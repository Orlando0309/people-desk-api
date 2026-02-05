package company

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles company settings requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new company handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// Get retrieves company settings
func (h *Handler) Get(c *gin.Context) {
	settings, err := h.repo.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get company settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// Update updates company settings (Admin only)
func (h *Handler) Update(c *gin.Context) {
	var input UpdateCompanySettingsRequest
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

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Only Admin can update company settings
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	// Get existing settings
	settings, err := h.repo.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get company settings"})
		return
	}

	// Update fields if provided
	if input.CompanyName != nil {
		settings.CompanyName = *input.CompanyName
	}
	if input.CompanyAddress != nil {
		settings.CompanyAddress = input.CompanyAddress
	}
	if input.CompanyNIF != nil {
		settings.CompanyNIF = input.CompanyNIF
	}
	if input.CompanySTAT != nil {
		settings.CompanySTAT = input.CompanySTAT
	}
	if input.CNAPSNumber != nil {
		settings.CNAPSNumber = input.CNAPSNumber
	}
	if input.OSTIENumber != nil {
		settings.OSTIENumber = input.OSTIENumber
	}
	if input.ContactEmail != nil {
		settings.ContactEmail = input.ContactEmail
	}
	if input.ContactPhone != nil {
		settings.ContactPhone = input.ContactPhone
	}
	if input.LogoURL != nil {
		settings.LogoURL = input.LogoURL
	}
	if input.Timezone != nil {
		settings.Timezone = *input.Timezone
	}
	if input.Currency != nil {
		settings.Currency = *input.Currency
	}
	if input.FiscalYearStart != nil {
		settings.FiscalYearStart = *input.FiscalYearStart
	}
	if input.WorkHoursPerDay != nil {
		settings.WorkHoursPerDay = *input.WorkHoursPerDay
	}
	if input.WorkDaysPerWeek != nil {
		settings.WorkDaysPerWeek = *input.WorkDaysPerWeek
	}
	if input.OvertimeWeekdayRate != nil {
		settings.OvertimeWeekdayRate = *input.OvertimeWeekdayRate
	}
	if input.OvertimeSaturdayRate != nil {
		settings.OvertimeSaturdayRate = *input.OvertimeSaturdayRate
	}
	if input.OvertimeSundayRate != nil {
		settings.OvertimeSundayRate = *input.OvertimeSundayRate
	}
	if input.AnnualLeaveDays != nil {
		settings.AnnualLeaveDays = *input.AnnualLeaveDays
	}
	if input.MinimumSalary != nil {
		settings.MinimumSalary = *input.MinimumSalary
	}

	// Set updated by
	settings.UpdatedBy = &userID

	if err := h.repo.Update(c.Request.Context(), settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UploadLogo uploads company logo (Admin only)
func (h *Handler) UploadLogo(c *gin.Context) {
	// Get user info from context
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

	// Only Admin can upload logo
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, PNG, and GIF are allowed"})
		return
	}

	// For now, we'll just return a mock URL
	// In a real implementation, you would upload to a cloud storage service
	logoURL := "https://example.com/logos/company-logo-" + uuid.New().String() + ext

	// Update logo URL in database
	if err := h.repo.UpdateLogo(c.Request.Context(), logoURL, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update logo"})
		return
	}

	response := UploadLogoResponse{
		LogoURL: logoURL,
	}

	c.JSON(http.StatusOK, response)
}

// ListHolidays retrieves company holidays
func (h *Handler) ListHolidays(c *gin.Context) {
	var query HolidaysListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default year to current year if not provided
	if query.Year == 0 && query.StartDate == "" && query.EndDate == "" {
		query.Year = time.Now().Year()
	}

	holidays, err := h.repo.ListHolidays(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list holidays"})
		return
	}

	response := HolidaysListResponse{
		Holidays: holidays,
	}

	c.JSON(http.StatusOK, response)
}

// CreateHoliday creates a new company holiday (Admin only)
func (h *Handler) CreateHoliday(c *gin.Context) {
	var input CreateHolidayRequest
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

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Only Admin can create holidays
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var description *string
	if input.Description != "" {
		description = &input.Description
	}

	holiday := &CompanyHoliday{
		Name:        input.Name,
		Date:        input.Date,
		IsRecurring: input.IsRecurring,
		Description: description,
		CreatedBy:   &userID,
	}

	if err := h.repo.CreateHoliday(c.Request.Context(), holiday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create holiday"})
		return
	}

	c.JSON(http.StatusCreated, holiday)
}

// UpdateHoliday updates a company holiday (Admin only)
func (h *Handler) UpdateHoliday(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid holiday ID"})
		return
	}

	var input UpdateHolidayRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Only Admin can update holidays
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	// Get existing holiday
	holiday, err := h.repo.GetHolidayByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Holiday not found"})
		return
	}

	// Update fields if provided
	if input.Name != nil {
		holiday.Name = *input.Name
	}
	if input.Date != nil {
		holiday.Date = *input.Date
	}
	if input.IsRecurring != nil {
		holiday.IsRecurring = *input.IsRecurring
	}
	if input.Description != nil {
		holiday.Description = input.Description
	}

	if err := h.repo.UpdateHoliday(c.Request.Context(), holiday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update holiday"})
		return
	}

	c.JSON(http.StatusOK, holiday)
}

// DeleteHoliday deletes a company holiday (Admin only)
func (h *Handler) DeleteHoliday(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid holiday ID"})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Only Admin can delete holidays
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	if err := h.repo.DeleteHoliday(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Holiday not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Holiday deleted successfully"})
}
