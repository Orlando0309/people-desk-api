package audit

import (
	"net/http"
	"strconv"
	"time"
	"fmt"
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles audit log requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new audit handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// LogAction is a helper function to log an action (can be called from other modules)
func (h *Handler) LogAction(ctx *gin.Context, actionType, module string, recordID *uuid.UUID, beforeValue, afterValue interface{}) error {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return err
	}

	userRole, err := middleware.GetUserRole(ctx)
	if err != nil {
		return err
	}

	log := &AuditLog{
		UserID:      userID,
		UserRole:    userRole,
		IPAddress:   ctx.ClientIP(),
		ActionType:  actionType,
		Module:      module,
		RecordID:    recordID,
		BeforeValue: beforeValue,
		AfterValue:  afterValue,
	}

	return h.repo.Create(ctx.Request.Context(), log)
}

// GetByID retrieves an audit log by ID (Admin only)
func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audit log ID"})
		return
	}

	log, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audit log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// List retrieves audit logs with filtering (Admin only)
func (h *Handler) List(c *gin.Context) {
	var query AuditLogListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logs, total, err := h.repo.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"total":  total,
		"limit":  query.Limit,
		"offset": query.Offset,
	})
}

// GetByUserID retrieves all audit logs for a specific user (Admin only)
func (h *Handler) GetByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit := 100
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil {
			limit = l
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil {
			offset = o
		}
	}

	logs, total, err := h.repo.GetByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetByModule retrieves all audit logs for a specific module (Admin only)
func (h *Handler) GetByModule(c *gin.Context) {
	module := c.Param("module")

	limit := 100
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil {
			limit = l
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil {
			offset = o
		}
	}

	logs, total, err := h.repo.GetByModule(c.Request.Context(), module, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get module audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetStats retrieves audit statistics (Admin only)
func (h *Handler) GetStats(c *gin.Context) {
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
		// Default to today
		endDate = time.Now()
	}

	stats, err := h.repo.GetStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get audit stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetRecordHistory retrieves all audit logs for a specific record (Admin only)
func (h *Handler) GetRecordHistory(c *gin.Context) {
	recordIDParam := c.Param("record_id")
	recordID, err := uuid.Parse(recordIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	logs, err := h.repo.GetRecordHistory(c.Request.Context(), recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get record history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// ExportLogs exports audit logs to CSV (Admin only)
func (h *Handler) ExportLogs(c *gin.Context) {
	var query AuditLogListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set a high limit for export
	query.Limit = 10000

	logs, _, err := h.repo.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export audit logs"})
		return
	}

	// Set headers for CSV download
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=audit_logs.csv")

	// Write CSV header
	c.Writer.Write([]byte("ID,User ID,User Role,IP Address,Action Type,Module,Record ID,Created At\n"))

	// Write CSV rows
	for _, log := range logs {
		recordID := ""
		if log.RecordID != nil {
			recordID = log.RecordID.String()
		}
		c.Writer.Write([]byte(fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s,%s\n",
			log.ID,
			log.UserID.String(),
			log.UserRole,
			log.IPAddress,
			log.ActionType,
			log.Module,
			recordID,
			log.CreatedAt.Format(time.RFC3339),
		)))
	}
}
