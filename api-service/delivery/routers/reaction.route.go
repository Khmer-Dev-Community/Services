package router

import (
	"net/http"

	HttpControllers "github.com/Khmer-Dev-Community/Services/api-service/delivery/http"
	"github.com/Khmer-Dev-Community/Services/api-service/delivery/rabbitmq"
	reactionServices "github.com/Khmer-Dev-Community/Services/api-service/lib/posts"
	"github.com/gin-gonic/gin"
)

// SetupReactionRouter sets up the Gin routes for Reactions.
// It now takes the eventPublisher instance as a dependency.
func SetupReactionRouter(r *gin.Engine, reactionService reactionServices.ReactionService, eventPublisher *rabbitmq.EventPublisher) {
	// Initialize the ReactionHandler with the service and the publisher.
	reactionHandler := HttpControllers.NewReactionHandler(reactionService, eventPublisher)

	api := r.Group("/api")
	postsGroup := api.Group("/m/:postID")
	{
		reactionsGroup := postsGroup.Group("/reactions")
		{
			reactionsGroup.POST("/", reactionHandler.CreateReaction)
		}
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	})
}
