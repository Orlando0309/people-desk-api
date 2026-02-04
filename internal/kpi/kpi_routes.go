package kpi

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers KPI and performance review routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	kpi := rg.Group("/kpi")
	kpi.Use(middleware.AuthMiddleware())
	{
		// KPI Routes (HR/Admin only for write, authenticated for read)
		kpi.POST("", middleware.RequireRole("admin", "hr"), handler.CreateKPI)
		kpi.GET("", handler.ListKPIs)
		kpi.GET("/:id", handler.GetKPIByID)
		kpi.PUT("/:id", middleware.RequireRole("admin", "hr"), handler.UpdateKPI)
		kpi.DELETE("/:id", middleware.RequireRole("admin", "hr"), handler.DeleteKPI)

		// Performance Review Routes
		kpi.POST("/reviews", middleware.RequireRole("admin", "hr"), handler.CreatePerformanceReview)
		kpi.GET("/reviews", handler.ListPerformanceReviews)
		kpi.GET("/reviews/:id", handler.GetPerformanceReviewByID)
		kpi.PUT("/reviews/:id", handler.UpdatePerformanceReview)
		kpi.DELETE("/reviews/:id", middleware.RequireRole("admin", "hr"), handler.DeletePerformanceReview)
		kpi.POST("/reviews/:id/calculate", middleware.RequireRole("admin", "hr"), handler.CalculateFinalScore)

		// Performance Report Routes
		kpi.GET("/reports/:employee_id", handler.GeneratePerformanceReport)
	}
}
