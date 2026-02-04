package leave

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Leave represents a leave request
type Leave struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EmployeeID      uuid.UUID      `gorm:"type:uuid;not null" json:"employee_id"`
	LeaveType       string         `gorm:"type:varchar(50);not null;check:leave_type IN ('annual', 'sick', 'maternity', 'exceptional', 'paternity', 'unpaid')" json:"leave_type"`
	StartDate       time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate         time.Time      `gorm:"type:date;not null" json:"end_date"`
	DaysRequested   float64        `gorm:"type:numeric(5,2);not null" json:"days_requested"`
	Status          string         `gorm:"type:varchar(20);default:'pending';not null;check:status IN ('pending', 'approved', 'rejected', 'cancelled')" json:"status"`
	ApproverID      *uuid.UUID     `gorm:"type:uuid" json:"approver_id,omitempty"`
	Reason          string         `gorm:"type:text" json:"reason,omitempty"`
	RejectionReason string         `gorm:"type:text" json:"rejection_reason,omitempty"`
	ApprovedAt      *time.Time     `json:"approved_at,omitempty"`
	CreatedAt       time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateLeaveRequest represents leave creation request
type CreateLeaveRequest struct {
	EmployeeID    uuid.UUID `json:"employee_id" binding:"required"`
	LeaveType     string    `json:"leave_type" binding:"required,oneof=annual sick maternity exceptional paternity unpaid"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required"`
	DaysRequested float64   `json:"days_requested" binding:"required,min=0.5"`
	Reason        string    `json:"reason,omitempty"`
}

// UpdateLeaveRequest represents leave update request
type UpdateLeaveRequest struct {
	LeaveType     *string    `json:"leave_type,omitempty" binding:"omitempty,oneof=annual sick maternity exceptional paternity unpaid"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	DaysRequested *float64   `json:"days_requested,omitempty" binding:"omitempty,min=0.5"`
	Reason        *string    `json:"reason,omitempty"`
}

// ApproveLeaveRequest represents leave approval request
type ApproveLeaveRequest struct {
	ApproverID uuid.UUID `json:"approver_id" binding:"required"`
}

// RejectLeaveRequest represents leave rejection request
type RejectLeaveRequest struct {
	ApproverID      uuid.UUID `json:"approver_id" binding:"required"`
	RejectionReason string    `json:"rejection_reason" binding:"required"`
}

// LeaveListQuery represents query parameters for listing leaves
type LeaveListQuery struct {
	EmployeeID *uuid.UUID `form:"employee_id"`
	LeaveType  string     `form:"leave_type" binding:"omitempty,oneof=annual sick maternity exceptional paternity unpaid"`
	Status     string     `form:"status" binding:"omitempty,oneof=pending approved rejected cancelled"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
	Limit      int        `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset     int        `form:"offset" binding:"omitempty,min=0"`
}

// LeaveBalance represents leave balance for an employee
type LeaveBalance struct {
	EmployeeID       uuid.UUID `json:"employee_id"`
	AnnualTotal      float64   `json:"annual_total"`
	AnnualUsed       float64   `json:"annual_used"`
	AnnualRemaining  float64   `json:"annual_remaining"`
	SickUsed         float64   `json:"sick_used"`
	MaternityUsed    float64   `json:"maternity_used"`
	ExceptionalUsed  float64   `json:"exceptional_used"`
	PaternityUsed    float64   `json:"paternity_used"`
	UnpaidUsed       float64   `json:"unpaid_used"`
}

// TableName specifies the table name for Leave model
func (Leave) TableName() string {
	return "leaves"
}
