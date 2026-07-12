package queue

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type Service interface {
	CreateQueue(userID, locationID int, preferredTime *time.Time) (Queue, error)
	GetQueue(id, userID int) (Queue, error)
	AdminListQueue(adminUserID int) ([]QueueWithUser, error)
	AdminVerifyQueue(adminUserID, queueID int) (Queue, error)
}

type service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(logger *slog.Logger, repo Repository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) CreateQueue(userID, locationID int, preferredTime *time.Time) (Queue, error) {
	exists, isOpen, openTime, closeTime, err := s.repo.LocationInfo(locationID)
	if err != nil {
		s.logger.Error("failed to fetch location info", "error", err, "location_id", locationID)
		return Queue{}, errors.New("gagal memeriksa data lokasi")
	}
	if !exists {
		return Queue{}, errors.New("lokasi tidak ditemukan")
	}
	if !isOpen {
		return Queue{}, errors.New("lokasi sedang tutup, tidak bisa request antrian")
	}

	if preferredTime != nil {
		if err := validateOperationalHours(*preferredTime, openTime, closeTime); err != nil {
			return Queue{}, err
		}
	}

	count, err := s.repo.CountTodayByLocation(locationID)
	if err != nil {
		s.logger.Error("failed to count today's queue", "error", err, "location_id", locationID)
		return Queue{}, errors.New("gagal membuat nomor antrian")
	}

	q := Queue{
		UserID:        userID,
		LocationID:    locationID,
		QueueNumber:   fmt.Sprintf("%03d", count+1),
		PreferredTime: preferredTime,
		Status:        StatusWaiting,
	}

	created, err := s.repo.Create(q)
	if err != nil {
		s.logger.Error("failed to create queue", "error", err, "user_id", userID)
		return Queue{}, errors.New("gagal membuat antrian, coba lagi")
	}

	return created, nil
}

func (s *service) GetQueue(id, userID int) (Queue, error) {
	q, err := s.repo.GetByID(id)
	if err != nil {
		return Queue{}, errors.New("antrian tidak ditemukan")
	}
	if q.UserID != userID {
		return Queue{}, errors.New("antrian ini bukan milik kamu")
	}
	return q, nil
}

func (s *service) AdminListQueue(adminUserID int) ([]QueueWithUser, error) {
	role, locationID, err := s.repo.AdminLocation(adminUserID)
	if err != nil {
		s.logger.Error("failed to resolve admin location", "error", err, "user_id", adminUserID)
		return nil, errors.New("gagal memverifikasi akses")
	}
	if role != "admin" || locationID == nil {
		return nil, errors.New("akses ditolak: hanya admin lokasi yang bisa melihat antrian")
	}

	list, err := s.repo.ListWaitingByLocation(*locationID)
	if err != nil {
		s.logger.Error("failed to list queue", "error", err, "location_id", *locationID)
		return nil, errors.New("gagal mengambil daftar antrian")
	}
	return list, nil
}

func (s *service) AdminVerifyQueue(adminUserID, queueID int) (Queue, error) {
	role, locationID, err := s.repo.AdminLocation(adminUserID)
	if err != nil {
		s.logger.Error("failed to resolve admin location", "error", err, "user_id", adminUserID)
		return Queue{}, errors.New("gagal memverifikasi akses")
	}
	if role != "admin" || locationID == nil {
		return Queue{}, errors.New("akses ditolak: hanya admin lokasi yang bisa verifikasi antrian")
	}

	q, err := s.repo.GetByID(queueID)
	if err != nil {
		return Queue{}, errors.New("antrian tidak ditemukan")
	}
	if q.LocationID != *locationID {
		return Queue{}, errors.New("antrian ini bukan di lokasi kamu")
	}
	if q.Status != StatusWaiting {
		return Queue{}, fmt.Errorf("antrian sudah berstatus %s, tidak bisa diverifikasi ulang", q.Status)
	}

	if err := s.repo.UpdateStatus(queueID, StatusVerified, adminUserID); err != nil {
		s.logger.Error("failed to verify queue", "error", err, "queue_id", queueID)
		return Queue{}, errors.New("gagal memverifikasi antrian")
	}

	q.Status = StatusVerified
	q.VerifiedBy = &adminUserID
	return q, nil
}

func validateOperationalHours(preferredTime time.Time, openTime, closeTime string) error {
	open, err := time.Parse("15:04", openTime)
	if err != nil {
		return errors.New("format jam operasional lokasi tidak valid")
	}
	closeT, err := time.Parse("15:04", closeTime)
	if err != nil {
		return errors.New("format jam operasional lokasi tidak valid")
	}

	reqMinutes := preferredTime.Hour()*60 + preferredTime.Minute()
	openMinutes := open.Hour()*60 + open.Minute()
	closeMinutes := closeT.Hour()*60 + closeT.Minute()

	if reqMinutes < openMinutes || reqMinutes > closeMinutes {
		return fmt.Errorf("jam kedatangan di luar jam operasional lokasi (%s - %s)", openTime, closeTime)
	}
	return nil
}
