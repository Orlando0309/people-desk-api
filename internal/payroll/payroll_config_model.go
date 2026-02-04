package payroll

import (
	"time"

	"github.com/google/uuid"
)

// PayrollConfiguration represents a configurable payroll parameter
type PayrollConfiguration struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Key         string     `gorm:"type:varchar(100);not null;uniqueIndex" json:"key"`
	Value       string     `gorm:"type:varchar(500);not null" json:"value"`
	Description string     `gorm:"type:text" json:"description"`
	DataType    string     `gorm:"type:varchar(20);not null;default:'string';check:data_type IN ('string', 'number', 'boolean')" json:"data_type"`
	Category    string     `gorm:"type:varchar(50);not null;default:'general'" json:"category"`
	IsActive    bool       `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt   time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:now()" json:"updated_at"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedBy   *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
}

// IRSATaxBracket represents an IRSA tax bracket
type IRSATaxBracket struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MinIncome     float64    `gorm:"type:numeric(15,2);not null;default:0;index" json:"min_income"`
	MaxIncome     *float64   `gorm:"type:numeric(15,2)" json:"max_income,omitempty"`
	TaxRate       float64    `gorm:"type:numeric(5,2);not null;default:0" json:"tax_rate"`
	MinTax        float64    `gorm:"type:numeric(15,2);not null;default:0" json:"min_tax"`
	BracketName   string     `gorm:"type:varchar(100);not null" json:"bracket_name"`
	IsActive      bool       `gorm:"not null;default:true;index" json:"is_active"`
	SortOrder     int        `gorm:"not null;default:0;index" json:"sort_order"`
	EffectiveDate time.Time  `gorm:"type:date;not null" json:"effective_date"`
	CreatedAt     time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"default:now()" json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedBy     *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
}

// CreatePayrollConfigurationRequest represents request to create a payroll configuration
type CreatePayrollConfigurationRequest struct {
	Key         string `json:"key" binding:"required"`
	Value       string `json:"value" binding:"required"`
	Description string `json:"description"`
	DataType    string `json:"data_type" binding:"required,oneof=string number boolean"`
	Category    string `json:"category" binding:"required"`
}

// UpdatePayrollConfigurationRequest represents request to update a payroll configuration
type UpdatePayrollConfigurationRequest struct {
	Value       *string `json:"value" binding:"omitempty"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

// CreateIRSATaxBracketRequest represents request to create an IRSA tax bracket
type CreateIRSATaxBracketRequest struct {
	MinIncome     float64   `json:"min_income" binding:"required,min=0"`
	MaxIncome     *float64  `json:"max_income"`
	TaxRate       float64   `json:"tax_rate" binding:"required,min=0,max=1"`
	MinTax        float64   `json:"min_tax" binding:"required,min=0"`
	BracketName   string    `json:"bracket_name" binding:"required"`
	SortOrder     int       `json:"sort_order"`
	EffectiveDate time.Time `json:"effective_date" binding:"required"`
}

// UpdateIRSATaxBracketRequest represents request to update an IRSA tax bracket
type UpdateIRSATaxBracketRequest struct {
	MinIncome     *float64   `json:"min_income" binding:"omitempty,min=0"`
	MaxIncome     *float64   `json:"max_income"`
	TaxRate       *float64   `json:"tax_rate" binding:"omitempty,min=0,max=1"`
	MinTax        *float64   `json:"min_tax" binding:"omitempty,min=0"`
	BracketName   *string    `json:"bracket_name"`
	SortOrder     *int       `json:"sort_order"`
	EffectiveDate *time.Time `json:"effective_date" binding:"omitempty"`
	IsActive      *bool      `json:"is_active"`
}

// PayrollConfigurationListQuery represents query parameters for listing payroll configurations
type PayrollConfigurationListQuery struct {
	Category string `form:"category"`
	Key      string `form:"key"`
	IsActive *bool  `form:"is_active"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset   int    `form:"offset" binding:"omitempty,min=0"`
}

// IRSATaxBracketListQuery represents query parameters for listing IRSA tax brackets
type IRSATaxBracketListQuery struct {
	IsActive *bool `form:"is_active"`
	Limit    int   `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset   int   `form:"offset" binding:"omitempty,min=0"`
}

// TableName specifies the table name for PayrollConfiguration model
func (PayrollConfiguration) TableName() string {
	return "payroll_configurations"
}

// TableName specifies the table name for IRSATaxBracket model
func (IRSATaxBracket) TableName() string {
	return "irsa_tax_brackets"
}
