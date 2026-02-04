package payroll

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterConfigRoutes registers payroll configuration routes
func RegisterConfigRoutes(router *gin.RouterGroup, handler *ConfigHandler) {
	// Payroll configuration routes (Admin only)
	configGroup := router.Group("/config")
	configGroup.Use(middleware.RequireRole("admin"))
	{
		configGroup.POST("", handler.CreatePayrollConfiguration)
		configGroup.GET("", handler.ListPayrollConfigurations)
		configGroup.GET("/:id", handler.GetPayrollConfiguration)
		configGroup.PUT("/:id", handler.UpdatePayrollConfiguration)
		configGroup.DELETE("/:id", handler.DeletePayrollConfiguration)
	}

	// IRSA tax bracket routes (Admin only)
	bracketGroup := router.Group("/irsa-brackets")
	bracketGroup.Use(middleware.RequireRole("admin"))
	{
		bracketGroup.POST("", handler.CreateIRSABracket)
		bracketGroup.GET("", handler.ListIRSABrackets)
		bracketGroup.GET("/:id", handler.GetIRSABracketByID)
		bracketGroup.PUT("/:id", handler.UpdateIRSABracket)
		bracketGroup.DELETE("/:id", handler.DeleteIRSABracket)
	}
}
