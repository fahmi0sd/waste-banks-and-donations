package user

const (
	RoleUser        = "user"
	RoleAdmin       = "admin"
	RoleMasterAdmin = "master_admin"
)

type Profile struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone,omitempty"`
	Role       string `json:"role"`
	LocationID *int   `json:"location_id,omitempty"`
	CreatedAt  string `json:"created_at"`
}

type UpdateProfileRequest struct {
	Name       *string `json:"name" validate:"omitempty,min=2,max=255"`
	Phone      *string `json:"phone" validate:"omitempty,min=8,max=20"`
	LocationID *int    `json:"location_id" validate:"omitempty,gt=0"`
}
