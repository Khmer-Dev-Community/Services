package http

import (
	"net/http"
	"strconv"

	service "github.com/Khmer-Dev-Community/Services/api-service/lib/clientfollow"
	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	FollowService       service.FollowService
	NotificationService service.NotificationService
}

func NewHandler(
	followService service.FollowService,
	notificationService service.NotificationService,
) *FollowHandler {
	return &FollowHandler{
		FollowService:       followService,
		NotificationService: notificationService,
	}
}

// NotificationHandler is an example of an API handler.
type NotificationHandler struct {
	Service service.NotificationService
}

// GetNotifications is the API handler that orchestrates the request.
// GET /api/notifications
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.MustGet("user_id").(uint) // Get user ID from middleware
	notifications, err := h.Service.GetUserNotifications(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get notifications"})
		return
	}
	notificationDTOs := service.ToNotificationDTOs(notifications)
	c.JSON(200, gin.H{"data": notificationDTOs})
}

// POST /api/users/:followedID/follow
func (h *FollowHandler) ToggleFollow(c *gin.Context) {
	// Extract IDs from the request
	followerID := c.MustGet("user_id").(uint) // Assuming a middleware sets this
	followedID, err := strconv.ParseUint(c.Param("followedID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid followed ID"})
		return
	}
	// Call the service layer to perform the business logic
	isNowFollowing, err := h.FollowService.ToggleFollow(followerID, uint(followedID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if isNowFollowing {
		c.JSON(http.StatusOK, gin.H{"message": "User followed successfully."})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "User unfollowed successfully."})
	}
}

// GET /api/users/:userID/followers
func (h *FollowHandler) GetUserFollowers(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the service layer
	followers, err := h.FollowService.GetUserFollowers(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve followers"})
		return
	}

	// Use the mapper to convert models to DTOs for the response
	followerDTOs := service.ToAccountFollowDTOs(followers)
	c.JSON(http.StatusOK, followerDTOs)
}

// GET /api/users/:userID/following
func (h *FollowHandler) GetUserFollowing(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the service layer
	following, err := h.FollowService.GetUserFollowing(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve following"})
		return
	}

	// Use the mapper to convert models to DTOs
	followingDTOs := service.ToAccountFollowDTOs(following)
	c.JSON(http.StatusOK, followingDTOs)
}

// PUT /api/notifications/:notificationID/read
func (h *FollowHandler) MarkNotificationAsRead(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("notificationID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Call the service layer
	err = h.NotificationService.MarkNotificationAsRead(uint(notificationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read."})
}
