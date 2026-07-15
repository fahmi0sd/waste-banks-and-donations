package user

import (
	"errors"
	"log/slog"
)

type Service interface {
	GetProfile(userID int) (Profile, error)
	UpdateProfile(userID int, req UpdateProfileRequest) (Profile, error)
}

type service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(logger *slog.Logger, repo Repository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) GetProfile(userID int) (Profile, error) {
	p, found, err := s.repo.FindByID(userID)
	if err != nil {
		s.logger.Error("failed to fetch profile", "error", err, "user_id", userID)
		return Profile{}, errors.New("gagal mengambil profil")
	}
	if !found {
		return Profile{}, errors.New("user tidak ditemukan")
	}
	return p, nil
}

func (s *service) UpdateProfile(userID int, req UpdateProfileRequest) (Profile, error) {
	if err := s.repo.UpdateProfile(userID, req.Name, req.Phone, req.LocationID); err != nil {
		s.logger.Error("failed to update profile", "error", err, "user_id", userID)
		return Profile{}, errors.New("gagal memperbarui profil")
	}
	return s.GetProfile(userID)
}
