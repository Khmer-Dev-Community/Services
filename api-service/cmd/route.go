package main

import (
	_ "github.com/Khmer-Dev-Community/Services/api-service/cmd/docs"
	"github.com/Khmer-Dev-Community/Services/api-service/config"
	routers "github.com/Khmer-Dev-Community/Services/api-service/delivery/routers"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SendMessageRequest struct {
	GroupID string `json:"groups"`
	Message string `json:"msg"`
	Phone   string `json:"phone"`
}

var whitelist = map[string]bool{
	"/api/auth/login":                          true,
	"/api/auth/logout":                         true,
	"/api/auth02/register":                     true,
	"/api/auth02/login":                        true,
	"/api/posts/list":                          true,
	"/api/posts/v/:id":                         true,
	"/api/posts/p/:slug":                       true,
	"/api/posts/:slug":                         true,
	"/api/client/profile/:username":            true,
	"/api/auth02/github/login":                 true,
	"/api/auth02/logout":                       true,
	"/api/account/auth02/github/login":         true,
	"/api/auth02/github/callback":              true,
	"/api/swagger/index.html":                  true,
	"/swagger/index.html":                      true,
	"/swagger/swagger-ui-bundle.js":            true,
	"/swagger/swagger-ui.css":                  true,
	"/swagger/swagger-ui-standalone-preset.js": true,
	"/swagger/doc.json":                        true,
	"/swagger/favicon-32x32.png":               true,
	"/swagger/favicon-16x16.png":               true,
}

func InitRoutes(cfg config.Config, s *Services) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:3000", "http://192.168.50.102:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Use(utils.AuthMiddlewareWithWhiteList(whitelist))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routers.SetupRouter(r, s.User)
	routers.SetupRouterAuth02(r, s.Auth02, &cfg.Github)
	routers.SetupPostRouter(r, s.Posts)
	routers.SetupCommentRouter(r, s.Comments)
	return r
}
