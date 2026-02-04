package payroll

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers payroll routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	// Create config handler for payroll configuration
	configRepo := NewConfigRepo(gormDB)
	configHandler := NewConfigHandler(configRepo)

	payroll := rg.Group("/payroll")
	payroll.Use(middleware.AuthMiddleware())
	{
		// HR Draft Routes (HR/Admin only)
		payroll.POST("/drafts", middleware.RequireRole("admin", "hr"), handler.CreateDraft)
		payroll.GET("/drafts", handler.ListDrafts)
		payroll.GET("/drafts/:id", handler.GetDraftByID)
		payroll.PUT("/drafts/:id", middleware.RequireRole("admin", "hr"), handler.UpdateDraft)
		payroll.DELETE("/drafts/:id", middleware.RequireRole("admin", "hr"), handler.DeleteDraft)

		// Accountant Approval Routes (Accountant/Admin only)
		payroll.PUT("/drafts/:id/approve", middleware.RequireRole("admin", "accountant"), handler.ApproveDraft)

		// Approved Payroll Routes (Accountant/Admin only for write, authenticated for read)
		payroll.GET("/approved", handler.ListApproved)
		payroll.GET("/approved/:id", handler.GetApprovedByID)
		payroll.GET("/approved/fiche/:fiche_paie_number", handler.GetApprovedByFichePaieNumber)
		payroll.GET("/approved/:id/fiche-paie", handler.GenerateFichePaie)

		// Reconciliation Report (Accountant/Admin only)
		payroll.GET("/reconciliation", middleware.RequireRole("admin", "accountant"), handler.GetReconciliationReport)
	}

	// Register payroll configuration routes
	RegisterConfigRoutes(rg, configHandler)
}
