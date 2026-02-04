package support

// Repository for CRUD operations on Support tickets

import (
	"gorm.io/gorm"
	"context"
	"time"
	"fmt"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(database *gorm.DB) *Repo {
	return &Repo{db: database}
}
func (r *Repo) Create(ctx context.Context, support *Support) (*Support, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := r.db.WithContext(ctx).Create(support).Error; err != nil {
		return nil, fmt.Errorf("create support: %w", err)
	}
	return support, nil
}

// func (r *Repo) CreateSupport(ctx context.Context, support *Support) (*Support, error) {
//     ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
//     defer cancel()

//     tx := r.db.WithContext(ctx).Begin()
//     if tx.Error != nil {
//         return nil, fmt.Errorf("begin tx: %w", tx.Error)
//     }

//     if err := tx.Create(support).Error; err != nil {
//         tx.Rollback()
//         return nil, fmt.Errorf("create support: %w", err)
//     }

//     if err := tx.Commit().Error; err != nil {
//         return nil, fmt.Errorf("commit tx: %w", err)
//     }

//     return support, nil
// }