package wallet

import (
	"errors"

	walletService "github.com/fahmi0sd/waste-banks-and-donations/service/wallet"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

// Compile-time check
var _ walletService.Repository = (*GormRepository)(nil)

func (r *GormRepository) GetWallet(userID int) (walletService.Wallet, error) {
	var wallet walletService.Wallet

	err := r.db.
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return walletService.Wallet{}, errors.New("wallet tidak ditemukan")
		}
		return walletService.Wallet{}, err
	}

	return wallet, nil
}

func (r *GormRepository) GetTransactions(userID int) ([]walletService.WalletTransaction, error) {
	var wallet walletService.Wallet

	err := r.db.
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet tidak ditemukan")
		}
		return nil, err
	}

	var transactions []walletService.WalletTransaction

	err = r.db.
		Where("wallet_id = ?", wallet.ID).
		Order("created_at DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *GormRepository) Withdraw(userID int, amount float64) (walletService.WalletTransaction, error) {
	var result walletService.WalletTransaction
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var wallet walletService.Wallet

		err := tx.Raw(`
			SELECT id, user_id, balance, updated_at
			FROM wallet
			WHERE user_id = ?
			FOR UPDATE
			`, userID).Scan(&wallet).Error

		if err != nil {
			return err
		}

		if wallet.ID == 0 {
			return errors.New("wallet tidak ditemukan")
		}

		if wallet.Balance < amount {
			return errors.New("saldo tidak mencukupi")
		}

		if err := tx.Exec(`
		UPDATE wallet
		SET balance = balance - ?, updated_at = NOW()
		WHERE id = ?
		`, amount, wallet.ID).Error; err != nil {
			return err
		}

		transaction := walletService.WalletTransaction{
			WalletID:      wallet.ID,
			Type:          "withdraw",
			Amount:        amount,
			ReferenceType: "withdraw",
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		result = transaction
		return nil
	})

	return result, err
}
