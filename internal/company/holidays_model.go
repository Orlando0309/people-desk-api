package company

import (
	"time"

	"github.com/google/uuid"
)

// CompanyHoliday represents a company holiday
type CompanyHoliday struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Date        string     `gorm:"type:date;not null" json:"date"`
	IsRecurring bool       `gorm:"default:false" json:"is_recurring"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy   *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
}

// HolidaysListQuery represents query parameters for listing holidays
type HolidaysListQuery struct {
	Year      int    `form:"year"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

// CreateHolidayRequest represents the request body for creating a holiday
type CreateHolidayRequest struct {
	Name        string `json:"name" binding:"required"`
	Date        string `json:"date" binding:"required"`
	IsRecurring bool   `json:"is_recurring"`
	Description string `json:"description,omitempty"`
}

// UpdateHolidayRequest represents the request body for updating a holiday
type UpdateHolidayRequest struct {
	Name        *string `json:"name,omitempty"`
	Date        *string `json:"date,omitempty"`
	IsRecurring *bool   `json:"is_recurring,omitempty"`
	Description *string `json:"description,omitempty"`
}

// HolidaysListResponse represents the response for listing holidays
type HolidaysListResponse struct {
	Holidays []CompanyHoliday `json:"holidays"`
}

// TableName specifies the table name for CompanyHoliday model
func (CompanyHoliday) TableName() string {
	return "company_holidays"
}
