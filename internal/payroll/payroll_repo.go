package payroll

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for payroll
type Repo struct {
	db         *gorm.DB
	configRepo *ConfigRepo
}

// NewRepo creates a new payroll repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{
		db:         database,
		configRepo: NewConfigRepo(database),
	}
}

// calculateCNAPS calculates CNAPS contributions using configurable rates from database
func (r *Repo) calculateCNAPS(ctx context.Context, grossSalary float64) (base, employee, employer float64, err error) {
	// Get configuration values from database
	ceiling, err := r.configRepo.GetConfigValueAsFloat(ctx, "cnaps_ostie_ceiling")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get CNAPS ceiling: %w", err)
	}

	employeeRate, err := r.configRepo.GetConfigValueAsFloat(ctx, "cnaps_employee_rate")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get CNAPS employee rate: %w", err)
	}

	employerRate, err := r.configRepo.GetConfigValueAsFloat(ctx, "cnaps_employer_rate")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get CNAPS employer rate: %w", err)
	}

	base = grossSalary
	if base > ceiling {
		base = ceiling
	}
	employee = base * employeeRate
	employer = base * employerRate
	return
}

// calculateOSTIE calculates OSTIE contributions using configurable rates from database
func (r *Repo) calculateOSTIE(ctx context.Context, grossSalary float64) (base, employee, employer float64, err error) {
	// Get configuration values from database
	ceiling, err := r.configRepo.GetConfigValueAsFloat(ctx, "cnaps_ostie_ceiling")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get CNAPS/OSTIE ceiling: %w", err)
	}

	employeeRate, err := r.configRepo.GetConfigValueAsFloat(ctx, "ostie_employee_rate")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get OSTIE employee rate: %w", err)
	}

	employerRate, err := r.configRepo.GetConfigValueAsFloat(ctx, "ostie_employer_rate")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get OSTIE employer rate: %w", err)
	}

	base = grossSalary
	if base > ceiling {
		base = ceiling
	}
	employee = base * employeeRate
	employer = base * employerRate
	return
}

// calculateIRSA calculates IRSA withholding based on configurable tax brackets from database
func (r *Repo) calculateIRSA(ctx context.Context, grossSalary, cnapsEmployee, ostieEmployee float64) (amount float64, bracket string, err error) {
	// Get IRSA tax brackets from database
	brackets, err := r.configRepo.GetActiveIRSABrackets(ctx)
	if err != nil {
		return 0, "", fmt.Errorf("get IRSA brackets: %w", err)
	}

	// Taxable income = gross - CNAPS employee - OSTIE employee
	taxableIncome := grossSalary - cnapsEmployee - ostieEmployee

	// Find the applicable bracket
	for _, b := range brackets {
		// Check if income falls within this bracket
		inRange := taxableIncome >= b.MinIncome
		if b.MaxIncome != nil {
			inRange = inRange && taxableIncome <= *b.MaxIncome
		}

		if inRange {
			amount = taxableIncome*b.TaxRate + b.MinTax
			bracket = b.BracketName
			return amount, bracket, nil
		}
	}

	// If no bracket found, return error
	return 0, "", fmt.Errorf("no applicable IRSA bracket found for income: %.2f", taxableIncome)
}

// CreateDraft creates a new payroll draft
func (r *Repo) CreateDraft(ctx context.Context, draft *PayrollDraft) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Check if draft already exists for this employee and period
	var existing PayrollDraft
	err := r.db.WithContext(ctx).Where("employee_id = ? AND period_start = ? AND period_end = ? AND deleted_at IS NULL",
		draft.EmployeeID, draft.PeriodStart, draft.PeriodEnd).First(&existing).Error
	if err == nil {
		return fmt.Errorf("payroll draft already exists for this employee and period")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("check existing draft: %w", err)
	}

	// Calculate CNAPS
	cnapsBase, cnapsEmployee, cnapsEmployer, err := r.calculateCNAPS(ctx, draft.GrossSalary)
	if err != nil {
		return fmt.Errorf("calculate CNAPS: %w", err)
	}

	// Calculate OSTIE
	ostieBase, ostieEmployee, ostieEmployer, err := r.calculateOSTIE(ctx, draft.GrossSalary)
	if err != nil {
		return fmt.Errorf("calculate OSTIE: %w", err)
	}

	// Calculate IRSA
	irsa, irsaBracket, err := r.calculateIRSA(ctx, draft.GrossSalary, cnapsEmployee, ostieEmployee)
	if err != nil {
		return fmt.Errorf("calculate IRSA: %w", err)
	}

	// Calculate net salary
	netSalary := draft.GrossSalary - cnapsEmployee - ostieEmployee - irsa

	// Set calculated values
	draft.CNAPSEmployee = cnapsEmployee
	draft.CNAPSEmployer = cnapsEmployer
	draft.CNAPSBase = cnapsBase
	draft.OSTIEEmployee = ostieEmployee
	draft.OSTIEEmployer = ostieEmployer
	draft.OSTIEBase = ostieBase
	draft.IRSA = irsa
	draft.IRSABracket = irsaBracket
	draft.NetSalary = netSalary

	if err := r.db.WithContext(ctx).Create(draft).Error; err != nil {
		return fmt.Errorf("create payroll draft: %w", err)
	}
	return nil
}

// GetDraftByID retrieves a payroll draft by ID
func (r *Repo) GetDraftByID(ctx context.Context, id uuid.UUID) (*PayrollDraft, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var draft PayrollDraft
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&draft).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payroll draft not found")
		}
		return nil, fmt.Errorf("get payroll draft: %w", err)
	}
	return &draft, nil
}

// UpdateDraft updates a payroll draft
func (r *Repo) UpdateDraft(ctx context.Context, draft *PayrollDraft) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Recalculate if gross salary changed
	if draft.GrossSalary > 0 {
		cnapsBase, cnapsEmployee, cnapsEmployer, err := r.calculateCNAPS(ctx, draft.GrossSalary)
		if err != nil {
			return fmt.Errorf("calculate CNAPS: %w", err)
		}

		ostieBase, ostieEmployee, ostieEmployer, err := r.calculateOSTIE(ctx, draft.GrossSalary)
		if err != nil {
			return fmt.Errorf("calculate OSTIE: %w", err)
		}

		irsa, irsaBracket, err := r.calculateIRSA(ctx, draft.GrossSalary, cnapsEmployee, ostieEmployee)
		if err != nil {
			return fmt.Errorf("calculate IRSA: %w", err)
		}

		netSalary := draft.GrossSalary - cnapsEmployee - ostieEmployee - irsa

		draft.CNAPSEmployee = cnapsEmployee
		draft.CNAPSEmployer = cnapsEmployer
		draft.CNAPSBase = cnapsBase
		draft.OSTIEEmployee = ostieEmployee
		draft.OSTIEEmployer = ostieEmployer
		draft.OSTIEBase = ostieBase
		draft.IRSA = irsa
		draft.IRSABracket = irsaBracket
		draft.NetSalary = netSalary
	}

	draft.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(draft).Error; err != nil {
		return fmt.Errorf("update payroll draft: %w", err)
	}
	return nil
}

// DeleteDraft soft deletes a payroll draft
func (r *Repo) DeleteDraft(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&PayrollDraft{})
	if result.Error != nil {
		return fmt.Errorf("delete payroll draft: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("payroll draft not found")
	}
	return nil
}

// ListDrafts retrieves payroll drafts with filtering
func (r *Repo) ListDrafts(ctx context.Context, query PayrollDraftListQuery) ([]PayrollDraft, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var drafts []PayrollDraft
	var total int64

	db := r.db.WithContext(ctx).Model(&PayrollDraft{}).Where("deleted_at IS NULL")

	// Apply filters
	if query.PeriodStart != nil {
		db = db.Where("period_start >= ?", *query.PeriodStart)
	}
	if query.PeriodEnd != nil {
		db = db.Where("period_end <= ?", *query.PeriodEnd)
	}
	if query.EmployeeID != nil {
		db = db.Where("employee_id = ?", *query.EmployeeID)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count payroll drafts: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("period_start DESC, created_at DESC").Find(&drafts).Error; err != nil {
		return nil, 0, fmt.Errorf("list payroll drafts: %w", err)
	}

	return drafts, total, nil
}

// GetDraftsByPeriod retrieves all drafts for a specific period
func (r *Repo) GetDraftsByPeriod(ctx context.Context, periodStart, periodEnd time.Time) ([]PayrollDraft, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var drafts []PayrollDraft
	if err := r.db.WithContext(ctx).Where("period_start = ? AND period_end = ? AND deleted_at IS NULL",
		periodStart, periodEnd).Find(&drafts).Error; err != nil {
		return nil, fmt.Errorf("get drafts by period: %w", err)
	}
	return drafts, nil
}

// CreateApproved creates a new approved payroll record
func (r *Repo) CreateApproved(ctx context.Context, approved *PayrollApproved) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Generate fiche paie number if not provided
	if approved.FichePaieNumber == "" {
		approved.FichePaieNumber = generateFichePaieNumber(approved.ApprovedAt)
	}

	if err := r.db.WithContext(ctx).Create(approved).Error; err != nil {
		return fmt.Errorf("create approved payroll: %w", err)
	}
	return nil
}

// GetApprovedByID retrieves an approved payroll by ID
func (r *Repo) GetApprovedByID(ctx context.Context, id uuid.UUID) (*PayrollApproved, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var approved PayrollApproved
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&approved).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("approved payroll not found")
		}
		return nil, fmt.Errorf("get approved payroll: %w", err)
	}
	return &approved, nil
}

// GetApprovedByFichePaieNumber retrieves an approved payroll by fiche paie number
func (r *Repo) GetApprovedByFichePaieNumber(ctx context.Context, fichePaieNumber string) (*PayrollApproved, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var approved PayrollApproved
	if err := r.db.WithContext(ctx).Where("fiche_paie_number = ?", fichePaieNumber).First(&approved).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("approved payroll not found")
		}
		return nil, fmt.Errorf("get approved payroll: %w", err)
	}
	return &approved, nil
}

// ListApproved retrieves approved payrolls with filtering
func (r *Repo) ListApproved(ctx context.Context, query PayrollApprovedListQuery) ([]PayrollApproved, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var approved []PayrollApproved
	var total int64

	db := r.db.WithContext(ctx).Model(&PayrollApproved{})

	// Apply filters
	if query.PeriodStart != nil {
		db = db.Where("created_at >= ?", *query.PeriodStart)
	}
	if query.PeriodEnd != nil {
		db = db.Where("created_at <= ?", *query.PeriodEnd)
	}
	if query.EmployeeID != nil {
		db = db.Joins("JOIN payroll_drafts ON payroll_approved.draft_id = payroll_drafts.id").
			Where("payroll_drafts.employee_id = ?", *query.EmployeeID)
	}
	if query.FichePaieNumber != "" {
		db = db.Where("fiche_paie_number LIKE ?", "%"+query.FichePaieNumber+"%")
	}
	if query.AccountantID != nil {
		db = db.Where("accountant_id = ?", *query.AccountantID)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count approved payrolls: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("approved_at DESC").Find(&approved).Error; err != nil {
		return nil, 0, fmt.Errorf("list approved payrolls: %w", err)
	}

	return approved, total, nil
}

// generateFichePaieNumber generates a unique fiche de paie number
func generateFichePaieNumber(approvedAt time.Time) string {
	// Format: FDPAIE-YYYYMMDD-Random(4)
	// This is a simple implementation. In production, you'd want to ensure uniqueness
	year := approvedAt.Year()
	month := int(approvedAt.Month())
	day := approvedAt.Day()
	return fmt.Sprintf("FDPAIE-%04d%02d%02d-%04d", year, month, day, 1000+int(approvedAt.Unix()%9000))
}

// GetReconciliationReport generates a reconciliation report for a period
func (r *Repo) GetReconciliationReport(ctx context.Context, periodStart, periodEnd time.Time) (*ReconciliationReport, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	report := &ReconciliationReport{
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
	}

	// Get all drafts for the period
	drafts, err := r.GetDraftsByPeriod(ctx, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("get drafts for reconciliation: %w", err)
	}

	// Calculate HR draft totals
	for _, draft := range drafts {
		report.HRDraftTotals.GrossSalary += draft.GrossSalary
		report.HRDraftTotals.CNAPSEmployee += draft.CNAPSEmployee
		report.HRDraftTotals.CNAPSEmployer += draft.CNAPSEmployer
		report.HRDraftTotals.OSTIEEmployee += draft.OSTIEEmployee
		report.HRDraftTotals.OSTIEEmployer += draft.OSTIEEmployer
		report.HRDraftTotals.IRSAWithheld += draft.IRSA
		report.HRDraftTotals.NetPayable += draft.NetSalary
	}

	// Get approved payrolls for the period
	var approvedPayrolls []PayrollApproved
	if err := r.db.WithContext(ctx).
		Joins("JOIN payroll_drafts ON payroll_approved.draft_id = payroll_drafts.id").
		Where("payroll_drafts.period_start = ? AND payroll_drafts.period_end = ?", periodStart, periodEnd).
		Find(&approvedPayrolls).Error; err != nil {
		return nil, fmt.Errorf("get approved payrolls for reconciliation: %w", err)
	}

	// Get draft IDs for approved payrolls
	draftIDs := make([]uuid.UUID, len(approvedPayrolls))
	for i, approved := range approvedPayrolls {
		draftIDs[i] = approved.DraftID
	}

	// Get approved drafts
	var approvedDrafts []PayrollDraft
	if len(draftIDs) > 0 {
		if err := r.db.WithContext(ctx).Where("id IN ?", draftIDs).Find(&approvedDrafts).Error; err != nil {
			return nil, fmt.Errorf("get approved drafts for reconciliation: %w", err)
		}
	}

	// Calculate accountant approved totals
	for _, draft := range approvedDrafts {
		report.AccountantApprovedTotals.GrossSalary += draft.GrossSalary
		report.AccountantApprovedTotals.CNAPSEmployee += draft.CNAPSEmployee
		report.AccountantApprovedTotals.CNAPSEmployer += draft.CNAPSEmployer
		report.AccountantApprovedTotals.OSTIEEmployee += draft.OSTIEEmployee
		report.AccountantApprovedTotals.OSTIEEmployer += draft.OSTIEEmployer
		report.AccountantApprovedTotals.IRSAWithheld += draft.IRSA
		report.AccountantApprovedTotals.NetPayable += draft.NetSalary
	}

	// Calculate GL recorded amounts from approved payrolls
	report.GLRecordedAmounts.Account431CNAPS = report.AccountantApprovedTotals.CNAPSEmployee + report.AccountantApprovedTotals.CNAPSEmployer
	report.GLRecordedAmounts.Account438OSTIE = report.AccountantApprovedTotals.OSTIEEmployee + report.AccountantApprovedTotals.OSTIEEmployer
	report.GLRecordedAmounts.Account437IRSA = report.AccountantApprovedTotals.IRSAWithheld

	// Calculate variance
	expectedTotal := report.HRDraftTotals.GrossSalary
	actualTotal := report.AccountantApprovedTotals.GrossSalary
	if expectedTotal > 0 {
		report.VariancePercentage = ((actualTotal - expectedTotal) / expectedTotal) * 100
	}

	// Determine status
	if report.VariancePercentage <= 0.1 && report.VariancePercentage >= -0.1 {
		report.Status = "RECONCILED"
	} else {
		report.Status = "VARIANCE DETECTED"
	}

	return report, nil
}
