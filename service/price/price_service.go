package price

import (
	"errors"
	"log/slog"
)

const (
	roleAdmin       = "admin"
	roleMasterAdmin = "master_admin"
)

type Service interface {
	SetPrice(requesterID, categoryID int, req SetPriceRequest) (Price, error)
	ActivePrice(categoryID int) (Price, error)
	History(categoryID int) ([]Price, error)
}

type service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(logger *slog.Logger, repo Repository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) requireAdmin(requesterID int) error {
	role, err := s.repo.RoleOf(requesterID)
	if err != nil {
		s.logger.Error("failed to resolve requester role", "error", err, "user_id", requesterID)
		return errors.New("gagal memverifikasi akses")
	}
	if role != roleAdmin && role != roleMasterAdmin {
		return errors.New("akses ditolak: hanya admin yang bisa mengatur harga sampah")
	}
	return nil
}

// SetPrice menutup harga aktif sebelumnya (jika ada) lalu menyimpan harga baru
// sehingga histori perubahan harga tetap tersimpan di tabel waste_price.
func (s *service) SetPrice(requesterID, categoryID int, req SetPriceRequest) (Price, error) {
	if err := s.requireAdmin(requesterID); err != nil {
		return Price{}, err
	}

	exists, err := s.repo.CategoryExists(categoryID)
	if err != nil {
		s.logger.Error("failed to check category", "error", err, "category_id", categoryID)
		return Price{}, errors.New("gagal memeriksa kategori sampah")
	}
	if !exists {
		return Price{}, errors.New("kategori sampah tidak ditemukan")
	}

	if err := s.repo.CloseActivePrice(categoryID); err != nil {
		s.logger.Error("failed to close active price", "error", err, "category_id", categoryID)
		return Price{}, errors.New("gagal menutup harga lama")
	}

	created, err := s.repo.InsertPrice(categoryID, req.PricePerUnit, requesterID)
	if err != nil {
		s.logger.Error("failed to insert price", "error", err, "category_id", categoryID)
		return Price{}, errors.New("gagal menyimpan harga baru")
	}
	return created, nil
}

func (s *service) ActivePrice(categoryID int) (Price, error) {
	p, found, err := s.repo.ActivePrice(categoryID)
	if err != nil {
		s.logger.Error("failed to fetch active price", "error", err, "category_id", categoryID)
		return Price{}, errors.New("gagal mengambil harga aktif")
	}
	if !found {
		return Price{}, errors.New("kategori ini belum punya harga aktif")
	}
	return p, nil
}

func (s *service) History(categoryID int) ([]Price, error) {
	list, err := s.repo.History(categoryID)
	if err != nil {
		s.logger.Error("failed to fetch price history", "error", err, "category_id", categoryID)
		return nil, errors.New("gagal mengambil histori harga")
	}
	return list, nil
}
