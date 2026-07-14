package wastetransaction

import (
	"errors"
	"log/slog"

	"github.com/fahmi0sd/waste-banks-and-donations/pkg"
	"gorm.io/gorm"
)

type Service interface {
	Create(adminID, queueID int, items []ItemInput) (WasteTransaction, error)
	GetByID(id, requesterID int) (WasteTransaction, error)
	MyHistory(userID int) ([]WasteTransaction, error)
	LocationHistory(adminID, locationID int) ([]WasteTransaction, error)
	Cancel(adminID, transactionID int) (WasteTransaction, error)
}
type NotificationSender interface {
	NotifyDepositUpdate(userID int, amount float64) error
}

type service struct {
	logger   *slog.Logger
	repo     Repository
	db       *gorm.DB
	notifier NotificationSender
}

func NewService(logger *slog.Logger, repo Repository, db *gorm.DB, notifier NotificationSender) Service {
	return &service{logger: logger, repo: repo, db: db, notifier: notifier}
}

func (s *service) Create(adminID, queueID int, items []ItemInput) (WasteTransaction, error) {
	identity, err := pkg.GetUserIdentity(s.db, adminID)
	if err != nil {
		s.logger.Error("failed to resolve identity", "error", err, "user_id", adminID)
		return WasteTransaction{}, errors.New("gagal memverifikasi akses")
	}
	if identity.Role != "admin" || identity.LocationID == nil {
		return WasteTransaction{}, errors.New("akses ditolak: hanya admin lokasi yang bisa membuat transaksi")
	}

	if len(items) == 0 {
		return WasteTransaction{}, errors.New("minimal 1 item sampah wajib diisi")
	}
	for _, item := range items {
		if item.Weight <= 0 {
			return WasteTransaction{}, errors.New("berat tiap item harus lebih dari 0")
		}
	}

	q, err := s.repo.QueueInfo(queueID)
	if err != nil {
		s.logger.Error("failed to fetch queue info", "error", err, "queue_id", queueID)
		return WasteTransaction{}, errors.New("gagal memeriksa data antrian")
	}
	if !q.Found {
		return WasteTransaction{}, errors.New("antrian tidak ditemukan")
	}
	if q.LocationID != *identity.LocationID {
		return WasteTransaction{}, errors.New("antrian ini bukan di lokasi kamu")
	}
	if q.Status != "verified" {
		return WasteTransaction{}, errors.New("antrian belum diverifikasi atau sudah diproses sebelumnya")
	}

	input := CreateInput{
		QueueID:    queueID,
		UserID:     q.UserID,
		LocationID: *identity.LocationID,
		AdminID:    adminID,
		Items:      items,
	}

	result, err := s.repo.CreateWithEffects(input)
	if err != nil {
		s.logger.Error("failed to create transaction", "error", err, "queue_id", queueID)
		return WasteTransaction{}, err // pesan dari repo sudah cukup jelas (lihat repository)
	}

	details, err := s.repo.GetDetailsByTransactionID(result.ID)
	if err == nil {
		result.Details = details
	}

	if s.notifier != nil {
		if notifErr := s.notifier.NotifyDepositUpdate(result.UserID, result.TotalRupiah); notifErr != nil {
			s.logger.Error("failed to send deposit notification", "error", notifErr, "user_id", result.UserID)
		}
	}

	return result, nil
}

func (s *service) GetByID(id, requesterID int) (WasteTransaction, error) {
	tx, err := s.repo.GetByID(id)
	if err != nil {
		return WasteTransaction{}, errors.New("transaksi tidak ditemukan")
	}

	identity, err := pkg.GetUserIdentity(s.db, requesterID)
	if err != nil {
		return WasteTransaction{}, errors.New("gagal memverifikasi akses")
	}

	allowed := false
	switch identity.Role {
	case "master_admin":
		allowed = true
	case "admin":
		allowed = identity.LocationID != nil && tx.LocationID == *identity.LocationID
	default:
		allowed = tx.UserID == requesterID
	}
	if !allowed {
		return WasteTransaction{}, errors.New("kamu tidak punya akses ke transaksi ini")
	}

	details, err := s.repo.GetDetailsByTransactionID(tx.ID)
	if err == nil {
		tx.Details = details
	}
	return tx, nil
}

func (s *service) MyHistory(userID int) ([]WasteTransaction, error) {
	list, err := s.repo.ListByUser(userID)
	if err != nil {
		s.logger.Error("failed to list user history", "error", err, "user_id", userID)
		return nil, errors.New("gagal mengambil riwayat transaksi")
	}
	return list, nil
}

func (s *service) LocationHistory(adminID, locationID int) ([]WasteTransaction, error) {
	identity, err := pkg.GetUserIdentity(s.db, adminID)
	if err != nil {
		return nil, errors.New("gagal memverifikasi akses")
	}
	if identity.Role != "admin" || identity.LocationID == nil || *identity.LocationID != locationID {
		return nil, errors.New("akses ditolak: hanya admin di lokasi tersebut yang bisa melihat riwayatnya")
	}

	list, err := s.repo.ListByLocation(locationID)
	if err != nil {
		s.logger.Error("failed to list location history", "error", err, "location_id", locationID)
		return nil, errors.New("gagal mengambil riwayat transaksi")
	}
	return list, nil
}

func (s *service) Cancel(adminID, transactionID int) (WasteTransaction, error) {
	identity, err := pkg.GetUserIdentity(s.db, adminID)
	if err != nil {
		return WasteTransaction{}, errors.New("gagal memverifikasi akses")
	}
	if identity.Role != "admin" || identity.LocationID == nil {
		return WasteTransaction{}, errors.New("akses ditolak: hanya admin lokasi yang bisa membatalkan transaksi")
	}

	tx, err := s.repo.GetByID(transactionID)
	if err != nil {
		return WasteTransaction{}, errors.New("transaksi tidak ditemukan")
	}
	if tx.LocationID != *identity.LocationID {
		return WasteTransaction{}, errors.New("transaksi ini bukan di lokasi kamu")
	}
	if tx.Status != StatusCompleted {
		return WasteTransaction{}, errors.New("transaksi dengan status " + tx.Status + " tidak bisa dibatalkan")
	}

	result, err := s.repo.CancelWithEffects(transactionID)
	if err != nil {
		s.logger.Error("failed to cancel transaction", "error", err, "transaction_id", transactionID)
		return WasteTransaction{}, err
	}
	return result, nil
}
