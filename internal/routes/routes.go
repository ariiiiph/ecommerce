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
}
