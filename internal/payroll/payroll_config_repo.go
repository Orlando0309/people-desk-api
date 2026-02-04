package payroll

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConfigRepo handles database operations for payroll configurations
type ConfigRepo struct {
	db *gorm.DB
}

// NewConfigRepo creates a new payroll configuration repository
func NewConfigRepo(database *gorm.DB) *ConfigRepo {
	return &ConfigRepo{db: database}
}

// GetConfigValue retrieves a configuration value by key
func (r *ConfigRepo) GetConfigValue(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var config PayrollConfiguration
	if err := r.db.WithContext(ctx).Where("key = ? AND is_active = ?", key, true).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("configuration key '%s' not found", key)
		}
		return "", fmt.Errorf("get configuration: %w", err)
	}
	return config.Value, nil
}

// GetConfigValueAsFloat retrieves a configuration value as a float64
func (r *ConfigRepo) GetConfigValueAsFloat(ctx context.Context, key string) (float64, error) {
	value, err := r.GetConfigValue(ctx, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(value, 64)
}

// GetConfigValueAsInt retrieves a configuration value as an int
func (r *ConfigRepo) GetConfigValueAsInt(ctx context.Context, key string) (int, error) {
	value, err := r.GetConfigValue(ctx, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

// GetConfigValueAsBool retrieves a configuration value as a bool
func (r *ConfigRepo) GetConfigValueAsBool(ctx context.Context, key string) (bool, error) {
	value, err := r.GetConfigValue(ctx, key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(value)
}

// CreateConfig creates a new payroll configuration
func (r *ConfigRepo) CreateConfig(ctx context.Context, config *PayrollConfiguration) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(config).Error; err != nil {
		return fmt.Errorf("create configuration: %w", err)
	}
	return nil
}

// UpdateConfig updates a payroll configuration
func (r *ConfigRepo) UpdateConfig(ctx context.Context, config *PayrollConfiguration) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	config.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(config).Error; err != nil {
		return fmt.Errorf("update configuration: %w", err)
	}
	return nil
}

// GetConfigByID retrieves a configuration by ID
func (r *ConfigRepo) GetConfigByID(ctx context.Context, id uuid.UUID) (*PayrollConfiguration, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var config PayrollConfiguration
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("configuration not found")
		}
		return nil, fmt.Errorf("get configuration: %w", err)
	}
	return &config, nil
}

// ListConfigs retrieves payroll configurations with filtering
func (r *ConfigRepo) ListConfigs(ctx context.Context, query PayrollConfigurationListQuery) ([]PayrollConfiguration, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var configs []PayrollConfiguration
	var total int64

	db := r.db.WithContext(ctx).Model(&PayrollConfiguration{})

	// Apply filters
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.Key != "" {
		db = db.Where("key LIKE ?", "%"+query.Key+"%")
	}
	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count configurations: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("category, key").Find(&configs).Error; err != nil {
		return nil, 0, fmt.Errorf("list configurations: %w", err)
	}

	return configs, total, nil
}

// DeleteConfig soft deletes a payroll configuration
func (r *ConfigRepo) DeleteConfig(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&PayrollConfiguration{})
	if result.Error != nil {
		return fmt.Errorf("delete configuration: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("configuration not found")
	}
	return nil
}

// GetActiveIRSABrackets retrieves all active IRSA tax brackets ordered by sort_order
func (r *ConfigRepo) GetActiveIRSABrackets(ctx context.Context) ([]IRSATaxBracket, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var brackets []IRSATaxBracket
	if err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC").
		Find(&brackets).Error; err != nil {
		return nil, fmt.Errorf("get IRSA brackets: %w", err)
	}
	return brackets, nil
}

// CreateIRSABracket creates a new IRSA tax bracket
func (r *ConfigRepo) CreateIRSABracket(ctx context.Context, bracket *IRSATaxBracket) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	bracket.CreatedAt = time.Now()
	bracket.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(bracket).Error; err != nil {
		return fmt.Errorf("create IRSA bracket: %w", err)
	}
	return nil
}

// UpdateIRSABracket updates an IRSA tax bracket
func (r *ConfigRepo) UpdateIRSABracket(ctx context.Context, bracket *IRSATaxBracket) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	bracket.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(bracket).Error; err != nil {
		return fmt.Errorf("update IRSA bracket: %w", err)
	}
	return nil
}

// GetIRSABracketByID retrieves an IRSA tax bracket by ID
func (r *ConfigRepo) GetIRSABracketByID(ctx context.Context, id uuid.UUID) (*IRSATaxBracket, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var bracket IRSATaxBracket
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&bracket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("IRSA bracket not found")
		}
		return nil, fmt.Errorf("get IRSA bracket: %w", err)
	}
	return &bracket, nil
}

// ListIRSABrackets retrieves IRSA tax brackets with filtering
func (r *ConfigRepo) ListIRSABrackets(ctx context.Context, query IRSATaxBracketListQuery) ([]IRSATaxBracket, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var brackets []IRSATaxBracket
	var total int64

	db := r.db.WithContext(ctx).Model(&IRSATaxBracket{})

	// Apply filters
	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count IRSA brackets: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("sort_order ASC").Find(&brackets).Error; err != nil {
		return nil, 0, fmt.Errorf("list IRSA brackets: %w", err)
	}

	return brackets, total, nil
}

// DeleteIRSABracket soft deletes an IRSA tax bracket
func (r *ConfigRepo) DeleteIRSABracket(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&IRSATaxBracket{})
	if result.Error != nil {
		return fmt.Errorf("delete IRSA bracket: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("IRSA bracket not found")
	}
	return nil
}
