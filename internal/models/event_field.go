package model

import (
	"time"

	"gorm.io/gorm"
)

type EventField struct {
	ID          string         `gorm:"primaryKey;column:id"`
	EventID     string         `gorm:"column:event_id;not null;index"`
	Name        string         `gorm:"column:name"`
	DataType    string         `gorm:"column:data_type"`
	IsRequired  bool           `gorm:"column:is_required;default:false"`
	Description string         `gorm:"column:description"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" sql:"index"`
}

func (EventField) TableName() string {
	return "tracking.event_fields"
}
