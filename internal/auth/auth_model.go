package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	ID                    uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email                 string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash          string         `gorm:"type:text;not null" json:"-"`
	Role                  string         `gorm:"type:varchar(20);not null;check:role IN ('admin', 'hr', 'accountant', 'employee')" json:"role"`
	EmployeeID            *uuid.UUID     `gorm:"type:uuid" json:"employee_id,omitempty"`
	Phone                 *string        `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Address               *string        `gorm:"type:text" json:"address,omitempty"`
	EmergencyContactName  *string        `gorm:"type:varchar(255)" json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string        `gorm:"type:varchar(50)" json:"emergency_contact_phone,omitempty"`
	IsActive              bool           `gorm:"default:true" json:"is_active"`
	FailedLoginAttempts   int            `gorm:"default:0" json:"-"`
	LockedUntil           *time.Time     `json:"-"`
	CreatedAt             time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"default:now()" json:"updated_at"`
	LastLogin             *time.Time     `json:"last_login,omitempty"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse represents login response with tokens
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID                    uuid.UUID  `json:"id"`
	Email                 string     `json:"email"`
	Role                  string     `json:"role"`
	EmployeeID            *uuid.UUID `json:"employee_id,omitempty"`
	Phone                 *string    `json:"phone,omitempty"`
	Address               *string    `json:"address,omitempty"`
	EmergencyContactName  *string    `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string    `json:"emergency_contact_phone,omitempty"`
	IsActive              bool       `json:"is_active"`
	LastLogin             *time.Time `json:"last_login,omitempty"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty"`
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email      string     `json:"email" binding:"required,email"`
	Password   string     `json:"password" binding:"required,min=8"`
	Role       string     `json:"role" binding:"required,oneof=admin hr accountant employee"`
	EmployeeID *uuid.UUID `json:"employee_id,omitempty"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Email                 *string `json:"email,omitempty"`
	Phone                 *string `json:"phone,omitempty"`
	Address               *string `json:"address,omitempty"`
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
