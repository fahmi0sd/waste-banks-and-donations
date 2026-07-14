package notification

import "time"

const (
	TypeDepositUpdate  = "deposit_update"
	TypeDonationUpdate = "donation_update"
)

const (
	ChannelEmail = "email"
	ChannelWA    = "wa"
)

const (
	StatusSent   = "sent"
	StatusFailed = "failed"
)

type NotificationLog struct {
	ID      int       `gorm:"primaryKey" json:"id"`
	UserID  int       `json:"user_id"`
	Type    string    `json:"type"`
	Channel string    `json:"channel"`
	Content string    `json:"content"`
	Status  string    `json:"status"`
	SentAt  time.Time `json:"sent_at"`
}

func (NotificationLog) TableName() string {
	return "notification_log"
}
