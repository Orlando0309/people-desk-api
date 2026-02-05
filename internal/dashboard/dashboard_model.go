package dashboard

import (
	"time"

	"github.com/google/uuid"
)

// DashboardStats represents aggregated statistics for the dashboard
type DashboardStats struct {
	TotalEmployees int     `json:"total_employees"`
	OnLeave        int     `json:"on_leave"`
	NewRequests    int     `json:"new_requests"`
	AttendanceRate float64 `json:"attendance_rate"`
	PendingLeaves  int     `json:"pending_leaves"`
	MonthlyPayroll float64 `json:"monthly_payroll"`
}

// AttendanceSummary represents daily or weekly attendance summary
type AttendanceSummary struct {
	Date           time.Time `json:"date"`
	TotalEmployees int       `json:"total_employees"`
	Present        int       `json:"present"`
	Late           int       `json:"late"`
	Absent         int       `json:"absent"`
	Overtime       int       `json:"overtime"`
	AttendanceRate float64   `json:"attendance_rate"`
}

// LeaveBalance represents leave balance for an employee
type LeaveBalance struct {
	EmployeeID        uuid.UUID `json:"employee_id"`
	Year              int       `json:"year"`
	AnnualEntitlement float64   `json:"annual_entitlement"`
	AnnualUsed        float64   `json:"annual_used"`
	AnnualRemaining   float64   `json:"annual_remaining"`
	SickUsed          float64   `json:"sick_used"`
	MaternityUsed     float64   `json:"maternity_used"`
	ExceptionalUsed   float64   `json:"exceptional_used"`
	PaternityUsed     float64   `json:"paternity_used"`
	UnpaidUsed        float64   `json:"unpaid_used"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	Role     *string `json:"role" binding:"omitempty,oneof=admin hr accountant employee"`
	IsActive *bool   `json:"is_active"`
}

// CalendarEvent represents a calendar event
type CalendarEvent struct {
	ID            uuid.UUID  `json:"id"`
	Type          string     `json:"type"`
	Title         string     `json:"title"`
	StartDate     string     `json:"start_date"`
	EndDate       string     `json:"end_date"`
	EmployeeID    *uuid.UUID `json:"employee_id,omitempty"`
	EmployeeName  string     `json:"employee_name,omitempty"`
	LeaveType     string     `json:"leave_type,omitempty"`
	Status        string     `json:"status,omitempty"`
	Color         string     `json:"color,omitempty"`
	IsCompanyWide bool       `json:"is_company_wide,omitempty"`
}

// CalendarEventsQuery represents query parameters for calendar events
type CalendarEventsQuery struct {
	StartDate  string     `form:"start_date" binding:"required"`
	EndDate    string     `form:"end_date" binding:"required"`
	EventTypes []string   `form:"event_types"`
	EmployeeID *uuid.UUID `form:"employee_id"`
}

// CalendarEventsResponse represents the response for calendar events
type CalendarEventsResponse struct {
	Events []CalendarEvent `json:"events"`
}

// TableName specifies the table name for models
func (DashboardStats) TableName() string {
	return "dashboard_stats"
}

func (AttendanceSummary) TableName() string {
	return "attendance_summary"
}

func (LeaveBalance) TableName() string {
	return "leave_balances"
}
