package gateway

import (
	"fmt"
	"websocket-service/utils"

	"github.com/go-playground/validator"
	socketio "github.com/googollee/go-socket.io"
)

// handleErrorAndRespond is a helper to centralize error logging and client response
func handleErrorAndRespond(s socketio.Conn, eventName string, data map[string]interface{}, err error, msg string, status string) {
	utils.ErrorLog(data, fmt.Sprintf("[Websocket] Error processing event '%s': %s - %v", eventName, msg, err))
	response := map[string]interface{}{
		"status":  status,
		"message": msg,
	}
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorsDetail []string
			for _, fieldErr := range validationErrors {
				errorsDetail = append(errorsDetail, fmt.Sprintf("Field '%s' failed on '%s' tag with value '%v'", fieldErr.Field(), fieldErr.Tag(), fieldErr.Value()))
			}
			response["errors"] = errorsDetail
		} else {
			response["errors"] = []string{err.Error()}
		}
	}
	s.Emit(fmt.Sprintf("ws:%sResponse", eventName), response) // Emit a response based on the original event name
}
