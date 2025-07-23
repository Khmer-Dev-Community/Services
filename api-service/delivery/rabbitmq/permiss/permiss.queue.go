package mq

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Khmer-Dev-Community/Services/api-service/delivery/rabbitmq"
	dto "github.com/Khmer-Dev-Community/Services/api-service/lib/permission/dto"
	service "github.com/Khmer-Dev-Community/Services/api-service/lib/permission/service"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"

	"github.com/streadway/amqp"
)

type PermissMQ struct {
	permisionService *service.PermissionService
}

func NewPermissMQ(us *service.PermissionService) *PermissMQ {
	return &PermissMQ{permisionService: us}
}

func (uc *PermissMQ) HandleMessage(msg amqp.Delivery, channel *amqp.Channel) {
	queueName := msg.RoutingKey

	switch queueName {
	case "permission_list_queue":
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
		userListResponse := uc.permisionService.HandlePermissionListRequest(requestDto)
		utils.SendPaginationResponse(userListResponse, userListResponse.Meta.Total, filter.Page, filter.Limit, msg.ReplyTo, channel)
	case "permission_findOne_queue":
		log.Printf("Receiced GetById: %v", msg.Body)
		var requestData struct {
			ID int `json:"id"`
		}
		err := json.Unmarshal(msg.Body, &requestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.permisionService.HandlePermissionByIdRequest(uint(requestData.ID))
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "permission_create_queue":

		userListResponse, err := uc.permisionService.HandlePermissionCreateRequest(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "permission_update_queue":
		var rquestData dto.PermissionUpdateDTO
		err := json.Unmarshal(msg.Body, &rquestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.permisionService.HandlePermissionUpdateRequest(uint(rquestData.ID), msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "permission_delete_queue":
		var rquestData dto.PermissionUpdateDTO
		err := json.Unmarshal(msg.Body, &rquestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.permisionService.HandlePermissionDeleteRequest(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "permission_role_update_queue":
		var rquestData dto.PermissionUpdateDTO
		err := json.Unmarshal(msg.Body, &rquestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.permisionService.UpdateRolePermissionByRoleID(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "permission_role_queue":
		var rquestData dto.PermissionUpdateDTO
		err := json.Unmarshal(msg.Body, &rquestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		userListResponse, err := uc.permisionService.HandleRolePermissionById(rquestData.ID)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to unmarshal JSON", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(userListResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	}

}

func (cc *PermissMQ) DeclareAndHandleExchanges(channel *amqp.Channel, permisionService_ *service.PermissionService) {
	rabbitmq.ExchnageDeclareAndHandleQueues(channel, "permission_exchange", "direct", []string{
		"permission_list_queue", "permission_byId_queue",
		"permission_create_queue", "permission_update_queue",
		"permission_delete_queue", "permission_role_update_queue", "permission_role_queue",
	},
		&PermissMQ{
			permisionService: permisionService_,
		},
	)
}
