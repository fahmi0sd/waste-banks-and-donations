package calculator

type SimulateItemRequest struct {
	CategoryID int     `json:"category_id" validate:"required"`
	Weight     float64 `json:"weight" validate:"required,gt=0"`
}

type SimulateRequest struct {
	Items []SimulateItemRequest `json:"items" validate:"required,min=1,dive"`
}

type SimulateItemResult struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Weight       float64 `json:"weight"`
	PricePerUnit float64 `json:"price_per_unit"`
	Subtotal     float64 `json:"subtotal"`
}

type SimulateResult struct {
	Items       []SimulateItemResult `json:"items"`
	TotalRupiah float64              `json:"total_rupiah"`
}
