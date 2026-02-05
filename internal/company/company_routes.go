package company

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers company routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	company := rg.Group("/company")
	company.Use(middleware.AuthMiddleware())
	{
		// Get company settings (all authenticated users)
		company.GET("/settings", handler.Get)

		// Update company settings (Admin only)
		company.PUT("/settings", middleware.RequireRole("admin"), handler.Update)

		// Upload company logo (Admin only)
		company.POST("/logo", middleware.RequireRole("admin"), handler.UploadLogo)

		// Company holidays routes
		company.GET("/holidays", handler.ListHolidays)
		company.POST("/holidays", middleware.RequireRole("admin"), handler.CreateHoliday)
		company.PUT("/holidays/:id", middleware.RequireRole("admin"), handler.UpdateHoliday)
		company.DELETE("/holidays/:id", middleware.RequireRole("admin"), handler.DeleteHoliday)
	}
}
