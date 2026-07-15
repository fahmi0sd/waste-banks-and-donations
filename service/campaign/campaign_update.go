package campaign

import (
	"errors"
	"time"
)

type CampaignUpdate struct {
	ID            int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CampaignID    int       `gorm:"column:campaign_id" json:"campaign_id"`
	Content       string    `gorm:"column:content" json:"content"`
	ReportFileURL string    `gorm:"column:report_file_url" json:"report_file_url"`
	PublishedAt   time.Time `gorm:"column:published_at" json:"published_at"`
}

func (CampaignUpdate) TableName() string {
	return "campaign_update"
}

type CreateUpdateRequest struct {
	Content       string `json:"content"`
	ReportFileURL string `json:"report_file_url"`
}

func (r CreateUpdateRequest) Validate() error {
	if r.Content == "" {
		return errors.New("content wajib diisi")
	}
	return nil
}
