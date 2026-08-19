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
	userRepository := repositories.NewUserRepository(dependencies.DB)
	roleRepository := repositories.NewRoleRepository(dependencies.DB)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(dependencies.DB)
	categoryRepository := repositories.NewCategoryRepository(dependencies.DB)
	brandRepository := repositories.NewBrandRepository(dependencies.DB)

	authService := services.NewAuthService(
		userRepository,
		roleRepository,
		refreshTokenRepository,
		dependencies.Config.JWT,
	)

	authHandler := handlers.NewAuthHandler(authService)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	brandService := services.NewBrandService(brandRepository)
	brandHandler := handlers.NewBrandHandler(brandService)

	dependencies.AuthHandler = authHandler
	dependencies.CategoryHandler = categoryHandler
	dependencies.BrandHandler = brandHandler

	return &App{
		Dependencies: dependencies,
	}
}
