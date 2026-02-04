package attendance

import (
	"net/http"
	"time"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles attendance requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new attendance handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// ClockIn handles employee clock-in
func (h *Handler) ClockIn(c *gin.Context) {
	var input ClockInRequest
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

	// Verify employee can only clock in for themselves (unless HR/Admin)
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" && input.EmployeeID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot clock in for another employee"})
		return
	}

	now := time.Now()
	attendance := &Attendance{
		EmployeeID:        input.EmployeeID,
		Date:              time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		ClockIn:           &now,
		IPAddress:         c.ClientIP(),
		DeviceFingerprint: c.GetHeader("User-Agent"),
		Status:            "present",
	}

	// Check if late (after 9:00 AM)
	if now.Hour() > 9 || (now.Hour() == 9 && now.Minute() > 0) {
		attendance.Status = "late"
	}

	if err := h.repo.ClockIn(c.Request.Context(), attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attendance)
}

// ClockOut handles employee clock-out
func (h *Handler) ClockOut(c *gin.Context) {
	var input ClockOutRequest
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

	// Verify employee can only clock out for themselves (unless HR/Admin)
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" && input.EmployeeID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot clock out for another employee"})
		return
	}

	now := time.Now()
	attendance, err := h.repo.ClockOut(c.Request.Context(), input.EmployeeID, now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

// GetByID retrieves an attendance record by ID
func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendance ID"})
		return
	}

	attendance, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance not found"})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

// List retrieves attendance records with filtering
func (h *Handler) List(c *gin.Context) {
	var query AttendanceListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If employee role, only show their own attendance
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" {
		userID, _ := middleware.GetUserID(c)
		query.EmployeeID = &userID
	}

	attendances, total, err := h.repo.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attendances": attendances,
		"total":       total,
		"limit":       query.Limit,
		"offset":      query.Offset,
	})
}

// GetTodayAttendance retrieves today's attendance for an employee
func (h *Handler) GetTodayAttendance(c *gin.Context) {
	employeeIDParam := c.Param("employee_id")
	employeeID, err := uuid.Parse(employeeIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Verify employee can only view their own attendance (unless HR/Admin)
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" {
		userID, _ := middleware.GetUserID(c)
		if employeeID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot view another employee's attendance"})
			return
		}
	}

	attendance, err := h.repo.GetTodayAttendance(c.Request.Context(), employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get today's attendance"})
		return
	}

	if attendance == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No attendance record for today"})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

// UpdateAttendance handles attendance correction (HR/Admin only)
func (h *Handler) UpdateAttendance(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendance ID"})
		return
	}

	var input AttendanceCorrectionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attendance, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance not found"})
		return
	}

	// Update fields if provided
	if input.ClockIn != nil {
		attendance.ClockIn = input.ClockIn
	}
	if input.ClockOut != nil {
		attendance.ClockOut = input.ClockOut
	}
	if input.Status != nil {
		attendance.Status = *input.Status
	}
	if input.Notes != nil {
		attendance.Notes = *input.Notes
	}

	// Recalculate total hours if both clock in and out are present
	if attendance.ClockIn != nil && attendance.ClockOut != nil {
		duration := attendance.ClockOut.Sub(*attendance.ClockIn)
		hours := duration.Hours()
		attendance.TotalHours = &hours

		// Calculate overtime (> 8 hours)
		if hours > 8 {
			attendance.OvertimeHours = hours - 8
		} else {
			attendance.OvertimeHours = 0
		}
	}

	if err := h.repo.Update(c.Request.Context(), attendance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

// GetStats retrieves attendance statistics for an employee
func (h *Handler) GetStats(c *gin.Context) {
	employeeIDParam := c.Param("employee_id")
	employeeID, err := uuid.Parse(employeeIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Verify employee can only view their own stats (unless HR/Admin)
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" {
		userID, _ := middleware.GetUserID(c)
		if employeeID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot view another employee's stats"})
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
		// Default to today
		endDate = time.Now()
	}

	stats, err := h.repo.GetStats(c.Request.Context(), employeeID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get attendance stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
