package employee

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Employee represents an employee in the system
type Employee struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID            uuid.UUID      `gorm:"type:uuid;not null" json:"company_id"`
	FirstName            string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName             string         `gorm:"type:varchar(100);not null" json:"last_name"`
	DateOfBirth          *time.Time     `gorm:"type:date" json:"date_of_birth,omitempty"`
	Gender               string         `gorm:"type:varchar(20);check:gender IN ('male', 'female', 'other')" json:"gender,omitempty"`
	Nationality          string         `gorm:"type:varchar(100);default:'Malagasy'" json:"nationality"`
	NationalID           string         `gorm:"type:varchar(50)" json:"national_id,omitempty"`
	Position             string         `gorm:"type:varchar(100)" json:"position,omitempty"`
	Department           string         `gorm:"type:varchar(100)" json:"department,omitempty"`
	HireDate             time.Time      `gorm:"type:date;not null" json:"hire_date"`
	ContractType         string         `gorm:"type:varchar(50);check:contract_type IN ('permanent', 'fixed_term', 'intern', 'contractor');default:'permanent'" json:"contract_type"`
	GrossSalary          float64        `gorm:"type:numeric(15,2);not null;check:gross_salary >= 200000" json:"gross_salary"`
	Status               string         `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'on_leave', 'terminated')" json:"status"`
	Address              string         `gorm:"type:text" json:"address,omitempty"`
	Phone                string         `gorm:"type:varchar(50)" json:"phone,omitempty"`
	EmergencyContactName string         `gorm:"type:varchar(100)" json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone string        `gorm:"type:varchar(50)" json:"emergency_contact_phone,omitempty"`
	ManagerID            *uuid.UUID     `gorm:"type:uuid" json:"manager_id,omitempty"`
	CreatedAt            time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateEmployeeRequest represents employee creation request
type CreateEmployeeRequest struct {
	CompanyID             uuid.UUID  `json:"company_id" binding:"required"`
	FirstName             string     `json:"first_name" binding:"required"`
	LastName              string     `json:"last_name" binding:"required"`
	DateOfBirth           *time.Time `json:"date_of_birth,omitempty"`
	Gender                string     `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Nationality           string     `json:"nationality,omitempty"`
	NationalID            string     `json:"national_id,omitempty"`
	Position              string     `json:"position,omitempty"`
	Department            string     `json:"department,omitempty"`
	HireDate              time.Time  `json:"hire_date" binding:"required"`
	ContractType          string     `json:"contract_type,omitempty" binding:"omitempty,oneof=permanent fixed_term intern contractor"`
	GrossSalary           float64    `json:"gross_salary" binding:"required,min=200000"`
	Address               string     `json:"address,omitempty"`
	Phone                 string     `json:"phone,omitempty"`
	EmergencyContactName  string     `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone string     `json:"emergency_contact_phone,omitempty"`
	ManagerID             *uuid.UUID `json:"manager_id,omitempty"`
}

// UpdateEmployeeRequest represents employee update request
type UpdateEmployeeRequest struct {
	FirstName             *string    `json:"first_name,omitempty"`
	LastName              *string    `json:"last_name,omitempty"`
	DateOfBirth           *time.Time `json:"date_of_birth,omitempty"`
	Gender                *string    `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Nationality           *string    `json:"nationality,omitempty"`
	NationalID            *string    `json:"national_id,omitempty"`
	Position              *string    `json:"position,omitempty"`
	Department            *string    `json:"department,omitempty"`
	ContractType          *string    `json:"contract_type,omitempty" binding:"omitempty,oneof=permanent fixed_term intern contractor"`
	GrossSalary           *float64   `json:"gross_salary,omitempty" binding:"omitempty,min=200000"`
	Status                *string    `json:"status,omitempty" binding:"omitempty,oneof=active on_leave terminated"`
	Address               *string    `json:"address,omitempty"`
	Phone                 *string    `json:"phone,omitempty"`
	EmergencyContactName  *string    `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string    `json:"emergency_contact_phone,omitempty"`
	ManagerID             *uuid.UUID `json:"manager_id,omitempty"`
}

// EmployeeListQuery represents query parameters for listing employees
type EmployeeListQuery struct {
	Search     string `form:"search"`
	Department string `form:"department"`
	Position   string `form:"position"`
	Status     string `form:"status" binding:"omitempty,oneof=active on_leave terminated"`
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset     int    `form:"offset" binding:"omitempty,min=0"`
}

// TableName specifies the table name for Employee model
func (Employee) TableName() string {
	return "employees"
}
