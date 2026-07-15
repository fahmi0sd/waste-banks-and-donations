package backup

import "time"

const (
	TriggeredManual    = "manual"
	TriggeredScheduler = "scheduler"

	StatusSuccess = "success"
	StatusFailed  = "failed"
)

type BackupLog struct {
	ID                int       `gorm:"primaryKey" json:"id"`
	TriggeredBy       string    `json:"triggered_by"`
	TriggeredByUserID *int      `json:"triggered_by_user_id,omitempty"`
	Status            string    `json:"status"`
	FileSizeBytes     int64     `json:"file_size_bytes"`
	StartedAt         time.Time `json:"started_at"`
	FinishedAt        time.Time `json:"finished_at"`
}

func (BackupLog) TableName() string {
	return "backup_log"
}
