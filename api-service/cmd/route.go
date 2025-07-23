package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "telegram-service/cmd/docs"
	"telegram-service/config"
	"telegram-service/delivery/middleware"
	routers "telegram-service/delivery/routers"
	telegrambot "telegram-service/telegram/bot"
	"telegram-service/telegram/dtos"
	"telegram-service/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type SendMessageRequest struct {
	GroupID string `json:"groups"`
	Message string `json:"msg"`
	Phone   string `json:"phone"`
}

func InitRoutes(cfg config.Config, s *Services) http.Handler {
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(utils.AuthMiddlewareWithWhiteList(whitelist))

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler())

	// Public routes
	routers.SetupAuthRouter(r, s.Auth)

	// Protected routes
	routers.SetupRouter(r, s.User)
	routers.SetupRoleRouter(r, s.Role)
	routers.SetupPermissionRouter(r, s.Permission)
	routers.SetupMenuRouter(r, s.Menu)
	routers.SetupEmployeeRouter(r, s.Employee)
	routers.SetupDepartmentRouter(r, s.Department)
	r.HandleFunc("/forward/", func(w http.ResponseWriter, r *http.Request) {
		var req SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		var userBotToUse *telegrambot.BotAccount
		var foundBotName string
		//lookup group
		filter := dtos.TelegramGroupFilter{
			GroupName: &req.GroupID,
		}
		tgGroup, err := s.TgGroup.TelegramGroupServiceGetByName(&filter)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to find Telegram group for ID %d in database: %v", req.GroupID, err)
			utils.ErrorLog(errMsg, fmt.Sprintf("Get Group from request not found: %s", req.Phone))
			http.Error(w, errMsg, http.StatusNotFound) // Use 404 Not Found if group is expected but not found
			return
		}
		if req.Phone != "" {
			userBotToUse = s.ActiveUserBots[req.Phone]
			if userBotToUse != nil {
				foundBotName = req.Phone
			}
		}
		if userBotToUse == nil {
			utils.ErrorLog(w, fmt.Sprintf("No active userbot found for phone number: %s", req.Phone))
			http.Error(w, fmt.Sprintf("No active userbot found for phone number: %s", req.Phone), http.StatusNotFound)
			return
		}
		sendCtx, sendCancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer sendCancel()
		err = userBotToUse.SendMessageToGroup(sendCtx, int64(tgGroup.GroupID), req.Message)
		if err != nil {
			log.Printf("ERROR sending message via bot '%s' (phone: %s) to group %d: %v", foundBotName, req.Phone, req.GroupID, err)
			http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message sent successfully."))
		log.Printf("Message successfully sent via bot '%s' (phone: %s) to group %d.", foundBotName, req.Phone, req.GroupID)
	}).Methods("POST")
	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	return corsHandler.Handler(r)
}

func StartHTTPServer(port string, handler http.Handler) error {
	return http.ListenAndServe(":"+port, handler)
}
