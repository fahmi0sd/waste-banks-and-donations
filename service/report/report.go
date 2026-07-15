package report

type TransactionSummary struct {
	LocationID   int     `json:"location_id"`
	LocationName string  `json:"location_name"`
	Date         string  `json:"date"`
	TotalTrx     int     `json:"total_transaksi"`
	TotalRupiah  float64 `json:"total_rupiah"`
}

type Filter struct {
	LocationID *int
	DateFrom   *string // format YYYY-MM-DD
	DateTo     *string // format YYYY-MM-DD
}
