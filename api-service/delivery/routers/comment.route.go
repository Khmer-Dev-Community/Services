package router

import (
	"net/http"

	commentControllers "github.com/Khmer-Dev-Community/Services/api-service/delivery/http"
	commentServices "github.com/Khmer-Dev-Community/Services/api-service/lib/posts"
	"github.com/gin-gonic/gin"
)

// CommentControllerWrapper is a wrapper for the CommentControllerHandler
type CommentControllerWrapper struct {
	controller *commentControllers.CommentHandler
}

// NewCommentControllerWrapper initializes the wrapper for CommentControllerHandler
func NewCommentControllerWrapper(ch *commentControllers.CommentHandler) *CommentControllerWrapper {
	return &CommentControllerWrapper{
		controller: ch,
	}
}

// SetupCommentRouter sets up the Gin routes for comments.
func SetupCommentRouter(r *gin.Engine, commentService commentServices.CommentService) {
	commentHandler := commentControllers.NewCommentHandler(commentService)
	wrappedCommentController := NewCommentControllerWrapper(commentHandler)
	api := r.Group("/api")
	postsGroup := api.Group("/m/:postID")
	{
		commentsGroup := postsGroup.Group("/comments")
		{
			publicCommentRoutes := commentsGroup.Group("/")
			{
				publicCommentRoutes.GET("/", wrappedCommentController.controller.GetCommentsByPostID)
			}
			authenticatedCommentRoutes := commentsGroup.Group("/")
			{
				authenticatedCommentRoutes.POST("/", wrappedCommentController.controller.CreateComment)
				authenticatedCommentRoutes.DELETE("/:commentID", wrappedCommentController.controller.DeleteComment)
			}
		}
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	})
}
