package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID            string         `gorm:"primaryKey;column:id"`
	ApplicationID string         `gorm:"column:application_id;not null;index"`
	PlatformID    int            `gorm:"column:platform_id;not null;index"`
	SessionKey    string         `gorm:"column:session_key;uniqueIndex"`
	UserID        *string        `gorm:"column:user_id"`
	UserAgent     *string        `gorm:"column:user_agent"`
	IPAddress     *string        `gorm:"column:ip_address"`
	StartedAt     time.Time      `gorm:"column:started_at"`
	EndedAt       *time.Time     `gorm:"column:ended_at"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" sql:"index"`
}

func (Session) TableName() string {
	return "tracking.sessions"
}
