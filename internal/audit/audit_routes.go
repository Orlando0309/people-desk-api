package audit

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers audit routes (Admin only)
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	audit := rg.Group("/audit")
	audit.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		// List and export audit logs
		audit.GET("/logs", handler.List)
		audit.GET("/logs/export", handler.ExportLogs)

		// Get audit log by ID
		audit.GET("/logs/:id", handler.GetByID)

		// Get audit logs by user
		audit.GET("/logs/user/:user_id", handler.GetByUserID)

		// Get audit logs by module
		audit.GET("/logs/module/:module", handler.GetByModule)

		// Get audit statistics
		audit.GET("/stats", handler.GetStats)

		// Get record history
		audit.GET("/history/:record_id", handler.GetRecordHistory)
	}
}
