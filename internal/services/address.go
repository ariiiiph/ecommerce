package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/models"
	"github.com/ariiiiph/ecommerce/internal/repositories"
)

type AddressService struct {
	addressRepo *repositories.AddressRepository
	userRepo    *repositories.UserRepository
}

func NewAddressService(addressRepo *repositories.AddressRepository, userRepo *repositories.UserRepository) *AddressService {
	return &AddressService{
		addressRepo: addressRepo,
		userRepo:    userRepo,
	}
}

func (s *AddressService) Create(ctx context.Context, req *dto.CreateAddressRequest) (*dto.AddressResponse, error) {
	if req.UserID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	title := strings.TrimSpace(req.Title)
	recipientName := strings.TrimSpace(req.RecipientName)
	phone := strings.TrimSpace(req.Phone)
	country := strings.TrimSpace(req.Country)
	city := strings.TrimSpace(req.City)
	addressLine := strings.TrimSpace(req.AddressLine)
	postalCode := strings.TrimSpace(req.PostalCode)

	if title == "" {
		return nil, apperror.BadRequest(
			"ADDRESS_TITLE_REQUIRED",
			"address title is required",
		)
	}

	if recipientName == "" {
		return nil, apperror.BadRequest(
			"RECIPIENT_NAME_REQUIRED",
			"recipient name is required",
		)
	}

	if phone == "" {
		return nil, apperror.BadRequest(
			"PHONE_REQUIRED",
			"phone is required",
		)
	}
	if country == "" {
		return nil, apperror.BadRequest(
			"COUNTRY_REQUIRED",
			"country is required",
		)
	}
	if city == "" {
		return nil, apperror.BadRequest(
			"CITY_REQUIRED",
			"city is required",
		)
	}

	if addressLine == "" {
		return nil, apperror.BadRequest(
			"ADDRESS_LINE_REQUIRED",
			"address line is required",
		)
	}

	if postalCode == "" {
		return nil, apperror.BadRequest(
			"POSTAL_CODE_REQUIRED",
			"postal code is required",
		)
	}

	_, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"USER_NOT_FOUND",
				"user not found",
			)
		}
		return nil, err
	}

	address := &models.Address{
		UserID:        req.UserID,
		Title:         title,
		RecipientName: recipientName,
		Phone:         phone,
		Country:       country,
		City:          city,
		AddressLine:   addressLine,
		PostalCode:    postalCode,
		IsDefault:     req.IsDefault,
	}

	if address.IsDefault {
		if err := s.addressRepo.CreateAsDefault(ctx, address); err != nil {
			return nil, err
		}
	} else {
		if err := s.addressRepo.Create(ctx, address); err != nil {
			return nil, err
		}
	}

	return toAddressResponse(address), nil

}

func (s *AddressService) GetByID(ctx context.Context, id int64) (*dto.AddressResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ADDRESS_ID",
			"invalid address id",
		)
	}

	address, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ADDRESS_NOT_FOUND",
				"address not found",
			)
		}
		return nil, err
	}
	return toAddressResponse(address), nil
}

func (s *AddressService) GetAllByUserID(ctx context.Context, userID int64) ([]*dto.AddressResponse, error) {
	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}
	addresses, err := s.addressRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.AddressResponse, 0, len(addresses))

	for _, address := range addresses {
		result = append(result, toAddressResponse(address))
	}

	return result, nil

}

func (s *AddressService) Update(ctx context.Context, id int64, req *dto.UpdateAddressRequest) (*dto.AddressResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ADDRESS_ID",
			"invalid address id",
		)
	}

	address, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ADDRESS_NOT_FOUND",
				"address not found",
			)
		}
		return nil, err
	}
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)

		if title == "" {
			return nil, apperror.BadRequest(
				"ADDRESS_TITLE_REQUIRED",
				"address title is required",
			)
		}
		address.Title = title
	}

	if req.RecipientName != nil {
		recipientName := strings.TrimSpace(*req.RecipientName)

		if recipientName == "" {
			return nil, apperror.BadRequest(
				"RECIPIENT_NAME_REQUIRED",
				"recipient name is required",
			)
		}
		address.RecipientName = recipientName
	}

	if req.Phone != nil {
		phone := strings.TrimSpace(*req.Phone)

		if phone == "" {
			return nil, apperror.BadRequest(
				"PHONE_REQUIRED",
				"phone is required",
			)
		}
		address.Phone = phone
	}

	if req.Country != nil {
		country := strings.TrimSpace(*req.Country)

		if country == "" {
			return nil, apperror.BadRequest(
				"COUNTRY_REQUIRED",
				"country is required",
			)
		}
		address.Country = country
	}

	if req.City != nil {
		city := strings.TrimSpace(*req.City)

		if city == "" {
			return nil, apperror.BadRequest(
				"CITY_REQUIRED",
				"city is required",
			)
		}
		address.City = city
	}

	if req.AddressLine != nil {
		addressLine := strings.TrimSpace(*req.AddressLine)

		if addressLine == "" {
			return nil, apperror.BadRequest(
				"ADDRESS_LINE_REQUIRED",
				"address line is required",
			)
		}
		address.AddressLine = addressLine
	}

	if req.PostalCode != nil {
		postalCode := strings.TrimSpace(*req.PostalCode)

		if postalCode == "" {
			return nil, apperror.BadRequest(
				"POSTAL_CODE_REQUIRED",
				"postal code is required",
			)
		}
		address.PostalCode = postalCode
	}

	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}

	if req.IsDefault != nil && *req.IsDefault {
		if err := s.addressRepo.UpdateAsDefault(ctx, address); err != nil {
			return nil, err
		}
	} else {
		if err := s.addressRepo.Update(ctx, address); err != nil {
			return nil, err
		}
	}
	return toAddressResponse(address), nil
}

func (s *AddressService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_ADDRESS_ID",
			"invalid address id",
		)
	}

	_, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"ADDRESS_NOT_FOUND",
				"address not found",
			)
		}
		return err
	}
	return s.addressRepo.Delete(ctx, id)
}

func toAddressResponse(address *models.Address) *dto.AddressResponse {
	return &dto.AddressResponse{
		ID:            address.ID,
		UserID:        address.UserID,
		Title:         address.Title,
		RecipientName: address.RecipientName,
		Phone:         address.Phone,
		Country:       address.Country,
		City:          address.City,
		AddressLine:   address.AddressLine,
		PostalCode:    address.PostalCode,
		IsDefault:     address.IsDefault,
		CreatedAt:     address.CreatedAt,
		UpdatedAt:     address.UpdatedAt,
	}
}
