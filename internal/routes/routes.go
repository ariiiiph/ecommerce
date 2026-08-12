package routes

import (
	"net/http"

	"github.com/ariiiiph/ecommerce/internal/app"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Setup(application *app.App) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/register", application.Dependencies.AuthHandler.Register)
	mux.HandleFunc("POST /api/auth/login", application.Dependencies.AuthHandler.Login)
	mux.HandleFunc("POST /api/auth/refresh", application.Dependencies.AuthHandler.Refresh)
	mux.HandleFunc("POST /api/auth/logout", application.Dependencies.AuthHandler.Logout)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
