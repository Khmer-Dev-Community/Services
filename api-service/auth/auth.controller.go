package auth

import (
	"fmt"
	"net/http"
	"time"

	redis "github.com/Khmer-Dev-Community/Services/api-service/config"
	users "github.com/Khmer-Dev-Community/Services/api-service/lib/users/models"
	utils "github.com/Khmer-Dev-Community/Services/api-service/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// POST /login
func (c *AuthController) Login(ctx *gin.Context) {
	var credentials LoginCredentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		utils.RespondWithError(ctx.Writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	utils.InfoLog(ctx.Request.Body, "Login Request :")

	user, err := c.authService.Login(&credentials)
	if err != nil {
		utils.RespondWithError(ctx.Writer, http.StatusUnauthorized, err.Error())
		return
	}

	if err := utils.ComparePassword(user.Password, credentials.Password); err != nil {
		utils.ErrorLog(err, "ComparePassword")
		utils.RespondWithError(ctx.Writer, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := c.authService.GenerateToken(user)
	if err != nil {
		utils.RespondWithError(ctx.Writer, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := struct {
		Token string      `json:"token"`
		User  *users.User `json:"user"`
	}{
		Token: token,
		User:  user,
	}

	expiration := time.Hour * 24
	key := fmt.Sprintf("user:%d", user.ID)
	user.Token = token

	if err := redis.SetWithExpiration(key, response, expiration); err != nil {
		utils.ErrorLog(err, "Failed Store User Info in Redis")
		utils.RespondWithError(ctx.Writer, http.StatusInternalServerError, "Failed to store user information in Redis")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// POST /logout
func (c *AuthController) Logout(ctx *gin.Context) {
	id, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	utils.InfoLog(id, "User logout id")

	key := fmt.Sprintf("user:%d", id)
	if err := redis.RemoveRedisKey(key); err != nil {
		utils.RespondWithError(ctx.Writer, http.StatusInternalServerError, "Failed to remove user information from Redis")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
