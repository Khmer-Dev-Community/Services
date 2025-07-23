package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	dto "telegram-service/lib/employee"
	service "telegram-service/lib/employee"
	"telegram-service/utils"

	"github.com/gorilla/mux"
)

type EmployeeController struct {
	employeeService service.EmployeeServiceInterface
}

func NewEmployeeController(es service.EmployeeServiceInterface) *EmployeeController {
	return &EmployeeController{employeeService: es}
}

type EmployeeControllerWrapper struct {
	employeeController *EmployeeController
}

func NewEmployeeControllerWrapper(ec *EmployeeController) *EmployeeControllerWrapper {
	return &EmployeeControllerWrapper{
		employeeController: ec,
	}
}

// EmployeeListHandler handles requests to list employees
// @Summary List employees
// @Description Get a list of employees
// @Tags employees
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param query query string false "Search query"
// @Success 200 {object} utils.PaginationResponse
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /employees/list [get]
func (ec *EmployeeController) EmployeeListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")

	pageInt, err := strconv.Atoi(pageStr)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil || limitInt <= 0 {
		limitInt = 10
	}

	requestDto := utils.PaginationRequestDTO{
		Page:  pageInt,
		Limit: limitInt,
		Query: query,
	}

	employeeListResponse := ec.employeeService.EmployeeServiceGetList(requestDto)
	utils.NewHttpPaginationResponse(w, employeeListResponse.Data, employeeListResponse.Meta.Total, employeeListResponse.Meta.Page, employeeListResponse.Meta.LastPage, int(employeeListResponse.Status), employeeListResponse.Message)
}

// EmployeeByIDHandler handles requests to get an employee by ID
// @Summary Get employee by ID
// @Description Get details of an employee by ID
// @Tags employees
// @Accept  json
// @Produce  json
// @Param id path int true "Employee ID"
// @Success 200 {object} EmployeeUpdateDTO
// @Failure 400 {string} string "Invalid employee ID"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /employees/{id} [get]
func (ec *EmployeeController) EmployeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	employeeID, err := strconv.Atoi(id)

	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	empDTO := &dto.EmployeeUpdateDTO{ID: uint(employeeID)}
	data, err := json.Marshal(empDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, "Failed to marshal request data")
		return
	}

	response, err := ec.employeeService.EmployeeServiceGetById(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// EmployeeCreateHandler handles requests to create a new employee
// @Summary Create a new employee
// @Description Create a new employee
// @Tags employees
// @Accept  json
// @Produce  json
// @Param employee body EmployeeUpdateDTO true "Employee data"
// @Success 200 {object} EmployeeUpdateDTO
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /employees/create [post]
func (ec *EmployeeController) EmployeeCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createEmployeeDTO dto.EmployeeUpdateDTO
	err := json.NewDecoder(r.Body).Decode(&createEmployeeDTO)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(createEmployeeDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	response, err := ec.employeeService.EmployeeServiceCreate(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// EmployeeUpdateHandler handles requests to update an existing employee
// @Summary Update an existing employee
// @Description Update an existing employee
// @Tags employees
// @Accept  json
// @Produce  json
// @Param employee body EmployeeUpdateDTO true "Employee data"
// @Success 200 {object} EmployeeUpdateDTO
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /employees/update [put]
func (ec *EmployeeController) EmployeeUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateEmployeeDTO dto.EmployeeUpdateDTO
	err := json.NewDecoder(r.Body).Decode(&updateEmployeeDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(updateEmployeeDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	response, err := ec.employeeService.EmployeeServiceUpdate(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// EmployeeDeleteHandler handles requests to delete an employee
// @Summary Delete an employee
// @Description Delete an employee
// @Tags employees
// @Accept  json
// @Produce  json
// @Param id query int true "Employee ID"
// @Success 200 {string} string "Employee deleted successfully"
// @Failure 400 {string} string "ID parameter is required"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /employees/delete [delete]
func (ec *EmployeeController) EmployeeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	employeeID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	empDTO := &dto.EmployeeUpdateDTO{ID: uint(employeeID)}
	data, err := json.Marshal(empDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	deleted, err := ec.employeeService.EmployeeServiceDelete(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, deleted, http.StatusOK, string(utils.SuccessMessage))
}
