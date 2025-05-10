package model

import (
	"time"
)

type EventLog struct {
	ID            string    `gorm:"primaryKey;column:id"`
	ApplicationID string    `gorm:"column:application_id;not null;index"`
	SessionID     string    `gorm:"column:session_id;not null;index"`
	EventID       string    `gorm:"column:event_id;not null;index"`
	PlatformID    int       `gorm:"column:platform_id"`
	Properties    JSONB     `gorm:"column:properties;type:jsonb"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

func (EventLog) TableName() string {
	return "tracking.event_logs"
}
