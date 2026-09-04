package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/models"
	"github.com/ariiiiph/ecommerce/internal/repositories"
	"github.com/shopspring/decimal"
)

type OrderService struct {
	db                 *sql.DB
	orderRepo          *repositories.OrderRepository
	orderItemRepo      *repositories.OrderItemRepository
	couponUsageRepo    *repositories.CouponUsageRepository
	cartRepo           *repositories.CartRepository
	cartItemRepo       *repositories.CartItemRepository
	inventoryRepo      *repositories.InventoryRepository
	productVariantRepo *repositories.ProductVariantRepository
	addressRepo        *repositories.AddressRepository
	couponRepo         *repositories.CouponRepository
}
type checkoutItem struct {
	CartItem    *models.CartItem
	Variant     *models.ProductVariant
	ProductName string
	UnitPrice   decimal.Decimal
	TotalPrice  decimal.Decimal
}

func NewOrderService(
	db *sql.DB,
	orderRepo *repositories.OrderRepository,
	orderItemRepo *repositories.OrderItemRepository,
	couponUsageRepo *repositories.CouponUsageRepository,
	cartRepo *repositories.CartRepository,
	cartItemRepo *repositories.CartItemRepository,
	inventoryRepo *repositories.InventoryRepository,
	productVariantRepo *repositories.ProductVariantRepository,
	addressRepo *repositories.AddressRepository,
	couponRepo *repositories.CouponRepository,
) *OrderService {
	return &OrderService{
		db:                 db,
		orderRepo:          orderRepo,
		orderItemRepo:      orderItemRepo,
		couponUsageRepo:    couponUsageRepo,
		cartRepo:           cartRepo,
		cartItemRepo:       cartItemRepo,
		inventoryRepo:      inventoryRepo,
		productVariantRepo: productVariantRepo,
		addressRepo:        addressRepo,
		couponRepo:         couponRepo,
	}
}

func (s *OrderService) Create(ctx context.Context, userID int64, req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	if req == nil {
		return nil, apperror.BadRequest(
			"INVALID_REQUEST_BODY",
			"request body is required",
		)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	cartRepo := repositories.NewCartRepository(tx)

	cart, err := cartRepo.GetByUserIDForUpdate(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CART_NOT_FOUND",
				"cart not found",
			)
		}
		return nil, err
	}
	cartItemRepo := repositories.NewCartItemRepository(tx)

	cartItems, err := cartItemRepo.GetAllByCartID(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, apperror.BadRequest(
			"EMPTY_CART",
			"cart is empty",
		)
	}

	sort.Slice(cartItems, func(i, j int) bool {
		return cartItems[i].VariantID < cartItems[j].VariantID
	})

	productVariantRepo := repositories.NewProductVariantRepository(tx)
	inventoryRepo := repositories.NewInventoryRepository(tx)

	checkoutItems := make([]checkoutItem, 0, len(cartItems))

	subtotal := decimal.Zero

	for _, item := range cartItems {
		if item.Quantity <= 0 {
			return nil, apperror.BadRequest(
				"INVALID_CART_ITEM_QUANTITY",
				"cart item quantity must be greater than zero",
			)
		}

		variantResult, err := productVariantRepo.GetByIDWithProduct(ctx, item.VariantID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.NotFound(
					"PRODUCT_VARIANT_NOT_FOUND",
					"product variant not found",
				)
			}
			return nil, err
		}

		inventory, err := inventoryRepo.GetByVariantIDForUpdate(ctx, item.VariantID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.NotFound(
					"INVENTORY_NOT_FOUND",
					"inventory not found",
				)
			}
			return nil, err
		}
		availableQuantity := inventory.Quantity - inventory.ReservedQuantity

		if item.Quantity > availableQuantity {
			return nil, apperror.BadRequest(
				"INSUFFICIENT_STOCK",
				"insufficient stock for product variant",
			)
		}

		unitPrice := decimal.NewFromFloat(variantResult.Variant.Price)

		if variantResult.Variant.DiscountPrice != nil {
			discountPrice := decimal.NewFromFloat(
				*variantResult.Variant.DiscountPrice,
			)
			if discountPrice.LessThan(unitPrice) {
				unitPrice = discountPrice
			}
		}

		totalPrice := unitPrice.Mul(
			decimal.NewFromInt(int64(item.Quantity)),
		)
		subtotal = subtotal.Add(totalPrice)

		checkoutItems = append(checkoutItems, checkoutItem{
			CartItem:    item,
			Variant:     variantResult.Variant,
			ProductName: variantResult.ProductName,
			UnitPrice:   unitPrice,
			TotalPrice:  totalPrice,
		})
	}

	addressRepo := repositories.NewAddressRepository(tx)

	var address *models.Address

	if req.AddressID != nil {
		if *req.AddressID <= 0 {
			return nil, apperror.BadRequest(
				"INVALID_ADDRESS_ID",
				"invalid address id",
			)
		}

		address, err = addressRepo.GetByID(ctx, *req.AddressID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.NotFound(
					"ADDRESS_NOT_FOUND",
					"address not found",
				)
			}
			return nil, err
		}
		if address.UserID != userID {
			return nil, apperror.NotFound(
				"ADDRESS_NOT_FOUND",
				"address not found",
			)
		}
	} else {
		address, err = addressRepo.GetDefaultByUserID(ctx, userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.BadRequest(
					"DEFAULT_ADDRESS_NOT_FOUND",
					"default address not found",
				)
			}
			return nil, err
		}
	}
	var coupon *models.Coupon
	discountAmount := decimal.Zero
	shippingAmount := decimal.Zero
	totalAmount := subtotal

	if req.CouponCode != nil {
		code := strings.TrimSpace(*req.CouponCode)

		if code == "" {
			return nil, apperror.BadRequest(
				"COUPON_CODE_REQUIRED",
				"coupon code is required",
			)
		}

		couponRepo := repositories.NewCouponRepository(tx)

		coupon, err = couponRepo.GetByCodeForUpdate(ctx, code)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.NotFound(
					"COUPON_NOT_FOUND",
					"coupon not found",
				)
			}
			return nil, err
		}
		if err := validateCoupon(coupon, subtotal, time.Now()); err != nil {
			return nil, err
		}

		if coupon.DiscountType == "percentage" {
			discountAmount = subtotal.
				Mul(coupon.DiscountValue).
				Div(decimal.NewFromInt(100))
		} else {
			discountAmount = coupon.DiscountValue
		}

		if coupon.MaximumDiscount != nil &&
			discountAmount.GreaterThan(*coupon.MaximumDiscount) {
			discountAmount = *coupon.MaximumDiscount
		}

		if discountAmount.GreaterThan(subtotal) {
			discountAmount = subtotal
		}

		totalAmount = subtotal.Sub(discountAmount).Add(shippingAmount)
	}
	orderNumber := generateOrderNumber()

	order := &models.Order{
		UserID:                userID,
		AddressID:             &address.ID,
		Status:                "pending",
		Subtotal:              subtotal,
		DiscountAmount:        discountAmount,
		ShippingAmount:        shippingAmount,
		TotalAmount:           totalAmount,
		ShippingRecipientName: address.RecipientName,
		ShippingPhone:         address.Phone,
		ShippingCountry:       address.Country,
		ShippingCity:          address.City,
		ShippingAddressLine:   address.AddressLine,
		ShippingPostalCode:    address.PostalCode,
		OrderNumber:           orderNumber,
	}

	if coupon != nil {
		order.CouponID = &coupon.ID
	}

	orderRepo := repositories.NewOrderRepository(tx)

	if err := orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	orderItemRepo := repositories.NewOrderItemRepository(tx)

	orderItems := make([]*models.OrderItem, 0, len(checkoutItems))

	for _, item := range checkoutItems {
		orderItem := &models.OrderItem{
			OrderID:     order.ID,
			VariantID:   item.Variant.ID,
			ProductName: item.ProductName,
			SKU:         item.Variant.SKU,
			Quantity:    item.CartItem.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
		}

		if err := orderItemRepo.Create(ctx, orderItem); err != nil {
			return nil, err
		}
	}

	for _, item := range checkoutItems {
		if err := inventoryRepo.DecreaseQuantity(
			ctx,
			item.Variant.ID,
			item.CartItem.Quantity,
		); err != nil {
			return nil, err
		}
	}
	if coupon != nil {
		couponUsage := &models.CouponUsage{
			CouponID: coupon.ID,
			UserID:   userID,
			OrderID:  order.ID,
		}
		couponUsageRepo := repositories.NewCouponUsageRepository(tx)

		if err := couponUsageRepo.Create(ctx, couponUsage); err != nil {
			return nil, err
		}
	}
	if coupon != nil {
		couponRepo := repositories.NewCouponRepository(tx)

		if err := couponRepo.IncrementUsedCount(ctx, coupon.ID); err != nil {
			return nil, err
		}
	}
	if err := cartItemRepo.DeleteByCartID(ctx, cart.ID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return toOrderResponse(order, orderItems), nil
}

func (s *OrderService) GetByID(ctx context.Context, id int64) (*dto.OrderResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ORDER_ID",
			"invalid order id",
		)
	}

	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ORDER_NOT_FOUND",
				"order not found",
			)
		}

		return nil, err
	}

	orderItems, err := s.orderItemRepo.GetAllByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return toOrderResponse(order, orderItems), nil
}

func (s *OrderService) GetAllByUserID(ctx context.Context, userID int64) ([]*dto.OrderResponse, error) {
	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	orders, err := s.orderRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.OrderResponse, 0, len(orders))

	for _, order := range orders {
		orderItems, err := s.orderItemRepo.GetAllByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, toOrderResponse(order, orderItems))
	}

	return result, nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, id int64, status string) (*dto.OrderResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ORDER_ID",
			"invalid order id",
		)
	}

	status = strings.TrimSpace(status)

	if status == "" {
		return nil, apperror.BadRequest(
			"ORDER_STATUS_REQUIRED",
			"order status is required",
		)
	}

	allowedStatuses := map[string]bool{
		"pending":    true,
		"confirmed":  true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
		"refunded":   true,
	}

	if !allowedStatuses[status] {
		return nil, apperror.BadRequest(
			"INVALID_ORDER_STATUS",
			"invalid order status",
		)
	}

	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ORDER_NOT_FOUND",
				"order not found",
			)
		}

		return nil, err
	}

	if err := s.orderRepo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}

	order.Status = status

	orderItems, err := s.orderItemRepo.GetAllByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return toOrderResponse(order, orderItems), nil
}

func generateOrderNumber() string {
	return fmt.Sprintf(
		"ORD-%d",
		time.Now().UnixNano(),
	)
}

func toOrderResponse(order *models.Order, orderItems []*models.OrderItem) *dto.OrderResponse {
	items := make([]dto.OrderItemResponse, 0, len(orderItems))

	for _, item := range orderItems {
		items = append(items, dto.OrderItemResponse{
			ID:          item.ID,
			OrderID:     item.OrderID,
			VariantID:   item.VariantID,
			ProductName: item.ProductName,
			SKU:         item.SKU,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
			CreatedAt:   item.CreatedAt,
		})
	}

	return &dto.OrderResponse{
		ID:                    order.ID,
		OrderNumber:           order.OrderNumber,
		UserID:                order.UserID,
		AddressID:             order.AddressID,
		CouponID:              order.CouponID,
		Status:                order.Status,
		Subtotal:              order.Subtotal,
		DiscountAmount:        order.DiscountAmount,
		ShippingAmount:        order.ShippingAmount,
		TotalAmount:           order.TotalAmount,
		ShippingRecipientName: order.ShippingRecipientName,
		ShippingPhone:         order.ShippingPhone,
		ShippingCountry:       order.ShippingCountry,
		ShippingCity:          order.ShippingCity,
		ShippingAddressLine:   order.ShippingAddressLine,
		ShippingPostalCode:    order.ShippingPostalCode,
		Items:                 items,
		CreatedAt:             order.CreatedAt,
		UpdatedtAt:            order.UpdatedAt,
	}
}
