package leave

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers leave routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	leaves := rg.Group("/leaves")
	leaves.Use(middleware.AuthMiddleware())
	{
		// List and create leaves
		leaves.GET("", handler.List)
		leaves.POST("", handler.Create)

		// Get pending leaves (HR/Admin only)
		leaves.GET("/pending", middleware.RequireRole("admin", "hr"), handler.GetPendingLeaves)

		// Get leave balance
		leaves.GET("/balance/:employee_id", handler.GetLeaveBalance)

		// Individual leave operations
		leaves.GET("/:id", handler.GetByID)
		leaves.PUT("/:id", handler.Update)
		leaves.DELETE("/:id", handler.Delete)

		// Approve/reject leaves (HR/Admin only)
		leaves.PUT("/:id/approve", middleware.RequireRole("admin", "hr"), handler.Approve)
		leaves.PUT("/:id/reject", middleware.RequireRole("admin", "hr"), handler.Reject)
	}
}
