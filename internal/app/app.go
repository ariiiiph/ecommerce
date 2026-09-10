package app

import (
	"github.com/ariiiiph/ecommerce/internal/handlers"
	"github.com/ariiiiph/ecommerce/internal/repositories"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type App struct {
	Dependencies *Dependencies
}

func New(dependencies *Dependencies) *App {
	// Repositories
	userRepository := repositories.NewUserRepository(dependencies.DB)
	roleRepository := repositories.NewRoleRepository(dependencies.DB)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(dependencies.DB)
	categoryRepository := repositories.NewCategoryRepository(dependencies.DB)
	brandRepository := repositories.NewBrandRepository(dependencies.DB)
	productRepository := repositories.NewProductRepository(dependencies.DB)
	productVariantRepository := repositories.NewProductVariantRepository(dependencies.DB)
	productImageRepository := repositories.NewProductImageRepository(dependencies.DB)
	attributeRepository := repositories.NewAttributeRepository(dependencies.DB)
	attributeValueRepository := repositories.NewAttributeValueRepository(dependencies.DB)
	variantAttributeValueRepository := repositories.NewVariantAttributeValueRepository(dependencies.DB)
	inventoryRepository := repositories.NewInventoryRepository(dependencies.DB)
	addressRepository := repositories.NewAddressRepository(dependencies.DB)
	cartRepository := repositories.NewCartRepository(dependencies.DB)
	cartItemRepository := repositories.NewCartItemRepository(dependencies.DB)
	wishlistRepo := repositories.NewWishlistRepository(dependencies.DB)
	wishlistItemRepo := repositories.NewWishlistItemRepository(dependencies.DB)
	couponRepo := repositories.NewCouponRepository(dependencies.DB)
	orderRepo := repositories.NewOrderRepository(dependencies.DB)
	orderItemRepo := repositories.NewOrderItemRepository(dependencies.DB)
	couponUsageRepo := repositories.NewCouponUsageRepository(dependencies.DB)
	paymentRepo := repositories.NewPaymentRepository(dependencies.DB)
	reviewRepository := repositories.NewReviewRepository(dependencies.DB)

	// Services
	authService := services.NewAuthService(
		userRepository,
		roleRepository,
		refreshTokenRepository,
		dependencies.Config.JWT,
	)

	categoryService := services.NewCategoryService(
		categoryRepository,
	)

	brandService := services.NewBrandService(
		brandRepository,
	)

	productService := services.NewProductService(
		productRepository,
		brandRepository,
		categoryRepository,
	)

	productVariantService := services.NewProductVariantService(
		productVariantRepository,
		productRepository,
	)

	productImageService := services.NewProductImageService(
		productImageRepository,
		productVariantRepository,
		productRepository,
	)

	attributeService := services.NewAttributeService(
		attributeRepository,
	)
	attributeValueService := services.NewAttributeValueService(
		attributeValueRepository,
		attributeRepository,
	)

	variantAttributeValueService := services.NewVariantAttributeValueService(
		variantAttributeValueRepository,
		productVariantRepository,
		attributeValueRepository,
	)

	inventoryService := services.NewInventoryService(
		inventoryRepository,
		productVariantRepository,
	)

	addressService := services.NewAddressService(
		addressRepository,
		userRepository,
	)

	cartService := services.NewCartService(
		cartRepository,
		userRepository,
	)

	cartItemService := services.NewCartItemService(
		cartItemRepository,
		cartRepository,
		productVariantRepository,
		dependencies.Redis,
	)

	wishlistService := services.NewWishlistService(
		wishlistRepo,
		userRepository,
	)

	wishlistItemService := services.NewWishlistItemService(
		wishlistItemRepo,
		wishlistRepo,
		productRepository,
	)

	couponService := services.NewCouponService(
		couponRepo,
	)

	orderService := services.NewOrderService(
		dependencies.DB,
		orderRepo,
		orderItemRepo,
		couponUsageRepo,
		cartRepository,
		cartItemRepository,
		inventoryRepository,
		productVariantRepository,
		addressRepository,
		couponRepo,
	)

	paymentService := services.NewPaymentService(
		dependencies.DB,
		paymentRepo,
		orderRepo,
	)

	reviewService := services.NewReviewService(
		reviewRepository,
		productRepository,
	)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)

	categoryHandler := handlers.NewCategoryHandler(
		categoryService,
	)

	brandHandler := handlers.NewBrandHandler(
		brandService,
	)

	productHandler := handlers.NewProductHandler(
		productService,
	)

	productVariantHandler := handlers.NewProductVariantHandler(productVariantService)

	productImageHandler := handlers.NewProductImageHandler(productImageService)

	attributeHandler := handlers.NewAttributeHandler(
		attributeService,
	)

	attributeValueHandler := handlers.NewAttributeValueHandler(
		attributeValueService,
	)

	variantAttributeValueHandler := handlers.NewVariantAttributeValueHandler(
		variantAttributeValueService,
	)

	inventoryHandler := handlers.NewInventoryHandler(
		inventoryService,
	)

	addressHandler := handlers.NewAddressHandler(addressService)

	cartHandler := handlers.NewCartHandler(cartService)

	cartItemHandler := handlers.NewCartItemHandler(cartItemService)

	wishlistHandler := handlers.NewWishlistHandler(wishlistService)

	wishlistItemHandler := handlers.NewWishlistItemHandler(wishlistItemService)

	couponHandler := handlers.NewCouponHandler(
		couponService,
	)
	orderHandler := handlers.NewOrderHandler(
		orderService,
	)

	paymentHandler := handlers.NewPaymentHandler(paymentService)

	reviewHandler := handlers.NewReviewHandler(
		reviewService,
	)

	// Dependencies
	dependencies.AuthHandler = authHandler
	dependencies.CategoryHandler = categoryHandler
	dependencies.BrandHandler = brandHandler
	dependencies.ProductHandler = productHandler
	dependencies.ProductVariantHandler = productVariantHandler
	dependencies.ProductImageHandler = productImageHandler
	dependencies.AttributeHandler = attributeHandler
	dependencies.AttributeValueHandler = attributeValueHandler
	dependencies.VariantAttributeValueHandler = variantAttributeValueHandler
	dependencies.InventoryHandler = inventoryHandler
	dependencies.AddressHandler = addressHandler
	dependencies.CartHandler = cartHandler
	dependencies.CartItemHandler = cartItemHandler
	dependencies.WishlistHandler = wishlistHandler
	dependencies.WishlistItemHandler = wishlistItemHandler
	dependencies.CouponHandler = couponHandler
	dependencies.OrderHandler = orderHandler
	dependencies.PaymentHandler = paymentHandler
	dependencies.ReviewHandler = reviewHandler

	return &App{
		Dependencies: dependencies,
	}
}
