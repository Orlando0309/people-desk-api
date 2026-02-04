package support

import (
	"gorm.io/gorm"
)

type Support struct {
	gorm.Model
	Message string `gorm:"type:text;not null" json:"message"`
	Email   string `gorm:"type:varchar(100);not null" json:"email"`
	Status string `gorm:"type:varchar(50);not null;default:'unread'" json:"status"`
}

type CreateSupportInput struct {
	Message string `json:"message" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
}

type UpdateSupportInput struct {
	Status string `json:"status" binding:"required"`
}