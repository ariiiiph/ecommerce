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
	mux.HandleFunc("POST /api/auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /api/auth/login", deps.AuthHandler.Login)
	mux.HandleFunc("POST /api/auth/refresh", deps.AuthHandler.Refresh)
	mux.Handle("POST /api/auth/logout", middleware.AuthMiddleware(deps.Config.JWT)(http.HandlerFunc(deps.AuthHandler.Logout)))

	mux.Handle("/swagger/", httpSwagger.WrapHandler)
}
