package http

import (
	"net/http"
	"strconv"

	"github.com/Khmer-Dev-Community/Services/api-service/delivery/rabbitmq"
	service "github.com/Khmer-Dev-Community/Services/api-service/lib/posts"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"
	"github.com/gin-gonic/gin"
)

// ReactionHandler handles HTTP requests related to Reactions.
type ReactionHandler struct {
	service   service.ReactionService
	publisher *rabbitmq.EventPublisher // Add a publisher field
}

// NewReactionHandler creates a new ReactionHandler instance.
// It now takes a publisher instance as an argument.
func NewReactionHandler(s service.ReactionService, p *rabbitmq.EventPublisher) *ReactionHandler {
	return &ReactionHandler{
		service:   s,
		publisher: p,
	}
}

// CreateReaction handles the creation of a new Reaction or reply.
func (h *ReactionHandler) CreateReaction(c *gin.Context) {
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		utils.ErrorLog(err, "CreateReaction")
		return
	}
	var newReaction service.Reaction
	if err := c.ShouldBindJSON(&newReaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		utils.ErrorLog(err, "CreateReaction")
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorLog(nil, "Unauthorized: UserID not found in context for post deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserID, ok := userID.(uint)
	if !ok {
		utils.ErrorLog(nil, "Internal Server Error: UserID in context is not of type uint")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	newReaction.PostID = uint(postID)
	newReaction.UserID = authUserID
	utils.InfoLog(newReaction, "data create reaction")
	createdReaction, err := h.service.CreateReaction(c.Request.Context(), newReaction, authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Reaction"})
		utils.ErrorLog(err.Error(), "CreateReaction")
		return
	}
	eventType := "reaction_added"
	exchangeName := "post_events"
	var increment int
	if createdReaction == nil {
		increment = -1
	} else {
		increment = 1
	}
	eventPayload := map[string]interface{}{
		"post_id":       newReaction.PostID,
		"user_id":       newReaction.UserID,
		"reaction_type": newReaction.ReactionType,
		"increment":     increment,
	}
	// Access the publisher through the handler's struct
	rabbitmq.PublishBroadcastEvent(eventType, exchangeName, eventPayload)

	utils.SuccessResponse(c, http.StatusOK, createdReaction, "success")
}
