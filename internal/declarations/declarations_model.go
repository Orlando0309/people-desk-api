package declarations

import (
	"time"

	"github.com/google/uuid"
)

// MonthlyDeclaration represents a monthly declaration for CNAPS, OSTIE, or IRSA
type MonthlyDeclaration struct {
	ID                         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DeclarationType            string     `gorm:"type:varchar(20);not null" json:"declaration_type"`
	DeclarationPeriodStart     time.Time  `gorm:"type:date;not null" json:"declaration_period_start"`
	DeclarationPeriodEnd       time.Time  `gorm:"type:date;not null" json:"declaration_period_end"`
	DeclarationNumber          string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"declaration_number"`
	CompanyName                string     `gorm:"type:varchar(255);not null" json:"company_name"`
	CompanyAddress             string     `gorm:"type:text" json:"company_address"`
	CompanyNIF                 string     `gorm:"type:varchar(50)" json:"company_nif"`
	TotalEmployees             int        `gorm:"not null" json:"total_employees"`
	TotalGrossSalary           float64    `gorm:"type:numeric(15,2);not null" json:"total_gross_salary"`
	TotalEmployeeContributions float64    `gorm:"type:numeric(15,2);not null" json:"total_employee_contributions"`
	TotalEmployerContributions float64    `gorm:"type:numeric(15,2);not null" json:"total_employer_contributions"`
	TotalAmountDue             float64    `gorm:"type:numeric(15,2);not null" json:"total_amount_due"`
	DeclarationData            string     `gorm:"type:jsonb;not null" json:"declaration_data"`
	Status                     string     `gorm:"type:varchar(20);default:'draft'" json:"status"`
	AccountantID               uuid.UUID  `gorm:"type:uuid;not null" json:"accountant_id"`
	SubmittedAt                *time.Time `json:"submitted_at"`
	PaidAt                     *time.Time `json:"paid_at"`
	CreatedAt                  time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt                  time.Time  `gorm:"default:now()" json:"updated_at"`
}

// IRSATaxBracket represents an IRSA tax bracket configuration
type IRSATaxBracket struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MinIncome     float64   `gorm:"type:numeric(15,2);not null" json:"min_income"`
	MaxIncome     *float64  `gorm:"type:numeric(15,2)" json:"max_income"`
	TaxRate       float64   `gorm:"type:numeric(5,2);not null" json:"tax_rate"`
	MinTax        float64   `gorm:"type:numeric(15,2);default:0" json:"min_tax"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	EffectiveDate time.Time `gorm:"type:date;not null" json:"effective_date"`
	CreatedBy     uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt     time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:now()" json:"updated_at"`
}

// CreateDeclarationRequest represents request to create a monthly declaration
type CreateDeclarationRequest struct {
	DeclarationType        string    `json:"declaration_type" binding:"required,oneof=cnaps ostie irsa"`
	DeclarationPeriodStart time.Time `json:"declaration_period_start" binding:"required"`
	DeclarationPeriodEnd   time.Time `json:"declaration_period_end" binding:"required"`
	CompanyName            string    `json:"company_name" binding:"required"`
	CompanyAddress         string    `json:"company_address"`
	CompanyNIF             string    `json:"company_nif"`
}

// UpdateDeclarationRequest represents request to update a monthly declaration
type UpdateDeclarationRequest struct {
	Status *string `json:"status" binding:"omitempty,oneof=draft submitted paid cancelled"`
}

// DeclarationListQuery represents query parameters for listing declarations
type DeclarationListQuery struct {
	DeclarationType        string     `form:"declaration_type" binding:"omitempty,oneof=cnaps ostie irsa"`
	DeclarationPeriodStart *time.Time `form:"declaration_period_start"`
	DeclarationPeriodEnd   *time.Time `form:"declaration_period_end"`
	Status                 string     `form:"status" binding:"omitempty,oneof=draft submitted paid cancelled"`
	AccountantID           *uuid.UUID `form:"accountant_id"`
	Limit                  int        `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset                 int        `form:"offset" binding:"omitempty,min=0"`
}

// CreateIRSATaxBracketRequest represents request to create an IRSA tax bracket
type CreateIRSATaxBracketRequest struct {
	MinIncome     float64   `json:"min_income" binding:"required,min=0"`
	MaxIncome     *float64  `json:"max_income" binding:"omitempty,min=0"`
	TaxRate       float64   `json:"tax_rate" binding:"required,min=0,max=100"`
	MinTax        float64   `json:"min_tax" binding:"omitempty,min=0"`
	EffectiveDate time.Time `json:"effective_date" binding:"required"`
}

// UpdateIRSATaxBracketRequest represents request to update an IRSA tax bracket
type UpdateIRSATaxBracketRequest struct {
	MinIncome     *float64   `json:"min_income" binding:"omitempty,min=0"`
	MaxIncome     *float64   `json:"max_income" binding:"omitempty,min=0"`
	TaxRate       *float64   `json:"tax_rate" binding:"omitempty,min=0,max=100"`
	MinTax        *float64   `json:"min_tax" binding:"omitempty,min=0"`
	IsActive      *bool      `json:"is_active"`
	EffectiveDate *time.Time `json:"effective_date"`
}

// IRSATaxBracketListQuery represents query parameters for listing IRSA tax brackets
type IRSATaxBracketListQuery struct {
	IsActive      *bool      `form:"is_active"`
	EffectiveDate *time.Time `form:"effective_date"`
	Limit         int        `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset        int        `form:"offset" binding:"omitempty,min=0"`
}

// DeclarationForm represents a declaration form for CNAPS, OSTIE, or IRSA
type DeclarationForm struct {
	DeclarationNumber          string              `json:"declaration_number"`
	DeclarationType            string              `json:"declaration_type"`
	DeclarationPeriodStart     time.Time           `json:"declaration_period_start"`
	DeclarationPeriodEnd       time.Time           `json:"declaration_period_end"`
	CompanyName                string              `json:"company_name"`
	CompanyAddress             string              `json:"company_address"`
	CompanyNIF                 string              `json:"company_nif"`
	TotalEmployees             int                 `json:"total_employees"`
	TotalGrossSalary           float64             `json:"total_gross_salary"`
	TotalEmployeeContributions float64             `json:"total_employee_contributions"`
	TotalEmployerContributions float64             `json:"total_employer_contributions"`
	TotalAmountDue             float64             `json:"total_amount_due"`
	EmployeeBreakdown          []EmployeeBreakdown `json:"employee_breakdown"`
	Status                     string              `json:"status"`
	AccountantName             string              `json:"accountant_name"`
	CreatedAt                  time.Time           `json:"created_at"`
	SubmittedAt                *time.Time          `json:"submitted_at"`
}

// EmployeeBreakdown represents employee contribution breakdown in a declaration
type EmployeeBreakdown struct {
	EmployeeID           uuid.UUID `json:"employee_id"`
	EmployeeName         string    `json:"employee_name"`
	EmployeeNIF          string    `json:"employee_nif"`
	GrossSalary          float64   `json:"gross_salary"`
	BaseAmount           float64   `json:"base_amount"`
	EmployeeContribution float64   `json:"employee_contribution"`
	EmployerContribution float64   `json:"employer_contribution"`
}

// TableName specifies the table name for MonthlyDeclaration model
func (MonthlyDeclaration) TableName() string {
	return "monthly_declarations"
}

// TableName specifies the table name for IRSATaxBracket model
func (IRSATaxBracket) TableName() string {
	return "irsa_tax_brackets"
}
