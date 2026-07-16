package campaign

import (
	"errors"
	"time"
)

const (
	StatusDraft     = "draft"
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusClosed    = "closed"
)

type Campaign struct {
	ID            int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title         string    `gorm:"column:title" json:"title"`
	Description   string    `gorm:"column:description" json:"description"`
	Organizer     string    `gorm:"column:organizer" json:"organizer"`
	TargetAmount  float64   `gorm:"column:target_amount" json:"target_amount"`
	CurrentAmount float64   `gorm:"column:current_amount" json:"current_amount"`
	Status        string    `gorm:"column:status" json:"status"`
	CreatedBy     *int      `gorm:"column:created_by" json:"created_by,omitempty"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Campaign) TableName() string {
	return "donation_campaign"
}

type CreateRequest struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Organizer    string  `json:"organizer"`
	TargetAmount float64 `json:"target_amount"`
	Status       string  `json:"status"`
}

type UpdateRequest struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Organizer    string  `json:"organizer"`
	TargetAmount float64 `json:"target_amount"`
	Status       string  `json:"status"`
}

func (r *CreateRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title wajib diisi")
	}

	if r.TargetAmount <= 0 {
		return errors.New("target amount harus lebih dari 0")
	}

	if r.Status == "" {
		r.Status = StatusDraft
	}
	return nil
}
