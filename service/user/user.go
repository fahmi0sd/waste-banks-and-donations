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

// CreateAccountRequest dipakai master_admin untuk membuat akun user/admin baru.
type CreateAccountRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=255"`
	Email      string `json:"email" validate:"required,email,max=255"`
	Phone      string `json:"phone" validate:"omitempty,min=8,max=20"`
	Password   string `json:"password" validate:"required,min=8,max=72"`
	Role       string `json:"role" validate:"required,oneof=user admin master_admin"`
	LocationID *int   `json:"location_id" validate:"omitempty,gt=0"`
}

// UpdateAccountRequest dipakai master_admin untuk mengubah akun user/admin.
type UpdateAccountRequest struct {
	Name       *string `json:"name" validate:"omitempty,min=2,max=255"`
	Phone      *string `json:"phone" validate:"omitempty,min=8,max=20"`
	Role       *string `json:"role" validate:"omitempty,oneof=user admin master_admin"`
	LocationID *int    `json:"location_id" validate:"omitempty,gt=0"`
}
