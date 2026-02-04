package payroll

import (
	"encoding/json"
	"net/http"
	"time"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles payroll requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new payroll handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// CreateDraft handles creation of a new payroll draft (HR only)
func (h *Handler) CreateDraft(c *gin.Context) {
	var input CreatePayrollDraftRequest
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

	// Verify HR role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "hr" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only HR can create payroll drafts"})
		return
	}

	draft := &PayrollDraft{
		PeriodStart: input.PeriodStart,
		PeriodEnd:   input.PeriodEnd,
		EmployeeID:  input.EmployeeID,
		GrossSalary: input.GrossSalary,
		CreatedBy:   userID,
	}

	if err := h.repo.CreateDraft(c.Request.Context(), draft); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, draft)
}

// GetDraftByID retrieves a payroll draft by ID
func (h *Handler) GetDraftByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
		return
	}

	draft, err := h.repo.GetDraftByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payroll draft not found"})
		return
	}

	c.JSON(http.StatusOK, draft)
}

// UpdateDraft handles update of a payroll draft (HR only)
func (h *Handler) UpdateDraft(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
		return
	}

	// Verify HR role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "hr" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only HR can update payroll drafts"})
		return
	}

	var input UpdatePayrollDraftRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draft, err := h.repo.GetDraftByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payroll draft not found"})
		return
	}

	// Update fields if provided
	if input.GrossSalary != nil {
		draft.GrossSalary = *input.GrossSalary
	}

	if err := h.repo.UpdateDraft(c.Request.Context(), draft); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payroll draft"})
		return
	}

	c.JSON(http.StatusOK, draft)
}

// DeleteDraft handles deletion of a payroll draft (HR only)
func (h *Handler) DeleteDraft(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
		return
	}

	// Verify HR role
	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	if userRole != "hr" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only HR can delete payroll drafts"})
		return
	}

	if err := h.repo.DeleteDraft(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payroll draft deleted successfully"})
}

// ListDrafts retrieves payroll drafts with filtering
func (h *Handler) ListDrafts(c *gin.Context) {
	var query PayrollDraftListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	drafts, total, err := h.repo.ListDrafts(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list payroll drafts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"drafts": drafts,
		"total":  total,
		"limit":  query.Limit,
		"offset": query.Offset,
	})
}

// ApproveDraft handles approval of a payroll draft (Accountant only)
func (h *Handler) ApproveDraft(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
		return
	}

	// Verify Accountant role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "accountant" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Accountant can approve payroll drafts"})
		return
	}

	accountantID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get the draft
	draft, err := h.repo.GetDraftByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payroll draft not found"})
		return
	}

	// Check if already approved by querying the approved table
	// This is a simplified check - in production you'd query the database properly

	// Generate OHADA-compliant GL entries
	glEntries := generateGLEntries(draft)
	glEntriesJSON, err := json.Marshal(glEntries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate GL entries"})
		return
	}

	// Generate digital signature
	digitalSignature := generateDigitalSignature(accountantID, id)

	// Create approved record
	approved := &PayrollApproved{
		DraftID:          id,
		AccountantID:     accountantID,
		GLEntries:        string(glEntriesJSON),
		DigitalSignature: digitalSignature,
	}

	if err := h.repo.CreateApproved(c.Request.Context(), approved); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve payroll draft"})
		return
	}

	c.JSON(http.StatusCreated, approved)
}

// generateGLEntries generates OHADA-compliant general ledger entries
func generateGLEntries(draft *PayrollDraft) []GLEntry {
	entries := []GLEntry{
		// Debit 641 (Salaires) | Credit 421 (Salaires à payer) - Net salary
		{
			AccountCode: "641",
			AccountName: "Salaires et traitements",
			Debit:       draft.NetSalary,
			Credit:      0,
			Description: "Salaire net à payer",
		},
		{
			AccountCode: "421",
			AccountName: "Salaires à payer",
			Debit:       0,
			Credit:      draft.NetSalary,
			Description: "Salaire net à payer",
		},
		// Debit 641 (Salaires) | Credit 431 (CNAPS à payer) - Employee CNAPS portion
		{
			AccountCode: "641",
			AccountName: "Salaires et traitements",
			Debit:       draft.CNAPSEmployee,
			Credit:      0,
			Description: "CNAPS employé",
		},
		{
			AccountCode: "431",
			AccountName: "CNAPS à payer",
			Debit:       0,
			Credit:      draft.CNAPSEmployee,
			Description: "CNAPS employé",
		},
		// Debit 646 (Charges sociales) | Credit 431 (CNAPS à payer) - Employer CNAPS portion
		{
			AccountCode: "646",
			AccountName: "Charges sociales",
			Debit:       draft.CNAPSEmployer,
			Credit:      0,
			Description: "CNAPS employeur",
		},
		{
			AccountCode: "431",
			AccountName: "CNAPS à payer",
			Debit:       0,
			Credit:      draft.CNAPSEmployer,
			Description: "CNAPS employeur",
		},
		// Debit 641 (Salaires) | Credit 438 (OSTIE à payer) - Employee OSTIE portion
		{
			AccountCode: "641",
			AccountName: "Salaires et traitements",
			Debit:       draft.OSTIEEmployee,
			Credit:      0,
			Description: "OSTIE employé",
		},
		{
			AccountCode: "438",
			AccountName: "OSTIE à payer",
			Debit:       0,
			Credit:      draft.OSTIEEmployee,
			Description: "OSTIE employé",
		},
		// Debit 646 (Charges sociales) | Credit 438 (OSTIE à payer) - Employer OSTIE portion
		{
			AccountCode: "646",
			AccountName: "Charges sociales",
			Debit:       draft.OSTIEEmployer,
			Credit:      0,
			Description: "OSTIE employeur",
		},
		{
			AccountCode: "438",
			AccountName: "OSTIE à payer",
			Debit:       0,
			Credit:      draft.OSTIEEmployer,
			Description: "OSTIE employeur",
		},
		// Debit 641 (Salaires) | Credit 437 (IRSA à verser) - IRSA withholding
		{
			AccountCode: "641",
			AccountName: "Salaires et traitements",
			Debit:       draft.IRSA,
			Credit:      0,
			Description: "IRSA à verser",
		},
		{
			AccountCode: "437",
			AccountName: "IRSA à verser",
			Debit:       0,
			Credit:      draft.IRSA,
			Description: "IRSA à verser",
		},
	}

	return entries
}

// generateDigitalSignature generates a digital signature for the approved payroll
func generateDigitalSignature(accountantID, draftID uuid.UUID) string {
	// In production, this would use proper cryptographic signing
	// For now, we'll create a simple signature string
	timestamp := time.Now().Unix()
	return "SIG-" + accountantID.String()[:8] + "-" + draftID.String()[:8] + "-" + string(rune(timestamp))
}

// GetApprovedByID retrieves an approved payroll by ID
func (h *Handler) GetApprovedByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid approved payroll ID"})
		return
	}

	approved, err := h.repo.GetApprovedByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Approved payroll not found"})
		return
	}

	c.JSON(http.StatusOK, approved)
}

// GetApprovedByFichePaieNumber retrieves an approved payroll by fiche paie number
func (h *Handler) GetApprovedByFichePaieNumber(c *gin.Context) {
	fichePaieNumber := c.Param("fiche_paie_number")

	approved, err := h.repo.GetApprovedByFichePaieNumber(c.Request.Context(), fichePaieNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Approved payroll not found"})
		return
	}

	c.JSON(http.StatusOK, approved)
}

// ListApproved retrieves approved payrolls with filtering
func (h *Handler) ListApproved(c *gin.Context) {
	var query PayrollApprovedListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	approved, total, err := h.repo.ListApproved(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list approved payrolls"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"approved": approved,
		"total":    total,
		"limit":    query.Limit,
		"offset":   query.Offset,
	})
}

// GetReconciliationReport generates a reconciliation report for a period
func (h *Handler) GetReconciliationReport(c *gin.Context) {
	// Parse date range from query params
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
			return
		}
	} else {
		// Default to current month
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
			return
		}
	} else {
		// Default to last day of current month
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())
	}

	report, err := h.repo.GetReconciliationReport(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reconciliation report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GenerateFichePaie generates a fiche de paie (payslip) for an approved payroll
func (h *Handler) GenerateFichePaie(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid approved payroll ID"})
		return
	}

	// Get approved payroll
	approved, err := h.repo.GetApprovedByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Approved payroll not found"})
		return
	}

	// Get the draft
	draft, err := h.repo.GetDraftByID(c.Request.Context(), approved.DraftID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payroll draft not found"})
		return
	}

	// TODO: Get employee details from employee service
	// For now, we'll use placeholder values
	fichePaie := &FichePaie{
		FichePaieNumber:    approved.FichePaieNumber,
		EmployeeID:         draft.EmployeeID,
		EmployeeName:       "Employee Name", // TODO: Get from employee service
		EmployeePosition:   "Position",      // TODO: Get from employee service
		EmployeeDepartment: "Department",    // TODO: Get from employee service
		PeriodStart:        draft.PeriodStart,
		PeriodEnd:          draft.PeriodEnd,
		GrossSalary:        draft.GrossSalary,
		CNAPSEmployee:      draft.CNAPSEmployee,
		CNAPSEmployer:      draft.CNAPSEmployer,
		OSTIEEmployee:      draft.OSTIEEmployee,
		OSTIEEmployer:      draft.OSTIEEmployer,
		IRSA:               draft.IRSA,
		IRSABracket:        draft.IRSABracket,
		NetSalary:          draft.NetSalary,
		AccountantName:     "Accountant Name", // TODO: Get from user service
		ApprovedAt:         approved.ApprovedAt,
		DigitalSignature:   approved.DigitalSignature,
	}

	c.JSON(http.StatusOK, fichePaie)
}
