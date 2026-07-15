package category

import (
	"errors"
	"log/slog"
)

const (
	roleMasterAdmin = "master_admin"
)

type Service interface {
	Create(requesterID int, req CreateRequest) (Category, error)
	List() ([]Category, error)
	Get(id int) (Category, error)
	Update(requesterID, id int, req UpdateRequest) (Category, error)
	Delete(requesterID, id int) error
}

type service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(logger *slog.Logger, repo Repository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) requireMasterAdmin(requesterID int) error {
	role, err := s.repo.RoleOf(requesterID)
	if err != nil {
		s.logger.Error("failed to resolve requester role", "error", err, "user_id", requesterID)
		return errors.New("gagal memverifikasi akses")
	}
	if role != roleMasterAdmin {
		return errors.New("akses ditolak: hanya master admin yang bisa mengelola kategori sampah")
	}
	return nil
}

func (s *service) Create(requesterID int, req CreateRequest) (Category, error) {
	if err := s.requireMasterAdmin(requesterID); err != nil {
		return Category{}, err
	}

	if req.ParentCategoryID != nil {
		if _, found, err := s.repo.GetByID(*req.ParentCategoryID); err != nil {
			return Category{}, errors.New("gagal memeriksa kategori induk")
		} else if !found {
			return Category{}, errors.New("kategori induk tidak ditemukan")
		}
	}

	unit := "kg"
	if req.Unit != "" {
		unit = req.Unit
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	created, err := s.repo.Create(Category{
		Name: req.Name, ParentCategoryID: req.ParentCategoryID,
		Unit: unit, IsActive: isActive,
	})
	if err != nil {
		s.logger.Error("failed to create category", "error", err)
		return Category{}, errors.New("gagal membuat kategori sampah")
	}
	return created, nil
}

func (s *service) List() ([]Category, error) {
	list, err := s.repo.List()
	if err != nil {
		s.logger.Error("failed to list categories", "error", err)
		return nil, errors.New("gagal mengambil daftar kategori sampah")
	}
	return list, nil
}

func (s *service) Get(id int) (Category, error) {
	c, found, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get category", "error", err, "id", id)
		return Category{}, errors.New("gagal mengambil kategori sampah")
	}
	if !found {
		return Category{}, errors.New("kategori sampah tidak ditemukan")
	}
	return c, nil
}

func (s *service) Update(requesterID, id int, req UpdateRequest) (Category, error) {
	if err := s.requireMasterAdmin(requesterID); err != nil {
		return Category{}, err
	}
	if _, found, err := s.repo.GetByID(id); err != nil {
		return Category{}, errors.New("gagal memeriksa kategori sampah")
	} else if !found {
		return Category{}, errors.New("kategori sampah tidak ditemukan")
	}

	if err := s.repo.Update(id, req.Name, req.ParentCategoryID, req.Unit, req.IsActive); err != nil {
		s.logger.Error("failed to update category", "error", err, "id", id)
		return Category{}, errors.New("gagal memperbarui kategori sampah")
	}
	return s.Get(id)
}

func (s *service) Delete(requesterID, id int) error {
	if err := s.requireMasterAdmin(requesterID); err != nil {
		return err
	}
	if _, found, err := s.repo.GetByID(id); err != nil {
		return errors.New("gagal memeriksa kategori sampah")
	} else if !found {
		return errors.New("kategori sampah tidak ditemukan")
	}

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("failed to delete category", "error", err, "id", id)
		return errors.New("gagal menghapus kategori sampah, kemungkinan masih dipakai data harga/transaksi")
	}
	return nil
}
