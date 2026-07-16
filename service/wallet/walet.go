package wallet

import "time"

type Wallet struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserID    int       `json:"user_id"`
	Balance   float64   `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Wallet) TableName() string {
	return "wallet"
}

type WalletTransaction struct {
	ID            int       `gorm:"primaryKey" json:"id"`
	WalletID      int       `json:"wallet_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   *int      `json:"reference_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transaction"
}

const (
	TransactionCreditWaste  = "credit_waste"
	TransactionDonationOut  = "donation_out"
	TransactionWithdraw     = "withdraw"
	TransactionCancellation = "cancellation"
)
