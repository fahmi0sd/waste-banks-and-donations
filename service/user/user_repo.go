package user

type Repository interface {
	FindByID(id int) (Profile, bool, error)
	UpdateProfile(id int, name, phone *string, locationID *int) error

	// Untuk keperluan CRUD akun oleh master_admin
	RoleOf(userID int) (string, error)
	ListAccounts() ([]Profile, error)
	CreateAccount(name, email, phone, passwordHash, role string, locationID *int) (Profile, error)
	EmailExists(email string) (bool, error)
	UpdateAccount(id int, name, phone, role *string, locationID *int) error
	DeleteAccount(id int) error
}
