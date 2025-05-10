package datastructure

type Platform struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type CreatePlatformRequest struct {
	Name string `json:"name" example:"Web" binding:"required"`
}
