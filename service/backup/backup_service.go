package backup

import (
	"errors"
	"log/slog"
)

type Service interface {
	Trigger(masterAdminID int) (BackupLog, error)
	History() ([]BackupLog, error)
}

type service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(
	logger *slog.Logger,
	repo Repository,
) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) Trigger(masterAdminID int) (BackupLog, error) {

	return BackupLog{}, errors.New("belum diimplementasikan")

}

func (s *service) History() ([]BackupLog, error) {

	list, err := s.repo.History()

	if err != nil {

		s.logger.Error(
			"failed get backup history",
			"error", err,
		)

		return nil, err

	}

	return list, nil
}
