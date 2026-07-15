package backup

type Repository interface {
	Create(log BackupLog) (BackupLog, error)
	History() ([]BackupLog, error)
}
