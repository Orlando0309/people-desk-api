package leave

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for leaves
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new leave repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// Create creates a new leave request
func (r *Repo) Create(ctx context.Context, leave *Leave) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(leave).Error; err != nil {
		return fmt.Errorf("create leave: %w", err)
	}
	return nil
}

// GetByID retrieves a leave by ID
func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*Leave, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var leave Leave
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&leave).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("leave not found")
		}
		return nil, fmt.Errorf("get leave by id: %w", err)
	}
	return &leave, nil
}

// Update updates a leave
func (r *Repo) Update(ctx context.Context, leave *Leave) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	leave.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(leave).Error; err != nil {
		return fmt.Errorf("update leave: %w", err)
	}
	return nil
}

// Delete soft deletes a leave
func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Delete(&Leave{}, id).Error; err != nil {
		return fmt.Errorf("delete leave: %w", err)
	}
	return nil
}

// List retrieves leaves with filtering and pagination
func (r *Repo) List(ctx context.Context, query LeaveListQuery) ([]Leave, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var leaves []Leave
	var total int64

	db := r.db.WithContext(ctx).Model(&Leave{})

	// Apply filters
	if query.EmployeeID != nil {
		db = db.Where("employee_id = ?", *query.EmployeeID)
	}

	if query.LeaveType != "" {
		db = db.Where("leave_type = ?", query.LeaveType)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.StartDate != nil {
		db = db.Where("start_date >= ?", *query.StartDate)
	}

	if query.EndDate != nil {
		db = db.Where("end_date <= ?", *query.EndDate)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count leaves: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("created_at DESC").Find(&leaves).Error; err != nil {
		return nil, 0, fmt.Errorf("list leaves: %w", err)
	}

	return leaves, total, nil
}

// GetPendingLeaves retrieves all pending leave requests
func (r *Repo) GetPendingLeaves(ctx context.Context) ([]Leave, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var leaves []Leave
	if err := r.db.WithContext(ctx).Where("status = ?", "pending").Order("created_at ASC").Find(&leaves).Error; err != nil {
		return nil, fmt.Errorf("get pending leaves: %w", err)
	}
	return leaves, nil
}

// GetLeaveBalance calculates leave balance for an employee
func (r *Repo) GetLeaveBalance(ctx context.Context, employeeID uuid.UUID, year int) (*LeaveBalance, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	balance := &LeaveBalance{
		EmployeeID:  employeeID,
		AnnualTotal: 30, // 30 days per year as per Madagascar law
	}

	// Calculate used leave days for each type
	startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	// Annual leave used
	var annualUsed *float64
	if err := r.db.WithContext(ctx).Model(&Leave{}).
		Where("employee_id = ? AND leave_type = ? AND status = ? AND start_date >= ? AND end_date <= ?",
			employeeID, "annual", "approved", startOfYear, endOfYear).
		Select("COALESCE(SUM(days_requested), 0)").Scan(&annualUsed).Error; err != nil {
		return nil, fmt.Errorf("calculate annual leave: %w", err)
	}
	if annualUsed != nil {
		balance.AnnualUsed = *annualUsed
	}

	// Sick leave used
	var sickUsed *float64
	if err := r.db.WithContext(ctx).Model(&Leave{}).
		Where("employee_id = ? AND leave_type = ? AND status = ? AND start_date >= ? AND end_date <= ?",
			employeeID, "sick", "approved", startOfYear, endOfYear).
		Select("COALESCE(SUM(days_requested), 0)").Scan(&sickUsed).Error; err != nil {
		return nil, fmt.Errorf("calculate sick leave: %w", err)
	}
	if sickUsed != nil {
		balance.SickUsed = *sickUsed
	}

	// Maternity leave used
	var maternityUsed *float64
	if err := r.db.WithContext(ctx).Model(&Leave{}).
		Where("employee_id = ? AND leave_type = ? AND status = ? AND start_date >= ? AND end_date <= ?",
			employeeID, "maternity", "approved", startOfYear, endOfYear).
		Select("COALESCE(SUM(days_requested), 0)").Scan(&maternityUsed).Error; err != nil {
		return nil, fmt.Errorf("calculate maternity leave: %w", err)
	}
	if maternityUsed != nil {
		balance.MaternityUsed = *maternityUsed
	}

	// Exceptional leave used
	var exceptionalUsed *float64
	if err := r.db.WithContext(ctx).Model(&Leave{}).
		Where("employee_id = ? AND leave_type = ? AND status = ? AND start_date >= ? AND end_date <= ?",
			employeeID, "exceptional", "approved", startOfYear, endOfYear).
		Select("COALESCE(SUM(days_requested), 0)").Scan(&exceptionalUsed).Error; err != nil {
		return nil, fmt.Errorf("calculate exceptional leave: %w", err)
	}
	if exceptionalUsed != nil {
		balance.ExceptionalUsed = *exceptionalUsed
	}

	// Paternity leave used
	var paternityUsed *float64
	if err := r.db.WithContext(ctx).Model(&Leave{}).
		Where("employee_id = ? AND leave_type = ? AND status = ? AND start_date >= ? AND end_date <= ?",
			employeeID, "paternity", "approved", startOfYear, endOfYear).
		Select("COALESCE(SUM(days_requested), 0)").Scan(&paternityUsed).Error; err != nil {
		return nil, fmt.Errorf("calculate paternity leave: %w", err)
	}
	if paternityUsed != nil {
		balance.PaternityUsed = *paternityUsed
	}

	// Unpaid leave used
	var unpaidUsed *float64
	if err := r.db.WithContext(ctx).Model(&Leave{}).
		Where("employee_id = ? AND leave_type = ? AND status = ? AND start_date >= ? AND end_date <= ?",
			employeeID, "unpaid", "approved", startOfYear, endOfYear).
		Select("COALESCE(SUM(days_requested), 0)").Scan(&unpaidUsed).Error; err != nil {
		return nil, fmt.Errorf("calculate unpaid leave: %w", err)
	}
	if unpaidUsed != nil {
		balance.UnpaidUsed = *unpaidUsed
	}

	// Calculate remaining annual leave
	balance.AnnualRemaining = balance.AnnualTotal - balance.AnnualUsed

	return balance, nil
}

// CheckOverlap checks if there's an overlapping leave request
func (r *Repo) CheckOverlap(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time, excludeID *uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := r.db.WithContext(ctx).Model(&Leave{}).
		Where("employee_id = ? AND status IN (?, ?)", employeeID, "pending", "approved").
		Where("(start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?)",
			endDate, startDate, startDate, endDate, startDate, endDate)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check overlap: %w", err)
	}

	return count > 0, nil
}
