package employee

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers employee routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	employees := rg.Group("/employees")
	employees.Use(middleware.AuthMiddleware())
	{
		// List and create employees (HR/Admin)
		employees.GET("", handler.List)
		employees.POST("", middleware.RequireRole("admin", "hr"), handler.Create)

		// Get departments and positions
		employees.GET("/departments", handler.GetDepartments)
		employees.GET("/positions", handler.GetPositions)

		// Individual employee operations
		employees.GET("/:id", handler.GetByID)
		employees.PUT("/:id", middleware.RequireRole("admin", "hr"), handler.Update)
		employees.DELETE("/:id", middleware.RequireRole("admin"), handler.Delete)

		// Get subordinates
		employees.GET("/:id/subordinates", handler.GetSubordinates)
	}
}
