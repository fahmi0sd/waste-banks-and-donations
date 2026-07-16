package backup

import (
	"errors"
	"log/slog"
	"time"

	"github.com/fahmi0sd/waste-banks-and-donations/pkg"
	"gorm.io/gorm"
)

type Service interface {
	Trigger(masterAdminID int) (BackupLog, error)
	TriggerScheduler() error
	History(requesterID int) ([]BackupLog, error)
}

type service struct {
	logger *slog.Logger
	repo   Repository
	db     *gorm.DB
}

func NewService(
	logger *slog.Logger,
	repo Repository,
	db *gorm.DB,
) Service {

	return &service{
		logger: logger,
		repo:   repo,
		db:     db,
	}
}

func (s *service) Trigger(masterAdminID int) (BackupLog, error) {

	identity, err := pkg.GetUserIdentity(s.db, masterAdminID)
	if err != nil {
		return BackupLog{}, errors.New("gagal memverifikasi user")
	}

	if identity.Role != "master_admin" {
		return BackupLog{}, errors.New("akses ditolak")
	}

	return s.run(
		TriggeredManual,
		&masterAdminID,
	)
}

func (s *service) TriggerScheduler() error {

	_, err := s.run(
		TriggeredScheduler,
		nil,
	)

	return err
}

func (s *service) History(requesterID int) ([]BackupLog, error) {

	identity, err := pkg.GetUserIdentity(s.db, requesterID)
	if err != nil {
		return nil, errors.New("gagal memverifikasi user")
	}

	if identity.Role != "master_admin" {
		return nil, errors.New("akses ditolak")
	}

	return s.repo.History()
}

func (s *service) run(
	triggeredBy string,
	userID *int,
) (BackupLog, error) {

	start := time.Now()

	filename, path, size, err := s.executeBackup()

	status := StatusSuccess

	if err != nil {

		status = StatusFailed

		s.logger.Error(
			"backup failed",
			"error", err,
		)
	}

	log := BackupLog{
		FileName:          filename,
		FilePath:          path,
		TriggeredBy:       triggeredBy,
		TriggeredByUserID: userID,
		Status:            status,
		FileSizeBytes:     size,
		StartedAt:         start,
		FinishedAt:        time.Now(),
	}

	result, saveErr := s.repo.Create(log)

	if saveErr != nil {
		return BackupLog{}, saveErr
	}

	if err != nil {
		return result, err
	}

	s.logger.Info(
		"backup created",
		"file", filename,
		"size", size,
	)

	return result, nil
}
