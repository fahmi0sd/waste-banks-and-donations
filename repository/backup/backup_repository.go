package backup

import (
	backupService "github.com/fahmi0sd/waste-banks-and-donations/service/backup"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

// Compile time check
var _ backupService.Repository = (*GormRepository)(nil)

func (r *GormRepository) Create(log backupService.BackupLog) (backupService.BackupLog, error) {

	if err := r.db.Create(&log).Error; err != nil {
		return backupService.BackupLog{}, err
	}

	return log, nil
}

func (r *GormRepository) History() ([]backupService.BackupLog, error) {

	var logs []backupService.BackupLog

	err := r.db.
		Order("started_at DESC").
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}
