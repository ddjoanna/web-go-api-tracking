package model

import (
	"time"

	"gorm.io/gorm"
)

type ApplicationApiKey struct {
	ID            string         `gorm:"primaryKey"`
	ApplicationID string         `gorm:"not null"`
	APIKey        string         `gorm:"unique;not null"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" sql:"index"`
}

func (ApplicationApiKey) TableName() string {
	return "tracking.applications_api_keys"
}
