package main

import (
	"net/http"

	_ "github.com/Khmer-Dev-Community/Services/api-service/cmd/docs"
	"github.com/Khmer-Dev-Community/Services/api-service/config"
	"github.com/Khmer-Dev-Community/Services/api-service/delivery/middleware"
	routers "github.com/Khmer-Dev-Community/Services/api-service/delivery/routers"

	"github.com/Khmer-Dev-Community/Services/api-service/utils"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type SendMessageRequest struct {
	GroupID string `json:"groups"`
	Message string `json:"msg"`
	Phone   string `json:"phone"`
}

func InitRoutes(cfg config.Config, s *Services) http.Handler {
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(utils.AuthMiddlewareWithWhiteList(whitelist))

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler())

	// Public routes
	routers.SetupAuthRouter(r, s.Auth)

	// Protected routes
	routers.SetupRouter(r, s.User)
	routers.SetupRoleRouter(r, s.Role)
	routers.SetupPermissionRouter(r, s.Permission)
	routers.SetupMenuRouter(r, s.Menu)
	routers.SetupRouterAuth02(r, s.Auth02)
	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	return corsHandler.Handler(r)
}

func StartHTTPServer(port string, handler http.Handler) error {
	return http.ListenAndServe(":"+port, handler)
}
