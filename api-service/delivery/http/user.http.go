package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	dto "github.com/Khmer-Dev-Community/Services/api-service/lib/users/dto"
	service "github.com/Khmer-Dev-Community/Services/api-service/lib/users/services"
	"github.com/Khmer-Dev-Community/Services/api-service/utils" // Your utility functions

	"github.com/gin-gonic/gin" // Gin framework
	// Removed: "github.com/gorilla/mux"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(us *service.UserService) *UserController {
	return &UserController{userService: us}
}

// UserControllerWrapper is typically not needed when directly passing *gin.Context
// to controller methods. Gin handlers are directly `func(*gin.Context)`.
// However, if you have other specific wrapping logic, you can keep it.
// For this conversion, we'll keep the struct but note its redundancy if only used for method access.
type UserControllerWrapper struct {
	userController *UserController
}

// NewUserControllerWrapper creates a new instance of UserControllerWrapper
// Note: This wrapper might be redundant if you're directly registering
// `uc.MethodName` as Gin handlers.
func NewUserControllerWrapper(us *UserController) *UserControllerWrapper {
	return &UserControllerWrapper{
		userController: us,
	}
}

// UserListHandler handles requests to list users.
// @Summary List users
// @Description Get a list of users
// @Tags users
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param query query string false "Search query"
// @Success 200 {object} utils.PaginationResponse // Assuming utils.PaginationResponse is a suitable struct
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users/list [get]
func (uc *UserController) UserListHandler(c *gin.Context) { // Changed signature to *gin.Context
	pageStr := c.Query("page")   // Use c.Query for query parameters
	limitStr := c.Query("limit") // Use c.Query for query parameters
	query := c.Query("query")    // Use c.Query for query parameters

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

	// Adapting NewHttpPaginationResponse to Gin's c.JSON
	// Assuming userListResponse.Data, Meta, Status, Message are accessible.
	c.JSON(int(userListResponse.Status), gin.H{ // Use c.JSON for response
		"data":      userListResponse.Data,
		"total":     userListResponse.Meta.Total,
		"page":      userListResponse.Meta.Page,
		"last_page": userListResponse.Meta.LastPage,
		"message":   userListResponse.Message,
	})
}

// UserListHandler1 (converted to Gin) - Consider consolidating with UserListHandler
func (uc *UserController) UserListHandler1(c *gin.Context) { // Changed signature to *gin.Context
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	query := c.Query("query")

	pageInt, err := strconv.Atoi(pageStr)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil || limitInt <= 0 {
		limitInt = 10
	}

	users, total, err := uc.userService.GetUserList1(pageInt, limitInt, query)
	if err != nil {
		log.Printf("Failed to retrieve users: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"}) // Use Gin's error response
		return
	}

	if users == nil {
		log.Print("No users found")
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No users found"}) // Use Gin's error response
		return
	}
	c.JSON(http.StatusOK, gin.H{ // Use c.JSON for response
		"data":      users,
		"total":     total,
		"page":      pageInt,
		"last_page": (total + limitInt - 1) / limitInt,
		"message":   "User list retrieved successfully",
	})
}

// UserByIDHandler handles requests to get a user by ID.
// @Summary Get user by ID
// @Description Get details of a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} UserDTO // Assuming UserDTO is visible to Swagger
// @Failure 400 {object} gin.H "Invalid user ID"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users/list/{id} [get]
func (uc *UserController) UserByIDHandler(c *gin.Context) { // Changed signature to *gin.Context
	idStr := c.Param("id") // Use c.Param for path parameters
	userID, err := strconv.Atoi(idStr)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid user ID: %v", err)}) // Use Gin's error response
		return
	}

	response, err := uc.userService.GetUserByID(uint(userID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Use Gin's error response
		return
	}

	if response == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": string(utils.NotFoundMessage)}) // Use Gin's error response
		return
	}

	c.JSON(http.StatusOK, gin.H{ // Use c.JSON for success response
		"data":    response,
		"message": string(utils.SuccessMessage),
	})
}

// UserProfileHandler (converted to Gin)
func (uc *UserController) UserProfileHandler(c *gin.Context) { // Changed signature to *gin.Context
	// Get user ID from Gin context, set by AuthMiddlewareWithWhiteList
	userIDAny, exists := c.Get("userID") // Assuming "userID" is set in context
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDAny.(uint) // Assert to uint
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type in context"})
		return
	}

	response, err := uc.userService.GetUserByID(userID) // Pass uint directly
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if response == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": string(utils.NotFoundMessage)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": string(utils.SuccessMessage),
	})
}

// UserCreateHandler handles requests to create a new user.
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body dto.UserDTO true "User data" // Ensure dto.UserDTO is visible to Swagger
// @Success 200 {object} gin.H "User created successfully"
// @Failure 400 {object} gin.H "Bad Request"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users/create [post]
func (uc *UserController) UserCreateHandler(c *gin.Context) { // Changed signature to *gin.Context
	var createUserDTO dto.UserDTO
	// Use c.ShouldBindJSON to bind the request body to the DTO.
	// This handles JSON decoding and error checking.
	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to decode request body: %v", err.Error())})
		return
	}
	data, err := json.Marshal(createUserDTO)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request data"})
		return
	}
	uc.userService.CreateUser(data)
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"}) // Send success response
}

// UserUpdateHandler handles requests to update an existing user.
// @Summary Update an existing user
// @Description Update an existing user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body dto.UserDTO true "User data"
// @Success 200 {object} gin.H "User updated successfully"
// @Failure 400 {object} gin.H "Bad Request"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users/update [put]
func (uc *UserController) UserUpdateHandler(c *gin.Context) { // Changed signature to *gin.Context
	var updateUserDTO dto.UserDTO
	if err := c.ShouldBindJSON(&updateUserDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to decode request body: %v", err.Error())})
		return
	}

	data, err := json.Marshal(updateUserDTO)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request data"})
		return
	}
	uc.userService.UpdateUser(uint(updateUserDTO.ID), data)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"}) // Send success response
}

// UserDeleteHandler handles requests to delete a user.
// @Summary Delete a user
// @Description Delete a user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id query int true "User ID"
// @Success 200 {object} gin.H "User deleted successfully"
// @Failure 400 {object} gin.H "Bad Request"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users/delete [delete]
func (uc *UserController) UserDeleteHandler(c *gin.Context) { // Changed signature to *gin.Context
	idStr := c.Query("id") // Use c.Query for query parameters

	if idStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Call the service method.
	// You might want to get a response from uc.userService.DeleteUser(uint(userID))
	// and use it to craft a success/failure message.
	uc.userService.DeleteUser(uint(userID))

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"}) // Send success response
}
