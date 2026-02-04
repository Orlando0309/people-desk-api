package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repo handles database operations for authentication
type Repo struct {
	db *gorm.DB
}

// NewRepo creates a new auth repository
func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}

// CreateUser creates a new user
func (r *Repo) CreateUser(ctx context.Context, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by email
func (r *Repo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user User
	if err := r.db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *Repo) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

// UpdateUser updates user information
func (r *Repo) UpdateUser(ctx context.Context, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// UpdateLastLogin updates the last login timestamp
func (r *Repo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"last_login":            now,
		"failed_login_attempts": 0,
		"locked_until":          nil,
	}).Error; err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	return nil
}

// IncrementFailedLoginAttempts increments failed login attempts
func (r *Repo) IncrementFailedLoginAttempts(ctx context.Context, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	user.FailedLoginAttempts++

	// Lock account after 5 failed attempts for 15 minutes
	if user.FailedLoginAttempts >= 5 {
		lockUntil := time.Now().Add(15 * time.Minute)
		user.LockedUntil = &lockUntil
	}

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return fmt.Errorf("increment failed login attempts: %w", err)
	}
	return nil
}

// IsAccountLocked checks if account is locked
func (r *Repo) IsAccountLocked(ctx context.Context, userID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return false, fmt.Errorf("get user: %w", err)
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return true, nil
	}

	// If lock period has expired, reset it
	if user.LockedUntil != nil && user.LockedUntil.Before(time.Now()) {
		user.LockedUntil = nil
		user.FailedLoginAttempts = 0
		if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
			return false, fmt.Errorf("reset lock: %w", err)
		}
	}

	return false, nil
}

// ListUsers retrieves all users (admin only)
func (r *Repo) ListUsers(ctx context.Context, limit, offset int) ([]User, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var users []User
	var total int64

	if err := r.db.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	return users, total, nil
}

// DeleteUser soft deletes a user
func (r *Repo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Delete(&User{}, id).Error; err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
