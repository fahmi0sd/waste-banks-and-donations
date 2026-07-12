package queue

import "time"

const (
	StatusWaiting  = "waiting"
	StatusVerified = "verified"
	StatusExpired  = "expired"
)

type Queue struct {
	ID            int        `gorm:"primaryKey" json:"id"`
	UserID        int        `json:"user_id"`
	LocationID    int        `json:"location_id"`
	QueueNumber   string     `json:"queue_number"`
	PreferredTime *time.Time `json:"preferred_time,omitempty"`
	Status        string     `json:"status"`
	VerifiedBy    *int       `json:"verified_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (Queue) TableName() string {
	return "queue"
}

type QueueWithUser struct {
	ID            int        `json:"id"`
	UserID        int        `json:"user_id"`
	UserName      string     `json:"user_name"`
	QueueNumber   string     `json:"queue_number"`
	PreferredTime *time.Time `json:"preferred_time,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
}
