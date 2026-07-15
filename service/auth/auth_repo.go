package auth

type UserRecord struct {
	ID           int
	Name         string
	Email        string
	Phone        string
	PasswordHash string
	Role         string
}

type Repository interface {
	Create(name, email, phone, passwordHash string) (UserRecord, error)
	FindByEmail(email string) (UserRecord, bool, error)
	EmailExists(email string) (bool, error)
}
