package audit

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an immutable audit log entry
type AuditLog struct {
	ID          int64     `gorm:"primary_key;autoIncrement" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	UserRole    string    `gorm:"type:varchar(20);not null" json:"user_role"`
	IPAddress   string    `gorm:"type:inet;not null" json:"ip_address"`
	ActionType  string    `gorm:"type:varchar(50);not null" json:"action_type"`
	Module      string    `gorm:"type:varchar(50);not null" json:"module"`
	RecordID    *uuid.UUID `gorm:"type:uuid" json:"record_id,omitempty"`
	BeforeValue interface{} `gorm:"type:jsonb" json:"before_value,omitempty"`
	AfterValue  interface{} `gorm:"type:jsonb" json:"after_value,omitempty"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
}

// CreateAuditLogRequest represents audit log creation request
type CreateAuditLogRequest struct {
	UserID      uuid.UUID   `json:"user_id" binding:"required"`
	UserRole    string      `json:"user_role" binding:"required"`
	IPAddress   string      `json:"ip_address" binding:"required"`
	ActionType  string      `json:"action_type" binding:"required"`
	Module      string      `json:"module" binding:"required"`
	RecordID    *uuid.UUID  `json:"record_id,omitempty"`
	BeforeValue interface{} `json:"before_value,omitempty"`
	AfterValue  interface{} `json:"after_value,omitempty"`
}

// AuditLogListQuery represents query parameters for listing audit logs
type AuditLogListQuery struct {
	UserID     *uuid.UUID `form:"user_id"`
	UserRole   string     `form:"user_role" binding:"omitempty,oneof=admin hr accountant employee"`
	ActionType string     `form:"action_type"`
	Module     string     `form:"module"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
	Limit      int        `form:"limit" binding:"omitempty,min=1,max=1000"`
	Offset     int        `form:"offset" binding:"omitempty,min=0"`
}

// AuditStats represents audit statistics
type AuditStats struct {
	TotalActions   int64            `json:"total_actions"`
	ActionsByType  map[string]int64 `json:"actions_by_type"`
	ActionsByModule map[string]int64 `json:"actions_by_module"`
	ActionsByUser  map[string]int64 `json:"actions_by_user"`
}

// TableName specifies the table name for AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}
