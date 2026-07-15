package location

type Location struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Address   string   `json:"address,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	IsOpen    bool     `json:"is_open"`
	OpenTime  string   `json:"open_time"`
	CloseTime string   `json:"close_time"`
	CreatedAt string   `json:"created_at"`
}

type CreateRequest struct {
	Name      string   `json:"name" validate:"required,min=2,max=255"`
	Address   string   `json:"address" validate:"omitempty,max=1000"`
	Latitude  *float64 `json:"latitude" validate:"omitempty,gte=-90,lte=90"`
	Longitude *float64 `json:"longitude" validate:"omitempty,gte=-180,lte=180"`
	OpenTime  string   `json:"open_time" validate:"omitempty,len=5"`
	CloseTime string   `json:"close_time" validate:"omitempty,len=5"`
	IsOpen    *bool    `json:"is_open"`
}

type UpdateRequest struct {
	Name      *string  `json:"name" validate:"omitempty,min=2,max=255"`
	Address   *string  `json:"address" validate:"omitempty,max=1000"`
	Latitude  *float64 `json:"latitude" validate:"omitempty,gte=-90,lte=90"`
	Longitude *float64 `json:"longitude" validate:"omitempty,gte=-180,lte=180"`
	OpenTime  *string  `json:"open_time" validate:"omitempty,len=5"`
	CloseTime *string  `json:"close_time" validate:"omitempty,len=5"`
	IsOpen    *bool    `json:"is_open"`
}
