package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	service "telegram-service/lib/menus"
	dto "telegram-service/lib/menus/dto"
	"telegram-service/utils"

	"github.com/gorilla/mux"
)

type MenuController struct {
	menuService *service.MenuService
}

func NewMenuController(us *service.MenuService) *MenuController {
	return &MenuController{menuService: us}
}

type MenuControllerWrapper struct {
	MenuController *MenuController
}

// NewMenuControllerWrapper creates a new instance of MenuControllerWrapper
func NewMenuControllerWrapper(us *MenuController) *MenuControllerWrapper {
	return &MenuControllerWrapper{
		MenuController: us,
	}
}

// MenuListSubHandler handles requests to list menu items by role ID.
// @Summary List menu items by role ID
// @Description Get menu items based on the role ID from headers
// @Tags menus
// @Accept  json
// @Produce  json
// @Param roleid header int true "Role ID"
// @Success 200 {object} interface{}
// @Failure 400 {string} string "Role ID not found in headers"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /menus/sub [get]
func (ctl *MenuController) MenuListByTenant(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("roleid")
	roleId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Role ID not found in headers", http.StatusBadRequest)
		return
	}
	repopnse, err := ctl.menuService.HandleSelectMenusByTenantRoleID(int(roleId))
	if err != nil {
		utils.HttpSuccessResponse(w, repopnse, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, repopnse, http.StatusOK, string(utils.SuccessMessage))

}
func (ctl *MenuController) MenuListSubHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("roleid")
	roleId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Role ID not found in headers", http.StatusBadRequest)
		return
	}
	repopnse, err := ctl.menuService.HandleMenuListSub(int(roleId))
	if err != nil {
		utils.HttpSuccessResponse(w, repopnse, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, repopnse, http.StatusOK, string(utils.SuccessMessage))

}

// MenuListHandler handles requests to list menu items.
// @Summary List menu items
// @Description Get a list of menu items
// @Tags menus
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param query query string false "Search query"
// @Success 200 {object} interface{}
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /menus/list [get]
func (ctl *MenuController) MenuListHandler(w http.ResponseWriter, r *http.Request) {
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

	menuListResponse := ctl.menuService.HandleMenuListRequest(requestDto)
	utils.NewHttpPaginationResponse(w, menuListResponse.Data, menuListResponse.Meta.Total, menuListResponse.Meta.Page, menuListResponse.Meta.LastPage, int(menuListResponse.Status), menuListResponse.Message)

}

// MenuListAllWithChild handles requests to list all menu items with their children.
// @Summary List all menu items with children
// @Description Get a hierarchical list of menu items with their children by ID
// @Tags menus
// @Accept  json
// @Produce  json
// @Param id path int true "Menu ID"
// @Success 200 {object} interface{}
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /menus/list-with-children/{id} [get]
func (ctl *MenuController) MenuListAllWithChild(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	pageInt, err := strconv.Atoi(id)
	if err != nil || pageInt <= 0 {
		pageInt = 0 // Default value
	}
	menuListResponse, err := ctl.menuService.HandleMenuListWithChild()
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusBadRequest, err.Error())
		return
	}
	utils.HttpSuccessResponse(w, menuListResponse, http.StatusOK, string(utils.SuccessMessage))

}

// MenuListAllWithChildLabel handles requests to list all menu items with their children by label.
// @Summary List all menu items with children by label
// @Description Get a hierarchical list of menu items with their children by ID for label
// @Tags menus
// @Accept  json
// @Produce  json
// @Param id path int true "Menu ID"
// @Success 200 {object} interface{}
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /menus/list-with-children-label/{id} [get]
func (ctl *MenuController) MenuListAllWithChildLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	pageInt, err := strconv.Atoi(id)
	if err != nil || pageInt <= 0 {
		pageInt = 0 // Default value
	}

	userListResponse, err := ctl.menuService.HandleMenuListWithChildToPermissionAuthorize(int(pageInt))
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusBadRequest, err.Error())
		return
	}
	utils.HttpSuccessResponse(w, userListResponse, http.StatusOK, string(utils.SuccessMessage))
	//utils.NewHttpPaginationResponse(w, userListResponse.Data, userListResponse.Meta.Total, userListResponse.Meta.Page, userListResponse.Meta.LastPage, int(userListResponse.Status), userListResponse.Message)

	//service.HandleAction(w, "List Menu as Child label", "menu_listWithChildLabel_queue", payload)
}

// MenuListAllWithChildById handles requests to list all menu items with their children by ID.
// @Summary List all menu items with children by ID
// @Description Get a hierarchical list of menu items with their children by ID
// @Tags menus
// @Accept  json
// @Produce  json
// @Param id path int true "Menu ID"
// @Success 200 {object} interface{}
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /menus/list-with-children-by-id/{id} [get]
func (ctl *MenuController) MenuListAllWithChildById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	pageInt, err := strconv.Atoi(id)
	if err != nil || pageInt <= 0 {
		pageInt = 0 // Default value
	}

	ListResponse, err := ctl.menuService.HandleMenuListWithChildByID(pageInt)
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusBadRequest, err.Error())
		return
	}
	utils.HttpSuccessResponse(w, ListResponse, http.StatusOK, string(utils.SuccessMessage))
	//service.HandleAction(w, "List Menu as Child", "menu_listWithChildById_queue", payload)
}

// MenuByIDHandler handles requests to get a menu item by ID.
// @Summary Get menu item by ID
// @Description Get details of a menu item by ID
// @Tags menus
// @Accept  json
// @Produce  json
// @Param id path int true "Menu ID"
// @Success 200 {object} interface{}
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /menus/list/{id} [get]
func (ctl *MenuController) MenuByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	menuID, err := strconv.Atoi(id)
	if err != nil {
		// Handle error
		http.Error(w, "Invalid menu ID", http.StatusBadRequest)
		return
	}

	ListResponse, err := ctl.menuService.HandleMenuByIdRequest(menuID)
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusBadRequest, err.Error())
		return
	}
	utils.HttpSuccessResponse(w, ListResponse, http.StatusOK, string(utils.SuccessMessage))
	//service.HandleAction(w, " Get by ID", "menu_findOne_queue", payload)
}

// MenuCreateHandler handles requests to create a menu item.
// @Summary Create menu item
// @Description Create a new menu item
// @Tags menus
// @Accept  json
// @Produce  json
// @Param body body MenuDTO true "Menu data"
// @Success 200 {string} string "Menu created successfully"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /menus/create [post]
func (ctl *MenuController) MenuCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createDto dto.MenuDTO
	err := json.NewDecoder(r.Body).Decode(&createDto)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(createDto)
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusInternalServerError, err.Error())
		return
	}
	ListResponse, err := ctl.menuService.HandleMenuCreateRequest(data)
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusBadRequest, err.Error())
		return
	}
	utils.HttpSuccessResponse(w, ListResponse, http.StatusOK, string(utils.SuccessMessage))

	//service.HandleAction(w, "Create Menu", "menu_create_queue", data)
}

// @Summary Update menu item
// @Description Update an existing menu item
// @Tags menus
// @Accept  json
// @Produce  json
// @Param body body MenuDTO true "Updated menu data"
// @Success 200 {string} string "Menu updated successfully"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /menus/update [put]
func (ctl *MenuController) MenuUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a MenuDTO object
	var updateDTO dto.MenuDTO
	err := json.NewDecoder(r.Body).Decode(&updateDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(updateDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusInternalServerError, err.Error())
		return
	}

	ListResponse, err := ctl.menuService.HandleMenuUpdateRequest(data)
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusInternalServerError, err.Error())
		return
	}
	utils.HttpSuccessResponse(w, ListResponse, http.StatusOK, string(utils.SuccessMessage))
	//service.HandleAction(w, "Update Menu", "menu_update_queue", data)
}

// MenuDeleteHandler handles requests to delete a menu item by ID.
// @Summary Delete menu item
// @Description Delete a menu item by ID
// @Tags menus
// @Accept  json
// @Produce  json
// @Param id query int true "Menu ID"
// @Success 200 {string} string "Menu deleted successfully"
// @Failure 400 {string} string "ID parameter is required"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /menus/delete [delete]
func (ctl *MenuController) MenuDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		utils.HttpSuccessResponse(w, "request missing", http.StatusBadGateway, "Request Not found")
		return
	}
	data := map[string]string{
		"id": id,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusInternalServerError, err.Error())
		return
	}
	ListResponse, err := ctl.menuService.HandleMenuDeleteRequest(payload)
	if err != nil {
		utils.HttpSuccessResponse(w, err, http.StatusBadRequest, err.Error())
		return
	}
	utils.HttpSuccessResponse(w, ListResponse, http.StatusOK, string(utils.SuccessMessage))
	//service.HandleAction(w, "Delete Menu", "menu_delete_queue", payload)
}
