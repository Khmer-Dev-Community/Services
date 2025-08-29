package http

import (
	"errors"
	"net/http"
	"strconv" // For converting string IDs to uint

	"github.com/Khmer-Dev-Community/Services/api-service/lib/posts"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"
	"github.com/gin-gonic/gin"
)

type PostControllerHandler struct {
	service posts.PostService
}

// NewPostControllerHandler creates a new PostControllerHandler instance.
func NewPostControllerHandler(service posts.PostService) *PostControllerHandler { // CORRECTED: Parameter also takes the interface directly
	return &PostControllerHandler{service: service}
}

func (h *PostControllerHandler) CreatePost(c *gin.Context) {
	var req posts.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LoggerRequest(map[string]interface{}{"error": err.Error()}, "Invalid CreatePostRequest", "Bad request body for post creation")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorLog(nil, "Unauthorized: UserID not found in context for post creation")
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	authorID, ok := userID.(uint)
	if !ok {
		utils.ErrorLog(nil, "Internal Server Error: UserID in context is not of type uint")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	postResponse, err := h.service.CreatePost(c.Request.Context(), req, authorID)
	if err != nil {
		if errors.Is(err, posts.ErrSlugGenerationFailed) {
			utils.ErrorLog(map[string]interface{}{"request": req, "error": err.Error()}, "Failed to create post due to slug generation error")
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate post slug")
			return
		}
		utils.ErrorLog(map[string]interface{}{"request": req, "author_id": authorID, "error": err.Error()}, "Service error during post creation")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	utils.InfoLog(postResponse, "Post created")
	utils.SuccessResponse(c, http.StatusOK, postResponse, "Post created successfully")
}

// GetPostByID handles GET /posts/:id requests.
func (h *PostControllerHandler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id") // Get ID from URL path parameter
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.LoggerRequest(map[string]interface{}{"id_param": idStr}, "Invalid Post ID format", "Bad request: Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	postResponse, err := h.service.GetPostByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, posts.ErrPostNotFound) {
			utils.LoggerRequest(map[string]interface{}{"post_id": id}, "Post not found by ID", "Resource not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		utils.ErrorLog(map[string]interface{}{"post_id": id, "error": err.Error()}, "Service error getting post by ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}

	utils.LoggerRequest(map[string]interface{}{"post_id": id, "title": postResponse.Title}, "Post retrieved by ID successfully", "Post retrieved")
	c.JSON(http.StatusOK, postResponse)
}

// GetPostBySlug handles GET /posts/slug/:slug requests.
func (h *PostControllerHandler) GetPostBySlug(c *gin.Context) {
	slugParam := c.Param("slug") // Get slug from URL path parameter
	if slugParam == "" {
		utils.LoggerRequest(map[string]interface{}{"slug_param": slugParam}, "Empty slug parameter", "Bad request: Slug is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug cannot be empty"})
		return
	}

	postResponse, err := h.service.GetPostBySlug(c.Request.Context(), slugParam)
	if err != nil {
		if errors.Is(err, posts.ErrPostNotFound) {
			utils.LoggerRequest(map[string]interface{}{"slug": slugParam}, "Post not found by slug", "Resource not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		utils.ErrorLog(map[string]interface{}{"slug": slugParam, "error": err.Error()}, "Service error getting post by slug")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}

	utils.LoggerRequest(map[string]interface{}{"slug": slugParam, "post_id": postResponse.ID}, "Post retrieved by slug successfully", "Post retrieved")
	c.JSON(http.StatusOK, postResponse)
}

// UpdatePost handles PUT/PATCH /posts/:id requests.
func (h *PostControllerHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.LoggerRequest(map[string]interface{}{"id_param": idStr}, "Invalid Post ID format for update", "Bad request: Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	var req posts.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LoggerRequest(map[string]interface{}{"error": err.Error()}, "Invalid UpdatePostRequest", "Bad request body for post update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorLog(nil, "Unauthorized: UserID not found in context for post update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserID, ok := userID.(uint)
	if !ok {
		utils.ErrorLog(nil, "Internal Server Error: UserID in context is not of type uint")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	postResponse, err := h.service.UpdatePost(c.Request.Context(), uint(id), req, authUserID)
	if err != nil {
		if errors.Is(err, posts.ErrPostNotFound) {
			utils.LoggerRequest(map[string]interface{}{"post_id": id, "user_id": authUserID}, "Post not found for update", "Resource not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if errors.Is(err, posts.ErrUnauthorizedPostAction) {
			utils.LoggerRequest(map[string]interface{}{"post_id": id, "user_id": authUserID}, "Unauthorized to update post", "Forbidden action")
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this post"})
			return
		}
		if errors.Is(err, posts.ErrSlugGenerationFailed) {
			utils.ErrorLog(map[string]interface{}{"request": req, "error": err.Error()}, "Failed to update post due to slug generation error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new post slug"})
			return
		}
		utils.ErrorLog(map[string]interface{}{"post_id": id, "user_id": authUserID, "request": req, "error": err.Error()}, "Service error during post update")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	utils.LoggerRequest(map[string]interface{}{"post_id": id, "user_id": authUserID}, "Post updated successfully", "Post update successful")
	c.JSON(http.StatusOK, postResponse)
}

// DeletePost handles DELETE /posts/:id requests.
func (h *PostControllerHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.LoggerRequest(map[string]interface{}{"id_param": idStr}, "Invalid Post ID format for delete", "Bad request: Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
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

	err = h.service.DeletePost(c.Request.Context(), uint(id), authUserID)
	if err != nil {
		if errors.Is(err, posts.ErrPostNotFound) {
			utils.LoggerRequest(map[string]interface{}{"post_id": id, "user_id": authUserID}, "Post not found for deletion", "Resource not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if errors.Is(err, posts.ErrUnauthorizedPostAction) {
			utils.LoggerRequest(map[string]interface{}{"post_id": id, "user_id": authUserID}, "Unauthorized to delete post", "Forbidden action")
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this post"})
			return
		}
		utils.ErrorLog(map[string]interface{}{"post_id": id, "user_id": authUserID, "error": err.Error()}, "Service error during post deletion")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	utils.LoggerRequest(map[string]interface{}{"post_id": id, "user_id": authUserID}, "Post deleted successfully", "Post deletion successful")
	// For 204 No Content, c.Status is more idiomatic than c.JSON(http.StatusNoContent, nil)
	c.Status(http.StatusNoContent)
}

// ListPosts handles GET /posts requests.
func (h *PostControllerHandler) ListPosts(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")
	status := c.DefaultQuery("status", "") // Empty string for no status filter
	tag := c.DefaultQuery("tag", "")       // Empty string for no tag filter

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		utils.LoggerRequest(map[string]interface{}{"offset_param": offsetStr}, "Invalid offset format", "Bad request: Invalid offset parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset format"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.LoggerRequest(map[string]interface{}{"limit_param": limitStr}, "Invalid limit format", "Bad request: Invalid limit parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit format"})
		return
	}

	posts, totalCount, err := h.service.ListPosts(c.Request.Context(), offset, limit, status, tag)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"offset": offset, "limit": limit, "status": status, "tag": tag, "error": err.Error()}, "Service error listing posts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	utils.LoggerRequest(map[string]interface{}{"offset": offset, "limit": limit, "status": status, "tag": tag, "count": len(posts), "total": totalCount}, "Posts listed successfully", "Posts retrieved")
	c.JSON(http.StatusOK, gin.H{
		"posts":      posts,
		"totalCount": totalCount,
		"offset":     offset,
		"limit":      limit,
	})
}
