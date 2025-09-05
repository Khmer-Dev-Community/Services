package http

import (
	"net/http"
	"strconv"

	service "github.com/Khmer-Dev-Community/Services/api-service/lib/posts"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"
	"github.com/gin-gonic/gin"
)

// CommentHandler handles HTTP requests related to comments.
type CommentHandler struct {
	service service.CommentService
}

// NewCommentHandler creates a new CommentHandler instance.
func NewCommentHandler(s service.CommentService) *CommentHandler {
	return &CommentHandler{service: s}
}

// CreateComment handles the creation of a new comment or reply.
func (h *CommentHandler) CreateComment(c *gin.Context) {
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		utils.ErrorLog(err, "CreateComment")
		return
	}

	var newComment service.Comment
	if err := c.ShouldBindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		utils.ErrorLog(err, "CreateComment")
		return
	}

	newComment.PostID = uint(postID)
	// Placeholder for authenticated user ID from context
	newComment.AuthorID = 1

	createdComment, err := h.service.CreateComment(c.Request.Context(), &newComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		utils.ErrorLog(err, "CreateComment")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, createdComment, "success")

}

// GetCommentsByPostID retrieves a list of comments for a given post.
func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		utils.ErrorLog(err, "GetCommentsByPostID")
		return
	}

	comments, err := h.service.GetCommentsByPostID(c.Request.Context(), uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		utils.ErrorLog(err, "GetCommentsByPostID")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, comments, "success")
}

// DeleteComment handles the deletion of a comment.
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("commentID")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		utils.ErrorLog(err, "DeleteComment")
		return
	}

	// Placeholder for authenticated user ID from context
	userID := uint(1)

	if err := h.service.DeleteComment(c.Request.Context(), userID, uint(commentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		utils.ErrorLog(err, "DeleteComment")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, commentID, "deleted success")
}
