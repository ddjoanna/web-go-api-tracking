package model

import (
	"time"

	"gorm.io/gorm"
)

type Application struct {
	ID          string `gorm:"primaryKey"`
	TenantID    string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time           `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time           `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt      `gorm:"column:deleted_at" sql:"index"`
	ApiKeys     []ApplicationApiKey `gorm:"foreignKey:ApplicationID;references:ID"`
}

func (Application) TableName() string {
	return "tracking.applications"
}
