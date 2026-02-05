package notifications

import (
	"time"

	"github.com/google/uuid"
)

// Notification represents a notification record
type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Type      string    `gorm:"type:varchar(50);not null;check:type IN ('info', 'warning', 'success', 'error')" json:"type"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	Link      *string   `gorm:"type:varchar(255)" json:"link,omitempty"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
}

// NotificationListQuery represents query parameters for listing notifications
type NotificationListQuery struct {
	IsRead *bool  `form:"is_read"`
	Type   string `form:"type" binding:"omitempty,oneof=info warning success error"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int    `form:"offset" binding:"omitempty,min=0"`
}

// NotificationListResponse represents the response for listing notifications
type NotificationListResponse struct {
	Notifications []Notification `json:"notifications"`
	Total         int64          `json:"total"`
	UnreadCount   int            `json:"unread_count"`
}

// UnreadCountResponse represents the response for unread count
type UnreadCountResponse struct {
	UnreadCount int `json:"unread_count"`
}

// MarkAsReadResponse represents the response for marking as read
type MarkAsReadResponse struct {
	ID     uuid.UUID `json:"id"`
	IsRead bool      `json:"is_read"`
}

// MarkAllReadResponse represents the response for marking all as read
type MarkAllReadResponse struct {
	UpdatedCount int `json:"updated_count"`
}

// TableName specifies the table name for Notification model
func (Notification) TableName() string {
	return "notifications"
}
