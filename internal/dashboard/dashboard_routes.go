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

		// Get monthly payroll summary (Admin/HR/Accountant only)
		dashboard.GET("/payroll/monthly", middleware.RequireRole("admin", "hr", "accountant"), handler.GetMonthlyPayroll)

		// Get compliance status (Admin/HR/Accountant only)
		dashboard.GET("/compliance", middleware.RequireRole("admin", "hr", "accountant"), handler.GetCompliance)

		// Get badge counts for navigation items
		dashboard.GET("/badges", handler.GetBadgeCounts)

		// Get weekly attendance summary
		dashboard.GET("/attendance/weekly", handler.GetWeeklyAttendance)

		// Get calendar events
		dashboard.GET("/calendar/events", handler.GetCalendarEvents)
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
