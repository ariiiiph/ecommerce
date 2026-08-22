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

	// Dependencies
	dependencies.AuthHandler = authHandler
	dependencies.CategoryHandler = categoryHandler
	dependencies.BrandHandler = brandHandler
	dependencies.ProductHandler = productHandler
	dependencies.ProductVariantHandler = productVariantHandler
	dependencies.ProductImageHandler = productImageHandler
	dependencies.AttributeHandler = attributeHandler
	dependencies.AttributeValueHandler = attributeValueHandler

	return &App{
		Dependencies: dependencies,
	}
}
