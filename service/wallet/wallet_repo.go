package wallet

type Repository interface {

	GetWallet(userID int) (Wallet, error)

	GetTransactions(userID int) ([]WalletTransaction, error)

	Withdraw(userID int, amount float64) (WalletTransaction, error)

}