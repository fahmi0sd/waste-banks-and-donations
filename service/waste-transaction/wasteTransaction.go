package wastetransaction

import "time"

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

type WasteTransaction struct {
	ID          int                      `gorm:"primaryKey" json:"id"`
	QueueID     *int                     `json:"queue_id,omitempty"`
	UserID      int                      `json:"user_id"`
	LocationID  int                      `json:"location_id"`
	AdminID     int                      `json:"admin_id"`
	TotalRupiah float64                  `json:"total_rupiah"`
	Status      string                   `json:"status"`
	CreatedAt   time.Time                `json:"created_at"`
	Details     []WasteTransactionDetail `gorm:"-" json:"details,omitempty"`
}

func (WasteTransaction) TableName() string {
	return "waste_transaction"
}

type WasteTransactionDetail struct {
	ID                   int     `gorm:"primaryKey" json:"id"`
	TransactionID        int     `json:"transaction_id"`
	CategoryID           int     `json:"category_id"`
	CategoryName         string  `gorm:"-" json:"category_name,omitempty"`
	Weight               float64 `json:"weight"`
	PricePerUnitSnapshot float64 `json:"price_per_unit_snapshot"`
	Subtotal             float64 `json:"subtotal"`
}

func (WasteTransactionDetail) TableName() string {
	return "waste_transaction_detail"
}
