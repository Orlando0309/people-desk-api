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

// GetMonthlyPayrollSummary retrieves monthly payroll summary
func (r *Repo) GetMonthlyPayrollSummary(ctx context.Context, month string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Parse month to get start and end dates
	parsedMonth, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("invalid month format: %w", err)
	}

	// Initialize result
	result := map[string]interface{}{
		"month":                  month,
		"total_gross_salary":     0.0,
		"total_net_salary":       0.0,
		"total_cnaps":            0.0,
		"total_ostie":            0.0,
		"total_irsa":             0.0,
		"total_employees_paid":   0,
		"pending_drafts":         0,
		"approved_payrolls":      0,
		"total_employer_charges": 0.0,
	}

	// Get payroll data
	var totalGross, totalNet, totalCNAPS, totalOSTIE, totalIRSA float64
	var totalEmployeesPaid, pendingDrafts, approvedPayrolls int64

	// Get approved payrolls for the month
	if err := r.db.WithContext(ctx).Table("payroll_approved").
		Where("month = ?", month).
		Select("COALESCE(SUM(gross_salary), 0)").
		Scan(&totalGross).Error; err != nil {
		return nil, fmt.Errorf("get total gross salary: %w", err)
	}

	if err := r.db.WithContext(ctx).Table("payroll_approved").
		Where("month = ?", month).
		Select("COALESCE(SUM(net_salary), 0)").
		Scan(&totalNet).Error; err != nil {
		return nil, fmt.Errorf("get total net salary: %w", err)
	}

	if err := r.db.WithContext(ctx).Table("payroll_approved").
		Where("month = ?", month).
		Select("COALESCE(SUM(cnaps), 0)").
		Scan(&totalCNAPS).Error; err != nil {
		return nil, fmt.Errorf("get total CNAPS: %w", err)
	}

	if err := r.db.WithContext(ctx).Table("payroll_approved").
		Where("month = ?", month).
		Select("COALESCE(SUM(ostie), 0)").
		Scan(&totalOSTIE).Error; err != nil {
		return nil, fmt.Errorf("get total OSTIE: %w", err)
	}

	if err := r.db.WithContext(ctx).Table("payroll_approved").
		Where("month = ?", month).
		Select("COALESCE(SUM(irsa), 0)").
		Scan(&totalIRSA).Error; err != nil {
		return nil, fmt.Errorf("get total IRSA: %w", err)
	}

	if err := r.db.WithContext(ctx).Table("payroll_approved").
		Where("month = ?", month).
		Count(&totalEmployeesPaid).Error; err != nil {
		return nil, fmt.Errorf("count approved payrolls: %w", err)
	}

	// Get pending drafts
	if err := r.db.WithContext(ctx).Table("payroll_drafts").
		Where("EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?",
			parsedMonth.Month(), parsedMonth.Year()).
		Count(&pendingDrafts).Error; err != nil {
		return nil, fmt.Errorf("count pending drafts: %w", err)
	}

	// Calculate employer charges (CNAPS + OSTIE)
	totalEmployerCharges := totalCNAPS + totalOSTIE

	// Update result
	result["total_gross_salary"] = totalGross
	result["total_net_salary"] = totalNet
	result["total_cnaps"] = totalCNAPS
	result["total_ostie"] = totalOSTIE
	result["total_irsa"] = totalIRSA
	result["total_employees_paid"] = totalEmployeesPaid
	result["pending_drafts"] = pendingDrafts
	result["approved_payrolls"] = approvedPayrolls
	result["total_employer_charges"] = totalEmployerCharges

	return result, nil
}

// GetComplianceStatus retrieves compliance status for all modules
func (r *Repo) GetComplianceStatus(ctx context.Context) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Initialize result
	result := map[string]interface{}{
		"cnaps": map[string]interface{}{
			"compliance_rate":     0.0,
			"compliant_employees": 0,
			"total_employees":     0,
			"issues":              []interface{}{},
		},
		"ostie": map[string]interface{}{
			"compliance_rate":     0.0,
			"compliant_employees": 0,
			"total_employees":     0,
			"issues":              []interface{}{},
		},
		"irsa": map[string]interface{}{
			"compliance_rate":     0.0,
			"compliant_employees": 0,
			"total_employees":     0,
			"issues":              []interface{}{},
		},
		"attendance": map[string]interface{}{
			"compliance_rate":     0.0,
			"days_tracked":        0,
			"total_required_days": 0,
			"issues":              []interface{}{},
		},
		"leave_balance": map[string]interface{}{
			"compliance_rate":         0.0,
			"employees_within_limits": 0,
			"total_employees":         0,
			"issues":                  []interface{}{},
		},
	}

	// Get total employees
	var totalEmployees int64
	if err := r.db.WithContext(ctx).Table("employees").
		Where("status = ?", "active").
		Count(&totalEmployees).Error; err != nil {
		return nil, fmt.Errorf("count total employees: %w", err)
	}

	// CNAPS compliance
	var cnapsCompliant int64
	if err := r.db.WithContext(ctx).Table("employees").
		Where("status = ? AND cnaps_number IS NOT NULL AND cnaps_number != ''", "active").
		Count(&cnapsCompliant).Error; err != nil {
		return nil, fmt.Errorf("count CNAPS compliant employees: %w", err)
	}

	cnapsRate := float64(cnapsCompliant) / float64(totalEmployees) * 100
	result["cnaps"].(map[string]interface{})["compliance_rate"] = cnapsRate
	result["cnaps"].(map[string]interface{})["compliant_employees"] = int(cnapsCompliant)
	result["cnaps"].(map[string]interface{})["total_employees"] = int(totalEmployees)

	// OSTIE compliance
	var ostieCompliant int64
	if err := r.db.WithContext(ctx).Table("employees").
		Where("status = ? AND ostie_number IS NOT NULL AND ostie_number != ''", "active").
		Count(&ostieCompliant).Error; err != nil {
		return nil, fmt.Errorf("count OSTIE compliant employees: %w", err)
	}

	ostieRate := float64(ostieCompliant) / float64(totalEmployees) * 100
	result["ostie"].(map[string]interface{})["compliance_rate"] = ostieRate
	result["ostie"].(map[string]interface{})["compliant_employees"] = int(ostieCompliant)
	result["ostie"].(map[string]interface{})["total_employees"] = int(totalEmployees)

	// IRSA compliance (assuming all employees are compliant for now)
	irsaRate := 100.0
	result["irsa"].(map[string]interface{})["compliance_rate"] = irsaRate
	result["irsa"].(map[string]interface{})["compliant_employees"] = int(totalEmployees)
	result["irsa"].(map[string]interface{})["total_employees"] = int(totalEmployees)

	// Attendance compliance (current month)
	now := time.Now()
	currentMonth := now.Format("2006-01")
	daysInMonth := getDaysInMonth(now.Year(), int(now.Month()))

	var daysTracked int64
	if err := r.db.WithContext(ctx).Table("attendance").
		Where("DATE_TRUNC('month', date) = ?", currentMonth).
		Count(&daysTracked).Error; err != nil {
		return nil, fmt.Errorf("count tracked days: %w", err)
	}

	attendanceRate := float64(daysTracked) / float64(daysInMonth) * 100
	result["attendance"].(map[string]interface{})["compliance_rate"] = attendanceRate
	result["attendance"].(map[string]interface{})["days_tracked"] = int(daysTracked)
	result["attendance"].(map[string]interface{})["total_required_days"] = daysInMonth

	// Leave balance compliance (assuming all employees are within limits for now)
	leaveRate := 100.0
	result["leave_balance"].(map[string]interface{})["compliance_rate"] = leaveRate
	result["leave_balance"].(map[string]interface{})["employees_within_limits"] = int(totalEmployees)
	result["leave_balance"].(map[string]interface{})["total_employees"] = int(totalEmployees)

	return result, nil
}

// GetWeeklyAttendanceSummary retrieves weekly attendance summary
func (r *Repo) GetWeeklyAttendanceSummary(ctx context.Context, employeeID *uuid.UUID, startDate time.Time) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Calculate week end (6 days after start)
	weekEnd := startDate.AddDate(0, 0, 6)

	// Initialize result
	result := map[string]interface{}{
		"week_start": startDate.Format("2006-01-02"),
		"week_end":   weekEnd.Format("2006-01-02"),
		"summary": map[string]interface{}{
			"present_rate":          0.0,
			"late_rate":             0.0,
			"absent_rate":           0.0,
			"total_employees":       0,
			"total_working_days":    5,
			"average_hours_per_day": 0.0,
			"total_overtime_hours":  0.0,
		},
		"daily_breakdown": []interface{}{},
	}

	// Build base query
	baseQuery := r.db.WithContext(ctx).Table("attendance").
		Where("date >= ? AND date <= ?", startDate.Format("2006-01-02"), weekEnd.Format("2006-01-02"))

	// If specific employee, filter by employee
	if employeeID != nil {
		baseQuery = baseQuery.Where("employee_id = ?", *employeeID)
	}

	// Get total employees for the period
	var totalEmployees int64
	employeesQuery := baseQuery
	if employeeID != nil {
		employeesQuery = employeesQuery.Select("COUNT(DISTINCT employee_id)")
	} else {
		employeesQuery = r.db.WithContext(ctx).Table("employees").
			Where("status = ?", "active")
	}

	if err := employeesQuery.Count(&totalEmployees).Error; err != nil {
		return nil, fmt.Errorf("count employees: %w", err)
	}

	// Get attendance counts
	var presentCount, lateCount, absentCount int64
	var totalHours, totalOvertime float64

	// Present count
	if err := baseQuery.Where("status = ?", "present").
		Count(&presentCount).Error; err != nil {
		return nil, fmt.Errorf("count present: %w", err)
	}

	// Late count
	if err := baseQuery.Where("status = ?", "late").
		Count(&lateCount).Error; err != nil {
		return nil, fmt.Errorf("count late: %w", err)
	}

	// Absent count
	if err := baseQuery.Where("status = ?", "absent").
		Count(&absentCount).Error; err != nil {
		return nil, fmt.Errorf("count absent: %w", err)
	}

	// Total hours
	if err := baseQuery.Select("COALESCE(SUM(total_hours), 0)").
		Scan(&totalHours).Error; err != nil {
		return nil, fmt.Errorf("sum total hours: %w", err)
	}

	// Total overtime
	if err := baseQuery.Select("COALESCE(SUM(overtime_hours), 0)").
		Scan(&totalOvertime).Error; err != nil {
		return nil, fmt.Errorf("sum overtime hours: %w", err)
	}

	// Calculate rates
	totalAttendance := presentCount + lateCount + absentCount
	presentRate := float64(presentCount) / float64(totalAttendance) * 100
	lateRate := float64(lateCount) / float64(totalAttendance) * 100
	absentRate := float64(absentCount) / float64(totalAttendance) * 100
	averageHours := totalHours / float64(totalAttendance)

	// Update summary
	result["summary"].(map[string]interface{})["present_rate"] = presentRate
	result["summary"].(map[string]interface{})["late_rate"] = lateRate
	result["summary"].(map[string]interface{})["absent_rate"] = absentRate
	result["summary"].(map[string]interface{})["total_employees"] = int(totalEmployees)
	result["summary"].(map[string]interface{})["average_hours_per_day"] = averageHours
	result["summary"].(map[string]interface{})["total_overtime_hours"] = totalOvertime

	// Get daily breakdown (simplified for now)
	dailyBreakdown := make([]interface{}, 0)
	for i := 0; i < 7; i++ {
		currentDate := startDate.AddDate(0, 0, i)
		if currentDate.After(weekEnd) {
			break
		}

		var dailyPresent, dailyLate, dailyAbsent, dailyOvertime int64
		dailyQuery := baseQuery.Where("date = ?", currentDate.Format("2006-01-02"))

		dailyQuery.Where("status = ?", "present").Count(&dailyPresent)
		dailyQuery.Where("status = ?", "late").Count(&dailyLate)
		dailyQuery.Where("status = ?", "absent").Count(&dailyAbsent)
		dailyQuery.Select("COALESCE(SUM(overtime_hours), 0)").Scan(&dailyOvertime)

		dailyBreakdown = append(dailyBreakdown, map[string]interface{}{
			"date":           currentDate.Format("2006-01-02"),
			"present":        int(dailyPresent),
			"late":           int(dailyLate),
			"absent":         int(dailyAbsent),
			"on_leave":       0, // Simplified for now
			"overtime_hours": dailyOvertime,
		})
	}

	result["daily_breakdown"] = dailyBreakdown

	return result, nil
}

// GetBadgeCounts retrieves badge counts for navigation items
func (r *Repo) GetBadgeCounts(ctx context.Context, userID uuid.UUID, userRole string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Initialize result based on role
	result := make(map[string]interface{})

	// Always include unread notifications
	var unreadNotifications int64
	if err := r.db.WithContext(ctx).Table("notifications").
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&unreadNotifications).Error; err != nil {
		return nil, fmt.Errorf("count unread notifications: %w", err)
	}
	result["unread_notifications"] = unreadNotifications

	// Role-specific counts
	switch userRole {
	case "admin", "hr":
		// Pending leaves
		var pendingLeaves int64
		if err := r.db.WithContext(ctx).Table("leaves").
			Where("status = ?", "pending").
			Count(&pendingLeaves).Error; err != nil {
			return nil, fmt.Errorf("count pending leaves: %w", err)
		}
		result["pending_leaves"] = pendingLeaves

		// Pending attendance corrections
		var pendingAttendanceCorrections int64
		if err := r.db.WithContext(ctx).Table("attendance_corrections").
			Where("status = ?", "pending").
			Count(&pendingAttendanceCorrections).Error; err != nil {
			return nil, fmt.Errorf("count pending attendance corrections: %w", err)
		}
		result["pending_attendance_corrections"] = pendingAttendanceCorrections

		// Pending payroll drafts
		var pendingPayrollDrafts int64
		if err := r.db.WithContext(ctx).Table("payroll_drafts").
			Where("status = ?", "draft").
			Count(&pendingPayrollDrafts).Error; err != nil {
			return nil, fmt.Errorf("count pending payroll drafts: %w", err)
		}
		result["pending_payroll_drafts"] = pendingPayrollDrafts

		// Pending declarations
		var pendingDeclarations int64
		if err := r.db.WithContext(ctx).Table("declarations").
			Where("status = ?", "pending").
			Count(&pendingDeclarations).Error; err != nil {
			return nil, fmt.Errorf("count pending declarations: %w", err)
		}
		result["pending_declarations"] = pendingDeclarations

		// Pending KPI reviews
		var pendingKPIReviews int64
		if err := r.db.WithContext(ctx).Table("kpis").
			Where("status = ?", "pending_review").
			Count(&pendingKPIReviews).Error; err != nil {
			return nil, fmt.Errorf("count pending KPI reviews: %w", err)
		}
		result["pending_kpi_reviews"] = pendingKPIReviews

		// Open support tickets
		var openSupportTickets int64
		if err := r.db.WithContext(ctx).Table("support_tickets").
			Where("status IN ?", []string{"open", "in_progress"}).
			Count(&openSupportTickets).Error; err != nil {
			return nil, fmt.Errorf("count open support tickets: %w", err)
		}
		result["open_support_tickets"] = openSupportTickets

	case "employee":
		// Employee's own pending leaves
		var pendingLeaves int64
		if err := r.db.WithContext(ctx).Table("leaves").
			Where("employee_id = ? AND status = ?", userID, "pending").
			Count(&pendingLeaves).Error; err != nil {
			return nil, fmt.Errorf("count pending leaves: %w", err)
		}
		result["pending_leaves"] = pendingLeaves

		// Employee's own open support tickets
		var openSupportTickets int64
		if err := r.db.WithContext(ctx).Table("support_tickets").
			Where("user_id = ? AND status IN ?", userID, []string{"open", "in_progress"}).
			Count(&openSupportTickets).Error; err != nil {
			return nil, fmt.Errorf("count open support tickets: %w", err)
		}
		result["open_support_tickets"] = openSupportTickets

	case "accountant":
		// Pending payroll approvals
		var pendingPayrollApprovals int64
		if err := r.db.WithContext(ctx).Table("payroll_drafts").
			Where("status = ?", "ready_for_approval").
			Count(&pendingPayrollApprovals).Error; err != nil {
			return nil, fmt.Errorf("count pending payroll approvals: %w", err)
		}
		result["pending_payroll_approvals"] = pendingPayrollApprovals

		// Pending declarations
		var pendingDeclarations int64
		if err := r.db.WithContext(ctx).Table("declarations").
			Where("status = ?", "pending").
			Count(&pendingDeclarations).Error; err != nil {
			return nil, fmt.Errorf("count pending declarations: %w", err)
		}
		result["pending_declarations"] = pendingDeclarations
	}

	return result, nil
}

// GetCalendarEvents retrieves calendar events for date range
func (r *Repo) GetCalendarEvents(ctx context.Context, query CalendarEventsQuery, userID uuid.UUID, userRole string) ([]CalendarEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var events []CalendarEvent

	// Determine which event types to include
	includeLeaves := len(query.EventTypes) == 0 || contains(query.EventTypes, "leave")
	includeHolidays := len(query.EventTypes) == 0 || contains(query.EventTypes, "holiday")

	// Get leave events if requested
	if includeLeaves {
		leavesQuery := r.db.WithContext(ctx).Table("leaves").
			Select(`
				leaves.id,
				'leave' as type,
				employees.name || ' - ' || leaves.leave_type as title,
				leaves.start_date,
				leaves.end_date,
				leaves.employee_id,
				employees.name as employee_name,
				leaves.leave_type,
				leaves.status
			`).
			Joins("JOIN employees ON leaves.employee_id = employees.id").
			Where("(leaves.start_date <= ? AND leaves.end_date >= ?) OR (leaves.start_date >= ? AND leaves.start_date <= ?)",
				query.EndDate, query.StartDate, query.StartDate, query.EndDate)

		// Filter by employee if provided
		if query.EmployeeID != nil {
			leavesQuery = leavesQuery.Where("leaves.employee_id = ?", *query.EmployeeID)
		}

		// For employees, only show their own leaves
		if userRole == "employee" {
			leavesQuery = leavesQuery.Where("leaves.employee_id = ?", userID)
		}

		var leaveEvents []CalendarEvent
		if err := leavesQuery.Scan(&leaveEvents).Error; err != nil {
			return nil, fmt.Errorf("get leave events: %w", err)
		}

		// Set colors based on leave status
		for i := range leaveEvents {
			switch leaveEvents[i].Status {
			case "approved":
				leaveEvents[i].Color = "#3b82f6" // Blue
			case "pending":
				leaveEvents[i].Color = "#f59e0b" // Yellow/Orange
			case "rejected":
				leaveEvents[i].Color = "#ef4444" // Red
			default:
				leaveEvents[i].Color = "#6b7280" // Gray
			}
		}

		events = append(events, leaveEvents...)
	}

	// Get company holidays if requested
	if includeHolidays {
		var holidays []CalendarEvent
		holidaysQuery := r.db.WithContext(ctx).Table("company_holidays").
			Select(`
				id,
				'holiday' as type,
				name as title,
				date as start_date,
				date as end_date,
				NULL as employee_id,
				NULL as employee_name,
				NULL as leave_type,
				NULL as status
			`).
			Where("date >= ? AND date <= ?", query.StartDate, query.EndDate)

		if err := holidaysQuery.Scan(&holidays).Error; err != nil {
			// If table doesn't exist, skip holidays
			if err != gorm.ErrRecordNotFound {
				// Log error but don't fail the request
				fmt.Printf("Warning: Could not get company holidays: %v\n", err)
			}
		} else {
			// Set holiday color
			for i := range holidays {
				holidays[i].Color = "#ef4444" // Red
				holidays[i].IsCompanyWide = true
			}
			events = append(events, holidays...)
		}
	}

	return events, nil
}

// contains checks if a string slice contains a specific value
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// getDaysInMonth returns the number of days in a month
func getDaysInMonth(year int, month int) int {
	// This is a simplified version
	// In a real implementation, you'd use a proper date library
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		// Simplified leap year check
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			return 29
		}
		return 28
	default:
		return 30
	}
}


