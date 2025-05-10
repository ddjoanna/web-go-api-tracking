package datastructure

type Application struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
}

type ApplicationAPIKey struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	APIKey        string `json:"api_key"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	DeletedAt     string `json:"deleted_at"`
}

type CreateApplicationRequest struct {
	TenantID    string `json:"tenant_id" example:"1231231123" binding:"required"`
	Name        string `json:"name" example:"My App" binding:"required"`
	Description string `json:"description" example:"My App Description" binding:"required"`
}

type UpdateApplicationRequest struct {
	TenantID    string `json:"tenant_id" example:"1231231123" binding:"required"`
	Name        string `json:"name" example:"My App" binding:"required"`
	Description string `json:"description" example:"My App Description" binding:"required"`
}
