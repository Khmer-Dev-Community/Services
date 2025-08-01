package router

import (
	// Your existing controller and service imports
	controllers "github.com/Khmer-Dev-Community/Services/api-service/delivery/http"
	services "github.com/Khmer-Dev-Community/Services/api-service/lib/users/services"
	"github.com/gin-gonic/gin"
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

func SetupRouter(r *gin.Engine, userService *services.UserService) { // Changed r to *gin.Engine
	// Initialize UserController with UserService
	userController := controllers.NewUserController(userService)
	userControllerWrapper := NewUserControllerWrapper(userController)
	api := r.Group("/api")            // Changed to Gin's Group method
	userRouter := api.Group("/users") // Changed to Gin's Group method

	// Define User routes using Gin's methods
	userRouter.GET("/list", userControllerWrapper.userController.UserListHandler)        // GET /api/users/list
	userRouter.GET("/list/:id", userControllerWrapper.userController.UserByIDHandler)    // GET /api/users/list/:id (Gin uses :id for path parameters)
	userRouter.POST("/create", userControllerWrapper.userController.UserCreateHandler)   // POST /api/users/create
	userRouter.PUT("/update", userControllerWrapper.userController.UserUpdateHandler)    // PUT /api/users/update
	userRouter.DELETE("/delete", userControllerWrapper.userController.UserDeleteHandler) // DELETE /api/users/delete
	userRouter.GET("/profile", userControllerWrapper.userController.UserProfileHandler)  // GET /api/users/profile

	// The Swagger documentation setup should be in your main InitRoutes function, not here.
	// r.PathPrefix("/swagger/").Handler(httpSwagger.Handler()) // REMOVED
}
