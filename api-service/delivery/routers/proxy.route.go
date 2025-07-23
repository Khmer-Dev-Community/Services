package router

import (
	"net/http"
	controllers "telegram-service/delivery/http"
	services "telegram-service/lib/proxy/services" // Adjust import path if necessary

	"github.com/gorilla/mux"
)

type ProxyListControllerWrapper struct {
	controller *controllers.ProxyListController
}

// NewProxyListControllerWrapper initializes the wrapper for the proxy list controller
func NewProxyListControllerWrapper(plc *controllers.ProxyListController) *ProxyListControllerWrapper {
	return &ProxyListControllerWrapper{
		controller: plc,
	}
}

func SetupProxyListRouter(r *mux.Router, service services.ProxyListServicer) {
	controller := controllers.NewProxyListController(service)
	controllerWrapper := NewProxyListControllerWrapper(controller)
	api := r.PathPrefix("/api").Subrouter()
	proxylists := api.PathPrefix("/proxylists").Subrouter()
	proxylists.Handle("",
		http.HandlerFunc(controllerWrapper.controller.ProxyListControllerGetList),
	).Methods("GET")

	// GET a single proxy list by ID
	proxylists.Handle("/{id}",
		http.HandlerFunc(controllerWrapper.controller.ProxyListControllerGetByID),
	).Methods("GET")
	// POST create a new proxy list
	proxylists.Handle("",
		http.HandlerFunc(controllerWrapper.controller.ProxyListControllerCreate),
	).Methods("POST")
	// PUT update an existing proxy list
	proxylists.Handle("",
		http.HandlerFunc(controllerWrapper.controller.ProxyListControllerUpdate),
	).Methods("PUT")
	// DELETE soft delete a proxy list by ID
	proxylists.Handle("/{id}",
		http.HandlerFunc(controllerWrapper.controller.ProxyListControllerDelete),
	).Methods("DELETE")
	// DELETE hard delete a proxy list by ID (use with extreme caution, typically for admin roles)
	proxylists.Handle("/hard/{id}",
		http.HandlerFunc(controllerWrapper.controller.ProxyListControllerHardDelete),
	).Methods("DELETE")
	proxylists.Handle("/generate",
		http.HandlerFunc(controllerWrapper.controller.ProxyListControllerGenerateBulk),
	).Methods("POST")
}
