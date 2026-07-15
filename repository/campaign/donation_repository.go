package campaign

import (
	"errors"

	campaignService "github.com/fahmi0sd/waste-banks-and-donations/service/campaign"
	walletService "github.com/fahmi0sd/waste-banks-and-donations/service/wallet"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Lock wallet user
func (r *GormRepository) FindWalletByUserID(tx *gorm.DB, userID int) (walletService.Wallet, error) {
	var wallet walletService.Wallet

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return walletService.Wallet{}, errors.New("wallet tidak ditemukan")
		}
		return walletService.Wallet{}, err
	}
	return wallet, nil
}

// Update saldo wallet
func (r *GormRepository) UpdateWallet(tx *gorm.DB, wallet walletService.Wallet) error {
	return tx.Model(&walletService.Wallet{}).Where("id = ?", wallet.ID).Update("balance", wallet.Balance).Error
}

// Insert wallet transaction
func (r *GormRepository) CreateWalletTransaction(tx *gorm.DB, transaction walletService.WalletTransaction) error {
	return tx.Create(&transaction).Error
}

// Insert donation
func (r *GormRepository) CreateDonation(tx *gorm.DB, donation campaignService.Donation) (campaignService.Donation, error) {
	if err := tx.Create(&donation).Error; err != nil {
		return campaignService.Donation{}, err
	}
	return donation, nil
}

// Tambah current_amount campaign
func (r *GormRepository) IncreaseCampaignAmount(tx *gorm.DB, campaignID int, amount float64) error {
	return tx.Model(&campaignService.Campaign{}).Where("id = ?", campaignID).Update("current_amount", gorm.Expr("current_amount + ?", amount)).Error
}

// History donasi user
func (r *GormRepository) MyDonations(userID int) ([]campaignService.Donation, error) {
	var donations []campaignService.Donation

	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&donations).Error
	if err != nil {
		return nil, err
	}
	return donations, nil
}
