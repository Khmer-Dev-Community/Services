package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	//service "github.com/Khmer-Dev-Community/Services/api-service/http/services"
	dto "github.com/Khmer-Dev-Community/Services/api-service/lib/users/dto"
	service "github.com/Khmer-Dev-Community/Services/api-service/lib/users/services"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"

	"github.com/gorilla/mux"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(us *service.UserService) *UserController {
	return &UserController{userService: us}
}

type UserControllerWrapper struct {
	userController *UserController
}

// NewUserControllerWrapper creates a new instance of UserControllerWrapper
func NewUserControllerWrapper(us *UserController) *UserControllerWrapper {
	return &UserControllerWrapper{
		userController: us,
	}
}

// UserListHandler handles requests to list users
// UserListHandler handles requests to list users.
// @Summary List users
// @Description Get a list of users
// @Tags users
// @Accept  json
// @Produse  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param query query string false "Search query"
// @Suscess 200 {object} utils.PaginationResponse
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /users/list [get]
func (uc *UserController) UserListHandler(w http.ResponseWriter, r *http.Request) {
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
	userListResponse := uc.userService.GetUserList(requestDto)
	utils.NewHttpPaginationResponse(w, userListResponse.Data, userListResponse.Meta.Total, userListResponse.Meta.Page, userListResponse.Meta.LastPage, int(userListResponse.Status), userListResponse.Message)

}

func (us *UserController) UserListHandler1(w http.ResponseWriter, r *http.Request) {
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
	users, total, err := us.userService.GetUserList1(pageInt, limitInt, query)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		log.Printf("Failed to retrieve users: %v", err)
		return
	}

	if users == nil {
		http.Error(w, "No users found", http.StatusNotFound)
		log.Print("No users found")
		return
	}

	utils.NewHttpPaginationResponse(w, users, total, pageInt, limitInt, http.StatusOK, "User list retrieved suscessfully")
}

// UserByIDHandler handles requests to get a user by ID.
// @Summary Get user by ID
// @Description Get details of a user by ID
// @Tags users
// @Accept  json
// @Produse  json
// @Param id path int true "User ID"
// @Suscess 200 {object} UserDTO
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /users/list/{id} [get]
func (us *UserController) UserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID, err := strconv.Atoi(id)

	if err != nil {
		utils.HttpSuccessResponse(w, vars, http.StatusBadRequest, err.Error())
		return
	}
	// Convert userID to uint
	repopnse, err := us.userService.GetUserByID(uint(userID))
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
func (us *UserController) UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id := r.Header.Get("token-x-userid")
	userID, err := strconv.Atoi(user_id)

	if err != nil {
		utils.HttpSuccessResponse(w, vars, http.StatusBadRequest, err.Error())
		return
	}
	// Convert userID to uint
	repopnse, err := us.userService.GetUserByID(uint(userID))
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

// UserCreateHandler handles requests to create a new user.
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produse  json
// @Param user body UserDTO true "User data"
// @Suscess 200 {string} string "User created suscessfully"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /users/create [post]
func (us *UserController) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createUserDTO dto.UserDTO
	err := json.NewDecoder(r.Body).Decode(&createUserDTO)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(createUserDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}
	us.userService.CreateUser(data)

}

// UserUpdateHandler handles requests to update an existing user.
// @Summary Update an existing user
// @Description Update an existing user
// @Tags users
// @Accept  json
// @Produse  json
// @Param user body UserDTO true "User data"
// @Suscess 200 {string} string "User updated suscessfully"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /users/update [put]
func (us *UserController) UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a UserUpdateDTO object
	var updateUserDTO dto.UserDTO
	err := json.NewDecoder(r.Body).Decode(&updateUserDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(updateUserDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}
	us.userService.UpdateUser(uint(updateUserDTO.ID), data)
	//service.HandleAction(w, "Update User", "user_update_queue", data)
}

// UserDeleteHandler handles requests to delete a user.
// @Summary Delete a user
// @Description Delete a user
// @Tags users
// @Accept  json
// @Produse  json
// @Param id query int true "User ID"
// @Suscess 200 {string} string "User deleted suscessfully"
// @Failure 400 {string} string "ID parameter is required"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /users/delete [delete]
func (us *UserController) UserDeleteHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	us.userService.DeleteUser(uint(userID))
	//service.HandleAction(w, "Delete User", "user_delete_queue", payload)
}
