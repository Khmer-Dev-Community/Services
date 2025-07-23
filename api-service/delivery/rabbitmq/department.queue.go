package rabbitmq

import (
	// Adjust the import path as necessary
	"encoding/json"
	"log"
	"net/http"
	service "telegram-service/lib/department" // Adjust the import path as necessary
	"telegram-service/utils"

	"github.com/streadway/amqp"
)

type DepartmentMQ struct {
	departmentService service.DepartmentServiceInterface
}

func NewDepartmentMQ(ds service.DepartmentServiceInterface) *DepartmentMQ {
	return &DepartmentMQ{departmentService: ds}
}

func (dc *DepartmentMQ) HandleMessage(msg amqp.Delivery, channel *amqp.Channel) {
	queueName := msg.RoutingKey

	switch queueName {
	case "department_list_queue":
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
		departmentListResponse := dc.departmentService.DepartmentServiceGetList(requestDto)
		utils.SendPaginationResponse(departmentListResponse, departmentListResponse.Meta.Total, filter.Page, filter.Limit, msg.ReplyTo, channel)
	case "department_byId_queue":

		departmentResponse, err := dc.departmentService.DepartmentServiceGetById(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to retrieve department", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(departmentResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "department_create_queue":

		departmentResponse, err := dc.departmentService.DepartmentServiceCreate(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to create department", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(departmentResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "department_update_queue":

		departmentResponse, err := dc.departmentService.DepartmentServiceUpdate(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to update department", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(departmentResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "department_delete_queue":

		departmentResponse, err := dc.departmentService.DepartmentServiceDelete(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to delete department", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(departmentResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	}
}

func (dc *DepartmentMQ) DepartmentDeclareAndHandleExchanges(channel *amqp.Channel, departmentService_ service.DepartmentServiceInterface) {
	ExchnageDeclareAndHandleQueues(channel, "department_exchange", "direct", []string{
		"department_list_queue", "department_byId_queue", "department_create_queue", "department_update_queue", "department_delete_queue",
	},
		&DepartmentMQ{
			departmentService: departmentService_,
		},
	)
}
