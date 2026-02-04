package kpi

import (
	"time"

	"github.com/google/uuid"
)

// KPI represents a Key Performance Indicator template
type KPI struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name             string    `gorm:"type:varchar(255);not null" json:"name"`
	Description      string    `gorm:"type:text" json:"description"`
	TargetValue      float64   `gorm:"type:numeric(15,2);not null" json:"target_value"`
	WeightPercentage float64   `gorm:"type:numeric(5,2);not null" json:"weight_percentage"`
	ScoringScale     string    `gorm:"type:varchar(20);default:'1_to_5'" json:"scoring_scale"`
	Department       string    `gorm:"type:varchar(100)" json:"department"`
	Position         string    `gorm:"type:varchar(100)" json:"position"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
	CreatedBy        uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt        time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:now()" json:"updated_at"`
}

// PerformanceReview represents a performance review for an employee
type PerformanceReview struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EmployeeID        uuid.UUID `gorm:"type:uuid;not null" json:"employee_id"`
	KPIID             uuid.UUID `gorm:"type:uuid;not null" json:"kpi_id"`
	ReviewPeriodStart time.Time `gorm:"type:date;not null" json:"review_period_start"`
	ReviewPeriodEnd   time.Time `gorm:"type:date;not null" json:"review_period_end"`
	SelfScore         *float64  `gorm:"type:numeric(5,2)" json:"self_score"`
	ManagerScore      *float64  `gorm:"type:numeric(5,2)" json:"manager_score"`
	FinalScore        *float64  `gorm:"type:numeric(5,2)" json:"final_score"`
	SelfAssessment    string    `gorm:"type:text" json:"self_assessment"`
	ManagerAssessment string    `gorm:"type:text" json:"manager_assessment"`
	ReviewerID        uuid.UUID `gorm:"type:uuid;not null" json:"reviewer_id"`
	Status            string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt         time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:now()" json:"updated_at"`
}

// CreateKPIRequest represents request to create a KPI
type CreateKPIRequest struct {
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	TargetValue      float64 `json:"target_value" binding:"required"`
	WeightPercentage float64 `json:"weight_percentage" binding:"required,min=0.01,max=100"`
	ScoringScale     string  `json:"scoring_scale" binding:"omitempty,oneof=1_to_5 1_to_10 custom"`
	Department       string  `json:"department"`
	Position         string  `json:"position"`
}

// UpdateKPIRequest represents request to update a KPI
type UpdateKPIRequest struct {
	Name             *string  `json:"name"`
	Description      *string  `json:"description"`
	TargetValue      *float64 `json:"target_value" binding:"omitempty,min=0"`
	WeightPercentage *float64 `json:"weight_percentage" binding:"omitempty,min=0.01,max=100"`
	ScoringScale     *string  `json:"scoring_scale" binding:"omitempty,oneof=1_to_5 1_to_10 custom"`
	Department       *string  `json:"department"`
	Position         *string  `json:"position"`
	IsActive         *bool    `json:"is_active"`
}

// KPIListQuery represents query parameters for listing KPIs
type KPIListQuery struct {
	Department string `form:"department"`
	Position   string `form:"position"`
	IsActive   *bool  `form:"is_active"`
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset     int    `form:"offset" binding:"omitempty,min=0"`
}

// CreatePerformanceReviewRequest represents request to create a performance review
type CreatePerformanceReviewRequest struct {
	EmployeeID        uuid.UUID `json:"employee_id" binding:"required"`
	KPIID             uuid.UUID `json:"kpi_id" binding:"required"`
	ReviewPeriodStart time.Time `json:"review_period_start" binding:"required"`
	ReviewPeriodEnd   time.Time `json:"review_period_end" binding:"required"`
	SelfScore         *float64  `json:"self_score" binding:"omitempty,min=0,max=10"`
	ManagerScore      *float64  `json:"manager_score" binding:"omitempty,min=0,max=10"`
	SelfAssessment    string    `json:"self_assessment"`
	ManagerAssessment string    `json:"manager_assessment"`
}

// UpdatePerformanceReviewRequest represents request to update a performance review
type UpdatePerformanceReviewRequest struct {
	SelfScore         *float64 `json:"self_score" binding:"omitempty,min=0,max=10"`
	ManagerScore      *float64 `json:"manager_score" binding:"omitempty,min=0,max=10"`
	FinalScore        *float64 `json:"final_score" binding:"omitempty,min=0,max=10"`
	SelfAssessment    *string  `json:"self_assessment"`
	ManagerAssessment *string  `json:"manager_assessment"`
	Status            *string  `json:"status" binding:"omitempty,oneof=pending completed approved"`
}

// PerformanceReviewListQuery represents query parameters for listing performance reviews
type PerformanceReviewListQuery struct {
	EmployeeID        *uuid.UUID `form:"employee_id"`
	KPIID             *uuid.UUID `form:"kpi_id"`
	ReviewerID        *uuid.UUID `form:"reviewer_id"`
	Status            string     `form:"status" binding:"omitempty,oneof=pending completed approved"`
	ReviewPeriodStart *time.Time `form:"review_period_start"`
	ReviewPeriodEnd   *time.Time `form:"review_period_end"`
	Limit             int        `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset            int        `form:"offset" binding:"omitempty,min=0"`
}

// PerformanceReport represents a performance report for an employee
type PerformanceReport struct {
	EmployeeID        uuid.UUID       `json:"employee_id"`
	EmployeeName      string          `json:"employee_name"`
	ReviewPeriodStart time.Time       `json:"review_period_start"`
	ReviewPeriodEnd   time.Time       `json:"review_period_end"`
	KPIs              []KPIReportItem `json:"kpis"`
	OverallScore      float64         `json:"overall_score"`
	TotalWeight       float64         `json:"total_weight"`
}

// KPIReportItem represents a KPI item in a performance report
type KPIReportItem struct {
	KPIName           string   `json:"kpi_name"`
	TargetValue       float64  `json:"target_value"`
	WeightPercentage  float64  `json:"weight_percentage"`
	SelfScore         *float64 `json:"self_score"`
	ManagerScore      *float64 `json:"manager_score"`
	FinalScore        *float64 `json:"final_score"`
	SelfAssessment    string   `json:"self_assessment"`
	ManagerAssessment string   `json:"manager_assessment"`
	Status            string   `json:"status"`
}

// TableName specifies the table name for KPI model
func (KPI) TableName() string {
	return "kpis"
}

// TableName specifies the table name for PerformanceReview model
func (PerformanceReview) TableName() string {
	return "performance_reviews"
}
