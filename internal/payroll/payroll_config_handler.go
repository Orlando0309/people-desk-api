package payroll

import (
	"net/http"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ConfigHandler handles payroll configuration requests
type ConfigHandler struct {
	configRepo *ConfigRepo
}

// NewConfigHandler creates a new ConfigHandler
func NewConfigHandler(configRepo *ConfigRepo) *ConfigHandler {
	return &ConfigHandler{
		configRepo: configRepo,
	}
}

// CreatePayrollConfiguration creates a new payroll configuration (Admin only)
func (h *ConfigHandler) CreatePayrollConfiguration(c *gin.Context) {
	var input CreatePayrollConfigurationRequest
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

	config := &PayrollConfiguration{
		Key:         input.Key,
		Value:       input.Value,
		Description: input.Description,
		DataType:    input.DataType,
		Category:    input.Category,
		IsActive:    true,
		CreatedBy:   userID,
	}

	if err := h.configRepo.CreateConfig(c.Request.Context(), config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, config)
}

// UpdatePayrollConfiguration updates an existing payroll configuration (Admin only)
func (h *ConfigHandler) UpdatePayrollConfiguration(c *gin.Context) {
	id := c.Param("id")
	var input UpdatePayrollConfigurationRequest
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

	config, err := h.configRepo.GetConfigByID(c.Request.Context(), parseUUID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	// Update fields if provided
	if input.Value != nil {
		config.Value = *input.Value
	}
	if input.Description != nil {
		config.Description = *input.Description
	}
	if input.IsActive != nil {
		config.IsActive = *input.IsActive
	}
	config.UpdatedBy = &userID

	if err := h.configRepo.UpdateConfig(c.Request.Context(), config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeletePayrollConfiguration deletes a payroll configuration (Admin only)
func (h *ConfigHandler) DeletePayrollConfiguration(c *gin.Context) {
	id := c.Param("id")
	if err := h.configRepo.DeleteConfig(c.Request.Context(), parseUUID(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListPayrollConfigurations lists all payroll configurations (Admin only)
func (h *ConfigHandler) ListPayrollConfigurations(c *gin.Context) {
	var query PayrollConfigurationListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configs, total, err := h.configRepo.ListConfigs(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list configurations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"configurations": configs,
		"total":          total,
		"limit":          query.Limit,
		"offset":         query.Offset,
	})
}

// GetPayrollConfiguration retrieves a payroll configuration by ID (Admin only)
func (h *ConfigHandler) GetPayrollConfiguration(c *gin.Context) {
	id := c.Param("id")
	config, err := h.configRepo.GetConfigByID(c.Request.Context(), parseUUID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// CreateIRSABracket creates a new IRSA tax bracket (Admin only)
func (h *ConfigHandler) CreateIRSABracket(c *gin.Context) {
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

	bracket := &IRSATaxBracket{
		MinIncome:     input.MinIncome,
		MaxIncome:     input.MaxIncome,
		TaxRate:       input.TaxRate,
		MinTax:        input.MinTax,
		BracketName:   input.BracketName,
		IsActive:      true,
		SortOrder:     input.SortOrder,
		EffectiveDate: input.EffectiveDate,
		CreatedBy:     userID,
	}

	if err := h.configRepo.CreateIRSABracket(c.Request.Context(), bracket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bracket)
}

// GetIRSABracketByID retrieves an IRSA tax bracket by ID
func (h *ConfigHandler) GetIRSABracketByID(c *gin.Context) {
	id := c.Param("id")
	bracket, err := h.configRepo.GetIRSABracketByID(c.Request.Context(), parseUUID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tax bracket not found"})
		return
	}

	c.JSON(http.StatusOK, bracket)
}

// UpdateIRSABracket updates an IRSA tax bracket (Admin only)
func (h *ConfigHandler) UpdateIRSABracket(c *gin.Context) {
	id := c.Param("id")
	var input UpdateIRSATaxBracketRequest
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

	bracket, err := h.configRepo.GetIRSABracketByID(c.Request.Context(), parseUUID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tax bracket not found"})
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
	if input.BracketName != nil {
		bracket.BracketName = *input.BracketName
	}
	if input.SortOrder != nil {
		bracket.SortOrder = *input.SortOrder
	}
	if input.EffectiveDate != nil {
		bracket.EffectiveDate = *input.EffectiveDate
	}
	if input.IsActive != nil {
		bracket.IsActive = *input.IsActive
	}
	bracket.UpdatedBy = &userID

	if err := h.configRepo.UpdateIRSABracket(c.Request.Context(), bracket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bracket)
}

// DeleteIRSABracket deletes an IRSA tax bracket (Admin only)
func (h *ConfigHandler) DeleteIRSABracket(c *gin.Context) {
	id := c.Param("id")
	if err := h.configRepo.DeleteIRSABracket(c.Request.Context(), parseUUID(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListIRSABrackets lists all IRSA tax brackets
func (h *ConfigHandler) ListIRSABrackets(c *gin.Context) {
	var query IRSATaxBracketListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	brackets, total, err := h.configRepo.ListIRSABrackets(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tax brackets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"brackets": brackets,
		"total":    total,
		"limit":    query.Limit,
		"offset":   query.Offset,
	})
}

// parseUUID converts a string to a UUID
func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
