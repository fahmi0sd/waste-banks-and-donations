package auth

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req RegisterRequest) (RegisteredUser, error)
	Login(req LoginRequest) (LoginResult, error)
}

type service struct {
	logger      *slog.Logger
	repo        Repository
	jwtSecret   string
	jwtTTLHours int
}

func NewService(logger *slog.Logger, repo Repository, jwtSecret string, jwtTTLHours int) Service {
	if jwtTTLHours <= 0 {
		jwtTTLHours = 24
	}
	return &service{logger: logger, repo: repo, jwtSecret: jwtSecret, jwtTTLHours: jwtTTLHours}
}

func (s *service) Register(req RegisterRequest) (RegisteredUser, error) {
	exists, err := s.repo.EmailExists(req.Email)
	if err != nil {
		s.logger.Error("failed to check email existence", "error", err, "email", req.Email)
		return RegisteredUser{}, errors.New("gagal memeriksa email")
	}
	if exists {
		return RegisteredUser{}, errors.New("email sudah terdaftar")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return RegisteredUser{}, errors.New("gagal memproses password")
	}

	created, err := s.repo.Create(req.Name, req.Email, req.Phone, string(hashed))
	if err != nil {
		s.logger.Error("failed to create user", "error", err, "email", req.Email)
		return RegisteredUser{}, errors.New("gagal mendaftarkan akun, coba lagi")
	}

	return RegisteredUser{ID: created.ID, Name: created.Name, Email: created.Email, Role: created.Role}, nil
}

func (s *service) Login(req LoginRequest) (LoginResult, error) {
	user, found, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		s.logger.Error("failed to find user by email", "error", err, "email", req.Email)
		return LoginResult{}, errors.New("gagal memproses login")
	}
	if !found {
		return LoginResult{}, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return LoginResult{}, errors.New("email atau password salah")
	}

	ttl := time.Duration(s.jwtTTLHours) * time.Hour
	now := time.Now()
	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"iat":  now.Unix(),
		"exp":  now.Add(ttl).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("failed to sign jwt", "error", err, "user_id", user.ID)
		return LoginResult{}, errors.New("gagal membuat token")
	}

	return LoginResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.jwtTTLHours * 3600,
	}, nil
}
