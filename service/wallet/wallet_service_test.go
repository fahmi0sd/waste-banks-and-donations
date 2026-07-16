package wallet_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/fahmi0sd/waste-banks-and-donations/service/wallet"
	"github.com/fahmi0sd/waste-banks-and-donations/service/wallet/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError + 1}))
}

func TestWithdraw_TolakNominalNolAtauNegatif(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := wallet.NewService(silentLogger(), repo)

	_, err := svc.Withdraw(1, 0)
	assert.Error(t, err, "withdraw dengan amount 0 harus ditolak")

	_, err = svc.Withdraw(1, -5000)
	assert.Error(t, err, "withdraw dengan amount negatif harus ditolak")

	repo.AssertNotCalled(t, "Withdraw", mock.Anything, mock.Anything)
}

func TestWithdraw_Sukses(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	expected := wallet.WalletTransaction{ID: 1, WalletID: 1, Type: wallet.TransactionWithdraw, Amount: 5000}
	repo.EXPECT().Withdraw(1, 5000.0).Return(expected, nil)

	svc := wallet.NewService(silentLogger(), repo)
	result, err := svc.Withdraw(1, 5000)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWithdraw_MeneruskanErrorDariRepository(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Withdraw(1, 100000.0).Return(wallet.WalletTransaction{}, errors.New("saldo tidak cukup"))

	svc := wallet.NewService(silentLogger(), repo)
	_, err := svc.Withdraw(1, 100000)

	assert.Error(t, err)
	assert.Equal(t, "saldo tidak cukup", err.Error())
}

func TestGetWallet_Sukses(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	expected := wallet.Wallet{ID: 1, UserID: 1, Balance: 17500}
	repo.EXPECT().GetWallet(1).Return(expected, nil)

	svc := wallet.NewService(silentLogger(), repo)
	result, err := svc.GetWallet(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetWallet_MeneruskanErrorDariRepository(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().GetWallet(1).Return(wallet.Wallet{}, errors.New("wallet tidak ditemukan"))

	svc := wallet.NewService(silentLogger(), repo)
	_, err := svc.GetWallet(1)

	assert.Error(t, err)
}

func TestGetTransactions_Sukses(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	expected := []wallet.WalletTransaction{
		{ID: 1, Type: wallet.TransactionCreditWaste, Amount: 17500},
		{ID: 2, Type: wallet.TransactionWithdraw, Amount: -5000},
	}
	repo.EXPECT().GetTransactions(1).Return(expected, nil)

	svc := wallet.NewService(silentLogger(), repo)
	result, err := svc.GetTransactions(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Len(t, result, 2)
}
