package router

import (
	HttpControllers "github.com/Khmer-Dev-Community/Services/api-service/delivery/http"
	saveServices "github.com/Khmer-Dev-Community/Services/api-service/lib/posts"
	"github.com/gin-gonic/gin"
)

// CommentControllerWrapper is a wrapper for the CommentControllerHandler
type SaveControllerWrapper struct {
	controller *HttpControllers.SavedPostHandler
}

// NewCommentControllerWrapper initializes the wrapper for CommentControllerHandler
func NewSaveControllerWrapper(ch *HttpControllers.SavedPostHandler) *SaveControllerWrapper {
	return &SaveControllerWrapper{
		controller: ch,
	}
}

func SetupSavePostRouter(r *gin.Engine, saveServices saveServices.SavedPostServiceImpl) {
	savetHandler := HttpControllers.NewSavedPostHandler(saveServices)
	saveHttp := NewSaveControllerWrapper(savetHandler)
	api := r.Group("/api")
	userRouter := api.Group("/saved-post")
	userRouter.GET("/:userID", saveHttp.controller.GetSavedPosts)
	userRouter.PUT("/", saveHttp.controller.UnsavePost)
	userRouter.POST("/", saveHttp.controller.SavePost)
}
