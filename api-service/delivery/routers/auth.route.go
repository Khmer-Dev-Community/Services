package router

import (
	controllers "github.com/Khmer-Dev-Community/Services/api-service/auth"
	services "github.com/Khmer-Dev-Community/Services/api-service/auth"

	"github.com/gin-gonic/gin"
)

type AuthControllerWrapper struct {
	authController *controllers.AuthController
}

func NewAuthControllerWrapper(ac *controllers.AuthController) *AuthControllerWrapper {
	return &AuthControllerWrapper{
		authController: ac,
	}
}

// SetupAuthRouter initializes the Gin router with auth routes
func SetupAuthRouter(r *gin.Engine, authService *services.AuthService) {
	authController := controllers.NewAuthController(authService)
	authWrapper := NewAuthControllerWrapper(authController)

	// Auth routes group
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/login", authWrapper.authController.Login)
		authGroup.POST("/logout", authWrapper.authController.Logout)
	}
}
