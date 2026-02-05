package company

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for company settings
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new company repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// Get retrieves company settings
func (r *Repo) Get(ctx context.Context) (*CompanySettings, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var settings CompanySettings
	if err := r.db.WithContext(ctx).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("company settings not found")
		}
		return nil, fmt.Errorf("get company settings: %w", err)
	}
	return &settings, nil
}

// Update updates company settings
func (r *Repo) Update(ctx context.Context, settings *CompanySettings) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	settings.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(settings).Error; err != nil {
		return fmt.Errorf("update company settings: %w", err)
	}
	return nil
}

// UpdateLogo updates the company logo URL
func (r *Repo) UpdateLogo(ctx context.Context, logoURL string, updatedBy uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&CompanySettings{}).
		Updates(map[string]interface{}{
			"logo_url":   logoURL,
			"updated_at": time.Now(),
			"updated_by": updatedBy,
		})

	if result.Error != nil {
		return fmt.Errorf("update company logo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("company settings not found")
	}

	return nil
}

// CreateHoliday creates a new company holiday
func (r *Repo) CreateHoliday(ctx context.Context, holiday *CompanyHoliday) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(holiday).Error; err != nil {
		return fmt.Errorf("create company holiday: %w", err)
	}
	return nil
}

// GetHolidayByID retrieves a company holiday by ID
func (r *Repo) GetHolidayByID(ctx context.Context, id uuid.UUID) (*CompanyHoliday, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var holiday CompanyHoliday
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&holiday).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("holiday not found")
		}
		return nil, fmt.Errorf("get holiday: %w", err)
	}
	return &holiday, nil
}

// ListHolidays retrieves company holidays with filtering
func (r *Repo) ListHolidays(ctx context.Context, query HolidaysListQuery) ([]CompanyHoliday, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var holidays []CompanyHoliday
	db := r.db.WithContext(ctx).Model(&CompanyHoliday{})

	// Apply filters
	if query.Year > 0 {
		db = db.Where("EXTRACT(YEAR FROM date) = ?", query.Year)
	}

	if query.StartDate != "" {
		db = db.Where("date >= ?", query.StartDate)
	}

	if query.EndDate != "" {
		db = db.Where("date <= ?", query.EndDate)
	}

	if err := db.Order("date ASC").Find(&holidays).Error; err != nil {
		return nil, fmt.Errorf("list holidays: %w", err)
	}

	return holidays, nil
}

// UpdateHoliday updates a company holiday
func (r *Repo) UpdateHoliday(ctx context.Context, holiday *CompanyHoliday) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Save(holiday).Error; err != nil {
		return fmt.Errorf("update holiday: %w", err)
	}
	return nil
}

// DeleteHoliday deletes a company holiday
func (r *Repo) DeleteHoliday(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&CompanyHoliday{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("delete holiday: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("holiday not found")
	}

	return nil
}
