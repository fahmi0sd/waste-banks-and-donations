package report

import (
	"errors"
	"log/slog"
)

const (
	roleAdmin       = "admin"
	roleMasterAdmin = "master_admin"
)

type Service interface {
	TransactionReport(requesterID int, filter Filter) ([]TransactionSummary, error)
}

type service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(logger *slog.Logger, repo Repository) Service {
	return &service{logger: logger, repo: repo}
}

// TransactionReport mengembalikan agregasi transaksi per lokasi/tanggal.
// master_admin bisa melihat semua lokasi; admin biasa hanya bisa melihat
// laporan lokasi tempat dia bertugas.
func (s *service) TransactionReport(requesterID int, filter Filter) ([]TransactionSummary, error) {
	role, locationID, err := s.repo.RoleAndLocation(requesterID)
	if err != nil {
		s.logger.Error("failed to resolve requester role", "error", err, "user_id", requesterID)
		return nil, errors.New("gagal memverifikasi akses")
	}
	if role != roleAdmin && role != roleMasterAdmin {
		return nil, errors.New("akses ditolak: hanya admin yang bisa melihat laporan transaksi")
	}

	if role == roleAdmin {
		if locationID == nil {
			return nil, errors.New("akun admin ini belum terdaftar di lokasi manapun")
		}
		if filter.LocationID != nil && *filter.LocationID != *locationID {
			return nil, errors.New("akses ditolak: admin hanya bisa melihat laporan lokasi sendiri")
		}
		filter.LocationID = locationID
	}

	list, err := s.repo.TransactionSummary(filter)
	if err != nil {
		s.logger.Error("failed to build transaction report", "error", err)
		return nil, errors.New("gagal mengambil laporan transaksi")
	}
	return list, nil
}
