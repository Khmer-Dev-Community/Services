package router

import (
	HttpControllers "github.com/Khmer-Dev-Community/Services/api-service/delivery/http"
	followServices "github.com/Khmer-Dev-Community/Services/api-service/lib/clientfollow"
	"github.com/gin-gonic/gin"
)

// NotificationControllerWrapper is a wrapper for the NotificationHandler.
type NotificationControllerWrapper struct {
	controller *HttpControllers.NotificationHandler
}

// NewNotificationControllerWrapper initializes the wrapper for NotificationHandler.
func NewNotificationControllerWrapper(nh *HttpControllers.NotificationHandler) *NotificationControllerWrapper {
	return &NotificationControllerWrapper{
		controller: nh,
	}
}

// FollowControllerWrapper is a wrapper for the FollowControllerHandler
type FollowControllerWrapper struct {
	controller *HttpControllers.FollowHandler
}

// NewFollowControllerWrapper initializes the wrapper for FollowControllerHandler
func NewFollowControllerWrapper(fh *HttpControllers.FollowHandler) *FollowControllerWrapper {
	return &FollowControllerWrapper{
		controller: fh,
	}
}

// SetupFollowRouter sets up the routes for the follow feature.
func SetupFollowRouter(r *gin.Engine, followService followServices.FollowService) {
	followHandler := HttpControllers.NewFollowHandler(followService)
	followHttp := NewFollowControllerWrapper(followHandler)

	api := r.Group("/api")
	followRouter := api.Group("/followers")
	{
		// To follow a user
		followRouter.POST("/:followedID", followHttp.controller.ToggleFollow)
		// To get a user's followers
		followRouter.GET("/:userID/followers", followHttp.controller.GetUserFollowers)
		// To get who a user is following
		followRouter.GET("/:userID/following", followHttp.controller.GetUserFollowing)
	}
}

func SetupNotificationRouter(r *gin.Engine, notificationService notificationServices.NotificationService) {
	notificationHandler := HttpControllers.NewHandler(notificationService)
	notificationHttp := NewNotificationControllerWrapper(notificationHandler)

	api := r.Group("/api")
	notificationRouter := api.Group("/notifications")
	{
		notificationRouter.GET("/", notificationHttp.controller.GetUserNotifications)
		notificationRouter.PUT("/:notificationID/read", notificationHttp.controller.MarkNotificationAsRead)
	}
}
