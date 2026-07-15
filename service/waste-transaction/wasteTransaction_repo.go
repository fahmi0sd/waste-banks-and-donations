package wastetransaction

type ItemInput struct {
	CategoryID int
	Weight     float64
}

type CreateInput struct {
	QueueID    int
	UserID     int
	LocationID int
	AdminID    int
	Items      []ItemInput
}

type QueueInfo struct {
	UserID     int
	LocationID int
	Status     string
	Found      bool
}

type Repository interface {
	QueueInfo(queueID int) (QueueInfo, error)
	CreateWithEffects(input CreateInput) (WasteTransaction, error)
	GetByID(id int) (WasteTransaction, error)
	GetDetailsByTransactionID(transactionID int) ([]WasteTransactionDetail, error)
	ListByUser(userID int) ([]WasteTransaction, error)
	ListByLocation(locationID int) ([]WasteTransaction, error)
	CancelWithEffects(transactionID int) (WasteTransaction, error)
}
