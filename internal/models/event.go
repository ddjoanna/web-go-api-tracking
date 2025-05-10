package model

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID            string         `gorm:"primaryKey;column:id"`
	ApplicationID string         `gorm:"column:application_id;not null;index"`
	PlatformID    int            `gorm:"column:platform_id"`
	Name          string         `gorm:"column:name"`
	Description   string         `gorm:"column:description"`
	IsActive      bool           `gorm:"column:is_active;default:true"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" sql:"index"`
	Application   *Application   `gorm:"foreignKey:ApplicationID"`
	Platform      *Platform      `gorm:"foreignKey:PlatformID"`
	Fields        []*EventField  `gorm:"foreignKey:EventID"`
}

func (Event) TableName() string {
	return "tracking.events"
}
