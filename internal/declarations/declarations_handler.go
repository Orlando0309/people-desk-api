package declarations

import (
	"encoding/json"
	"net/http"
	"time"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles declaration requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new declarations handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// CreateDeclaration handles creation of a new monthly declaration (Accountant only)
func (h *Handler) CreateDeclaration(c *gin.Context) {
	var input CreateDeclarationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	accountantID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify Accountant role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "accountant" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Accountant can create declarations"})
		return
	}

	declaration := &MonthlyDeclaration{
		DeclarationType:        input.DeclarationType,
		DeclarationPeriodStart: input.DeclarationPeriodStart,
		DeclarationPeriodEnd:   input.DeclarationPeriodEnd,
		CompanyName:            input.CompanyName,
		CompanyAddress:         input.CompanyAddress,
		CompanyNIF:             input.CompanyNIF,
		AccountantID:           accountantID,
		Status:                 "draft",
		DeclarationData:        "[]", // Empty JSON array for employee breakdown
	}

	if err := h.repo.CreateDeclaration(c.Request.Context(), declaration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, declaration)
}

// GetDeclarationByID retrieves a declaration by ID
func (h *Handler) GetDeclarationByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid declaration ID"})
		return
	}

	declaration, err := h.repo.GetDeclarationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Declaration not found"})
		return
	}

	c.JSON(http.StatusOK, declaration)
}

// GetDeclarationByNumber retrieves a declaration by declaration number
func (h *Handler) GetDeclarationByNumber(c *gin.Context) {
	declarationNumber := c.Param("declaration_number")

	declaration, err := h.repo.GetDeclarationByNumber(c.Request.Context(), declarationNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Declaration not found"})
		return
	}

	c.JSON(http.StatusOK, declaration)
}

// UpdateDeclaration handles update of a declaration (Accountant only)
func (h *Handler) UpdateDeclaration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid declaration ID"})
		return
	}

	// Verify Accountant role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "accountant" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Accountant can update declarations"})
		return
	}

	var input UpdateDeclarationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	declaration, err := h.repo.GetDeclarationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Declaration not found"})
		return
	}

	// Update fields if provided
	if input.Status != nil {
		declaration.Status = *input.Status
		now := time.Now()
		if *input.Status == "submitted" && declaration.SubmittedAt == nil {
			declaration.SubmittedAt = &now
		}
		if *input.Status == "paid" && declaration.PaidAt == nil {
			declaration.PaidAt = &now
		}
	}

	if err := h.repo.UpdateDeclaration(c.Request.Context(), declaration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update declaration"})
		return
	}

	c.JSON(http.StatusOK, declaration)
}

// DeleteDeclaration handles deletion of a declaration (Accountant only)
func (h *Handler) DeleteDeclaration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid declaration ID"})
		return
	}

	// Verify Accountant role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "accountant" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Accountant can delete declarations"})
		return
	}

	if err := h.repo.DeleteDeclaration(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Declaration deleted successfully"})
}

// ListDeclarations retrieves declarations with filtering
func (h *Handler) ListDeclarations(c *gin.Context) {
	var query DeclarationListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	declarations, total, err := h.repo.ListDeclarations(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list declarations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"declarations": declarations,
		"total":        total,
		"limit":        query.Limit,
		"offset":       query.Offset,
	})
}

// GenerateDeclarationForm generates a declaration form for CNAPS, OSTIE, or IRSA
func (h *Handler) GenerateDeclarationForm(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid declaration ID"})
		return
	}

	form, err := h.repo.GenerateDeclarationForm(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate declaration form"})
		return
	}

	c.JSON(http.StatusOK, form)
}

// CreateIRSATaxBracket handles creation of a new IRSA tax bracket (Admin only)
func (h *Handler) CreateIRSATaxBracket(c *gin.Context) {
	var input CreateIRSATaxBracketRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify Admin role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Admin can create IRSA tax brackets"})
		return
	}

	bracket := &IRSATaxBracket{
		MinIncome:     input.MinIncome,
		MaxIncome:     input.MaxIncome,
		TaxRate:       input.TaxRate,
		MinTax:        input.MinTax,
		EffectiveDate: input.EffectiveDate,
		CreatedBy:     userID,
		IsActive:      true,
	}

	if err := h.repo.CreateIRSATaxBracket(c.Request.Context(), bracket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bracket)
}

// GetIRSATaxBracketByID retrieves an IRSA tax bracket by ID
func (h *Handler) GetIRSATaxBracketByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IRSA tax bracket ID"})
		return
	}

	bracket, err := h.repo.GetIRSATaxBracketByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "IRSA tax bracket not found"})
		return
	}

	c.JSON(http.StatusOK, bracket)
}

// UpdateIRSATaxBracket handles update of an IRSA tax bracket (Admin only)
func (h *Handler) UpdateIRSATaxBracket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IRSA tax bracket ID"})
		return
	}

	// Verify Admin role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Admin can update IRSA tax brackets"})
		return
	}

	var input UpdateIRSATaxBracketRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bracket, err := h.repo.GetIRSATaxBracketByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "IRSA tax bracket not found"})
		return
	}

	// Update fields if provided
	if input.MinIncome != nil {
		bracket.MinIncome = *input.MinIncome
	}
	if input.MaxIncome != nil {
		bracket.MaxIncome = input.MaxIncome
	}
	if input.TaxRate != nil {
		bracket.TaxRate = *input.TaxRate
	}
	if input.MinTax != nil {
		bracket.MinTax = *input.MinTax
	}
	if input.IsActive != nil {
		bracket.IsActive = *input.IsActive
	}
	if input.EffectiveDate != nil {
		bracket.EffectiveDate = *input.EffectiveDate
	}

	if err := h.repo.UpdateIRSATaxBracket(c.Request.Context(), bracket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update IRSA tax bracket"})
		return
	}

	c.JSON(http.StatusOK, bracket)
}

// DeleteIRSATaxBracket handles deletion of an IRSA tax bracket (Admin only)
func (h *Handler) DeleteIRSATaxBracket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IRSA tax bracket ID"})
		return
	}

	// Verify Admin role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Admin can delete IRSA tax brackets"})
		return
	}

	if err := h.repo.DeleteIRSATaxBracket(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "IRSA tax bracket deleted successfully"})
}

// ListIRSATaxBrackets retrieves IRSA tax brackets with filtering
func (h *Handler) ListIRSATaxBrackets(c *gin.Context) {
	var query IRSATaxBracketListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	brackets, total, err := h.repo.ListIRSATaxBrackets(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list IRSA tax brackets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"brackets": brackets,
		"total":    total,
		"limit":    query.Limit,
		"offset":   query.Offset,
	})
}

// GenerateCNAPSDeclaration generates a CNAPS declaration form for a specific month
func (h *Handler) GenerateCNAPSDeclaration(c *gin.Context) {
	// Parse month from query params (format: YYYY-MM)
	monthStr := c.Query("month")
	if monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month parameter is required (format: YYYY-MM)"})
		return
	}

	monthTime, err := time.Parse("2006-01", monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format (use YYYY-MM)"})
		return
	}

	// Calculate period start and end
	periodStart := time.Date(monthTime.Year(), monthTime.Month(), 1, 0, 0, 0, 0, monthTime.Location())
	periodEnd := time.Date(monthTime.Year(), monthTime.Month()+1, 0, 0, 0, 0, 0, monthTime.Location())

	// Get user info
	accountantID, _ := middleware.GetUserID(c)

	// Check if declaration already exists
	var existing MonthlyDeclaration
	err = h.repo.db.WithContext(c.Request.Context()).
		Where("declaration_type = ? AND declaration_period_start = ? AND declaration_period_end = ?",
			"cnaps", periodStart, periodEnd).First(&existing).Error

	if err == nil {
		// Declaration exists, return it
		form, err := h.repo.GenerateDeclarationForm(c.Request.Context(), existing.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate declaration form"})
			return
		}
		c.JSON(http.StatusOK, form)
		return
	}

	// Create new declaration
	declaration := &MonthlyDeclaration{
		DeclarationType:        "cnaps",
		DeclarationPeriodStart: periodStart,
		DeclarationPeriodEnd:   periodEnd,
		CompanyName:            "Company Name",    // TODO: Get from company settings
		CompanyAddress:         "Company Address", // TODO: Get from company settings
		CompanyNIF:             "NIF",             // TODO: Get from company settings
		AccountantID:           accountantID,
		Status:                 "draft",
		DeclarationData:        "[]",
	}

	if err := h.repo.CreateDeclaration(c.Request.Context(), declaration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CNAPS declaration"})
		return
	}

	form, err := h.repo.GenerateDeclarationForm(c.Request.Context(), declaration.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate declaration form"})
		return
	}

	c.JSON(http.StatusOK, form)
}

// GenerateOSTIEDeclaration generates an OSTIE declaration form for a specific month
func (h *Handler) GenerateOSTIEDeclaration(c *gin.Context) {
	// Parse month from query params (format: YYYY-MM)
	monthStr := c.Query("month")
	if monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month parameter is required (format: YYYY-MM)"})
		return
	}

	monthTime, err := time.Parse("2006-01", monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format (use YYYY-MM)"})
		return
	}

	// Calculate period start and end
	periodStart := time.Date(monthTime.Year(), monthTime.Month(), 1, 0, 0, 0, 0, monthTime.Location())
	periodEnd := time.Date(monthTime.Year(), monthTime.Month()+1, 0, 0, 0, 0, 0, monthTime.Location())

	// Get user info
	accountantID, _ := middleware.GetUserID(c)

	// Check if declaration already exists
	var existing MonthlyDeclaration
	err = h.repo.db.WithContext(c.Request.Context()).
		Where("declaration_type = ? AND declaration_period_start = ? AND declaration_period_end = ?",
			"ostie", periodStart, periodEnd).First(&existing).Error

	if err == nil {
		// Declaration exists, return it
		form, err := h.repo.GenerateDeclarationForm(c.Request.Context(), existing.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate declaration form"})
			return
		}
		c.JSON(http.StatusOK, form)
		return
	}

	// Create new declaration
	declaration := &MonthlyDeclaration{
		DeclarationType:        "ostie",
		DeclarationPeriodStart: periodStart,
		DeclarationPeriodEnd:   periodEnd,
		CompanyName:            "Company Name",    // TODO: Get from company settings
		CompanyAddress:         "Company Address", // TODO: Get from company settings
		CompanyNIF:             "NIF",             // TODO: Get from company settings
		AccountantID:           accountantID,
		Status:                 "draft",
		DeclarationData:        "[]",
	}

	if err := h.repo.CreateDeclaration(c.Request.Context(), declaration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OSTIE declaration"})
		return
	}

	form, err := h.repo.GenerateDeclarationForm(c.Request.Context(), declaration.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate declaration form"})
		return
	}

	c.JSON(http.StatusOK, form)
}

// GenerateIRSADeclaration generates an IRSA declaration form for a specific month
func (h *Handler) GenerateIRSADeclaration(c *gin.Context) {
	// Parse month from query params (format: YYYY-MM)
	monthStr := c.Query("month")
	if monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month parameter is required (format: YYYY-MM)"})
		return
	}

	monthTime, err := time.Parse("2006-01", monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format (use YYYY-MM)"})
		return
	}

	// Calculate period start and end
	periodStart := time.Date(monthTime.Year(), monthTime.Month(), 1, 0, 0, 0, 0, monthTime.Location())
	periodEnd := time.Date(monthTime.Year(), monthTime.Month()+1, 0, 0, 0, 0, 0, monthTime.Location())

	// Get user info
	accountantID, _ := middleware.GetUserID(c)

	// Check if declaration already exists
	var existing MonthlyDeclaration
	err = h.repo.db.WithContext(c.Request.Context()).
		Where("declaration_type = ? AND declaration_period_start = ? AND declaration_period_end = ?",
			"irsa", periodStart, periodEnd).First(&existing).Error

	if err == nil {
		// Declaration exists, return it
		form, err := h.repo.GenerateDeclarationForm(c.Request.Context(), existing.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate declaration form"})
			return
		}
		c.JSON(http.StatusOK, form)
		return
	}

	// Create new declaration
	declaration := &MonthlyDeclaration{
		DeclarationType:        "irsa",
		DeclarationPeriodStart: periodStart,
		DeclarationPeriodEnd:   periodEnd,
		CompanyName:            "Company Name",    // TODO: Get from company settings
		CompanyAddress:         "Company Address", // TODO: Get from company settings
		CompanyNIF:             "NIF",             // TODO: Get from company settings
		AccountantID:           accountantID,
		Status:                 "draft",
		DeclarationData:        "[]",
	}

	if err := h.repo.CreateDeclaration(c.Request.Context(), declaration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create IRSA declaration"})
		return
	}

	form, err := h.repo.GenerateDeclarationForm(c.Request.Context(), declaration.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate declaration form"})
		return
	}

	c.JSON(http.StatusOK, form)
}

// PopulateDeclarationData populates declaration data from approved payrolls
func (h *Handler) PopulateDeclarationData(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid declaration ID"})
		return
	}

	// Verify Accountant role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "accountant" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Accountant can populate declaration data"})
		return
	}

	declaration, err := h.repo.GetDeclarationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Declaration not found"})
		return
	}

	// TODO: Query approved payrolls for the period and populate employee breakdown
	// This is a placeholder implementation
	employeeBreakdown := []EmployeeBreakdown{}

	// Marshal to JSON
	declarationData, err := json.Marshal(employeeBreakdown)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal declaration data"})
		return
	}

	declaration.DeclarationData = string(declarationData)
	declaration.TotalEmployees = len(employeeBreakdown)

	if err := h.repo.UpdateDeclaration(c.Request.Context(), declaration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update declaration"})
		return
	}

	c.JSON(http.StatusOK, declaration)
}
