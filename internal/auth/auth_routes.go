package auth

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers authentication routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	auth := rg.Group("/auth")
	{
		// Public routes
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)

		// Protected routes
		auth.POST("/register", middleware.AuthMiddleware(), middleware.RequireRole("admin"), handler.Register)
		auth.GET("/profile", middleware.AuthMiddleware(), handler.GetProfile)
		auth.PUT("/change-password", middleware.AuthMiddleware(), handler.ChangePassword)
		auth.GET("/users", middleware.AuthMiddleware(), middleware.RequireRole("admin"), handler.ListUsers)
	}
}
