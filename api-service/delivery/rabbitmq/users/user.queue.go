package mq

import (
	"encoding/json"
	"log"
	"net/http"
	"telegram-service/delivery/rabbitmq"
	"telegram-service/lib/users/dto"
	service "telegram-service/lib/users/services"
	"telegram-service/utils"

	"github.com/streadway/amqp"
)

type UserMQ struct {
	userService *service.UserService
}

func NewUserMQ(us *service.UserService) *UserMQ {
	return &UserMQ{userService: us}
}

func (uc *UserMQ) HandleMessage(msg amqp.Delivery, channel *amqp.Channel) {
	queueName := msg.RoutingKey

	switch queueName {
	case "user_list_queue":
		var filter utils.PaginationRequestDTO
		if err := json.Unmarshal(msg.Body, &filter); err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		var requestData struct {
			Page  int    `json:"page"`
			Limit int    `json:"limit"`
			Query string `json:"query"`
		}
		err := json.Unmarshal(msg.Body, &requestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		requestDto := utils.PaginationRequestDTO{
			Page:  requestData.Page,
			Limit: requestData.Limit,
			Query: requestData.Query,
		}
		userListResponse := uc.userService.GetUserList(requestDto)
		utils.SendPaginationResponse(userListResponse, userListResponse.Meta.Total, filter.Page, filter.Limit, msg.ReplyTo, channel)
	case "user_byId_queue":
		log.Printf("Receiced GetById: %v", msg.Body)
		var requestData struct {
			ID int `json:"id"`
		}
		err := json.Unmarshal(msg.Body, &requestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.userService.GetUserByID(uint(requestData.ID))
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "user_create_queue":

		userListResponse, err := uc.userService.CreateUser(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "user_update_queue":
		var rquestData dto.UserUpdateDTO
		err := json.Unmarshal(msg.Body, &rquestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.userService.UpdateUser(uint(rquestData.ID), msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "user_delete_queue":
		var rquestData dto.UserUpdateDTO
		err := json.Unmarshal(msg.Body, &rquestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.userService.DeleteUser(uint(rquestData.ID))
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	}

}

func (cc *UserMQ) DeclareAndHandleExchanges(channel *amqp.Channel, userService_ *service.UserService) {
	rabbitmq.ExchnageDeclareAndHandleQueues(channel, "user_exchange", "direct", []string{
		"user_list_queue", "user_byId_queue", "user_create_queue", "user_update_queue", "user_delete_queue",
	},
		&UserMQ{
			userService: userService_,
		},
	)
}
