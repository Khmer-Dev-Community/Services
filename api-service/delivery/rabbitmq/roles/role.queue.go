package mq

import (
	"encoding/json"
	"log"
	"net/http"
	"telegram-service/delivery/rabbitmq"
	dto "telegram-service/lib/roles/dto"
	service "telegram-service/lib/roles/services"
	"telegram-service/utils"

	"github.com/streadway/amqp"
)

type RoleMQ struct {
	roleService *service.RoleService
}

func NewRoleMQ(us *service.RoleService) *RoleMQ {
	return &RoleMQ{roleService: us}
}

func (uc *RoleMQ) HandleMessage(msg amqp.Delivery, channel *amqp.Channel) {
	queueName := msg.RoutingKey

	switch queueName {
	case "role_list_queue":
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
		userListResponse := uc.roleService.HandleRoleListRequest(requestDto)
		utils.SendPaginationResponse(userListResponse, userListResponse.Meta.Total, filter.Page, filter.Limit, msg.ReplyTo, channel)
	case "role_byId_queue":
		log.Printf("Receiced GetById: %v", msg.Body)
		var requestData struct {
			ID int `json:"id"`
		}
		err := json.Unmarshal(msg.Body, &requestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.roleService.HandleRoleByIdRequest(int(requestData.ID))
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "role_create_queue":

		userListResponse, err := uc.roleService.HandleRoleCreateRequest(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "role_update_queue":
		var rquestData dto.RoleUpdateDTO
		err := json.Unmarshal(msg.Body, &rquestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.roleService.HandleRoleUpdateRequest(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "role_delete_queue":
		var rquestData dto.RoleUpdateDTO
		err := json.Unmarshal(msg.Body, &rquestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.roleService.HandleRoleDeleteRequest(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	}

}

func (cc *RoleMQ) DeclareAndHandleExchanges(channel *amqp.Channel, roleService_ *service.RoleService) {
	rabbitmq.ExchnageDeclareAndHandleQueues(channel, "role_exchange", "direct", []string{
		"role_list_queue", "role_byId_queue", "role_create_queue", "role_update_queue", "role_delete_queue",
	},
		&RoleMQ{
			roleService: roleService_,
		},
	)
}
