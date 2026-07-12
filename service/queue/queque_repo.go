package queue

type Repository interface {
	Create(q Queue) (Queue, error)
	GetByID(id int) (Queue, error)
	CountTodayByLocation(locationID int) (int, error)
	ListWaitingByLocation(locationID int) ([]QueueWithUser, error)
	UpdateStatus(id int, status string, verifiedBy int) error
	LocationInfo(locationID int) (exists bool, isOpen bool, openTime string, closeTime string, err error)
	AdminLocation(userID int) (role string, locationID *int, err error)
}
