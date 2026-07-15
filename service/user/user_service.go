package user

import (
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetProfile(userID int) (Profile, error)
	UpdateProfile(userID int, req UpdateProfileRequest) (Profile, error)

	ListAccounts(requesterID int) ([]Profile, error)
	CreateAccount(requesterID int, req CreateAccountRequest) (Profile, error)
	UpdateAccount(requesterID, targetID int, req UpdateAccountRequest) (Profile, error)
	DeleteAccount(requesterID, targetID int) error
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

func (s *service) requireMasterAdmin(requesterID int) error {
	role, err := s.repo.RoleOf(requesterID)
	if err != nil {
		s.logger.Error("failed to resolve requester role", "error", err, "user_id", requesterID)
		return errors.New("gagal memverifikasi akses")
	}
	if role != RoleMasterAdmin {
		return errors.New("akses ditolak: hanya master admin yang bisa mengelola akun")
	}
	return nil
}

func (s *service) ListAccounts(requesterID int) ([]Profile, error) {
	if err := s.requireMasterAdmin(requesterID); err != nil {
		return nil, err
	}
	list, err := s.repo.ListAccounts()
	if err != nil {
		s.logger.Error("failed to list accounts", "error", err)
		return nil, errors.New("gagal mengambil daftar akun")
	}
	return list, nil
}

func (s *service) CreateAccount(requesterID int, req CreateAccountRequest) (Profile, error) {
	if err := s.requireMasterAdmin(requesterID); err != nil {
		return Profile{}, err
	}

	exists, err := s.repo.EmailExists(req.Email)
	if err != nil {
		s.logger.Error("failed to check email existence", "error", err, "email", req.Email)
		return Profile{}, errors.New("gagal memeriksa email")
	}
	if exists {
		return Profile{}, errors.New("email sudah terdaftar")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return Profile{}, errors.New("gagal memproses password")
	}

	created, err := s.repo.CreateAccount(req.Name, req.Email, req.Phone, string(hashed), req.Role, req.LocationID)
	if err != nil {
		s.logger.Error("failed to create account", "error", err, "email", req.Email)
		return Profile{}, errors.New("gagal membuat akun")
	}
	return created, nil
}

func (s *service) UpdateAccount(requesterID, targetID int, req UpdateAccountRequest) (Profile, error) {
	if err := s.requireMasterAdmin(requesterID); err != nil {
		return Profile{}, err
	}

	if _, found, err := s.repo.FindByID(targetID); err != nil {
		s.logger.Error("failed to check target account", "error", err, "target_id", targetID)
		return Profile{}, errors.New("gagal memeriksa akun")
	} else if !found {
		return Profile{}, errors.New("akun tidak ditemukan")
	}

	if err := s.repo.UpdateAccount(targetID, req.Name, req.Phone, req.Role, req.LocationID); err != nil {
		s.logger.Error("failed to update account", "error", err, "target_id", targetID)
		return Profile{}, errors.New("gagal memperbarui akun")
	}

	p, _, err := s.repo.FindByID(targetID)
	if err != nil {
		return Profile{}, errors.New("gagal mengambil akun terbaru")
	}
	return p, nil
}

func (s *service) DeleteAccount(requesterID, targetID int) error {
	if err := s.requireMasterAdmin(requesterID); err != nil {
		return err
	}
	if requesterID == targetID {
		return errors.New("tidak bisa menghapus akun sendiri")
	}

	if _, found, err := s.repo.FindByID(targetID); err != nil {
		return errors.New("gagal memeriksa akun")
	} else if !found {
		return errors.New("akun tidak ditemukan")
	}

	if err := s.repo.DeleteAccount(targetID); err != nil {
		s.logger.Error("failed to delete account", "error", err, "target_id", targetID)
		return errors.New("gagal menghapus akun")
	}
	return nil
}
