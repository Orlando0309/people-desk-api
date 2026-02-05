package server

import (
	"go-server/internal/attendance"
	"go-server/internal/audit"
	"go-server/internal/auth"
	"go-server/internal/dashboard"
	"go-server/internal/declarations"
	"go-server/internal/employee"
	"go-server/internal/kpi"
	"go-server/internal/leave"
	"go-server/internal/payroll"
	"go-server/internal/support"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(gormDB *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Go Server!",
			"ok": true, "status": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		// Register routes
		auth.RegisterRoutes(api, gormDB)
		employee.RegisterRoutes(api, gormDB)
		attendance.RegisterRoutes(api, gormDB)
		leave.RegisterRoutes(api, gormDB)
		audit.RegisterRoutes(api, gormDB)
		support.RegisterRoutes(api, gormDB)
		payroll.RegisterRoutes(api, gormDB)
		kpi.RegisterRoutes(api, gormDB)
		declarations.RegisterRoutes(api, gormDB)
		dashboard.RegisterRoutes(api, gormDB)
	}

	return r
}
