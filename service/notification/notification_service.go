package notification

import (
	"fmt"
	"log/slog"
)

type Mailer interface {
	SendEmail(toEmail, toName, subject, textContent string) error
}

type Service interface {
	NotifyDepositUpdate(userID int, amount float64) error
	NotifyDonationUpdate(userID int, campaignTitle, updateContent string) error
	MyNotifications(userID int) ([]NotificationLog, error)
}

type service struct {
	logger *slog.Logger
	repo   Repository
	mailer Mailer
}

func NewService(logger *slog.Logger, repo Repository, mailer Mailer) Service {
	return &service{logger: logger, repo: repo, mailer: mailer}
}

func (s *service) NotifyDepositUpdate(userID int, amount float64) error {
	name, email, err := s.repo.UserContact(userID)
	if err != nil {
		s.logger.Error("failed to resolve user contact", "error", err, "user_id", userID)
		return err
	}

	subject := "Deposit Bertambah - Bank Sampah & Donasi"
	content := fmt.Sprintf(
		"Halo %s,\n\nDeposit sebesar Rp%.2f baru saja masuk ke wallet kamu dari hasil tukar sampah.\n\nTerima kasih sudah menjaga lingkungan bersama kami!\n\nSalam,\nBank Sampah & Donasi",
		name, amount,
	)

	return s.send(userID, TypeDepositUpdate, email, name, subject, content)
}

func (s *service) NotifyDonationUpdate(userID int, campaignTitle, updateContent string) error {
	name, email, err := s.repo.UserContact(userID)
	if err != nil {
		s.logger.Error("failed to resolve user contact", "error", err, "user_id", userID)
		return err
	}

	subject := "Update Penyaluran Donasi - " + campaignTitle
	content := fmt.Sprintf(
		"Halo %s,\n\nAda update baru untuk kampanye \"%s\" yang kamu dukung:\n\n%s\n\nTerima kasih atas kepedulianmu!\n\nSalam,\nBank Sampah & Donasi",
		name, campaignTitle, updateContent,
	)

	return s.send(userID, TypeDonationUpdate, email, name, subject, content)
}

func (s *service) send(userID int, notifType, email, name, subject, content string) error {
	sendErr := s.mailer.SendEmail(email, name, subject, content)

	status := StatusSent
	if sendErr != nil {
		status = StatusFailed
		s.logger.Error("failed to send email", "error", sendErr, "user_id", userID, "type", notifType)
	}

	logErr := s.repo.Log(NotificationLog{
		UserID:  userID,
		Type:    notifType,
		Channel: ChannelEmail,
		Content: content,
		Status:  status,
	})
	if logErr != nil {
		s.logger.Error("failed to write notification log", "error", logErr, "user_id", userID)
	}

	return sendErr
}

func (s *service) MyNotifications(userID int) ([]NotificationLog, error) {
	list, err := s.repo.ListByUser(userID)
	if err != nil {
		s.logger.Error("failed to list notifications", "error", err, "user_id", userID)
		return nil, err
	}
	return list, nil
}
