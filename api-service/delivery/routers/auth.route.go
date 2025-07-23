package router

import (
	controllers "telegram-service/auth"
	services "telegram-service/auth"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// UserControllerWrapper is a wrapper for the user controller
type AuthControllerWrapper struct {
	authController *controllers.AuthController
}

// NewUserControllerWrapper initializes the wrapper for UserController
func NewAuthControllerWrapper(uc *controllers.AuthController) *AuthControllerWrapper {
	return &AuthControllerWrapper{
		authController: uc,
	}
}

// SetupRouter initializes and configures the router
func SetupAuthRouter(r *mux.Router, authService *services.AuthService) {
	// Initialize UserController with UserService
	authController := controllers.NewAuthController(authService)

	// Create a new instance of the wrapper for UserController
	authControllerWrapper := NewAuthControllerWrapper(authController)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// User routes
	userRouter := api.PathPrefix("/auth").Subrouter()
	userRouter.HandleFunc("/login", authControllerWrapper.authController.Login).Methods("POST")
	userRouter.HandleFunc("/logout", authControllerWrapper.authController.Logout).Methods("GET")

}
