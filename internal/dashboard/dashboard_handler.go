package dashboard

import (
	"net/http"
	"strconv"
	"time"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles dashboard requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new dashboard handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// GetStats retrieves dashboard statistics
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.repo.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetAttendanceSummary retrieves daily or weekly attendance summary
func (h *Handler) GetAttendanceSummary(c *gin.Context) {
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	period := c.DefaultQuery("period", "day")

	summary, err := h.repo.GetAttendanceSummary(c.Request.Context(), date, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get attendance summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetLeaveBalances retrieves leave balances for multiple employees
func (h *Handler) GetLeaveBalances(c *gin.Context) {
	// Parse query parameters
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
		return
	}

	department := c.Query("department")
	employeeIDsStr := c.Query("employee_ids")

	var employeeIDs []string
	if employeeIDsStr != "" {
		ids := splitString(employeeIDsStr, ",")
		employeeIDs = ids
	}

	// Convert string IDs to UUIDs
	var uuids []uuid.UUID
	for _, id := range employeeIDs {
		uuid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID format"})
			return
		}
		uuids = append(uuids, uuid)
	}

	balances, err := h.repo.GetLeaveBalances(c.Request.Context(), year, department, uuids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get leave balances"})
		return
	}

	c.JSON(http.StatusOK, balances)
}

// UpdateUser updates a user
func (h *Handler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input UpdateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user
	user, err := h.repo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields if provided
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	if err := h.repo.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is trying to delete themselves
	currentUserID, _ := middleware.GetUserID(c)
	if currentUserID == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	// Soft delete the user
	if err := h.repo.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetBadgeCounts retrieves badge counts for navigation items
func (h *Handler) GetBadgeCounts(c *gin.Context) {
	// Get user info
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	badgeCounts, err := h.repo.GetBadgeCounts(c.Request.Context(), userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get badge counts"})
		return
	}

	c.JSON(http.StatusOK, badgeCounts)
}

// GetCalendarEvents retrieves calendar events for date range
func (h *Handler) GetCalendarEvents(c *gin.Context) {
	var query CalendarEventsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Employees can only see their own data unless employee_id is explicitly set
	if userRole == "employee" && query.EmployeeID == nil {
		query.EmployeeID = &userID
	}

	events, err := h.repo.GetCalendarEvents(c.Request.Context(), query, userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get calendar events"})
		return
	}

	response := CalendarEventsResponse{
		Events: events,
	}

	c.JSON(http.StatusOK, response)
}

// splitString splits a comma-separated string into a slice
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	result := make([]string, 0)
	for _, part := range rangeSplit(s, sep) {
		result = append(result, part)
	}
	return result
}

// rangeSplit is a helper function to split a string by separator
func rangeSplit(s, sep string) []string {
	result := make([]string, 0)
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

// GetMonthlyPayroll retrieves monthly payroll summary
func (h *Handler) GetMonthlyPayroll(c *gin.Context) {
	// Parse month parameter
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))

	// Validate month format
	if _, err := time.Parse("2006-01", month); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format. Use YYYY-MM"})
		return
	}

	// Get user role
	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Only Admin, HR, and Accountant can access payroll data
	if userRole != "admin" && userRole != "hr" && userRole != "accountant" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	payrollSummary, err := h.repo.GetMonthlyPayrollSummary(c.Request.Context(), month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get monthly payroll summary"})
		return
	}

	c.JSON(http.StatusOK, payrollSummary)
}

// GetCompliance retrieves compliance status for all modules
func (h *Handler) GetCompliance(c *gin.Context) {
	// Get user role
	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Only Admin, HR, and Accountant can access compliance data
	if userRole != "admin" && userRole != "hr" && userRole != "accountant" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	compliance, err := h.repo.GetComplianceStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get compliance status"})
		return
	}

	c.JSON(http.StatusOK, compliance)
}

// GetWeeklyAttendance retrieves weekly attendance summary
func (h *Handler) GetWeeklyAttendance(c *gin.Context) {
	// Parse query parameters
	employeeIDParam := c.Query("employee_id")
	startDateParam := c.DefaultQuery("start_date", time.Now().Format("2006-01-02"))

	// Validate date format
	startDate, err := time.Parse("2006-01-02", startDateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	// Get user info
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Parse employee ID if provided
	var employeeID *uuid.UUID
	if employeeIDParam != "" {
		id, err := uuid.Parse(employeeIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
			return
		}
		employeeID = &id
	}

	// Employees can only see their own data
	if userRole == "employee" && employeeID != nil && *employeeID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot view another employee's attendance"})
		return
	}

	// If employee and no employee_id specified, use their own ID
	if userRole == "employee" && employeeID == nil {
		employeeID = &userID
	}

	weeklyAttendance, err := h.repo.GetWeeklyAttendanceSummary(c.Request.Context(), employeeID, startDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weekly attendance summary"})
		return
	}

	c.JSON(http.StatusOK, weeklyAttendance)
}
