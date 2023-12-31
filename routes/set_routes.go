package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetRoutes(router *gin.Engine, db *gorm.DB) {
	UserRoutes(router.Group("/users"), db)
	FriendshipRoutes(router.Group("/friendships"), db)
	AuthRoutes(router.Group("/auth"), db)
	PostRoutes(router.Group("/posts"), db)
	TimelineRoutes(router.Group("/timeline"), db)
	CommentRoutes(router.Group("/comments"), db)
	NotificationRoutes(router.Group("/notifications"), db)
}
