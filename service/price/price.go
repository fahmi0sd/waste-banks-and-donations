package price

type Price struct {
	ID             int     `json:"id"`
	CategoryID     int     `json:"category_id"`
	PricePerUnit   float64 `json:"price_per_unit"`
	EffectiveFrom  string  `json:"effective_from"`
	EffectiveUntil *string `json:"effective_until,omitempty"`
	SetBy          *int    `json:"set_by,omitempty"`
}

type SetPriceRequest struct {
	PricePerUnit float64 `json:"price_per_unit" validate:"required,gt=0"`
}
