package router

import (
	_ "github.com/Khmer-Dev-Community/Services/api-service/cmd/docs"
	controllers "github.com/Khmer-Dev-Community/Services/api-service/delivery/http"
	services "github.com/Khmer-Dev-Community/Services/api-service/lib/users/services"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// UserControllerWrapper is a wrapper for the user controller
type UserControllerWrapper struct {
	userController *controllers.UserController
}

// NewUserControllerWrapper initializes the wrapper for UserController
func NewUserControllerWrapper(uc *controllers.UserController) *UserControllerWrapper {
	return &UserControllerWrapper{
		userController: uc,
	}
}

// SetupRouter initializes and configures the router
func SetupRouter(r *mux.Router, userService *services.UserService) {
	// Initialize UserController with UserService
	userController := controllers.NewUserController(userService)

	// Create a new instance of the wrapper for UserController
	userControllerWrapper := NewUserControllerWrapper(userController)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler())

	// User routes
	userRouter := api.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/list", userControllerWrapper.userController.UserListHandler).Methods("GET")
	userRouter.HandleFunc("/list/{id}", userControllerWrapper.userController.UserByIDHandler).Methods("GET")
	userRouter.HandleFunc("/create", userControllerWrapper.userController.UserCreateHandler).Methods("POST")
	userRouter.HandleFunc("/update", userControllerWrapper.userController.UserUpdateHandler).Methods("PUT")
	userRouter.HandleFunc("/delete", userControllerWrapper.userController.UserDeleteHandler).Methods("DELETE")
	userRouter.HandleFunc("/profile", userControllerWrapper.userController.UserProfileHandler).Methods("GET")

}
