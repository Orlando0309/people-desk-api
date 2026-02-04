package declarations

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for declarations
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new declarations repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// CreateDeclaration creates a new monthly declaration
func (r *Repo) CreateDeclaration(ctx context.Context, declaration *MonthlyDeclaration) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Check if declaration already exists for this type and period
	var existing MonthlyDeclaration
	err := r.db.WithContext(ctx).Where("declaration_type = ? AND declaration_period_start = ? AND declaration_period_end = ?",
		declaration.DeclarationType, declaration.DeclarationPeriodStart, declaration.DeclarationPeriodEnd).First(&existing).Error
	if err == nil {
		return fmt.Errorf("declaration already exists for this type and period")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("check existing declaration: %w", err)
	}

	// Generate declaration number if not provided
	if declaration.DeclarationNumber == "" {
		declaration.DeclarationNumber = generateDeclarationNumber(declaration.DeclarationType, declaration.DeclarationPeriodStart)
	}

	if err := r.db.WithContext(ctx).Create(declaration).Error; err != nil {
		return fmt.Errorf("create declaration: %w", err)
	}
	return nil
}

// GetDeclarationByID retrieves a declaration by ID
func (r *Repo) GetDeclarationByID(ctx context.Context, id uuid.UUID) (*MonthlyDeclaration, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var declaration MonthlyDeclaration
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&declaration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("declaration not found")
		}
		return nil, fmt.Errorf("get declaration: %w", err)
	}
	return &declaration, nil
}

// GetDeclarationByNumber retrieves a declaration by declaration number
func (r *Repo) GetDeclarationByNumber(ctx context.Context, declarationNumber string) (*MonthlyDeclaration, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var declaration MonthlyDeclaration
	if err := r.db.WithContext(ctx).Where("declaration_number = ?", declarationNumber).First(&declaration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("declaration not found")
		}
		return nil, fmt.Errorf("get declaration: %w", err)
	}
	return &declaration, nil
}

// UpdateDeclaration updates a declaration
func (r *Repo) UpdateDeclaration(ctx context.Context, declaration *MonthlyDeclaration) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	declaration.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(declaration).Error; err != nil {
		return fmt.Errorf("update declaration: %w", err)
	}
	return nil
}

// DeleteDeclaration deletes a declaration
func (r *Repo) DeleteDeclaration(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&MonthlyDeclaration{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete declaration: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("declaration not found")
	}
	return nil
}

// ListDeclarations retrieves declarations with filtering
func (r *Repo) ListDeclarations(ctx context.Context, query DeclarationListQuery) ([]MonthlyDeclaration, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var declarations []MonthlyDeclaration
	var total int64

	db := r.db.WithContext(ctx).Model(&MonthlyDeclaration{})

	// Apply filters
	if query.DeclarationType != "" {
		db = db.Where("declaration_type = ?", query.DeclarationType)
	}
	if query.DeclarationPeriodStart != nil {
		db = db.Where("declaration_period_start >= ?", *query.DeclarationPeriodStart)
	}
	if query.DeclarationPeriodEnd != nil {
		db = db.Where("declaration_period_end <= ?", *query.DeclarationPeriodEnd)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.AccountantID != nil {
		db = db.Where("accountant_id = ?", *query.AccountantID)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count declarations: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("declaration_period_start DESC, created_at DESC").Find(&declarations).Error; err != nil {
		return nil, 0, fmt.Errorf("list declarations: %w", err)
	}

	return declarations, total, nil
}

// CreateIRSATaxBracket creates a new IRSA tax bracket
func (r *Repo) CreateIRSATaxBracket(ctx context.Context, bracket *IRSATaxBracket) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(bracket).Error; err != nil {
		return fmt.Errorf("create IRSA tax bracket: %w", err)
	}
	return nil
}

// GetIRSATaxBracketByID retrieves an IRSA tax bracket by ID
func (r *Repo) GetIRSATaxBracketByID(ctx context.Context, id uuid.UUID) (*IRSATaxBracket, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var bracket IRSATaxBracket
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&bracket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("IRSA tax bracket not found")
		}
		return nil, fmt.Errorf("get IRSA tax bracket: %w", err)
	}
	return &bracket, nil
}

// UpdateIRSATaxBracket updates an IRSA tax bracket
func (r *Repo) UpdateIRSATaxBracket(ctx context.Context, bracket *IRSATaxBracket) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	bracket.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(bracket).Error; err != nil {
		return fmt.Errorf("update IRSA tax bracket: %w", err)
	}
	return nil
}

// DeleteIRSATaxBracket deletes an IRSA tax bracket
func (r *Repo) DeleteIRSATaxBracket(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&IRSATaxBracket{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete IRSA tax bracket: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("IRSA tax bracket not found")
	}
	return nil
}

// ListIRSATaxBrackets retrieves IRSA tax brackets with filtering
func (r *Repo) ListIRSATaxBrackets(ctx context.Context, query IRSATaxBracketListQuery) ([]IRSATaxBracket, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var brackets []IRSATaxBracket
	var total int64

	db := r.db.WithContext(ctx).Model(&IRSATaxBracket{})

	// Apply filters
	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}
	if query.EffectiveDate != nil {
		db = db.Where("effective_date <= ?", *query.EffectiveDate)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count IRSA tax brackets: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("min_income ASC").Find(&brackets).Error; err != nil {
		return nil, 0, fmt.Errorf("list IRSA tax brackets: %w", err)
	}

	return brackets, total, nil
}

// GetActiveIRSATaxBrackets retrieves all active IRSA tax brackets for a given date
func (r *Repo) GetActiveIRSATaxBrackets(ctx context.Context, effectiveDate time.Time) ([]IRSATaxBracket, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var brackets []IRSATaxBracket
	if err := r.db.WithContext(ctx).Where("is_active = ? AND effective_date <= ?", true, effectiveDate).
		Order("min_income ASC").Find(&brackets).Error; err != nil {
		return nil, fmt.Errorf("get active IRSA tax brackets: %w", err)
	}
	return brackets, nil
}

// generateDeclarationNumber generates a unique declaration number
func generateDeclarationNumber(declarationType string, periodStart time.Time) string {
	// Format: TYPE-YYYYMM-XXXX
	// e.g., CNAPS-202601-0001, OSTIE-202601-0001, IRSA-202601-0001
	year := periodStart.Year()
	month := int(periodStart.Month())
	return fmt.Sprintf("%s-%04d%02d-%04d", declarationType, year, month, 1000+int(periodStart.Unix()%9000))
}

// GenerateDeclarationForm generates a declaration form for CNAPS, OSTIE, or IRSA
func (r *Repo) GenerateDeclarationForm(ctx context.Context, declarationID uuid.UUID) (*DeclarationForm, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	declaration, err := r.GetDeclarationByID(ctx, declarationID)
	if err != nil {
		return nil, err
	}

	// Parse declaration data
	var employeeBreakdown []EmployeeBreakdown
	if declaration.DeclarationData != "" {
		if err := json.Unmarshal([]byte(declaration.DeclarationData), &employeeBreakdown); err != nil {
			// If parsing fails, continue with empty breakdown
			employeeBreakdown = []EmployeeBreakdown{}
		}
	}

	form := &DeclarationForm{
		DeclarationNumber:          declaration.DeclarationNumber,
		DeclarationType:            declaration.DeclarationType,
		DeclarationPeriodStart:     declaration.DeclarationPeriodStart,
		DeclarationPeriodEnd:       declaration.DeclarationPeriodEnd,
		CompanyName:                declaration.CompanyName,
		CompanyAddress:             declaration.CompanyAddress,
		CompanyNIF:                 declaration.CompanyNIF,
		TotalEmployees:             declaration.TotalEmployees,
		TotalGrossSalary:           declaration.TotalGrossSalary,
		TotalEmployeeContributions: declaration.TotalEmployeeContributions,
		TotalEmployerContributions: declaration.TotalEmployerContributions,
		TotalAmountDue:             declaration.TotalAmountDue,
		EmployeeBreakdown:          employeeBreakdown,
		Status:                     declaration.Status,
		AccountantName:             "Accountant Name", // TODO: Get from user service
		CreatedAt:                  declaration.CreatedAt,
		SubmittedAt:                declaration.SubmittedAt,
	}

	return form, nil
}
