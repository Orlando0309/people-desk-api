package payroll

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PayrollDraft represents a payroll draft created by HR
type PayrollDraft struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PeriodStart   time.Time      `gorm:"type:date;not null" json:"period_start"`
	PeriodEnd     time.Time      `gorm:"type:date;not null" json:"period_end"`
	EmployeeID    uuid.UUID      `gorm:"type:uuid;not null" json:"employee_id"`
	GrossSalary   float64        `gorm:"type:numeric(15,2);not null" json:"gross_salary"`
	CNAPSEmployee float64        `gorm:"type:numeric(15,2);not null" json:"cnaps_employee"`
	CNAPSEmployer float64        `gorm:"type:numeric(15,2);not null" json:"cnaps_employer"`
	OSTIEEmployee float64        `gorm:"type:numeric(15,2);not null" json:"ostie_employee"`
	OSTIEEmployer float64        `gorm:"type:numeric(15,2);not null" json:"ostie_employer"`
	IRSA          float64        `gorm:"type:numeric(15,2);not null" json:"irsa"`
	NetSalary     float64        `gorm:"type:numeric(15,2);not null" json:"net_salary"`
	CNAPSBase     float64        `gorm:"type:numeric(15,2);not null" json:"cnaps_base"`
	OSTIEBase     float64        `gorm:"type:numeric(15,2);not null" json:"ostie_base"`
	IRSABracket   string         `gorm:"type:varchar(50)" json:"irsa_bracket"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt     time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// PayrollApproved represents an approved payroll record by Accountant
type PayrollApproved struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DraftID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"draft_id"`
	FichePaieNumber  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"fiche_paie_number"`
	AccountantID     uuid.UUID `gorm:"type:uuid;not null" json:"accountant_id"`
	GLEntries        string    `gorm:"type:jsonb;not null" json:"gl_entries"`
	ApprovedAt       time.Time `gorm:"default:now()" json:"approved_at"`
	DigitalSignature string    `gorm:"type:text;not null" json:"digital_signature"`
	CreatedAt        time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:now()" json:"updated_at"`
}

// CreatePayrollDraftRequest represents request to create a payroll draft
type CreatePayrollDraftRequest struct {
	PeriodStart time.Time `json:"period_start" binding:"required"`
	PeriodEnd   time.Time `json:"period_end" binding:"required"`
	EmployeeID  uuid.UUID `json:"employee_id" binding:"required"`
	GrossSalary float64   `json:"gross_salary" binding:"required,min=200000"`
}

// UpdatePayrollDraftRequest represents request to update a payroll draft
type UpdatePayrollDraftRequest struct {
	GrossSalary *float64 `json:"gross_salary" binding:"omitempty,min=200000"`
}

// ApprovePayrollDraftRequest represents request to approve a payroll draft
type ApprovePayrollDraftRequest struct {
	Comment string `json:"comment"`
}

// PayrollDraftListQuery represents query parameters for listing payroll drafts
type PayrollDraftListQuery struct {
	PeriodStart *time.Time `form:"period_start"`
	PeriodEnd   *time.Time `form:"period_end"`
	EmployeeID  *uuid.UUID `form:"employee_id"`
	Status      string     `form:"status" binding:"omitempty,oneof=draft approved"`
	Limit       int        `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset      int        `form:"offset" binding:"omitempty,min=0"`
}

// PayrollApprovedListQuery represents query parameters for listing approved payrolls
type PayrollApprovedListQuery struct {
	PeriodStart     *time.Time `form:"period_start"`
	PeriodEnd       *time.Time `form:"period_end"`
	EmployeeID      *uuid.UUID `form:"employee_id"`
	FichePaieNumber string     `form:"fiche_paie_number"`
	AccountantID    *uuid.UUID `form:"accountant_id"`
	Limit           int        `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset          int        `form:"offset" binding:"omitempty,min=0"`
}

// FichePaie represents a payslip (fiche de paie)
type FichePaie struct {
	FichePaieNumber    string    `json:"fiche_paie_number"`
	EmployeeID         uuid.UUID `json:"employee_id"`
	EmployeeName       string    `json:"employee_name"`
	EmployeePosition   string    `json:"employee_position"`
	EmployeeDepartment string    `json:"employee_department"`
	PeriodStart        time.Time `json:"period_start"`
	PeriodEnd          time.Time `json:"period_end"`
	GrossSalary        float64   `json:"gross_salary"`
	CNAPSEmployee      float64   `json:"cnaps_employee"`
	CNAPSEmployer      float64   `json:"cnaps_employer"`
	OSTIEEmployee      float64   `json:"ostie_employee"`
	OSTIEEmployer      float64   `json:"ostie_employer"`
	IRSA               float64   `json:"irsa"`
	IRSABracket        string    `json:"irsa_bracket"`
	NetSalary          float64   `json:"net_salary"`
	AccountantName     string    `json:"accountant_name"`
	ApprovedAt         time.Time `json:"approved_at"`
	DigitalSignature   string    `json:"digital_signature"`
}

// GLEntry represents a general ledger entry (OHADA compliant)
type GLEntry struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Description string  `json:"description"`
}

// ReconciliationReport represents payroll reconciliation report
type ReconciliationReport struct {
	PeriodStart   time.Time `json:"period_start"`
	PeriodEnd     time.Time `json:"period_end"`
	HRDraftTotals struct {
		GrossSalary   float64 `json:"gross_salary"`
		CNAPSEmployee float64 `json:"cnaps_employee"`
		CNAPSEmployer float64 `json:"cnaps_employer"`
		OSTIEEmployee float64 `json:"ostie_employee"`
		OSTIEEmployer float64 `json:"ostie_employer"`
		IRSAWithheld  float64 `json:"irsa_withheld"`
		NetPayable    float64 `json:"net_payable"`
	} `json:"hr_draft_totals"`
	AccountantApprovedTotals struct {
		GrossSalary   float64 `json:"gross_salary"`
		CNAPSEmployee float64 `json:"cnaps_employee"`
		CNAPSEmployer float64 `json:"cnaps_employer"`
		OSTIEEmployee float64 `json:"ostie_employee"`
		OSTIEEmployer float64 `json:"ostie_employer"`
		IRSAWithheld  float64 `json:"irsa_withheld"`
		NetPayable    float64 `json:"net_payable"`
	} `json:"accountant_approved_totals"`
	GLRecordedAmounts struct {
		Account431CNAPS float64 `json:"account_431_cnaps"`
		Account438OSTIE float64 `json:"account_438_ostie"`
		Account437IRSA  float64 `json:"account_437_irsa"`
	} `json:"gl_recorded_amounts"`
	VariancePercentage float64 `json:"variance_percentage"`
	Status             string  `json:"status"`
}

// TableName specifies the table name for PayrollDraft model
func (PayrollDraft) TableName() string {
	return "payroll_drafts"
}

// TableName specifies the table name for PayrollApproved model
func (PayrollApproved) TableName() string {
	return "payroll_approved"
}
