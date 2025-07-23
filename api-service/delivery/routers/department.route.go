package router

import (
	controllers "telegram-service/delivery/http"
	services "telegram-service/lib/department" // Adjust the import path as necessary

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// DepartmentControllerWrapper is a wrapper for the department controller
type DepartmentControllerWrapper struct {
	controller *controllers.DepartmentController
}

// NewDepartmentControllerWrapper initializes the wrapper for the department controller
func NewDepartmentControllerWrapper(dc *controllers.DepartmentController) *DepartmentControllerWrapper {
	return &DepartmentControllerWrapper{
		controller: dc,
	}
}

// SetupDepartmentRouter initializes and configures the router for department-related routes
func SetupDepartmentRouter(r *mux.Router, service services.DepartmentServiceInterface) {
	// Initialize controller with DepartmentService
	controller := controllers.NewDepartmentController(service)

	// Create a new instance of the wrapper for the controller
	controllerWrapper := NewDepartmentControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Department routes
	department := api.PathPrefix("/department").Subrouter()
	department.HandleFunc("/list", controllerWrapper.controller.DepartmentListHandler).Methods("GET")
	department.HandleFunc("/create", controllerWrapper.controller.DepartmentCreateHandler).Methods("POST")
	department.HandleFunc("/update", controllerWrapper.controller.DepartmentUpdateHandler).Methods("PUT")
	department.HandleFunc("/delete", controllerWrapper.controller.DepartmentDeleteHandler).Methods("DELETE")
	department.HandleFunc("/{id}", controllerWrapper.controller.DepartmentByIDHandler).Methods("GET")
}
