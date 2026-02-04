package attendance

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Attendance represents an attendance record
type Attendance struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EmployeeID      uuid.UUID      `gorm:"type:uuid;not null" json:"employee_id"`
	Date            time.Time      `gorm:"type:date;not null" json:"date"`
	ClockIn         *time.Time     `json:"clock_in,omitempty"`
	ClockOut        *time.Time     `json:"clock_out,omitempty"`
	IPAddress       string         `gorm:"type:inet" json:"ip_address,omitempty"`
	DeviceFingerprint string       `gorm:"type:text" json:"device_fingerprint,omitempty"`
	Status          string         `gorm:"type:varchar(20);default:'present';check:status IN ('present', 'absent', 'late', 'overtime', 'half_day')" json:"status"`
	TotalHours      *float64       `gorm:"type:numeric(5,2)" json:"total_hours,omitempty"`
	OvertimeHours   float64        `gorm:"type:numeric(5,2);default:0" json:"overtime_hours"`
	Notes           string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt       time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// ClockInRequest represents clock-in request
type ClockInRequest struct {
	EmployeeID uuid.UUID `json:"employee_id" binding:"required"`
}

// ClockOutRequest represents clock-out request
type ClockOutRequest struct {
	EmployeeID uuid.UUID `json:"employee_id" binding:"required"`
}

// AttendanceListQuery represents query parameters for listing attendance
type AttendanceListQuery struct {
	EmployeeID *uuid.UUID `form:"employee_id"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
	Status     string     `form:"status" binding:"omitempty,oneof=present absent late overtime half_day"`
	Limit      int        `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset     int        `form:"offset" binding:"omitempty,min=0"`
}

// AttendanceCorrectionRequest represents attendance correction request
type AttendanceCorrectionRequest struct {
	ClockIn  *time.Time `json:"clock_in,omitempty"`
	ClockOut *time.Time `json:"clock_out,omitempty"`
	Status   *string    `json:"status,omitempty" binding:"omitempty,oneof=present absent late overtime half_day"`
	Notes    *string    `json:"notes,omitempty"`
}

// AttendanceStats represents attendance statistics
type AttendanceStats struct {
	TotalDays      int     `json:"total_days"`
	PresentDays    int     `json:"present_days"`
	AbsentDays     int     `json:"absent_days"`
	LateDays       int     `json:"late_days"`
	TotalHours     float64 `json:"total_hours"`
	OvertimeHours  float64 `json:"overtime_hours"`
	AttendanceRate float64 `json:"attendance_rate"`
}

// TableName specifies the table name for Attendance model
func (Attendance) TableName() string {
	return "attendance"
}
