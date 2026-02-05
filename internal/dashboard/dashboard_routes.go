package dashboard

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers dashboard routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	dashboard := rg.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware())
	{
		// Get dashboard statistics
		dashboard.GET("/stats", handler.GetStats)

		// Get attendance summary
		dashboard.GET("/attendance/summary", handler.GetAttendanceSummary)

		// Get leave balances
		dashboard.GET("/leaves/balances", handler.GetLeaveBalances)
	}

	// User management routes (admin only)
	users := rg.Group("/auth/users")
	users.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		// Update user
		users.PUT("/:id", handler.UpdateUser)

		// Delete user
		users.DELETE("/:id", handler.DeleteUser)
	}
}
