package rabbitmq

import (
	"encoding/json"
	"log"
	"net/http"
	dto "telegram-service/lib/employee"
	service "telegram-service/lib/employee"
	"telegram-service/utils"

	"github.com/streadway/amqp"
)

type EmployeeMQ struct {
	employeeService service.EmployeeServiceInterface
}

func NewEmployeeMQ(es service.EmployeeServiceInterface) *EmployeeMQ {
	return &EmployeeMQ{employeeService: es}
}

func (ec *EmployeeMQ) HandleMessage(msg amqp.Delivery, channel *amqp.Channel) {
	queueName := msg.RoutingKey

	switch queueName {
	case "employee_list_queue":
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
		employeeListResponse := ec.employeeService.EmployeeServiceGetList(requestDto)
		utils.SendPaginationResponse(employeeListResponse, employeeListResponse.Meta.Total, filter.Page, filter.Limit, msg.ReplyTo, channel)
	case "employee_byId_queue":
		log.Printf("Received GetById: %v", msg.Body)
		var requestData struct {
			ID int `json:"id"`
		}
		err := json.Unmarshal(msg.Body, &requestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		employeeResponse, err := ec.employeeService.EmployeeServiceGetById(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to retrieve employee", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(employeeResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "employee_create_queue":
		employeeResponse, err := ec.employeeService.EmployeeServiceCreate(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to create employee", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(employeeResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "employee_update_queue":
		var requestData dto.EmployeeUpdateDTO
		err := json.Unmarshal(msg.Body, &requestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		employeeResponse, err := ec.employeeService.EmployeeServiceUpdate(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to update employee", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(employeeResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	case "employee_delete_queue":
		var requestData dto.EmployeeUpdateDTO
		err := json.Unmarshal(msg.Body, &requestData)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			return
		}
		employeeResponse, err := ec.employeeService.EmployeeServiceDelete(msg.Body)
		if err != nil {
			utils.SendExchnageErrorResponse("Failed to delete employee", msg.ReplyTo, channel)
			return
		}
		utils.SendExchageResponseQuene(employeeResponse, http.StatusOK, http.StatusText(http.StatusOK), msg.ReplyTo, channel)
	}
}

func (ec *EmployeeMQ) EmployeeDeclareAndHandleExchanges(channel *amqp.Channel, employeeService_ service.EmployeeServiceInterface) {
	ExchnageDeclareAndHandleQueues(channel, "employee_exchange", "direct", []string{
		"employee_list_queue", "employee_byId_queue", "employee_create_queue", "employee_update_queue", "employee_delete_queue",
	},
		&EmployeeMQ{
			employeeService: employeeService_,
		},
	)
}
