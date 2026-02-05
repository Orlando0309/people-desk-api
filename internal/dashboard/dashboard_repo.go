package dashboard

import (
	"context"
	"fmt"
	"go-server/internal/auth"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for dashboard
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new dashboard repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// GetDashboardStats retrieves aggregated dashboard statistics
func (r *Repo) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stats := &DashboardStats{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var err error

	// Helper function to safely execute queries concurrently
	execQuery := func(name string, query func() (*int64, error)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count, queryErr := query()
			if queryErr != nil {
				mu.Lock()
				if err == nil {
					err = fmt.Errorf("%s: %w", name, queryErr)
				}
				mu.Unlock()
				return
			}
			if count != nil {
				mu.Lock()
				switch name {
				case "total_employees":
					stats.TotalEmployees = int(*count)
				case "on_leave":
					stats.OnLeave = int(*count)
				case "pending_leaves":
					stats.PendingLeaves = int(*count)
				case "new_requests":
					stats.NewRequests = int(*count)
				}
				mu.Unlock()
			}
		}()
	}

	// Execute all queries concurrently
	execQuery("total_employees", func() (*int64, error) {
		var count int64
		result := r.db.WithContext(ctx).Model(&struct{}{}).Table("employees").Count(&count)
		return &count, result.Error
	})

	execQuery("on_leave", func() (*int64, error) {
		var count int64
		result := r.db.WithContext(ctx).Model(&struct{}{}).Table("employees").
			Where("status = ?", "on_leave").Count(&count)
		return &count, result.Error
	})

	execQuery("pending_leaves", func() (*int64, error) {
		var count int64
		result := r.db.WithContext(ctx).Model(&struct{}{}).Table("leaves").
			Where("status = ?", "pending").Count(&count)
		return &count, result.Error
	})

	execQuery("new_requests", func() (*int64, error) {
		var count int64
		result := r.db.WithContext(ctx).Model(&struct{}{}).Table("leaves").
			Where("status = ?", "pending").Count(&count)
		return &count, result.Error
	})

	// Wait for all queries to complete
	wg.Wait()

	if err != nil {
		return nil, err
	}

	// Get attendance rate for today concurrently
	var totalAttendance, presentAttendance, lateAttendance int64
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := r.db.WithContext(ctx).Model(&struct{}{}).Table("attendance").
			Joins("JOIN employees ON attendance.employee_id = employees.id").
			Where("attendance.date = ?", time.Now().Format("2006-01-02")).
			Count(&totalAttendance)
		if result != nil && result.Error != nil {
			mu.Lock()
			if err == nil {
				err = fmt.Errorf("count total attendance: %w", result.Error)
			}
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := r.db.WithContext(ctx).Model(&struct{}{}).Table("attendance").
			Joins("JOIN employees ON attendance.employee_id = employees.id").
			Where("attendance.date = ? AND attendance.status = ?", time.Now().Format("2006-01-02"), "present").
			Count(&presentAttendance)
		if result != nil && result.Error != nil {
			mu.Lock()
			if err == nil {
				err = fmt.Errorf("count present attendance: %w", result.Error)
			}
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := r.db.WithContext(ctx).Model(&struct{}{}).Table("attendance").
			Joins("JOIN employees ON attendance.employee_id = employees.id").
			Where("attendance.date = ? AND attendance.status = ?", time.Now().Format("2006-01-02"), "late").
			Count(&lateAttendance)
		if result != nil && result.Error != nil {
			mu.Lock()
			if err == nil {
				err = fmt.Errorf("count late attendance: %w", result.Error)
			}
			mu.Unlock()
		}
	}()

	wg.Wait()

	if err != nil {
		return nil, err
	}

	// Calculate attendance rate
	if totalAttendance > 0 {
		stats.AttendanceRate = float64(presentAttendance+lateAttendance) / float64(totalAttendance) * 100
	} else {
		stats.AttendanceRate = 0
	}

	// Get monthly payroll concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := r.db.WithContext(ctx).Model(&struct{}{}).Table("employees").
			Where("status = ?", "active").
			Select("COALESCE(SUM(gross_salary), 0)").
			Scan(&stats.MonthlyPayroll)
		if result != nil && result.Error != nil {
			mu.Lock()
			if err == nil {
				err = fmt.Errorf("calculate monthly payroll: %w", result.Error)
			}
			mu.Unlock()
		}
	}()

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetAttendanceSummary retrieves daily or weekly attendance summary
func (r *Repo) GetAttendanceSummary(ctx context.Context, date string, period string) (*AttendanceSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	summary := &AttendanceSummary{}

	// Parse date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	summary.Date = parsedDate

	// Get total employees concurrently
	var totalEmployees int64
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		r.db.WithContext(ctx).Model(&struct{}{}).Table("employees").Count(&totalEmployees)
	}()

	// Get attendance counts concurrently
	var present, late, absent, overtime int64
	wg.Add(4)
	go func() {
		defer wg.Done()
		r.db.WithContext(ctx).Model(&struct{}{}).Table("attendance").
			Joins("JOIN employees ON attendance.employee_id = employees.id").
			Where("attendance.date = ?", date).
			Where("attendance.status = ?", "present").
			Count(&present)
	}()
	go func() {
		defer wg.Done()
		r.db.WithContext(ctx).Model(&struct{}{}).Table("attendance").
			Joins("JOIN employees ON attendance.employee_id = employees.id").
			Where("attendance.date = ?", date).
			Where("attendance.status = ?", "late").
			Count(&late)
	}()
	go func() {
		defer wg.Done()
		r.db.WithContext(ctx).Model(&struct{}{}).Table("attendance").
			Joins("JOIN employees ON attendance.employee_id = employees.id").
			Where("attendance.date = ?", date).
			Where("attendance.status = ?", "absent").
			Count(&absent)
	}()
	go func() {
		defer wg.Done()
		r.db.WithContext(ctx).Model(&struct{}{}).Table("attendance").
			Joins("JOIN employees ON attendance.employee_id = employees.id").
			Where("attendance.date = ?", date).
			Where("attendance.status = ?", "overtime").
			Count(&overtime)
	}()

	// Wait for all queries to complete
	wg.Wait()

	summary.TotalEmployees = int(totalEmployees)
	summary.Present = int(present)
	summary.Late = int(late)
	summary.Absent = int(absent)
	summary.Overtime = int(overtime)

	// Calculate attendance rate
	if totalEmployees > 0 {
		summary.AttendanceRate = float64(present+late) / float64(totalEmployees) * 100
	} else {
		summary.AttendanceRate = 0
	}

	return summary, nil
}

// GetLeaveBalances retrieves leave balances for multiple employees
func (r *Repo) GetLeaveBalances(ctx context.Context, year int, department string, employeeIDs []uuid.UUID) ([]LeaveBalance, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var balances []LeaveBalance

	// Build query
	query := r.db.WithContext(ctx).Table("leaves").
		Select(`
			employee_id,
			COUNT(CASE WHEN leave_type = 'annual' AND status = 'approved' THEN 1 END) as annual_used,
			COUNT(CASE WHEN leave_type = 'sick' AND status = 'approved' THEN 1 END) as sick_used,
			COUNT(CASE WHEN leave_type = 'maternity' AND status = 'approved' THEN 1 END) as maternity_used,
			COUNT(CASE WHEN leave_type = 'exceptional' AND status = 'approved' THEN 1 END) as exceptional_used,
			COUNT(CASE WHEN leave_type = 'paternity' AND status = 'approved' THEN 1 END) as paternity_used,
			COUNT(CASE WHEN leave_type = 'unpaid' AND status = 'approved' THEN 1 END) as unpaid_used
		`).
		Where("status = ?", "approved").
		Where("DATE_TRUNC('year', start_date) = ?", year)

	// Apply filters
	if department != "" {
		query = query.Joins("JOIN employees ON leaves.employee_id = employees.id").
			Where("employees.department = ?", department)
	}

	if len(employeeIDs) > 0 {
		query = query.Where("employee_id IN ?", employeeIDs)
	}

	// Group by employee_id
	query = query.Group("employee_id")

	// Get results
	if err := query.Scan(&balances).Error; err != nil {
		return nil, fmt.Errorf("get leave balances: %w", err)
	}

	// Get annual entitlement for each employee
	for i := range balances {
		var annualEntitlement float64
		if err := r.db.WithContext(ctx).Table("payroll_config").
			Where("config_key = ?", "annual_leave_entitlement").
			Scan(&annualEntitlement).Error; err != nil {
			return nil, fmt.Errorf("get annual entitlement: %w", err)
		}
		balances[i].AnnualEntitlement = annualEntitlement
		balances[i].Year = year
		balances[i].AnnualRemaining = annualEntitlement - balances[i].AnnualUsed
	}

	return balances, nil
}

// GetUserByID retrieves a user by ID
func (r *Repo) GetUserByID(ctx context.Context, id uuid.UUID) (*auth.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user auth.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

// UpdateUser updates user information
func (r *Repo) UpdateUser(ctx context.Context, user *auth.User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// DeleteUser soft deletes a user
func (r *Repo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Delete(&auth.User{}, id).Error; err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
