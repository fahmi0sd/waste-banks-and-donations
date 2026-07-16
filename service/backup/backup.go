package backup

import "time"

const (
	TriggeredManual    = "manual"
	TriggeredScheduler = "scheduler"

	StatusSuccess = "success"
	StatusFailed  = "failed"
)

type BackupLog struct {
	ID                int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	FileName          string    `gorm:"column:file_name" json:"file_name"`
	FilePath          string    `gorm:"column:file_path" json:"file_path"`
	TriggeredBy       string    `gorm:"column:triggered_by" json:"triggered_by"`
	TriggeredByUserID *int      `gorm:"column:triggered_by_user_id" json:"triggered_by_user_id,omitempty"`
	Status            string    `gorm:"column:status" json:"status"`
	FileSizeBytes     int64     `gorm:"column:file_size_bytes" json:"file_size_bytes"`
	StartedAt         time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt        time.Time `gorm:"column:finished_at" json:"finished_at"`
}

func (BackupLog) TableName() string {
	return "backup_log"
}
