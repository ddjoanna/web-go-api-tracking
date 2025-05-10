package model

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID          string         `gorm:"primaryKey;column:id"`
	Name        string         `gorm:"column:name;uniqueIndex"`
	Description string         `gorm:"column:description"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" sql:"index"`
}

func (Tenant) TableName() string {
	return "tracking.tenants"
}
