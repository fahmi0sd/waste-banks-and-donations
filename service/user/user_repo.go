package user

type Repository interface {
	FindByID(id int) (Profile, bool, error)
	UpdateProfile(id int, name, phone *string, locationID *int) error
}
