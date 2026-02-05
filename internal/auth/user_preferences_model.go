package auth

import (
	"time"

	"github.com/google/uuid"
)

// UserPreferences represents user preferences
type UserPreferences struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID             uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	EmailNotifications bool      `gorm:"default:true" json:"email_notifications"`
	PushNotifications  bool      `gorm:"default:false" json:"push_notifications"`
	LeaveUpdates       bool      `gorm:"default:true" json:"leave_updates"`
	PayrollUpdates     bool      `gorm:"default:true" json:"payroll_updates"`
	SystemUpdates      bool      `gorm:"default:false" json:"system_updates"`
	Theme              string    `gorm:"type:varchar(20);default:'light'" json:"theme"`
	Language           string    `gorm:"type:varchar(10);default:'en'" json:"language"`
	DateFormat         string    `gorm:"type:varchar(20);default:'DD/MM/YYYY'" json:"date_format"`
	UpdatedAt          time.Time `gorm:"default:now()" json:"updated_at"`
}

// UpdatePreferencesRequest represents the request body for updating preferences
type UpdatePreferencesRequest struct {
	EmailNotifications *bool   `json:"email_notifications,omitempty"`
	PushNotifications  *bool   `json:"push_notifications,omitempty"`
	LeaveUpdates       *bool   `json:"leave_updates,omitempty"`
	PayrollUpdates     *bool   `json:"payroll_updates,omitempty"`
	SystemUpdates      *bool   `json:"system_updates,omitempty"`
	Theme              *string `json:"theme,omitempty"`
	Language           *string `json:"language,omitempty"`
	DateFormat         *string `json:"date_format,omitempty"`
}

// PreferencesResponse represents the response for user preferences
type PreferencesResponse struct {
	Notifications struct {
		EmailEnabled   bool `json:"email_enabled"`
		PushEnabled    bool `json:"push_enabled"`
		LeaveUpdates   bool `json:"leave_updates"`
		PayrollUpdates bool `json:"payroll_updates"`
		SystemUpdates  bool `json:"system_updates"`
	} `json:"notifications"`
	Display struct {
		Theme      string `json:"theme"`
		Language   string `json:"language"`
		DateFormat string `json:"date_format"`
	} `json:"display"`
}

// TableName specifies the table name for UserPreferences model
func (UserPreferences) TableName() string {
	return "user_preferences"
}
