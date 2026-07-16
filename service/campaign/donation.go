package campaign

import (
	"errors"
	"time"
)

type Donation struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     int       `gorm:"column:user_id" json:"user_id"`
	CampaignID int       `gorm:"column:campaign_id" json:"campaign_id"`
	Amount     float64   `gorm:"column:amount" json:"amount"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Donation) TableName() string {
	return "donation"
}

type DonateRequest struct {
	Amount float64 `json:"amount"`
}

func (r DonateRequest) Validate() error {
	if r.Amount <= 0 {
		return errors.New("amount harus lebih dari 0")
	}
	return nil
}
