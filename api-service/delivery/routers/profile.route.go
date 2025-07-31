package router

import (
	// For health check/root route JSON

	// Import your swagger docs if you have them generated for mux
	_ "github.com/Khmer-Dev-Community/Services/api-service/cmd/docs"

	"github.com/gorilla/mux"

	"github.com/Khmer-Dev-Community/Services/api-service/auth02" // The auth02 package (controller)
	// The new oauth_config package
)

// ClientProfileControllerWrapper is a wrapper for the client auth controller
type ClientProfileControllerWrapper struct {
	ClientProfileController *auth02.ClientAuthController
}

func NewClientProfileControllerWrapper(cac *auth02.ClientAuthController) *ClientProfileControllerWrapper {
	return &ClientProfileControllerWrapper{
		ClientProfileController: cac,
	}
}

func SetupRouterClientAuth(r *mux.Router, ClientProfileCtrl *auth02.ClientAuthController) { // <--- CHANGED SIGNATURE

	ClientProfileControllerWrapper := NewClientProfileControllerWrapper(ClientProfileCtrl)
	// Define API routes
	api := r.PathPrefix("/api").Subrouter()

	ClientProfileRouter := api.PathPrefix("/kdc").Subrouter()
	ClientProfileRouter.HandleFunc("/profile", ClientProfileControllerWrapper.ClientProfileController.GetClientProfile).Methods("GET")
	ClientProfileRouter.HandleFunc("/profile", ClientProfileControllerWrapper.ClientProfileController.UpdateClientProfile).Methods("PUT")

}
