package attendance

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers attendance routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	attendance := rg.Group("/attendance")
	attendance.Use(middleware.AuthMiddleware())
	{
		// Clock in/out (all authenticated users)
		attendance.POST("/clock-in", handler.ClockIn)
		attendance.POST("/clock-out", handler.ClockOut)

		// Get today's attendance
		attendance.GET("/today/:employee_id", handler.GetTodayAttendance)

		// List attendance records
		attendance.GET("", handler.List)

		// Get attendance by ID
		attendance.GET("/:id", handler.GetByID)

		// Update attendance (HR/Admin only)
		attendance.PUT("/:id", middleware.RequireRole("admin", "hr"), handler.UpdateAttendance)

		// Get attendance statistics
		attendance.GET("/stats/:employee_id", handler.GetStats)
	}
}
