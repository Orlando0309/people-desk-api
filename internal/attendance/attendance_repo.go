package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for attendance
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new attendance repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// ClockIn records employee clock-in
func (r *Repo) ClockIn(ctx context.Context, attendance *Attendance) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Check if already clocked in today
	var existing Attendance
	err := r.db.WithContext(ctx).Where("employee_id = ? AND date = ?", attendance.EmployeeID, attendance.Date).First(&existing).Error
	if err == nil {
		return fmt.Errorf("already clocked in today")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("check existing attendance: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(attendance).Error; err != nil {
		return fmt.Errorf("clock in: %w", err)
	}
	return nil
}

// ClockOut records employee clock-out
func (r *Repo) ClockOut(ctx context.Context, employeeID uuid.UUID, clockOut time.Time) (*Attendance, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Use the same date format as ClockIn handler for consistency
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var attendance Attendance
	if err := r.db.WithContext(ctx).Where("employee_id = ? AND date = ?", employeeID, today).First(&attendance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no clock-in record found for today")
		}
		return nil, fmt.Errorf("get attendance: %w", err)
	}

	if attendance.ClockOut != nil {
		return nil, fmt.Errorf("already clocked out")
	}

	attendance.ClockOut = &clockOut

	// Calculate total hours
	if attendance.ClockIn != nil {
		duration := clockOut.Sub(*attendance.ClockIn)
		hours := duration.Hours()
		attendance.TotalHours = &hours

		// Calculate overtime (> 8 hours)
		if hours > 8 {
			attendance.OvertimeHours = hours - 8
			attendance.Status = "overtime"
		}
	}

	attendance.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&attendance).Error; err != nil {
		return nil, fmt.Errorf("clock out: %w", err)
	}

	return &attendance, nil
}

// GetByID retrieves an attendance record by ID
func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*Attendance, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var attendance Attendance
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&attendance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("attendance not found")
		}
		return nil, fmt.Errorf("get attendance: %w", err)
	}
	return &attendance, nil
}

// Update updates an attendance record
func (r *Repo) Update(ctx context.Context, attendance *Attendance) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	attendance.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(attendance).Error; err != nil {
		return fmt.Errorf("update attendance: %w", err)
	}
	return nil
}

// List retrieves attendance records with filtering
func (r *Repo) List(ctx context.Context, query AttendanceListQuery) ([]Attendance, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var attendances []Attendance
	var total int64

	db := r.db.WithContext(ctx).Model(&Attendance{})

	// Apply filters
	if query.EmployeeID != nil {
		db = db.Where("employee_id = ?", *query.EmployeeID)
	}

	if query.StartDate != nil {
		db = db.Where("date >= ?", *query.StartDate)
	}

	if query.EndDate != nil {
		db = db.Where("date <= ?", *query.EndDate)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count attendance: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("date DESC, clock_in DESC").Find(&attendances).Error; err != nil {
		return nil, 0, fmt.Errorf("list attendance: %w", err)
	}

	return attendances, total, nil
}

// GetTodayAttendance retrieves today's attendance for an employee
func (r *Repo) GetTodayAttendance(ctx context.Context, employeeID uuid.UUID) (*Attendance, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	today := time.Now().Format("2006-01-02")
	var attendance Attendance
	if err := r.db.WithContext(ctx).Where("employee_id = ? AND date = ?", employeeID, today).First(&attendance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get today attendance: %w", err)
	}
	return &attendance, nil
}

// GetStats retrieves attendance statistics for an employee
func (r *Repo) GetStats(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (*AttendanceStats, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var stats AttendanceStats

	// Count total days
	var totalDays int64
	if err := r.db.WithContext(ctx).Model(&Attendance{}).
		Where("employee_id = ? AND date >= ? AND date <= ?", employeeID, startDate, endDate).
		Count(&totalDays).Error; err != nil {
		return nil, fmt.Errorf("count total days: %w", err)
	}
	stats.TotalDays = int(totalDays)

	// Count present days
	var presentDays int64
	if err := r.db.WithContext(ctx).Model(&Attendance{}).
		Where("employee_id = ? AND date >= ? AND date <= ? AND status = ?", employeeID, startDate, endDate, "present").
		Count(&presentDays).Error; err != nil {
		return nil, fmt.Errorf("count present days: %w", err)
	}
	stats.PresentDays = int(presentDays)

	// Count absent days
	var absentDays int64
	if err := r.db.WithContext(ctx).Model(&Attendance{}).
		Where("employee_id = ? AND date >= ? AND date <= ? AND status = ?", employeeID, startDate, endDate, "absent").
		Count(&absentDays).Error; err != nil {
		return nil, fmt.Errorf("count absent days: %w", err)
	}
	stats.AbsentDays = int(absentDays)

	// Count late days
	var lateDays int64
	if err := r.db.WithContext(ctx).Model(&Attendance{}).
		Where("employee_id = ? AND date >= ? AND date <= ? AND status = ?", employeeID, startDate, endDate, "late").
		Count(&lateDays).Error; err != nil {
		return nil, fmt.Errorf("count late days: %w", err)
	}
	stats.LateDays = int(lateDays)

	// Sum total hours
	var totalHours *float64
	if err := r.db.WithContext(ctx).Model(&Attendance{}).
		Where("employee_id = ? AND date >= ? AND date <= ?", employeeID, startDate, endDate).
		Select("SUM(total_hours)").Scan(&totalHours).Error; err != nil {
		return nil, fmt.Errorf("sum total hours: %w", err)
	}
	if totalHours != nil {
		stats.TotalHours = *totalHours
	}

	// Sum overtime hours
	var overtimeHours *float64
	if err := r.db.WithContext(ctx).Model(&Attendance{}).
		Where("employee_id = ? AND date >= ? AND date <= ?", employeeID, startDate, endDate).
		Select("SUM(overtime_hours)").Scan(&overtimeHours).Error; err != nil {
		return nil, fmt.Errorf("sum overtime hours: %w", err)
	}
	if overtimeHours != nil {
		stats.OvertimeHours = *overtimeHours
	}

	// Calculate attendance rate
	if stats.TotalDays > 0 {
		stats.AttendanceRate = float64(stats.PresentDays) / float64(stats.TotalDays) * 100
	}

	return &stats, nil
}
