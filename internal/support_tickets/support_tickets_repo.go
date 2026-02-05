package support_tickets

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for support tickets
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new support tickets repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// generateTicketNumber generates a unique ticket number
func (r *Repo) generateTicketNumber(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	year := time.Now().Year()

	// Get the count of tickets for this year
	var count int64
	if err := r.db.WithContext(ctx).Model(&SupportTicket{}).
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Count(&count).Error; err != nil {
		return "", fmt.Errorf("count tickets for year: %w", err)
	}

	// Generate ticket number: TICKET-YYYY-NNN (padded with leading zeros)
	ticketNumber := fmt.Sprintf("TICKET-%d-%03d", year, count+1)
	return ticketNumber, nil
}

// Create creates a new support ticket
func (r *Repo) Create(ctx context.Context, ticket *SupportTicket) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Generate ticket number
	ticketNumber, err := r.generateTicketNumber(ctx)
	if err != nil {
		return fmt.Errorf("generate ticket number: %w", err)
	}
	ticket.TicketNumber = ticketNumber

	if err := r.db.WithContext(ctx).Create(ticket).Error; err != nil {
		return fmt.Errorf("create support ticket: %w", err)
	}
	return nil
}

// GetByID retrieves a support ticket by ID with user info
func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*SupportTicket, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var ticket SupportTicket

	// Join with users table to get user info
	query := `
		SELECT 
			st.id, st.ticket_number, st.user_id, st.subject, st.description, 
			st.category, st.priority, st.status, st.assigned_to_id,
			st.created_at, st.updated_at, st.resolved_at,
			u.name as user_name, u.email as user_email,
			au.name as assigned_to_name
		FROM support_tickets st
		LEFT JOIN users u ON st.user_id = u.id
		LEFT JOIN users au ON st.assigned_to_id = au.id
		WHERE st.id = ?
	`

	if err := r.db.WithContext(ctx).Raw(query, id).Scan(&ticket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("support ticket not found")
		}
		return nil, fmt.Errorf("get support ticket: %w", err)
	}

	// Get replies count
	var repliesCount int64
	if err := r.db.WithContext(ctx).Model(&SupportTicketReply{}).
		Where("ticket_id = ?", id).
		Count(&repliesCount).Error; err != nil {
		return nil, fmt.Errorf("count ticket replies: %w", err)
	}
	ticket.RepliesCount = int(repliesCount)

	return &ticket, nil
}

// List retrieves support tickets with filtering
func (r *Repo) List(ctx context.Context, userID uuid.UUID, userRole string, query TicketListQuery) ([]SupportTicket, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var tickets []SupportTicket
	var total int64

	// Build base query
	db := r.db.WithContext(ctx).Model(&SupportTicket{})

	// Apply role-based filtering
	if userRole == "employee" {
		// Employees can only see their own tickets
		db = db.Where("user_id = ?", userID)
	}
	// Admin/HR can see all tickets, but can filter by user_id if provided

	// Apply filters
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}

	if query.Priority != "" {
		db = db.Where("priority = ?", query.Priority)
	}

	if query.UserID != nil && (userRole == "admin" || userRole == "hr") {
		db = db.Where("user_id = ?", *query.UserID)
	}

	if query.AssignedTo != nil && (userRole == "admin" || userRole == "hr") {
		db = db.Where("assigned_to_id = ?", *query.AssignedTo)
	}

	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(subject) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count support tickets: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	// Get tickets with user info
	queryStr := `
		SELECT 
			st.id, st.ticket_number, st.user_id, st.subject, st.description, 
			st.category, st.priority, st.status, st.assigned_to_id,
			st.created_at, st.updated_at, st.resolved_at,
			u.name as user_name, u.email as user_email,
			au.name as assigned_to_name,
			(SELECT COUNT(*) FROM support_ticket_replies str WHERE str.ticket_id = st.id) as replies_count
		FROM support_tickets st
		LEFT JOIN users u ON st.user_id = u.id
		LEFT JOIN users au ON st.assigned_to_id = au.id
		ORDER BY st.created_at DESC
		LIMIT ? OFFSET ?
	`

	if err := r.db.WithContext(ctx).Raw(queryStr, limit, query.Offset).Scan(&tickets).Error; err != nil {
		return nil, 0, fmt.Errorf("list support tickets: %w", err)
	}

	return tickets, total, nil
}

// GetTicketWithReplies retrieves a ticket with all its replies
func (r *Repo) GetTicketWithReplies(ctx context.Context, id uuid.UUID) (*TicketDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Get ticket
	ticket, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get replies
	var replies []SupportTicketReply
	query := `
		SELECT 
			str.id, str.ticket_id, str.user_id, str.message, 
			str.is_internal_note, str.created_at,
			u.name as user_name, u.role as user_role
		FROM support_ticket_replies str
		LEFT JOIN users u ON str.user_id = u.id
		WHERE str.ticket_id = ?
		ORDER BY str.created_at ASC
	`

	if err := r.db.WithContext(ctx).Raw(query, id).Scan(&replies).Error; err != nil {
		return nil, fmt.Errorf("get ticket replies: %w", err)
	}

	response := &TicketDetailResponse{
		SupportTicket: *ticket,
		Replies:       replies,
	}

	return response, nil
}

// Update updates a support ticket
func (r *Repo) Update(ctx context.Context, ticket *SupportTicket) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ticket.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(ticket).Error; err != nil {
		return fmt.Errorf("update support ticket: %w", err)
	}
	return nil
}

// AddReply adds a reply to a support ticket
func (r *Repo) AddReply(ctx context.Context, reply *SupportTicketReply) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(reply).Error; err != nil {
		return fmt.Errorf("add ticket reply: %w", err)
	}
	return nil
}

// Resolve resolves a support ticket
func (r *Repo) Resolve(ctx context.Context, id uuid.UUID, resolutionNote string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	now := time.Now()
	result := r.db.WithContext(ctx).Model(&SupportTicket{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "resolved",
			"resolved_at": &now,
			"updated_at":  now,
		})

	if result.Error != nil {
		return fmt.Errorf("resolve support ticket: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("support ticket not found")
	}

	// Add resolution note as a reply
	reply := &SupportTicketReply{
		TicketID:       id,
		UserID:         uuid.Nil, // This should be set by the handler
		Message:        resolutionNote,
		IsInternalNote: false,
	}

	if err := r.db.WithContext(ctx).Create(reply).Error; err != nil {
		return fmt.Errorf("add resolution note: %w", err)
	}

	return nil
}

// GetCategories returns the list of support categories
func (r *Repo) GetCategories(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	categories := []string{
		"technical",
		"payroll",
		"leave",
		"attendance",
		"account",
		"other",
	}

	return categories, nil
}
