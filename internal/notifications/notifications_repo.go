package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for notifications
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new notifications repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// Create creates a new notification
func (r *Repo) Create(ctx context.Context, notification *Notification) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

// GetByID retrieves a notification by ID
func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var notification Notification
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&notification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return &notification, nil
}

// List retrieves notifications with filtering
func (r *Repo) List(ctx context.Context, userID uuid.UUID, query NotificationListQuery) ([]Notification, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var notifications []Notification
	var total int64

	db := r.db.WithContext(ctx).Model(&Notification{}).Where("user_id = ?", userID)

	// Apply filters
	if query.IsRead != nil {
		db = db.Where("is_read = ?", *query.IsRead)
	}

	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count notifications: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 50
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, fmt.Errorf("list notifications: %w", err)
	}

	return notifications, total, nil
}

// GetUnreadCount retrieves the count of unread notifications for a user
func (r *Repo) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var count int64
	if err := r.db.WithContext(ctx).Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count unread notifications: %w", err)
	}

	return int(count), nil
}

// MarkAsRead marks a notification as read
func (r *Repo) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&Notification{}).
		Where("id = ?", id).
		Update("is_read", true)

	if result.Error != nil {
		return fmt.Errorf("mark notification as read: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (r *Repo) MarkAllAsRead(ctx context.Context, userID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	if result.Error != nil {
		return 0, fmt.Errorf("mark all notifications as read: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}

// Delete deletes a notification
func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&Notification{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("delete notification: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// CreateForUser creates a notification for a specific user
func (r *Repo) CreateForUser(ctx context.Context, userID uuid.UUID, title, message, notificationType string, link *string) error {
	notification := &Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    notificationType,
		IsRead:  false,
		Link:    link,
	}

	return r.Create(ctx, notification)
}

// CreateForMultipleUsers creates notifications for multiple users
func (r *Repo) CreateForMultipleUsers(ctx context.Context, userIDs []uuid.UUID, title, message, notificationType string, link *string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var notifications []Notification
	for _, userID := range userIDs {
		notifications = append(notifications, Notification{
			UserID:  userID,
			Title:   title,
			Message: message,
			Type:    notificationType,
			IsRead:  false,
			Link:    link,
		})
	}

	if err := r.db.WithContext(ctx).CreateInBatches(notifications, 100).Error; err != nil {
		return fmt.Errorf("create notifications for multiple users: %w", err)
	}

	return nil
}
