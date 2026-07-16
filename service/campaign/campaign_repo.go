package campaign

import (
	walletService "github.com/fahmi0sd/waste-banks-and-donations/service/wallet"
	"gorm.io/gorm"
)

type Repository interface {
	List() ([]Campaign, error)
	Get(id int) (Campaign, error)
	Create(c Campaign) (Campaign, error)
	Update(c Campaign) (Campaign, error)
	Delete(id int) error

	CreateUpdate(update CampaignUpdate) (CampaignUpdate, error)
	ListUpdate(campaignID int) ([]CampaignUpdate, error)
	FindWalletByUserID(tx *gorm.DB, userID int) (walletService.Wallet, error)
	UpdateWallet(tx *gorm.DB, wallet walletService.Wallet) error
	CreateWalletTransaction(tx *gorm.DB, transaction walletService.WalletTransaction) error
	CreateDonation(tx *gorm.DB, donation Donation) (Donation, error)
	IncreaseCampaignAmount(tx *gorm.DB, campaignID int, amount float64) error
	MyDonations(userID int) ([]Donation, error)

	DonorUserIDs(campaignID int) ([]int, error)
}
