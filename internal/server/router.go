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
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(gormDB *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Configure CORS from environment variable CORS_ALLOWED_ORIGINS
	// Example: CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
	corsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowOrigins []string
	if corsEnv != "" {
		for _, o := range strings.Split(corsEnv, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				allowOrigins = append(allowOrigins, trimmed)
			}
		}
	}
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
