package router

import (
	controllers "github.com/Khmer-Dev-Community/Services/api-service/delivery/http"
	services "github.com/Khmer-Dev-Community/Services/api-service/lib/roles/services"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// controllerWrapper is a wrapper for the user controller
type RoleControllerWrapper struct {
	controller *controllers.RoleController
}

// NewcontrollerWrapper initializes the wrapper for controller
func NewRoleControllerWrapper(uc *controllers.RoleController) *RoleControllerWrapper {
	return &RoleControllerWrapper{
		controller: uc,
	}
}

// SetupRouter initializes and configures the router
func SetupRoleRouter(r *mux.Router, service *services.RoleService) {
	// Initialize controller with UserService
	controller := controllers.NewRoleController(service)

	// Create a new instance of the wrapper for controller
	controllerWrapper := NewRoleControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// User routes
	userRouter := api.PathPrefix("/roles").Subrouter()
	userRouter.HandleFunc("/list", controllerWrapper.controller.RoleListHandler).Methods("GET")

	userRouter.HandleFunc("/list/{id}", controllerWrapper.controller.RoleByIDHandler).Methods("GET")
	userRouter.HandleFunc("/create", controllerWrapper.controller.RoleCreateHandler).Methods("POST")
	userRouter.HandleFunc("/update", controllerWrapper.controller.RoleUpdateHandler).Methods("PUT")
	userRouter.HandleFunc("/delete", controllerWrapper.controller.RoleDeleteHandler).Methods("DELETE")
}
