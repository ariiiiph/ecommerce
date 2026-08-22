package routes

import (
	"net/http"

	"github.com/ariiiiph/ecommerce/internal/app"
	"github.com/ariiiiph/ecommerce/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(
	mux *http.ServeMux,
	deps *app.Dependencies,
) {
	//swager
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Auth routes
	mux.HandleFunc("POST /api/auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /api/auth/login", deps.AuthHandler.Login)
	mux.HandleFunc("POST /api/auth/refresh", deps.AuthHandler.Refresh)
	mux.Handle("POST /api/auth/logout", middleware.AuthMiddleware(deps.Config.JWT)(http.HandlerFunc(deps.AuthHandler.Logout)))

	// Category routes
	mux.Handle("POST /api/categories", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.CategoryHandler.Create))))
	mux.Handle("PUT /api/categories/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.CategoryHandler.Update))))
	mux.Handle("DELETE /api/categories/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.CategoryHandler.Delete))))
	mux.HandleFunc("GET /api/categories", deps.CategoryHandler.GetAll)
	mux.HandleFunc("GET /api/categories/{id}", deps.CategoryHandler.GetByID)
	mux.HandleFunc("GET /api/categories/{id}/children", deps.CategoryHandler.GetChildren)
	mux.HandleFunc("GET /api/categories/{id}/tree", deps.CategoryHandler.GetTree)

	//Brand routes
	mux.Handle("POST /api/brands", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.BrandHandler.Create))))
	mux.Handle("PUT /api/brands/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.BrandHandler.Update))))
	mux.Handle("DELETE /api/brands/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.BrandHandler.Delete))))
	mux.HandleFunc("GET /api/brands", deps.BrandHandler.GetAll)
	mux.HandleFunc("GET /api/brands/{id}", deps.BrandHandler.GetByID)

	// Product routes
	mux.Handle("POST /api/products", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.ProductHandler.Create))))
	mux.Handle("PUT /api/products/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.ProductHandler.Update))))
	mux.Handle("DELETE /api/products/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.ProductHandler.Delete))))
	mux.HandleFunc("GET /api/products", deps.ProductHandler.GetAll)
	mux.HandleFunc("GET /api/products/{id}", deps.ProductHandler.GetByID)

	// Product Variant routes
	mux.Handle("POST /api/product-variants", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.ProductVariantHandler.Create))))
	mux.Handle("PUT /api/product-variants/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.ProductVariantHandler.Update))))
	mux.Handle("DELETE /api/product-variants/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.ProductVariantHandler.Delete))))
	mux.HandleFunc("GET /api/product-variants/{id}", deps.ProductVariantHandler.GetByID)
	mux.HandleFunc("GET /api/products/{id}/variants", deps.ProductVariantHandler.GetAllByProductID)

	//Product Image routes
	mux.Handle("POST /api/product-images", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.ProductImageHandler.Create))))
	mux.Handle("PUT /api/product-images/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.ProductImageHandler.Update))))
	mux.Handle("DELETE /api/product-images/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.ProductImageHandler.Delete))))
	mux.HandleFunc("GET /api/product-images/{id}", deps.ProductImageHandler.GetByID)
	mux.HandleFunc("GET /api/products/{id}/images", deps.ProductImageHandler.GetAllByProductID)

	// Attribute routes
	mux.Handle("POST /api/attributes", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.AttributeHandler.Create))))
	mux.Handle("PUT /api/attributes/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.AttributeHandler.Update))))
	mux.Handle("DELETE /api/attributes/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.AttributeHandler.Delete))))
	mux.HandleFunc("GET /api/attributes/{id}", deps.AttributeHandler.GetByID)
	mux.HandleFunc("GET /api/attributes", deps.AttributeHandler.GetAll)

	// Attribute Value routes
	mux.Handle("POST /api/attribute-values", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.AttributeValueHandler.Create))))
	mux.Handle("DELETE /api/attribute-values/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.AttributeValueHandler.Delete))))
	mux.Handle("PUT /api/attribute-values/{id}", middleware.AuthMiddleware(deps.Config.JWT)(middleware.RequireRole("admin")(http.HandlerFunc(deps.AttributeValueHandler.Update))))
	mux.HandleFunc("GET /api/attribute-values/{id}", deps.AttributeValueHandler.GetByID)
	mux.HandleFunc("GET /api/attribute-values", deps.AttributeValueHandler.GetAll)

}
