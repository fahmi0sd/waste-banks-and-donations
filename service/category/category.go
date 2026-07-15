package category

type Category struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ParentCategoryID *int   `json:"parent_category_id,omitempty"`
	Unit             string `json:"unit"`
	IsActive         bool   `json:"is_active"`
}

type CreateRequest struct {
	Name             string `json:"name" validate:"required,min=2,max=100"`
	ParentCategoryID *int   `json:"parent_category_id" validate:"omitempty,gt=0"`
	Unit             string `json:"unit" validate:"omitempty,oneof=kg liter"`
	IsActive         *bool  `json:"is_active"`
}

type UpdateRequest struct {
	Name             *string `json:"name" validate:"omitempty,min=2,max=100"`
	ParentCategoryID *int    `json:"parent_category_id" validate:"omitempty,gt=0"`
	Unit             *string `json:"unit" validate:"omitempty,oneof=kg liter"`
	IsActive         *bool   `json:"is_active"`
}
