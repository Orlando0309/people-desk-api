package kpi

import (
	"net/http"
	"time"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles KPI and performance review requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new KPI handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// CreateKPI handles creation of a new KPI (HR/Admin only)
func (h *Handler) CreateKPI(c *gin.Context) {
	var input CreateKPIRequest
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

	// Verify HR or Admin role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "hr" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only HR and Admin can create KPIs"})
		return
	}

	kpi := &KPI{
		Name:             input.Name,
		Description:      input.Description,
		TargetValue:      input.TargetValue,
		WeightPercentage: input.WeightPercentage,
		ScoringScale:     input.ScoringScale,
		Department:       input.Department,
		Position:         input.Position,
		CreatedBy:        userID,
	}

	if kpi.ScoringScale == "" {
		kpi.ScoringScale = "1_to_5"
	}

	if err := h.repo.CreateKPI(c.Request.Context(), kpi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, kpi)
}

// GetKPIByID retrieves a KPI by ID
func (h *Handler) GetKPIByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KPI ID"})
		return
	}

	kpi, err := h.repo.GetKPIByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "KPI not found"})
		return
	}

	c.JSON(http.StatusOK, kpi)
}

// UpdateKPI handles update of a KPI (HR/Admin only)
func (h *Handler) UpdateKPI(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KPI ID"})
		return
	}

	// Verify HR or Admin role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "hr" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only HR and Admin can update KPIs"})
		return
	}

	var input UpdateKPIRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpi, err := h.repo.GetKPIByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "KPI not found"})
		return
	}

	// Update fields if provided
	if input.Name != nil {
		kpi.Name = *input.Name
	}
	if input.Description != nil {
		kpi.Description = *input.Description
	}
	if input.TargetValue != nil {
		kpi.TargetValue = *input.TargetValue
	}
	if input.WeightPercentage != nil {
		kpi.WeightPercentage = *input.WeightPercentage
	}
	if input.ScoringScale != nil {
		kpi.ScoringScale = *input.ScoringScale
	}
	if input.Department != nil {
		kpi.Department = *input.Department
	}
	if input.Position != nil {
		kpi.Position = *input.Position
	}
	if input.IsActive != nil {
		kpi.IsActive = *input.IsActive
	}

	if err := h.repo.UpdateKPI(c.Request.Context(), kpi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update KPI"})
		return
	}

	c.JSON(http.StatusOK, kpi)
}

// DeleteKPI handles deletion of a KPI (HR/Admin only)
func (h *Handler) DeleteKPI(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KPI ID"})
		return
	}

	// Verify HR or Admin role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "hr" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only HR and Admin can delete KPIs"})
		return
	}

	if err := h.repo.DeleteKPI(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KPI deleted successfully"})
}

// ListKPIs retrieves KPIs with filtering
func (h *Handler) ListKPIs(c *gin.Context) {
	var query KPIListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpis, total, err := h.repo.ListKPIs(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list KPIs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kpis":   kpis,
		"total":  total,
		"limit":  query.Limit,
		"offset": query.Offset,
	})
}

// CreatePerformanceReview handles creation of a new performance review (HR/Admin only)
func (h *Handler) CreatePerformanceReview(c *gin.Context) {
	var input CreatePerformanceReviewRequest
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

	// Verify HR or Admin role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "hr" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only HR and Admin can create performance reviews"})
		return
	}

	review := &PerformanceReview{
		EmployeeID:        input.EmployeeID,
		KPIID:             input.KPIID,
		ReviewPeriodStart: input.ReviewPeriodStart,
		ReviewPeriodEnd:   input.ReviewPeriodEnd,
		SelfScore:         input.SelfScore,
		ManagerScore:      input.ManagerScore,
		SelfAssessment:    input.SelfAssessment,
		ManagerAssessment: input.ManagerAssessment,
		ReviewerID:        userID,
		Status:            "pending",
	}

	if err := h.repo.CreatePerformanceReview(c.Request.Context(), review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetPerformanceReviewByID retrieves a performance review by ID
func (h *Handler) GetPerformanceReviewByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid performance review ID"})
		return
	}

	review, err := h.repo.GetPerformanceReviewByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Performance review not found"})
		return
	}

	// Verify employee can only view their own reviews
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" {
		userID, _ := middleware.GetUserID(c)
		if review.EmployeeID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot view another employee's performance review"})
			return
		}
	}

	c.JSON(http.StatusOK, review)
}

// UpdatePerformanceReview handles update of a performance review
func (h *Handler) UpdatePerformanceReview(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid performance review ID"})
		return
	}

	var input UpdatePerformanceReviewRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.repo.GetPerformanceReviewByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Performance review not found"})
		return
	}

	// Verify permissions
	userRole, _ := middleware.GetUserRole(c)
	userID, _ := middleware.GetUserID(c)

	if userRole == "employee" {
		// Employee can only update their self-assessment and self-score
		if review.EmployeeID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update another employee's performance review"})
			return
		}
		if input.SelfAssessment != nil {
			review.SelfAssessment = *input.SelfAssessment
		}
		if input.SelfScore != nil {
			review.SelfScore = input.SelfScore
		}
	} else if userRole == "hr" || userRole == "admin" {
		// HR/Admin can update all fields
		if input.SelfScore != nil {
			review.SelfScore = input.SelfScore
		}
		if input.ManagerScore != nil {
			review.ManagerScore = input.ManagerScore
		}
		if input.FinalScore != nil {
			review.FinalScore = input.FinalScore
		}
		if input.SelfAssessment != nil {
			review.SelfAssessment = *input.SelfAssessment
		}
		if input.ManagerAssessment != nil {
			review.ManagerAssessment = *input.ManagerAssessment
		}
		if input.Status != nil {
			review.Status = *input.Status
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	if err := h.repo.UpdatePerformanceReview(c.Request.Context(), review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update performance review"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// DeletePerformanceReview handles deletion of a performance review (HR/Admin only)
func (h *Handler) DeletePerformanceReview(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid performance review ID"})
		return
	}

	// Verify HR or Admin role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "hr" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only HR and Admin can delete performance reviews"})
		return
	}

	if err := h.repo.DeletePerformanceReview(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Performance review deleted successfully"})
}

// ListPerformanceReviews retrieves performance reviews with filtering
func (h *Handler) ListPerformanceReviews(c *gin.Context) {
	var query PerformanceReviewListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If employee role, only show their own reviews
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" {
		userID, _ := middleware.GetUserID(c)
		query.EmployeeID = &userID
	}

	reviews, total, err := h.repo.ListPerformanceReviews(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list performance reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"total":   total,
		"limit":   query.Limit,
		"offset":  query.Offset,
	})
}

// CalculateFinalScore calculates the final score for a performance review (HR/Admin only)
func (h *Handler) CalculateFinalScore(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid performance review ID"})
		return
	}

	// Verify HR or Admin role
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "hr" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only HR and Admin can calculate final scores"})
		return
	}

	review, err := h.repo.CalculateFinalScore(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate final score"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// GeneratePerformanceReport generates a performance report for an employee
func (h *Handler) GeneratePerformanceReport(c *gin.Context) {
	employeeIDParam := c.Param("employee_id")
	employeeID, err := uuid.Parse(employeeIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Verify employee can only view their own report (unless HR/Admin)
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" {
		userID, _ := middleware.GetUserID(c)
		if employeeID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot view another employee's performance report"})
			return
		}
	}

	// Parse date range from query params
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
			return
		}
	} else {
		// Default to current quarter
		now := time.Now()
		quarter := (now.Month() - 1) / 3
		startDate = time.Date(now.Year(), quarter*3+1, 1, 0, 0, 0, 0, now.Location())
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
			return
		}
	} else {
		// Default to end of current quarter
		now := time.Now()
		quarter := (now.Month() - 1) / 3
		endDate = time.Date(now.Year(), quarter*3+4, 0, 0, 0, 0, 0, now.Location())
	}

	report, err := h.repo.GeneratePerformanceReport(c.Request.Context(), employeeID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate performance report"})
		return
	}

	c.JSON(http.StatusOK, report)
}
