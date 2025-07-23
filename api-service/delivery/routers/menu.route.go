package router

import (
	controllers "telegram-service/delivery/http"
	services "telegram-service/lib/menus"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// controllerWrapper is a wrapper for the user controller
type MenuControllerWrapper struct {
	controller *controllers.MenuController
}

// NewcontrollerWrapper initializes the wrapper for controller
func NewMenuControllerWrapper(uc *controllers.MenuController) *MenuControllerWrapper {
	return &MenuControllerWrapper{
		controller: uc,
	}
}

// SetupRouter initializes and configures the router
func SetupMenuRouter(r *mux.Router, service *services.MenuService) {
	// Initialize controller with UserService
	controller := controllers.NewMenuController(service)

	// Create a new instance of the wrapper for controller
	controllerWrapper := NewMenuControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// User routes

	menu := api.PathPrefix("/menus").Subrouter()
	menu.HandleFunc("/list-paginate", controllerWrapper.controller.MenuListHandler).Methods("GET")
	menu.HandleFunc("/create", controllerWrapper.controller.MenuCreateHandler).Methods("POST")
	menu.HandleFunc("/update", controllerWrapper.controller.MenuUpdateHandler).Methods("PUT")
	menu.HandleFunc("/delete", controllerWrapper.controller.MenuDeleteHandler).Methods("DELETE")
	menu.HandleFunc("/menu", controllerWrapper.controller.MenuListSubHandler).Methods("GET")
	menu.HandleFunc("/id/{id}", controllerWrapper.controller.MenuByIDHandler).Methods("GET")
	//HandleMenuByIdRequest
	menu.HandleFunc("/list", controllerWrapper.controller.MenuListAllWithChild).Methods("GET")
	menu.HandleFunc("/list-permission", controllerWrapper.controller.MenuListByTenant).Methods("GET")

	menu.HandleFunc("/roleMenuTreeselect/{id}", controllerWrapper.controller.MenuListAllWithChildLabel).Methods("GET")
	menu.HandleFunc("/list/{id}", controllerWrapper.controller.MenuListAllWithChildById).Methods("GET")
}
