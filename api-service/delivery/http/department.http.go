package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	department "telegram-service/lib/department"
	//service "telegram-service/lib/department"
	"telegram-service/utils"

	"github.com/gorilla/mux"
)

type DepartmentController struct {
	departmentService department.DepartmentServiceInterface
}

func NewDepartmentController(ds department.DepartmentServiceInterface) *DepartmentController {
	return &DepartmentController{departmentService: ds}
}

// DepartmentListHandler handles requests to list departments
// @Summary List departments
// @Description Get a list of departments
// @Tags departments
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param query query string false "Search query"
// @Success 200 {object} utils.PaginationResponse
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /departments/list [get]
func (dc *DepartmentController) DepartmentListHandler(w http.ResponseWriter, r *http.Request) {
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

	departmentListResponse := dc.departmentService.DepartmentServiceGetList(requestDto)
	utils.NewHttpPaginationResponse(w, departmentListResponse.Data, departmentListResponse.Meta.Total, departmentListResponse.Meta.Page, departmentListResponse.Meta.LastPage, int(departmentListResponse.Status), departmentListResponse.Message)
}

// DepartmentByIDHandler handles requests to get a department by ID
// @Summary Get department by ID
// @Description Get details of a department by ID
// @Tags departments
// @Accept  json
// @Produce  json
// @Param id path int true "Department ID"
// @Success 200 {object} department.DepartmentDTO
// @Failure 400 {string} string "Invalid department ID"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /departments/{id} [get]
func (dc *DepartmentController) DepartmentByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	departmentID, err := strconv.Atoi(id)

	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid department ID")
		return
	}

	deptDTO := &department.DepartmentDTO{ID: int(departmentID)}
	data, err := json.Marshal(deptDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, "Failed to marshal request data")
		return
	}

	response, err := dc.departmentService.DepartmentServiceGetById(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// DepartmentCreateHandler handles requests to create a new department
// @Summary Create a new department
// @Description Create a new department
// @Tags departments
// @Accept  json
// @Produce  json
// @Param department body department.DepartmentDTO true "Department data"
// @Success 200 {object} department.DepartmentDTO
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /departments/create [post]
func (dc *DepartmentController) DepartmentCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createDepartmentDTO department.DepartmentDTO
	err := json.NewDecoder(r.Body).Decode(&createDepartmentDTO)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(createDepartmentDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	response, err := dc.departmentService.DepartmentServiceCreate(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// DepartmentUpdateHandler handles requests to update an existing department
// @Summary Update an existing department
// @Description Update an existing department
// @Tags departments
// @Accept  json
// @Produce  json
// @Param department body department.DepartmentDTO true "Department data"
// @Success 200 {object} department.DepartmentDTO
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /departments/update [put]
func (dc *DepartmentController) DepartmentUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateDepartmentDTO department.DepartmentDTO
	err := json.NewDecoder(r.Body).Decode(&updateDepartmentDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(updateDepartmentDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	response, err := dc.departmentService.DepartmentServiceUpdate(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// DepartmentDeleteHandler handles requests to delete a department
// @Summary Delete a department
// @Description Delete a department
// @Tags departments
// @Accept  json
// @Produce  json
// @Param id query int true "Department ID"
// @Success 200 {string} string "Department deleted successfully"
// @Failure 400 {string} string "ID parameter is required"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /departments/delete [delete]
func (dc *DepartmentController) DepartmentDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	departmentID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	deptDTO := &department.DepartmentDTO{ID: int(departmentID)}
	data, err := json.Marshal(deptDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	deleted, err := dc.departmentService.DepartmentServiceDelete(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, deleted, http.StatusOK, string(utils.SuccessMessage))
}
