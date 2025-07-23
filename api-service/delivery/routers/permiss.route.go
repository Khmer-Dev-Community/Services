package router

import (
	controllers "telegram-service/delivery/http"
	services "telegram-service/lib/permission/service"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// controllerWrapper is a wrapper for the user controller
type PermissControllerWrapper struct {
	controller *controllers.PermissController
}

// NewcontrollerWrapper initializes the wrapper for controller
func NewPermissControllerWrapper(uc *controllers.PermissController) *PermissControllerWrapper {
	return &PermissControllerWrapper{
		controller: uc,
	}
}

// SetupRouter initializes and configures the router
func SetupPermissionRouter(r *mux.Router, service *services.PermissionService) {
	// Initialize controller with UserService
	controller := controllers.NewPermissController(service)

	// Create a new instance of the wrapper for controller
	controllerWrapper := NewPermissControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// User routes
	permissRouter := api.PathPrefix("/permiss").Subrouter()
	permissRouter.HandleFunc("/list", controllerWrapper.controller.PermissionListHandler).Methods("GET")
	permissRouter.HandleFunc("/list/{id}", controllerWrapper.controller.PermissionByIDHandler).Methods("GET")
	permissRouter.HandleFunc("/create", controllerWrapper.controller.PermissionCreateHandler).Methods("POST")
	permissRouter.HandleFunc("/update/{id}", controllerWrapper.controller.PermissionRoleByIDHandler).Methods("PUT")
	permissRouter.HandleFunc("/update", controllerWrapper.controller.PermissionUpdateHandler).Methods("PUT")
	permissRouter.HandleFunc("/delete", controllerWrapper.controller.PermissionDeleteHandler).Methods("DELETE")
	permissRouter.HandleFunc("/rolesacess", controllerWrapper.controller.PermissionRoleUpdateHandler).Methods("POST") // Update Role Permission

	permissRouter.HandleFunc("/info", controllerWrapper.controller.PermissionRoleByIDHandler).Methods("GET")
}
