package server

import (
	"net/http"
	"go-server/internal/support"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(gormDB *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Go Server!",
		"ok":true,"status": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		// Register other routes here
		support.RegisterRoutes(api, gormDB)
	}

	

	
	return r
}