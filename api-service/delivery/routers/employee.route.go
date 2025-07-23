package router

import (
	controllers "telegram-service/delivery/http"
	services "telegram-service/lib/employee"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// EmployeeControllerWrapper is a wrapper for the employee controller
type EmployeeControllerWrapper struct {
	controller *controllers.EmployeeController
}

// NewEmployeeControllerWrapper initializes the wrapper for the employee controller
func NewEmployeeControllerWrapper(ec *controllers.EmployeeController) *EmployeeControllerWrapper {
	return &EmployeeControllerWrapper{
		controller: ec,
	}
}

// SetupEmployeeRouter initializes and configures the router for employee-related routes
func SetupEmployeeRouter(r *mux.Router, service services.EmployeeServiceInterface) {
	// Initialize controller with EmployeeService
	controller := controllers.NewEmployeeController(service)

	// Create a new instance of the wrapper for the controller
	controllerWrapper := NewEmployeeControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Employee routes
	employee := api.PathPrefix("/employees").Subrouter()
	employee.HandleFunc("/list", controllerWrapper.controller.EmployeeListHandler).Methods("GET")
	employee.HandleFunc("/create", controllerWrapper.controller.EmployeeCreateHandler).Methods("POST")
	employee.HandleFunc("/update", controllerWrapper.controller.EmployeeUpdateHandler).Methods("PUT")
	employee.HandleFunc("/delete", controllerWrapper.controller.EmployeeDeleteHandler).Methods("DELETE")
	employee.HandleFunc("/{id}", controllerWrapper.controller.EmployeeByIDHandler).Methods("GET")
}
