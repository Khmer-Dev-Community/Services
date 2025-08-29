package router

import (
	// Your controller and service imports remain the same
	"net/http"

	postControllers "github.com/Khmer-Dev-Community/Services/api-service/delivery/http"
	postServices "github.com/Khmer-Dev-Community/Services/api-service/lib/posts"
	"github.com/gin-gonic/gin"
)

// PostControllerWrapper is a wrapper for the PostControllerHandler
type PostControllerWrapper struct {
	controller *postControllers.PostControllerHandler
}

// NewPostControllerWrapper initializes the wrapper for PostControllerHandler
func NewPostControllerWrapper(ph *postControllers.PostControllerHandler) *PostControllerWrapper {
	return &PostControllerWrapper{
		controller: ph,
	}
}
func SetupPostRouter(r *gin.Engine, postService postServices.PostService) { // Changed r to *gin.Engine
	postHandler := postControllers.NewPostControllerHandler(postService)
	wrappedPostController := NewPostControllerWrapper(postHandler)
	api := r.Group("/api")                  // Changed to Gin's Group method
	publicPostRoutes := api.Group("/posts") // Changed to Gin's Group method
	publicPostRoutes.GET("/", wrappedPostController.controller.ListPosts)
	publicPostRoutes.GET("/:id", wrappedPostController.controller.GetPostByID) // Path parameters use :id
	publicPostRoutes.GET("/slug/:slug", wrappedPostController.controller.GetPostBySlug)
	authenticatedPostRoutes := api.Group("/posts")

	authenticatedPostRoutes.POST("/create", wrappedPostController.controller.CreatePost)
	authenticatedPostRoutes.PUT("/:id", wrappedPostController.controller.UpdatePost)
	authenticatedPostRoutes.PATCH("/:id", wrappedPostController.controller.UpdatePost)
	// DELETE /api/posts/:id - Delete Post
	authenticatedPostRoutes.DELETE("/:id", wrappedPostController.controller.DeletePost)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	})
}
