package datastructure

type Event struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	PlatformID    int    `json:"platform_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	IsActive      bool   `json:"is_active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	DeletedAt     string `json:"deleted_at"`
}

type EventField struct {
	ID          string `json:"id"`
	EventID     string `json:"event_id"`
	Name        string `json:"name"`
	DataType    string `json:"data_type"`
	IsRequired  bool   `json:"is_required"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
}

type EventLog struct {
	ID            string                 `json:"id"`
	ApplicationID string                 `json:"application_id"`
	SessionID     string                 `json:"session_id"`
	EventID       string                 `json:"event_id"`
	PlatformID    int                    `json:"platform_id"`
	Properties    map[string]interface{} `json:"properties"`
	CreatedAt     string                 `json:"created_at"`
}

type EventResponse struct {
	Event
	Fields []EventField `json:"fields"`
}

type CreateEventRequest struct {
	ApplicationID string `json:"application_id" example:"1231231123" binding:"omitempty"`
	PlatformID    int    `json:"platform_id" example:"1" binding:"required"`
	Name          string `json:"name" example:"click_button" binding:"required"`
	Description   string `json:"description" example:"Click on CTA" binding:"required"`
	IsActive      bool   `json:"is_active" example:"true" binding:"required"`
}

type UpdateEventRequest struct {
	ApplicationID string `json:"application_id" example:"1231231123" binding:"omitempty"`
	PlatformID    int    `json:"platform_id" example:"1" binding:"required"`
	Name          string `json:"name" example:"click_button" binding:"required"`
	Description   string `json:"description" example:"Click on CTA" binding:"required"`
	IsActive      bool   `json:"is_active" example:"true" binding:"required"`
}

type CreateEventFieldRequest struct {
	EventID     string `json:"event_id" example:"1231231123" binding:"omitempty"`
	Name        string `json:"name" example:"button_id" binding:"required"`
	DataType    string `json:"data_type" example:"string" binding:"required,oneof=string int float boolean datetime json"`
	IsRequired  bool   `json:"is_required" example:"true"`
	Description string `json:"description" example:"Button ID" binding:"omitempty"`
}

type UpdateEventFieldRequest struct {
	EventID     string `json:"event_id" example:"1231231123" binding:"omitempty"`
	Name        string `json:"name" example:"button_id" binding:"required"`
	DataType    string `json:"data_type" example:"string" binding:"required,oneof=string int float boolean datetime json"`
	IsRequired  bool   `json:"is_required" example:"true"`
	Description string `json:"description" example:"Button ID" binding:"omitempty"`
}

type CreateEventLogRequest struct {
	ApplicationID string                 `json:"application_id" example:"1231231123" binding:"omitempty"`
	PlatformID    int                    `json:"platform_id" example:"1" binding:"required"`
	EventID       string                 `json:"event_id" example:"1231231123" binding:"omitempty"`
	SessionID     string                 `json:"session_id" example:"1231231123" binding:"required"`
	Properties    map[string]interface{} `json:"properties" binding:"required"`
}
