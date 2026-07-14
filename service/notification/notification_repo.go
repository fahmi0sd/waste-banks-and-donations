package notification

type Repository interface {
	UserContact(userID int) (name string, email string, err error)
	Log(n NotificationLog) error
	ListByUser(userID int) ([]NotificationLog, error)
}
