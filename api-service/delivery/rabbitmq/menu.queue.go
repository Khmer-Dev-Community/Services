package rabbitmq

import (
	"encoding/json"
	"log"
	service "telegram-service/lib/menus"
	"telegram-service/utils"

	"github.com/streadway/amqp"
)

type MenuMQ struct {
	menuService *service.MenuService
}

func NewMenuMQ(us *service.MenuService) *MenuMQ {
	return &MenuMQ{menuService: us}
}

func (uc *MenuMQ) HandleMessage(msg amqp.Delivery, channel *amqp.Channel) {
	queueName := msg.RoutingKey

	switch queueName {
	case "menu_list_queue":
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
		userListResponse := uc.menuService.HandleMenuListRequest(requestDto)
		utils.SendPaginationResponse(userListResponse, userListResponse.Meta.Total, filter.Page, filter.Limit, msg.ReplyTo, channel)

	}

}

func (cc *MenuMQ) MenuDeclareAndHandleExchanges(channel *amqp.Channel, menuService_ *service.MenuService) {
	ExchnageDeclareAndHandleQueues(channel, "role_exchange", "direct", []string{
		"role_list_queue", "role_byId_queue", "role_create_queue", "role_update_queue", "role_delete_queue",
	},
		&MenuMQ{
			menuService: menuService_,
		},
	)
}
