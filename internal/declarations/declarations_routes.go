package declarations

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers declarations routes
func RegisterRoutes(rg *gin.RouterGroup, gormDB *gorm.DB) {
	repo := NewRepo(gormDB)
	handler := NewHandler(repo)

	declarations := rg.Group("/declarations")
	declarations.Use(middleware.AuthMiddleware())
	{
		// Monthly Declaration Routes (Accountant only for write, authenticated for read)
		declarations.POST("", middleware.RequireRole("admin", "accountant"), handler.CreateDeclaration)
		declarations.GET("", handler.ListDeclarations)
		declarations.GET("/:id", handler.GetDeclarationByID)
		declarations.GET("/number/:declaration_number", handler.GetDeclarationByNumber)
		declarations.PUT("/:id", middleware.RequireRole("admin", "accountant"), handler.UpdateDeclaration)
		declarations.DELETE("/:id", middleware.RequireRole("admin", "accountant"), handler.DeleteDeclaration)
		declarations.GET("/:id/form", handler.GenerateDeclarationForm)
		declarations.POST("/:id/populate", middleware.RequireRole("admin", "accountant"), handler.PopulateDeclarationData)

		// CNAPS Declaration Generation (Accountant only)
		declarations.GET("/cnaps/generate", middleware.RequireRole("admin", "accountant"), handler.GenerateCNAPSDeclaration)

		// OSTIE Declaration Generation (Accountant only)
		declarations.GET("/ostie/generate", middleware.RequireRole("admin", "accountant"), handler.GenerateOSTIEDeclaration)

		// IRSA Declaration Generation (Accountant only)
		declarations.GET("/irsa/generate", middleware.RequireRole("admin", "accountant"), handler.GenerateIRSADeclaration)

		// IRSA Tax Bracket Routes (Admin only)
		declarations.POST("/irsa-brackets", middleware.RequireRole("admin"), handler.CreateIRSATaxBracket)
		declarations.GET("/irsa-brackets", handler.ListIRSATaxBrackets)
		declarations.GET("/irsa-brackets/:id", handler.GetIRSATaxBracketByID)
		declarations.PUT("/irsa-brackets/:id", middleware.RequireRole("admin"), handler.UpdateIRSATaxBracket)
		declarations.DELETE("/irsa-brackets/:id", middleware.RequireRole("admin"), handler.DeleteIRSATaxBracket)
	}
}
