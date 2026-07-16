package campaign

import (
	"errors"
	"log/slog"
	"time"

	"github.com/fahmi0sd/waste-banks-and-donations/pkg"
	walletService "github.com/fahmi0sd/waste-banks-and-donations/service/wallet"
	"gorm.io/gorm"
)

type Service interface {
	List() ([]Campaign, error)
	Get(id int) (Campaign, error)
	Create(masterAdminID int, req CreateRequest) (Campaign, error)
	Update(masterAdminID int, id int, req UpdateRequest) (Campaign, error)
	Delete(masterAdminID int, id int) error
	CreateUpdate(masterAdminID int, campaignID int, req CreateUpdateRequest) (CampaignUpdate, error)
	ListUpdate(campaignID int) ([]CampaignUpdate, error)
	Donate(userID int, campaignID int, req DonateRequest) error
	MyDonations(userID int) ([]Donation, error)
}

type NotificationSender interface {
	NotifyDonationUpdate(userID int, campaignTitle, updateContent string) error
}

type service struct {
	logger   *slog.Logger
	repo     Repository
	db       *gorm.DB
	notifier NotificationSender
}

func NewService(logger *slog.Logger, repo Repository, db *gorm.DB, notifier NotificationSender) Service {
	return &service{
		logger:   logger,
		repo:     repo,
		db:       db,
		notifier: notifier,
	}
}

func (s *service) List() ([]Campaign, error) {
	return s.repo.List()
}

func (s *service) Get(id int) (Campaign, error) {
	return s.repo.Get(id)
}

func (s *service) Create(masterAdminID int, req CreateRequest) (Campaign, error) {
	identity, err := pkg.GetUserIdentity(s.db, masterAdminID)
	if err != nil {
		return Campaign{}, errors.New("gagal memverifikasi user")
	}

	if identity.Role != "master_admin" {
		return Campaign{}, errors.New("akses ditolak")
	}

	if err := req.Validate(); err != nil {
		return Campaign{}, err
	}

	status := req.Status
	if status == "" {
		status = StatusDraft
	}

	campaign := Campaign{
		Title:         req.Title,
		Description:   req.Description,
		Organizer:     req.Organizer,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: 0,
		Status:        status,
		CreatedBy:     &masterAdminID,
	}

	result, err := s.repo.Create(campaign)
	if err != nil {
		s.logger.Error("failed create campaign", "error", err)
		return Campaign{}, err
	}

	return result, nil
}

func (s *service) Update(masterAdminID int, id int, req UpdateRequest) (Campaign, error) {
	identity, err := pkg.GetUserIdentity(s.db, masterAdminID)
	if err != nil {
		return Campaign{}, errors.New("gagal memverifikasi user")
	}

	if identity.Role != "master_admin" {
		return Campaign{}, errors.New("akses ditolak")
	}

	campaign, err := s.repo.Get(id)
	if err != nil {
		return Campaign{}, err
	}

	if req.Title != "" {
		campaign.Title = req.Title
	}

	if req.Description != "" {
		campaign.Description = req.Description
	}

	if req.Organizer != "" {
		campaign.Organizer = req.Organizer
	}

	if req.TargetAmount > 0 {
		campaign.TargetAmount = req.TargetAmount
	}

	if req.Status != "" {
		campaign.Status = req.Status
	}

	result, err := s.repo.Update(campaign)
	if err != nil {
		s.logger.Error("failed update campaign", "error", err)
		return Campaign{}, err
	}
	return result, nil
}

func (s *service) Delete(masterAdminID int, id int) error {
	identity, err := pkg.GetUserIdentity(s.db, masterAdminID)
	if err != nil {
		return errors.New("gagal memverifikasi user")
	}

	if identity.Role != "master_admin" {
		return errors.New("akses ditolak")
	}
	return s.repo.Delete(id)
}

func (s *service) ListUpdate(campaignID int) ([]CampaignUpdate, error) {
	return s.repo.ListUpdate(campaignID)
}

func (s *service) CreateUpdate(masterAdminID int, campaignID int, req CreateUpdateRequest) (CampaignUpdate, error) {
	identity, err := pkg.GetUserIdentity(s.db, masterAdminID)
	if err != nil {
		return CampaignUpdate{},
			errors.New("gagal memverifikasi user")
	}

	if identity.Role != "master_admin" {
		return CampaignUpdate{}, errors.New("akses ditolak")
	}

	if err := req.Validate(); err != nil {
		return CampaignUpdate{}, err
	}

	campaign, err := s.repo.Get(campaignID)
	if err != nil {
		return CampaignUpdate{}, err
	}

	update := CampaignUpdate{
		CampaignID:    campaignID,
		Content:       req.Content,
		ReportFileURL: req.ReportFileURL,
		PublishedAt:   time.Now(),
	}

	result, err := s.repo.CreateUpdate(update)
	if err != nil {
		s.logger.Error("failed create campaign update", "error", err)
		return CampaignUpdate{}, err
	}

	if s.notifier != nil {
		donorIDs, err := s.repo.DonorUserIDs(campaignID)
		if err != nil {
			s.logger.Error("failed to fetch donor list for notification", "error", err, "campaign_id", campaignID)
		} else {
			for _, donorID := range donorIDs {
				if notifErr := s.notifier.NotifyDonationUpdate(donorID, campaign.Title, req.Content); notifErr != nil {
					s.logger.Error("failed to send donation update notification", "error", notifErr, "user_id", donorID, "campaign_id", campaignID)
				}
			}
		}
	}

	return result, nil
}

func (s *service) Donate(userID int, campaignID int, req DonateRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	campaign, err := s.repo.Get(campaignID)
	if err != nil {
		return err
	}
	if campaign.Status != StatusActive {
		return errors.New("campaign ini sedang tidak menerima donasi")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Lock wallet
		wallet, err := s.repo.FindWalletByUserID(tx, userID)
		if err != nil {
			return err
		}

		// cek saldo
		if wallet.Balance < req.Amount {
			return errors.New("saldo wallet tidak mencukupi")
		}

		// pengurangan saldo
		wallet.Balance -= req.Amount
		if err := s.repo.UpdateWallet(tx, wallet); err != nil {
			return err
		}

		// simpan donation
		donation, err := s.repo.CreateDonation(tx, Donation{
			UserID:     userID,
			CampaignID: campaignID,
			Amount:     req.Amount,
			CreatedAt:  time.Now(),
		})

		if err != nil {
			return err
		}

		// update campaign
		if err := s.repo.IncreaseCampaignAmount(tx, campaignID, req.Amount); err != nil {
			return err
		}

		// wallet history
		if err := s.repo.CreateWalletTransaction(tx, walletService.WalletTransaction{
			WalletID:      wallet.ID,
			Type:          walletService.TransactionDonationOut,
			Amount:        req.Amount,
			ReferenceType: "donation",
			ReferenceID:   &donation.ID,
			CreatedAt:     time.Now(),
		},
		); err != nil {
			return err
		}

		s.logger.Info(
			"campaign donation success",
			"user_id", userID,
			"campaign_id", campaignID,
			"amount", req.Amount,
		)

		return nil
	})
}

func (s *service) MyDonations(userID int) ([]Donation, error) {
	return s.repo.MyDonations(userID)
}
