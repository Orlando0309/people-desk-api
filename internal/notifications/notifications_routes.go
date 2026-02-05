package notifications

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers notification routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	notifications := rg.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware())
	{
		// Get all notifications for current user
		notifications.GET("", handler.List)

		// Get unread count
		notifications.GET("/unread-count", handler.GetUnreadCount)

		// Mark notification as read
		notifications.PUT("/:id/read", handler.MarkAsRead)

		// Mark all notifications as read
		notifications.PUT("/mark-all-read", handler.MarkAllAsRead)

		// Delete notification
		notifications.DELETE("/:id", handler.Delete)

		// Create notification (Admin/HR only)
		notifications.POST("", middleware.RequireRole("admin", "hr"), handler.Create)
	}
}
