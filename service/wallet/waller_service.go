package wallet

import (
	"errors"
	"log/slog"
)

type Service interface {
	GetWallet(userID int) (Wallet, error)
	GetTransactions(userID int) ([]WalletTransaction, error)
	Withdraw(userID int, amount float64) (WalletTransaction, error)
}

type service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(logger *slog.Logger, repo Repository) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) GetWallet(userID int) (Wallet, error) {
	wallet, err := s.repo.GetWallet(userID)
	if err != nil {
		s.logger.Error(
			"failed get wallet",
			"error", err,
			"user_id", userID,
		)
		return Wallet{}, err
	}

	return wallet, nil
}

func (s *service) GetTransactions(userID int) ([]WalletTransaction, error) {
	transactions, err := s.repo.GetTransactions(userID)
	if err != nil {
		s.logger.Error(
			"failed get transactions",
			"error", err,
			"user_id", userID,
		)
		return nil, err
	}

	return transactions, nil
}

func (s *service) Withdraw(userID int, amount float64) (WalletTransaction, error) {
	if amount <= 0 {
		return WalletTransaction{}, errors.New("nominal withdraw harus lebih dari 0")
	}

	transaction, err := s.repo.Withdraw(userID, amount)
	if err != nil {
		s.logger.Error(
			"failed withdraw wallet",
			"error", err,
			"user_id", userID,
			"amount", amount,
		)
		return WalletTransaction{}, err
	}

	return transaction, nil
}
