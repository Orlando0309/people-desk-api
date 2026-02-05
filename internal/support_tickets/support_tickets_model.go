package support_tickets

import (
	"time"

	"github.com/google/uuid"
)

// SupportTicket represents a support ticket
type SupportTicket struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TicketNumber   string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"ticket_number"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	UserName       string     `gorm:"-" json:"user_name"`
	UserEmail      string     `gorm:"-" json:"user_email"`
	Subject        string     `gorm:"type:varchar(255);not null" json:"subject"`
	Description    string     `gorm:"type:text;not null" json:"description"`
	Category       string     `gorm:"type:varchar(50);not null;check:category IN ('technical', 'payroll', 'leave', 'attendance', 'account', 'other')" json:"category"`
	Priority       string     `gorm:"type:varchar(20);default:'medium';check:priority IN ('low', 'medium', 'high', 'urgent')" json:"priority"`
	Status         string     `gorm:"type:varchar(20);default:'open';check:status IN ('open', 'in_progress', 'resolved', 'closed')" json:"status"`
	AssignedToID   *uuid.UUID `gorm:"type:uuid" json:"assigned_to_id,omitempty"`
	AssignedToName string     `gorm:"-" json:"assigned_to_name,omitempty"`
	CreatedAt      time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"default:now()" json:"updated_at"`
	ResolvedAt     *time.Time `gorm:"type:timestamptz" json:"resolved_at,omitempty"`
	RepliesCount   int        `gorm:"-" json:"replies_count"`
}

// SupportTicketReply represents a reply to a support ticket
type SupportTicketReply struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TicketID       uuid.UUID `gorm:"type:uuid;not null" json:"ticket_id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	UserName       string    `gorm:"-" json:"user_name"`
	UserRole       string    `gorm:"-" json:"user_role"`
	Message        string    `gorm:"type:text;not null" json:"message"`
	IsInternalNote bool      `gorm:"default:false" json:"is_internal_note"`
	CreatedAt      time.Time `gorm:"default:now()" json:"created_at"`
}

// TicketListQuery represents query parameters for listing tickets
type TicketListQuery struct {
	Status     string     `form:"status" binding:"omitempty,oneof=open in_progress resolved closed"`
	Category   string     `form:"category" binding:"omitempty,oneof=technical payroll leave attendance account other"`
	Priority   string     `form:"priority" binding:"omitempty,oneof=low medium high urgent"`
	UserID     *uuid.UUID `form:"user_id"`
	AssignedTo *uuid.UUID `form:"assigned_to"`
	Search     string     `form:"search"`
	Limit      int        `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset     int        `form:"offset" binding:"omitempty,min=0"`
}

// TicketListResponse represents the response for listing tickets
type TicketListResponse struct {
	Tickets []SupportTicket `json:"tickets"`
	Total   int64           `json:"total"`
}

// TicketDetailResponse represents the response for ticket details
type TicketDetailResponse struct {
	SupportTicket
	Replies []SupportTicketReply `json:"replies"`
}

// CreateTicketRequest represents the request body for creating a ticket
type CreateTicketRequest struct {
	Subject     string `json:"subject" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required,oneof=technical payroll leave attendance account other"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
}

// UpdateTicketRequest represents the request body for updating a ticket
type UpdateTicketRequest struct {
	Status       string     `json:"status" binding:"omitempty,oneof=open in_progress resolved closed"`
	Priority     string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	AssignedToID *uuid.UUID `json:"assigned_to_id"`
	Category     string     `json:"category" binding:"omitempty,oneof=technical payroll leave attendance account other"`
}

// ReplyTicketRequest represents the request body for replying to a ticket
type ReplyTicketRequest struct {
	Message        string `json:"message" binding:"required"`
	IsInternalNote bool   `json:"is_internal_note"`
}

// ResolveTicketRequest represents the request body for resolving a ticket
type ResolveTicketRequest struct {
	ResolutionNote string `json:"resolution_note" binding:"required"`
}

// CategoriesResponse represents the response for support categories
type CategoriesResponse struct {
	Categories []string `json:"categories"`
}

// TableName specifies the table name for SupportTicket model
func (SupportTicket) TableName() string {
	return "support_tickets"
}

// TableName specifies the table name for SupportTicketReply model
func (SupportTicketReply) TableName() string {
	return "support_ticket_replies"
}
