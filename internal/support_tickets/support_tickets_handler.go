package support_tickets

import (
	"net/http"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles support ticket requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new support tickets handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// List retrieves support tickets
func (h *Handler) List(c *gin.Context) {
	var query TicketListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	tickets, total, err := h.repo.List(c.Request.Context(), userID, userRole, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list support tickets"})
		return
	}

	response := TicketListResponse{
		Tickets: tickets,
		Total:   total,
	}

	c.JSON(http.StatusOK, response)
}

// GetByID retrieves a support ticket with replies
func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	// Get user info from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	ticket, err := h.repo.GetTicketWithReplies(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Support ticket not found"})
		return
	}

	// Check access permissions
	if userRole == "employee" && ticket.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot access another user's ticket"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// Create creates a new support ticket
func (h *Handler) Create(c *gin.Context) {
	var input CreateTicketRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Set default priority if not provided
	priority := input.Priority
	if priority == "" {
		priority = "medium"
	}

	ticket := &SupportTicket{
		UserID:      userID,
		Subject:     input.Subject,
		Description: input.Description,
		Category:    input.Category,
		Priority:    priority,
		Status:      "open",
	}

	if err := h.repo.Create(c.Request.Context(), ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create support ticket"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":            ticket.ID,
		"ticket_number": ticket.TicketNumber,
		"user_id":       ticket.UserID,
		"subject":       ticket.Subject,
		"status":        ticket.Status,
		"created_at":    ticket.CreatedAt,
	})
}

// Update updates a support ticket
func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var input UpdateTicketRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	_, err = middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Only Admin/HR can update tickets
	if userRole != "admin" && userRole != "hr" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	// Get existing ticket
	ticket, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Support ticket not found"})
		return
	}

	// Update fields if provided
	if input.Status != "" {
		ticket.Status = input.Status
	}
	if input.Priority != "" {
		ticket.Priority = input.Priority
	}
	if input.AssignedToID != nil {
		ticket.AssignedToID = input.AssignedToID
	}
	if input.Category != "" {
		ticket.Category = input.Category
	}

	if err := h.repo.Update(c.Request.Context(), ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update support ticket"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// Reply adds a reply to a support ticket
func (h *Handler) Reply(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var input ReplyTicketRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Get ticket to check permissions
	ticket, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Support ticket not found"})
		return
	}

	// Check access permissions
	if userRole == "employee" && ticket.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot reply to another user's ticket"})
		return
	}

	// Only Admin/HR can add internal notes
	if input.IsInternalNote && userRole != "admin" && userRole != "hr" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot add internal notes"})
		return
	}

	reply := &SupportTicketReply{
		TicketID:       id,
		UserID:         userID,
		Message:        input.Message,
		IsInternalNote: input.IsInternalNote,
	}

	if err := h.repo.AddReply(c.Request.Context(), reply); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reply"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         reply.ID,
		"ticket_id":  reply.TicketID,
		"user_id":    reply.UserID,
		"message":    reply.Message,
		"created_at": reply.CreatedAt,
	})
}

// Resolve resolves a support ticket
func (h *Handler) Resolve(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var input ResolveTicketRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	_, err = middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Only Admin/HR can resolve tickets
	if userRole != "admin" && userRole != "hr" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	if err := h.repo.Resolve(c.Request.Context(), id, input.ResolutionNote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket resolved successfully"})
}

// GetCategories returns the list of support categories
func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}

	response := CategoriesResponse{
		Categories: categories,
	}

	c.JSON(http.StatusOK, response)
}
