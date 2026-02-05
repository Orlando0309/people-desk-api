package company

import (
	"time"

	"github.com/google/uuid"
)

// CompanySettings represents company settings
type CompanySettings struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyName          string     `gorm:"type:varchar(255);not null" json:"company_name"`
	CompanyAddress       *string    `gorm:"type:text" json:"company_address,omitempty"`
	CompanyNIF           *string    `gorm:"type:varchar(50)" json:"company_nif,omitempty"`
	CompanySTAT          *string    `gorm:"type:varchar(50)" json:"company_stat,omitempty"`
	CNAPSNumber          *string    `gorm:"type:varchar(50)" json:"cnaps_number,omitempty"`
	OSTIENumber          *string    `gorm:"type:varchar(50)" json:"ostie_number,omitempty"`
	ContactEmail         *string    `gorm:"type:varchar(255)" json:"contact_email,omitempty"`
	ContactPhone         *string    `gorm:"type:varchar(50)" json:"contact_phone,omitempty"`
	LogoURL              *string    `gorm:"type:text" json:"logo_url,omitempty"`
	Timezone             string     `gorm:"type:varchar(100);default:'Indian/Antananarivo'" json:"timezone"`
	Currency             string     `gorm:"type:varchar(10);default:'MGA'" json:"currency"`
	FiscalYearStart      string     `gorm:"type:varchar(5);default:'01-01'" json:"fiscal_year_start"`
	WorkHoursPerDay      float64    `gorm:"type:decimal(4,2);default:8.00" json:"work_hours_per_day"`
	WorkDaysPerWeek      int        `gorm:"default:5" json:"work_days_per_week"`
	OvertimeWeekdayRate  float64    `gorm:"type:decimal(4,2);default:1.25" json:"overtime_weekday_rate"`
	OvertimeSaturdayRate float64    `gorm:"type:decimal(4,2);default:1.50" json:"overtime_saturday_rate"`
	OvertimeSundayRate   float64    `gorm:"type:decimal(4,2);default:2.00" json:"overtime_sunday_rate"`
	AnnualLeaveDays      int        `gorm:"default:30" json:"annual_leave_days"`
	MinimumSalary        float64    `gorm:"type:decimal(15,2);default:200000" json:"minimum_salary"`
	UpdatedAt            time.Time  `gorm:"default:now()" json:"updated_at"`
	UpdatedBy            *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
}

// UpdateCompanySettingsRequest represents the request body for updating company settings
type UpdateCompanySettingsRequest struct {
	CompanyName          *string  `json:"company_name,omitempty"`
	CompanyAddress       *string  `json:"company_address,omitempty"`
	CompanyNIF           *string  `json:"company_nif,omitempty"`
	CompanySTAT          *string  `json:"company_stat,omitempty"`
	CNAPSNumber          *string  `json:"cnaps_number,omitempty"`
	OSTIENumber          *string  `json:"ostie_number,omitempty"`
	ContactEmail         *string  `json:"contact_email,omitempty"`
	ContactPhone         *string  `json:"contact_phone,omitempty"`
	LogoURL              *string  `json:"logo_url,omitempty"`
	Timezone             *string  `json:"timezone,omitempty"`
	Currency             *string  `json:"currency,omitempty"`
	FiscalYearStart      *string  `json:"fiscal_year_start,omitempty"`
	WorkHoursPerDay      *float64 `json:"work_hours_per_day,omitempty"`
	WorkDaysPerWeek      *int     `json:"work_days_per_week,omitempty"`
	OvertimeWeekdayRate  *float64 `json:"overtime_weekday_rate,omitempty"`
	OvertimeSaturdayRate *float64 `json:"overtime_saturday_rate,omitempty"`
	OvertimeSundayRate   *float64 `json:"overtime_sunday_rate,omitempty"`
	AnnualLeaveDays      *int     `json:"annual_leave_days,omitempty"`
	MinimumSalary        *float64 `json:"minimum_salary,omitempty"`
}

// UploadLogoResponse represents the response for logo upload
type UploadLogoResponse struct {
	LogoURL string `json:"logo_url"`
}

// TableName specifies the table name for CompanySettings model
func (CompanySettings) TableName() string {
	return "company_settings"
}
