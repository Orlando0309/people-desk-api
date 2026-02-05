package support_tickets

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers support ticket routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	tickets := rg.Group("/support/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		// List all support tickets
		tickets.GET("", handler.List)

		// Get ticket details with replies
		tickets.GET("/:id", handler.GetByID)

		// Create a new support ticket
		tickets.POST("", handler.Create)

		// Update ticket (Admin/HR only)
		tickets.PUT("/:id", middleware.RequireRole("admin", "hr"), handler.Update)

		// Add reply to ticket
		tickets.POST("/:id/reply", handler.Reply)

		// Resolve ticket (Admin/HR only)
		tickets.PUT("/:id/resolve", middleware.RequireRole("admin", "hr"), handler.Resolve)
	}

	// Support categories (no role restriction)
	rg.GET("/support/categories", middleware.AuthMiddleware(), handler.GetCategories)
}
