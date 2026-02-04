package support

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)
	api := rg.Group("/support")
	{
		api.POST("", handler.Create)
	}
}
