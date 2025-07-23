package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	service "telegram-service/lib/permission/service"

	dto "telegram-service/lib/permission/dto"
	"telegram-service/utils"

	"github.com/gorilla/mux"
)

type PermissController struct {
	permisService *service.PermissionService
}

func NewPermissController(us *service.PermissionService) *PermissController {
	return &PermissController{permisService: us}
}

type PermissControllerWrapper struct {
	PermissController *PermissController
}

// NewPermissControllerWrapper creates a new instance of PermissControllerWrapper
func NewPermissControllerWrapper(us *PermissController) *PermissControllerWrapper {
	return &PermissControllerWrapper{
		PermissController: us,
	}
}

// PermissionListHandler handles requests to list permissions.
// @Summary List permissions
// @Description Get a list of permissions
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param query query string false "Search query"
// @Success 200 {object} utils.PaginationResponse
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /permissions/list [get]
func (ctl *PermissController) PermissionListHandler(w http.ResponseWriter, r *http.Request) {
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
	userListResponse := ctl.permisService.HandlePermissionListRequest(requestDto)
	utils.NewHttpPaginationResponse(w, userListResponse.Data, userListResponse.Meta.Total, userListResponse.Meta.Page, userListResponse.Meta.LastPage, int(userListResponse.Status), userListResponse.Message)

}

// PermissionRoleByIDHandler handles requests to get permissions by role ID.
// @Summary Get permissions by role ID
// @Description Get permissions assigned to a specific role by ID
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param roleid header int true "Role ID"
// @Success 200 {object} PermissionDTO
// @Failure 400 {string} string "Role ID not found in headers"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /permissions/role [get]
func (ctl *PermissController) PermissionRoleByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("PermissionRoleByIDHandler called")
	id := r.Header.Get("roleid")
	if id == "" {
		utils.RespondWithJSON(w, http.StatusOK, "logout")
		http.Error(w, "Role ID not found in headers", http.StatusBadRequest)
		return
	}

	roleId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	repopnse, err := ctl.permisService.HandleRolePermissionById(roleId)
	if err != nil {
		utils.HttpSuccessResponse(w, repopnse, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, repopnse, http.StatusOK, string(utils.SuccessMessage))
	//service.HandleAction(w, "Get by ID", "permission_role_queue", payload)
}

// PermissionByIDHandler handles requests to get a permission by ID.
// @Summary Get permission by ID
// @Description Get details of a permission by ID
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param id path int true "Permission ID"
// @Success 200 {object} PermissionDTO
// @Failure 400 {string} string "Invalid permission ID"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /permissions/list/{id} [get]
func (ctl *PermissController) PermissionByIDHandler1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	permissionID, err := strconv.Atoi(id)
	if err != nil {
		// Handle error
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}
	ctl.permisService.HandlePermissionByIdRequest(uint(permissionID))
	//service.HandleAction(w, " Get by ID", "permission_findOne_queue", payload)
}
func (ctl *PermissController) PermissionByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	permissionID, err := strconv.Atoi(id)

	if err != nil {
		utils.HttpSuccessResponse(w, vars, http.StatusBadRequest, err.Error())
		return
	}
	// Convert userID to uint
	repopnse, err := ctl.permisService.HandlePermissionByIdRequest(uint(permissionID))
	if err != nil {
		utils.HttpSuccessResponse(w, repopnse, http.StatusBadRequest, err.Error())
		return
	}

	if repopnse == nil {
		utils.HttpSuccessResponse(w, repopnse, http.StatusBadRequest, string(utils.NotFoundMessage))
		return
	}
	utils.HttpSuccessResponse(w, repopnse, http.StatusOK, string(utils.SuccessMessage))

}

// PermissionCreateHandler handles requests to create a permission.
// @Summary Create permission
// @Description Create a new permission
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param body body PermissionDTO true "Permission data"
// @Success 200 {string} string "Permission created successfully"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /permissions/create [post]
func (ctl *PermissController) PermissionCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createDto dto.PermissionCreateDTO
	err := json.NewDecoder(r.Body).Decode(&createDto)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(createDto)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}
	ctl.permisService.HandlePermissionCreateRequest(data)
	//service.HandleAction(w, "Create Permission", "permission_create_queue", data)

}

// PermissionUpdateHandler handles requests to update a permission.
// @Summary Update permission
// @Description Update an existing permission
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param body body PermissionDTO true "Permission data"
// @Success 200 {string} string "Permission updated successfully"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /permissions/update [put]
func (ctl *PermissController) PermissionUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a PermissionDTO object
	var updateDTO dto.PermissionUpdateDTO
	err := json.NewDecoder(r.Body).Decode(&updateDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(updateDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}
	ctl.permisService.HandlePermissionUpdateRequest(uint(updateDTO.ID), data)
	//service.HandleAction(w, "Update Permission", "permission_update_queue", data)
}

// PermissionDeleteHandler handles requests to delete a permission.
// @Summary Delete permission
// @Description Delete an existing permission
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param id query int true "Permission ID"
// @Success 200 {string} string "Permission deleted successfully"
// @Failure 400 {string} string "ID parameter is required"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /permissions/delete [delete]
func (ctl *PermissController) PermissionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}
	data := map[string]string{
		"action": "Delete Permission",
		"id":     id,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		log.Printf("Failed to encode JSON: %v", err)
		return
	}
	ctl.permisService.HandlePermissionDeleteRequest(payload)
	//service.HandleAction(w, "Delete Permission", "permission_delete_queue", payload)
}

// PermissionRoleUpdateHandler handles requests to update permissions for a role.
// @Summary Update role permissions
// @Description Update permissions assigned to a role
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param body body UpdateRolePermissionsMessage true "Role permissions data"
// @Success 200 {string} string "Role permissions updated successfully"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /permissions/role/update [put]
func (ctl *PermissController) PermissionRoleUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateDTO dto.UpdateRolePermissionsMessage
	err := json.NewDecoder(r.Body).Decode(&updateDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		utils.ErrorLog(err.Error(), fmt.Sprintf("Failed to decode request body: %v", err))
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}
	utils.InfoLog(updateDTO, "Request")
	data, err := json.Marshal(updateDTO)
	if err != nil {
		utils.ErrorLog(err.Error(), fmt.Sprintf("Failed to marshal request data: %v", err))
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}
	ctl.permisService.UpdateRolePermissionByRoleID(data)
	utils.HttpSuccessResponse(w, updateDTO, http.StatusOK, string(utils.SuccessMessage))

}
