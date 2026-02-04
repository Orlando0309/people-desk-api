package leave

import (
	"net/http"
	"time"

	"go-server/internal/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles leave requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new leave handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// Create handles leave request creation
func (h *Handler) Create(c *gin.Context) {
	var input CreateLeaveRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate dates
	if input.EndDate.Before(input.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// Check for overlapping leave requests
	overlap, err := h.repo.CheckOverlap(c.Request.Context(), input.EmployeeID, input.StartDate, input.EndDate, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check overlapping leaves"})
		return
	}
	if overlap {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave request overlaps with existing leave"})
		return
	}

	// Check leave balance for annual leave
	if input.LeaveType == "annual" {
		balance, err := h.repo.GetLeaveBalance(c.Request.Context(), input.EmployeeID, time.Now().Year())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check leave balance"})
			return
		}
		if balance.AnnualRemaining < input.DaysRequested {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient annual leave balance"})
			return
		}
	}

	leave := &Leave{
		EmployeeID:    input.EmployeeID,
		LeaveType:     input.LeaveType,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		DaysRequested: input.DaysRequested,
		Reason:        input.Reason,
		Status:        "pending",
	}

	if err := h.repo.Create(c.Request.Context(), leave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create leave request"})
		return
	}

	c.JSON(http.StatusCreated, leave)
}

// GetByID retrieves a leave by ID
func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave ID"})
		return
	}

	leave, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave not found"})
		return
	}

	c.JSON(http.StatusOK, leave)
}

// Update updates a leave request (only if pending)
func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave ID"})
		return
	}

	var input UpdateLeaveRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leave, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave not found"})
		return
	}

	// Only allow updates to pending leaves
	if leave.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update non-pending leave"})
		return
	}

	// Update fields if provided
	if input.LeaveType != nil {
		leave.LeaveType = *input.LeaveType
	}
	if input.StartDate != nil {
		leave.StartDate = *input.StartDate
	}
	if input.EndDate != nil {
		leave.EndDate = *input.EndDate
	}
	if input.DaysRequested != nil {
		leave.DaysRequested = *input.DaysRequested
	}
	if input.Reason != nil {
		leave.Reason = *input.Reason
	}

	// Validate dates
	if leave.EndDate.Before(leave.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// Check for overlapping leave requests (excluding current leave)
	overlap, err := h.repo.CheckOverlap(c.Request.Context(), leave.EmployeeID, leave.StartDate, leave.EndDate, &id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check overlapping leaves"})
		return
	}
	if overlap {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave request overlaps with existing leave"})
		return
	}

	if err := h.repo.Update(c.Request.Context(), leave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leave"})
		return
	}

	c.JSON(http.StatusOK, leave)
}

// Delete cancels a leave request
func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave ID"})
		return
	}

	leave, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave not found"})
		return
	}

	// Only allow cancellation of pending or approved leaves
	if leave.Status != "pending" && leave.Status != "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel this leave"})
		return
	}

	leave.Status = "cancelled"
	if err := h.repo.Update(c.Request.Context(), leave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel leave"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leave cancelled successfully"})
}

// List retrieves leaves with filtering
func (h *Handler) List(c *gin.Context) {
	var query LeaveListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If employee role, only show their own leaves
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" {
		userID, _ := middleware.GetUserID(c)
		query.EmployeeID = &userID
	}

	leaves, total, err := h.repo.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list leaves"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leaves": leaves,
		"total":  total,
		"limit":  query.Limit,
		"offset": query.Offset,
	})
}

// Approve approves a leave request (HR/Admin only)
func (h *Handler) Approve(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave ID"})
		return
	}

	approverID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	leave, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave not found"})
		return
	}

	if leave.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave is not pending"})
		return
	}

	now := time.Now()
	leave.Status = "approved"
	leave.ApproverID = &approverID
	leave.ApprovedAt = &now

	if err := h.repo.Update(c.Request.Context(), leave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve leave"})
		return
	}

	c.JSON(http.StatusOK, leave)
}

// Reject rejects a leave request (HR/Admin only)
func (h *Handler) Reject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave ID"})
		return
	}

	var input RejectLeaveRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	approverID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	leave, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave not found"})
		return
	}

	if leave.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave is not pending"})
		return
	}

	leave.Status = "rejected"
	leave.ApproverID = &approverID
	leave.RejectionReason = input.RejectionReason

	if err := h.repo.Update(c.Request.Context(), leave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject leave"})
		return
	}

	c.JSON(http.StatusOK, leave)
}

// GetPendingLeaves retrieves all pending leave requests (HR/Admin only)
func (h *Handler) GetPendingLeaves(c *gin.Context) {
	leaves, err := h.repo.GetPendingLeaves(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending leaves"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaves": leaves})
}

// GetLeaveBalance retrieves leave balance for an employee
func (h *Handler) GetLeaveBalance(c *gin.Context) {
	employeeIDParam := c.Param("employee_id")
	employeeID, err := uuid.Parse(employeeIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Verify employee can only view their own balance (unless HR/Admin)
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" {
		userID, _ := middleware.GetUserID(c)
		if employeeID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot view another employee's leave balance"})
			return
		}
	}

	year := time.Now().Year()
	if yearParam := c.Query("year"); yearParam != "" {
		if _, err := fmt.Sscanf(yearParam, "%d", &year); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
			return
		}
	}

	balance, err := h.repo.GetLeaveBalance(c.Request.Context(), employeeID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get leave balance"})
		return
	}

	c.JSON(http.StatusOK, balance)
}
