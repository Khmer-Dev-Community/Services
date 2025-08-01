package main

import (
	"net/http"

	_ "github.com/Khmer-Dev-Community/Services/api-service/cmd/docs"
	"github.com/Khmer-Dev-Community/Services/api-service/config"
	routers "github.com/Khmer-Dev-Community/Services/api-service/delivery/routers"
	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SendMessageRequest struct {
	GroupID string `json:"groups"`
	Message string `json:"msg"`
	Phone   string `json:"phone"`
}

func InitRoutes(cfg config.Config, s *Services) *gin.Engine {
	r := gin.Default()

	// Middleware
	//r.Use(middleware.LoggingMiddleware)
	//r.Use(utils.AuthMiddlewareWithWhiteList(whitelist))
	//r.Use(utils.AuthMiddlewareWithWhiteList(whitelist))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//routers.SetupAuthRouter(r, s.Auth)

	// Protected routes
	routers.SetupRouter(r, s.User)
	//routers.SetupRoleRouter(r, s.Role)
	//routers.SetupPermissionRouter(r, s.Permission)
	//routers.SetupMenuRouter(r, s.Menu)
	routers.SetupRouterAuth02(r, s.Auth02)
	routers.SetupPostRouter(r, s.Posts)
	// CORS configuration
	r.Use(ginCors.New(ginCors.Config{
		AllowOrigins:     []string{"http://localhost", "*"},                                                                                // Be more specific in production, e.g., "http://localhost:3000"
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},                                             // Added HEAD and PATCH
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "X-Requested-With"}, // Added common headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	return r
}

func StartHTTPServer(port string, handler http.Handler) error {
	return http.ListenAndServe(":"+port, handler)
}
