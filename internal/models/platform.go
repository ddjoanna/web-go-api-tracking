package model

import (
	"time"

	"gorm.io/gorm"
)

type Platform struct {
	ID        int            `gorm:"primaryKey;column:id"`
	Name      string         `gorm:"column:name;uniqueIndex"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" sql:"index"`
}

func (Platform) TableName() string {
	return "tracking.platforms"
}
