package kpi

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for KPIs and performance reviews
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new KPI repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// CreateKPI creates a new KPI
func (r *Repo) CreateKPI(ctx context.Context, kpi *KPI) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(kpi).Error; err != nil {
		return fmt.Errorf("create KPI: %w", err)
	}
	return nil
}

// GetKPIByID retrieves a KPI by ID
func (r *Repo) GetKPIByID(ctx context.Context, id uuid.UUID) (*KPI, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var kpi KPI
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&kpi).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("KPI not found")
		}
		return nil, fmt.Errorf("get KPI: %w", err)
	}
	return &kpi, nil
}

// UpdateKPI updates a KPI
func (r *Repo) UpdateKPI(ctx context.Context, kpi *KPI) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	kpi.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(kpi).Error; err != nil {
		return fmt.Errorf("update KPI: %w", err)
	}
	return nil
}

// DeleteKPI deletes a KPI
func (r *Repo) DeleteKPI(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&KPI{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete KPI: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("KPI not found")
	}
	return nil
}

// ListKPIs retrieves KPIs with filtering
func (r *Repo) ListKPIs(ctx context.Context, query KPIListQuery) ([]KPI, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var kpis []KPI
	var total int64

	db := r.db.WithContext(ctx).Model(&KPI{})

	// Apply filters
	if query.Department != "" {
		db = db.Where("department = ?", query.Department)
	}
	if query.Position != "" {
		db = db.Where("position = ?", query.Position)
	}
	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count KPIs: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("created_at DESC").Find(&kpis).Error; err != nil {
		return nil, 0, fmt.Errorf("list KPIs: %w", err)
	}

	return kpis, total, nil
}

// CreatePerformanceReview creates a new performance review
func (r *Repo) CreatePerformanceReview(ctx context.Context, review *PerformanceReview) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Check if review already exists for this employee, KPI, and period
	var existing PerformanceReview
	err := r.db.WithContext(ctx).Where("employee_id = ? AND kpi_id = ? AND review_period_start = ? AND review_period_end = ?",
		review.EmployeeID, review.KPIID, review.ReviewPeriodStart, review.ReviewPeriodEnd).First(&existing).Error
	if err == nil {
		return fmt.Errorf("performance review already exists for this employee, KPI, and period")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("check existing review: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		return fmt.Errorf("create performance review: %w", err)
	}
	return nil
}

// GetPerformanceReviewByID retrieves a performance review by ID
func (r *Repo) GetPerformanceReviewByID(ctx context.Context, id uuid.UUID) (*PerformanceReview, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var review PerformanceReview
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&review).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance review not found")
		}
		return nil, fmt.Errorf("get performance review: %w", err)
	}
	return &review, nil
}

// UpdatePerformanceReview updates a performance review
func (r *Repo) UpdatePerformanceReview(ctx context.Context, review *PerformanceReview) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	review.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(review).Error; err != nil {
		return fmt.Errorf("update performance review: %w", err)
	}
	return nil
}

// DeletePerformanceReview deletes a performance review
func (r *Repo) DeletePerformanceReview(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&PerformanceReview{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete performance review: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance review not found")
	}
	return nil
}

// ListPerformanceReviews retrieves performance reviews with filtering
func (r *Repo) ListPerformanceReviews(ctx context.Context, query PerformanceReviewListQuery) ([]PerformanceReview, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var reviews []PerformanceReview
	var total int64

	db := r.db.WithContext(ctx).Model(&PerformanceReview{})

	// Apply filters
	if query.EmployeeID != nil {
		db = db.Where("employee_id = ?", *query.EmployeeID)
	}
	if query.KPIID != nil {
		db = db.Where("kpi_id = ?", *query.KPIID)
	}
	if query.ReviewerID != nil {
		db = db.Where("reviewer_id = ?", *query.ReviewerID)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.ReviewPeriodStart != nil {
		db = db.Where("review_period_start >= ?", *query.ReviewPeriodStart)
	}
	if query.ReviewPeriodEnd != nil {
		db = db.Where("review_period_end <= ?", *query.ReviewPeriodEnd)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count performance reviews: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("review_period_start DESC, created_at DESC").Find(&reviews).Error; err != nil {
		return nil, 0, fmt.Errorf("list performance reviews: %w", err)
	}

	return reviews, total, nil
}

// GetPerformanceReviewsByEmployee retrieves all performance reviews for an employee in a period
func (r *Repo) GetPerformanceReviewsByEmployee(ctx context.Context, employeeID uuid.UUID, periodStart, periodEnd time.Time) ([]PerformanceReview, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var reviews []PerformanceReview
	if err := r.db.WithContext(ctx).Where("employee_id = ? AND review_period_start >= ? AND review_period_end <= ?",
		employeeID, periodStart, periodEnd).Find(&reviews).Error; err != nil {
		return nil, fmt.Errorf("get performance reviews by employee: %w", err)
	}
	return reviews, nil
}

// GetKPIsByDepartment retrieves all KPIs for a department
func (r *Repo) GetKPIsByDepartment(ctx context.Context, department string) ([]KPI, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var kpis []KPI
	if err := r.db.WithContext(ctx).Where("department = ? AND is_active = ?", department, true).Find(&kpis).Error; err != nil {
		return nil, fmt.Errorf("get KPIs by department: %w", err)
	}
	return kpis, nil
}

// GetKPIsByPosition retrieves all KPIs for a position
func (r *Repo) GetKPIsByPosition(ctx context.Context, position string) ([]KPI, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var kpis []KPI
	if err := r.db.WithContext(ctx).Where("position = ? AND is_active = ?", position, true).Find(&kpis).Error; err != nil {
		return nil, fmt.Errorf("get KPIs by position: %w", err)
	}
	return kpis, nil
}

// CalculateFinalScore calculates the final score for a performance review
func (r *Repo) CalculateFinalScore(ctx context.Context, reviewID uuid.UUID) (*PerformanceReview, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	review, err := r.GetPerformanceReviewByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}

	// Calculate final score as average of self and manager scores
	// If only one score is provided, use that score
	if review.SelfScore != nil && review.ManagerScore != nil {
		avgScore := (*review.SelfScore + *review.ManagerScore) / 2
		review.FinalScore = &avgScore
	} else if review.SelfScore != nil {
		review.FinalScore = review.SelfScore
	} else if review.ManagerScore != nil {
		review.FinalScore = review.ManagerScore
	}

	review.Status = "completed"
	review.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(review).Error; err != nil {
		return nil, fmt.Errorf("calculate final score: %w", err)
	}

	return review, nil
}

// GeneratePerformanceReport generates a performance report for an employee
func (r *Repo) GeneratePerformanceReport(ctx context.Context, employeeID uuid.UUID, periodStart, periodEnd time.Time) (*PerformanceReport, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	report := &PerformanceReport{
		EmployeeID:        employeeID,
		EmployeeName:      "Employee Name", // TODO: Get from employee service
		ReviewPeriodStart: periodStart,
		ReviewPeriodEnd:   periodEnd,
	}

	// Get all performance reviews for the employee in the period
	reviews, err := r.GetPerformanceReviewsByEmployee(ctx, employeeID, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("get performance reviews for report: %w", err)
	}

	// Get KPI details for each review
	var totalScore, totalWeight float64
	for _, review := range reviews {
		kpi, err := r.GetKPIByID(ctx, review.KPIID)
		if err != nil {
			continue
		}

		item := KPIReportItem{
			KPIName:           kpi.Name,
			TargetValue:       kpi.TargetValue,
			WeightPercentage:  kpi.WeightPercentage,
			SelfScore:         review.SelfScore,
			ManagerScore:      review.ManagerScore,
			FinalScore:        review.FinalScore,
			SelfAssessment:    review.SelfAssessment,
			ManagerAssessment: review.ManagerAssessment,
			Status:            review.Status,
		}

		report.KPIs = append(report.KPIs, item)

		// Calculate weighted score
		if review.FinalScore != nil {
			totalScore += *review.FinalScore * (kpi.WeightPercentage / 100)
			totalWeight += kpi.WeightPercentage
		}
	}

	// Calculate overall score
	if totalWeight > 0 {
		report.OverallScore = totalScore / (totalWeight / 100)
	}
	report.TotalWeight = totalWeight

	return report, nil
}
