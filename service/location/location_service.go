package location

import (
	"errors"
	"log/slog"
)

const (
	roleAdmin       = "admin"
	roleMasterAdmin = "master_admin"
)

type Service interface {
	Create(requesterID int, req CreateRequest) (Location, error)
	List() ([]Location, error)
	Get(id int) (Location, error)
	Update(requesterID, id int, req UpdateRequest) (Location, error)
	Delete(requesterID, id int) error
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
		return errors.New("akses ditolak: hanya admin yang bisa mengelola lokasi")
	}
	return nil
}

func (s *service) Create(requesterID int, req CreateRequest) (Location, error) {
	if err := s.requireAdmin(requesterID); err != nil {
		return Location{}, err
	}

	l := Location{
		Name:      req.Name,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		IsOpen:    true,
		OpenTime:  "08:00",
		CloseTime: "16:00",
	}
	if req.OpenTime != "" {
		l.OpenTime = req.OpenTime
	}
	if req.CloseTime != "" {
		l.CloseTime = req.CloseTime
	}
	if req.IsOpen != nil {
		l.IsOpen = *req.IsOpen
	}

	created, err := s.repo.Create(l)
	if err != nil {
		s.logger.Error("failed to create location", "error", err)
		return Location{}, errors.New("gagal membuat lokasi")
	}
	return created, nil
}

func (s *service) List() ([]Location, error) {
	list, err := s.repo.List()
	if err != nil {
		s.logger.Error("failed to list locations", "error", err)
		return nil, errors.New("gagal mengambil daftar lokasi")
	}
	return list, nil
}

func (s *service) Get(id int) (Location, error) {
	l, found, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get location", "error", err, "id", id)
		return Location{}, errors.New("gagal mengambil lokasi")
	}
	if !found {
		return Location{}, errors.New("lokasi tidak ditemukan")
	}
	return l, nil
}

func (s *service) Update(requesterID, id int, req UpdateRequest) (Location, error) {
	if err := s.requireAdmin(requesterID); err != nil {
		return Location{}, err
	}
	if _, found, err := s.repo.GetByID(id); err != nil {
		return Location{}, errors.New("gagal memeriksa lokasi")
	} else if !found {
		return Location{}, errors.New("lokasi tidak ditemukan")
	}

	if err := s.repo.Update(id, req.Name, req.Address, req.Latitude, req.Longitude, req.OpenTime, req.CloseTime, req.IsOpen); err != nil {
		s.logger.Error("failed to update location", "error", err, "id", id)
		return Location{}, errors.New("gagal memperbarui lokasi")
	}
	return s.Get(id)
}

func (s *service) Delete(requesterID, id int) error {
	if err := s.requireAdmin(requesterID); err != nil {
		return err
	}
	if _, found, err := s.repo.GetByID(id); err != nil {
		return errors.New("gagal memeriksa lokasi")
	} else if !found {
		return errors.New("lokasi tidak ditemukan")
	}

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("failed to delete location", "error", err, "id", id)
		return errors.New("gagal menghapus lokasi")
	}
	return nil
}
