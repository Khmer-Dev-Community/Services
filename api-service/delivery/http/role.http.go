package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	dto "github.com/Khmer-Dev-Community/Services/api-service/lib/roles/dto"
	service "github.com/Khmer-Dev-Community/Services/api-service/lib/roles/services"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"

	"github.com/gorilla/mux"
)

type RoleController struct {
	roleService *service.RoleService
}

func NewRoleController(srv *service.RoleService) *RoleController {
	return &RoleController{roleService: srv}
}

type RoleControllerWrapper struct {
	roleController *RoleController
}

// NewUserControllerWrapper creates a new instance of UserControllerWrapper
func NewRoleControllerWrapper(srv *RoleController) *RoleControllerWrapper {
	return &RoleControllerWrapper{
		roleController: srv,
	}
}

// RoleListHandler handles requests to list roles.
// @Summary List roles
// @Description Get a list of roles
// @Tags roles
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param query query string false "Search query"
// @Success 200 {object} utils.PaginationResponse
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /roles/list [get]
func (srv *RoleController) RoleListHandler(w http.ResponseWriter, r *http.Request) {
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
	userListResponse := srv.roleService.HandleRoleListRequest(requestDto)
	utils.NewHttpPaginationResponse(w, userListResponse.Data, userListResponse.Meta.Total, userListResponse.Meta.Page, userListResponse.Meta.LastPage, int(userListResponse.Status), userListResponse.Message)

}

// RoleByIDHandler handles requests to get a role by ID.
// @Summary Get role by ID
// @Description Get details of a role by ID
// @Tags roles
// @Accept  json
// @Produce  json
// @Param id path int true "Role ID"
// @Success 200 {object} RoleDTO
// @Failure 400 {string} string "Invalid role ID"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /roles/list/{id} [get]
func (srv *RoleController) RoleByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	roleId, err := strconv.Atoi(id)
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusBadRequest, err.Error())
		return
	}
	repopnse, err := srv.roleService.HandleRoleByIdRequest(int(roleId))
	if err != nil {
		utils.HttpSuccessResponse(w, repopnse, http.StatusBadRequest, err.Error())
		return
	}
	utils.HttpSuccessResponse(w, repopnse, http.StatusOK, string(utils.SuccessMessage))

}

// RoleCreateHandler handles requests to create a role.
// @Summary Create role
// @Description Create a new role
// @Tags roles
// @Accept  json
// @Produce  json
// @Param body body RoleDTO true "Role data"
// @Success 200 {string} string "Role created successfully"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /roles/create [post]
func (srv *RoleController) RoleCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createDTO dto.RoleDTO
	err := json.NewDecoder(r.Body).Decode(&createDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusBadRequest, err.Error())
		return
	}
	repopnse, err := json.Marshal(createDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, repopnse, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusBadRequest, err.Error())
		return
	}
	utils.HttpSuccessResponse(w, repopnse, http.StatusOK, string(utils.SuccessMessage))
}

// RoleUpdateHandler handles requests to update a role.
// @Summary Update role
// @Description Update an existing role
// @Tags roles
// @Accept  json
// @Produce  json
// @Param body body RoleDTO true "Role data"
// @Success 200 {string} string "Role updated successfully"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /roles/update [put]
func (srv *RoleController) RoleUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateDTO dto.RoleDTO
	err := json.NewDecoder(r.Body).Decode(&updateDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, err.Error())
		return
	}
	data, err := json.Marshal(updateDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, err.Error())
		return
	}
	utils.InfoLog(data, "RoleUpdateHandler")
	response, err := srv.roleService.HandleRoleUpdateRequest(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))

}

// RoleDeleteHandler handles requests to delete a role.
// @Summary Delete role
// @Description Delete an existing role
// @Tags roles
// @Accept  json
// @Produce  json
// @Param id query int true "Role ID"
// @Success 200 {string} string "Role deleted successfully"
// @Failure 400 {string} string "ID parameter is required"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /roles/delete [delete]
func (srv *RoleController) RoleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, string(utils.ErrorMessage))
		return
	}
	data := map[string]string{
		"action": "Delete Role",
		"id":     id,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, string(utils.ErrorMessage))
		return
	}
	response, err := srv.roleService.HandleRoleDeleteRequest(payload)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, string(utils.ErrorMessage))
		return
	}
	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}
