package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for employees
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new employee repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// Create creates a new employee
func (r *Repo) Create(ctx context.Context, employee *Employee) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(employee).Error; err != nil {
		return fmt.Errorf("create employee: %w", err)
	}
	return nil
}

// GetByID retrieves an employee by ID
func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var employee Employee
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("get employee by id: %w", err)
	}
	return &employee, nil
}

// Update updates an employee
func (r *Repo) Update(ctx context.Context, employee *Employee) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	employee.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(employee).Error; err != nil {
		return fmt.Errorf("update employee: %w", err)
	}
	return nil
}

// Delete soft deletes an employee
func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Delete(&Employee{}, id).Error; err != nil {
		return fmt.Errorf("delete employee: %w", err)
	}
	return nil
}

// List retrieves employees with filtering and pagination
func (r *Repo) List(ctx context.Context, query EmployeeListQuery) ([]Employee, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var employees []Employee
	var total int64

	db := r.db.WithContext(ctx).Model(&Employee{})

	// Apply filters
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("first_name ILIKE ? OR last_name ILIKE ? OR national_id ILIKE ?", 
			searchPattern, searchPattern, searchPattern)
	}

	if query.Department != "" {
		db = db.Where("department = ?", query.Department)
	}

	if query.Position != "" {
		db = db.Where("position = ?", query.Position)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count employees: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("created_at DESC").Find(&employees).Error; err != nil {
		return nil, 0, fmt.Errorf("list employees: %w", err)
	}

	return employees, total, nil
}

// GetByCompanyID retrieves all employees for a company
func (r *Repo) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var employees []Employee
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("get employees by company id: %w", err)
	}
	return employees, nil
}

// GetByDepartment retrieves all employees in a department
func (r *Repo) GetByDepartment(ctx context.Context, department string) ([]Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var employees []Employee
	if err := r.db.WithContext(ctx).Where("department = ? AND status = ?", department, "active").Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("get employees by department: %w", err)
	}
	return employees, nil
}

// GetSubordinates retrieves all employees reporting to a manager
func (r *Repo) GetSubordinates(ctx context.Context, managerID uuid.UUID) ([]Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var employees []Employee
	if err := r.db.WithContext(ctx).Where("manager_id = ? AND status = ?", managerID, "active").Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("get subordinates: %w", err)
	}
	return employees, nil
}

// GetDepartments retrieves all unique departments
func (r *Repo) GetDepartments(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var departments []string
	if err := r.db.WithContext(ctx).Model(&Employee{}).Distinct("department").Pluck("department", &departments).Error; err != nil {
		return nil, fmt.Errorf("get departments: %w", err)
	}
	return departments, nil
}

// GetPositions retrieves all unique positions
func (r *Repo) GetPositions(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var positions []string
	if err := r.db.WithContext(ctx).Model(&Employee{}).Distinct("position").Pluck("position", &positions).Error; err != nil {
		return nil, fmt.Errorf("get positions: %w", err)
	}
	return positions, nil
}
