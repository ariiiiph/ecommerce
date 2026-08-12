package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ariiiiph/ecommerce/internal/config"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/models"
	"github.com/ariiiiph/ecommerce/internal/repositories"
	"github.com/ariiiiph/ecommerce/internal/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type AuthService struct {
	userRepo         *repositories.UserRepository
	roleRepo         *repositories.RoleRepository
	refreshTokenRepo *repositories.RefreshTokenRepository
	jwtConfig        config.JWTConfig
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	roleRepo *repositories.RoleRepository,
	refreshTokenRepo *repositories.RefreshTokenRepository,
	jwtConfig config.JWTConfig,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		roleRepo:         roleRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtConfig:        jwtConfig,
	}
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && user != nil {
		return nil, ErrEmailAlreadyExists
	}

	role, err := s.roleRepo.FindByName(ctx, "customer")
	if err != nil {
		return nil, err
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user = &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		RoleID:       role.ID,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(
		s.jwtConfig,
		user.ID,
		user.Email,
		role.Name,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashRefreshToken(refreshToken),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshTokenModel); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  role.Name,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	role, err := s.roleRepo.FindByID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(
		s.jwtConfig,
		user.ID,
		user.Email,
		role.Name,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashRefreshToken(refreshToken),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshTokenModel); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  role.Name,
		},
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, rawToken string) (*dto.AuthResponse, error) {
	tokenHash := hashRefreshToken(rawToken)

	refreshToken, err := s.refreshTokenRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if refreshToken.RevokedAt != nil {
		return nil, ErrInvalidCredentials
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.FindByID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(
		s.jwtConfig,
		user.ID,
		user.Email,
		role.Name,
	)
	if err != nil {
		return nil, err
	}

	newRawToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	newRefreshToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashRefreshToken(newRawToken),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.refreshTokenRepo.Create(ctx, newRefreshToken); err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.SetReplacement(
		ctx,
		refreshToken.ID,
		newRefreshToken.ID,
	); err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.Revoke(
		ctx,
		refreshToken.ID,
	); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRawToken,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  role.Name,
		},
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	tokenHash := hashRefreshToken(rawToken)

	refreshToken, err := s.refreshTokenRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		return ErrInvalidCredentials
	}

	if refreshToken.RevokedAt != nil {
		return nil
	}

	return s.refreshTokenRepo.Revoke(ctx, refreshToken.ID)
}
