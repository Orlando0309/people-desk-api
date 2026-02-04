package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for audit logs
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new audit repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// Create creates a new audit log entry (immutable)
func (r *Repo) Create(ctx context.Context, log *AuditLog) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

// GetByID retrieves an audit log by ID
func (r *Repo) GetByID(ctx context.Context, id int64) (*AuditLog, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var log AuditLog
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("audit log not found")
		}
		return nil, fmt.Errorf("get audit log: %w", err)
	}
	return &log, nil
}

// List retrieves audit logs with filtering and pagination
func (r *Repo) List(ctx context.Context, query AuditLogListQuery) ([]AuditLog, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var logs []AuditLog
	var total int64

	db := r.db.WithContext(ctx).Model(&AuditLog{})

	// Apply filters
	if query.UserID != nil {
		db = db.Where("user_id = ?", *query.UserID)
	}

	if query.UserRole != "" {
		db = db.Where("user_role = ?", query.UserRole)
	}

	if query.ActionType != "" {
		db = db.Where("action_type = ?", query.ActionType)
	}

	if query.Module != "" {
		db = db.Where("module = ?", query.Module)
	}

	if query.StartDate != nil {
		db = db.Where("created_at >= ?", *query.StartDate)
	}

	if query.EndDate != nil {
		db = db.Where("created_at <= ?", *query.EndDate)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}

	// Apply pagination
	limit := query.Limit
	if limit == 0 {
		limit = 100
	}

	if err := db.Limit(limit).Offset(query.Offset).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}

	return logs, total, nil
}

// GetByUserID retrieves all audit logs for a specific user
func (r *Repo) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]AuditLog, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var logs []AuditLog
	var total int64

	if err := r.db.WithContext(ctx).Model(&AuditLog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count user audit logs: %w", err)
	}

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("get user audit logs: %w", err)
	}

	return logs, total, nil
}

// GetByModule retrieves all audit logs for a specific module
func (r *Repo) GetByModule(ctx context.Context, module string, limit, offset int) ([]AuditLog, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var logs []AuditLog
	var total int64

	if err := r.db.WithContext(ctx).Model(&AuditLog{}).Where("module = ?", module).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count module audit logs: %w", err)
	}

	if err := r.db.WithContext(ctx).Where("module = ?", module).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("get module audit logs: %w", err)
	}

	return logs, total, nil
}

// GetStats retrieves audit statistics
func (r *Repo) GetStats(ctx context.Context, startDate, endDate time.Time) (*AuditStats, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stats := &AuditStats{
		ActionsByType:   make(map[string]int64),
		ActionsByModule: make(map[string]int64),
		ActionsByUser:   make(map[string]int64),
	}

	// Count total actions
	if err := r.db.WithContext(ctx).Model(&AuditLog{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Count(&stats.TotalActions).Error; err != nil {
		return nil, fmt.Errorf("count total actions: %w", err)
	}

	// Count actions by type
	var actionTypes []struct {
		ActionType string
		Count      int64
	}
	if err := r.db.WithContext(ctx).Model(&AuditLog{}).
		Select("action_type, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("action_type").Scan(&actionTypes).Error; err != nil {
		return nil, fmt.Errorf("count actions by type: %w", err)
	}
	for _, at := range actionTypes {
		stats.ActionsByType[at.ActionType] = at.Count
	}

	// Count actions by module
	var modules []struct {
		Module string
		Count  int64
	}
	if err := r.db.WithContext(ctx).Model(&AuditLog{}).
		Select("module, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("module").Scan(&modules).Error; err != nil {
		return nil, fmt.Errorf("count actions by module: %w", err)
	}
	for _, m := range modules {
		stats.ActionsByModule[m.Module] = m.Count
	}

	// Count actions by user
	var users []struct {
		UserID string
		Count  int64
	}
	if err := r.db.WithContext(ctx).Model(&AuditLog{}).
		Select("user_id::text as user_id, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("user_id").Scan(&users).Error; err != nil {
		return nil, fmt.Errorf("count actions by user: %w", err)
	}
	for _, u := range users {
		stats.ActionsByUser[u.UserID] = u.Count
	}

	return stats, nil
}

// GetRecordHistory retrieves all audit logs for a specific record
func (r *Repo) GetRecordHistory(ctx context.Context, recordID uuid.UUID) ([]AuditLog, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var logs []AuditLog
	if err := r.db.WithContext(ctx).Where("record_id = ?", recordID).
		Order("created_at ASC").Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("get record history: %w", err)
	}
	return logs, nil
}
