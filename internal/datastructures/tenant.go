package datastructure

type Tenant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
}

type CreateTenantRequest struct {
	Name        string `json:"name" example:"My Shop" binding:"required"`
	Description string `json:"description" example:"Shopping Description" binding:"required"`
}

type UpdateTenantRequest struct {
	Name        string `json:"name" example:"My Shop" binding:"required"`
	Description string `json:"description" example:"Shopping Description" binding:"required"`
}
