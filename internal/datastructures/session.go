package datastructure

type Session struct {
	ID            string  `json:"id"`
	ApplicationID string  `json:"application_id"`
	PlatformID    int     `json:"platform_id"`
	SessionKey    string  `json:"session_key"`
	UserID        *string `json:"user_id"`
	UserAgent     *string `json:"user_agent"`
	IPAddress     *string `json:"ip_address"`
	StartedAt     string  `json:"started_at"`
	EndedAt       *string `json:"ended_at"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     string  `json:"deleted_at"`
}

type CreateSessionRequest struct {
	ApplicationID string  `json:"application_id" example:"1231231123" binding:"omitempty"`
	PlatformID    int     `json:"platform_id" example:"1" binding:"required"`
	SessionKey    string  `json:"session_key" example:"1231231123" binding:"required"`
	UserID        *string `json:"user_id" example:"1231231123" binding:"omitempty"`
	UserAgent     *string `json:"user_agent" example:"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36" binding:"omitempty"`
	IPAddress     *string `json:"ip_address" example:"127.0.0.1" binding:"omitempty"`
	StartedAt     string  `json:"started_at" example:"2006-01-02 15:04:05" binding:"required,datetime_format"`
	EndedAt       *string `json:"ended_at" example:"2007-01-02 15:04:05" binding:"omitempty,datetime_format"`
}

type UpdateSessionRequest struct {
	UserID  string `json:"user_id" binding:"omitempty"`
	EndedAt string `json:"ended_at" binding:"omitempty,datetime_format"`
}
