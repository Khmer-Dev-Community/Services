package http

import (
	"net/http"
	"strconv"

	service "github.com/Khmer-Dev-Community/Services/api-service/lib/posts"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"
	"github.com/gin-gonic/gin"
)

// SavedPostHandler handles HTTP requests for saved posts using Gin.
type SavedPostHandler struct {
	savedPostService service.SavedPostServiceImpl
}

// NewSavedPostHandler creates a new instance of SavedPostHandler.
func NewSavedPostHandler(savedPostService service.SavedPostServiceImpl) *SavedPostHandler {
	return &SavedPostHandler{savedPostService: savedPostService}
}

func (h *SavedPostHandler) GetSavedPosts(c *gin.Context) {
	userIDStr := c.Param("userID")
	if userIDStr == "" {
		utils.ErrorResponse(c, http.StatusFound, "request not found", nil)
		return
	}
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid offset format", err)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid limit format", err)
		return
	}

	var viewerID uint
	viewerIDVal, exists := c.Get("userID")
	if exists {
		viewerID = viewerIDVal.(uint)
	}
	savedPosts, err := h.savedPostService.GetSavedPosts(c.Request.Context(), userIDStr, viewerID, offset, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve saved posts", err)
		return
	}
	c.JSON(http.StatusOK, savedPosts)
}

// POST /saved-posts
func (h *SavedPostHandler) SavePost(c *gin.Context) {
	var req service.SavePostData
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	// Get the authenticated user ID from the context.
	_, exist := c.Get("username")
	if !exist {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	authUserID, ok := userID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}
	err := h.savedPostService.SavePost(c.Request.Context(), authUserID, req.PostID, req.Username)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save post", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post saved successfully"})
}

// UnsavePost handles a request to unsave a post.
// DELETE /saved-posts
func (h *SavedPostHandler) UnsavePost(c *gin.Context) {
	var req service.SavePostData
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	// Get the authenticated user ID from the context.
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	authUserID, ok := userID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}
	// Call the service layer to perform the business logic.
	err := h.savedPostService.UnsavePost(c.Request.Context(), authUserID, req.PostID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to unsave post", err)
		return
	}

	c.Status(http.StatusNoContent)
}
