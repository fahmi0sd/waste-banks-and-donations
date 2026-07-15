package location

type Repository interface {
	Create(l Location) (Location, error)
	List() ([]Location, error)
	GetByID(id int) (Location, bool, error)
	Update(id int, name, address *string, latitude, longitude *float64, openTime, closeTime *string, isOpen *bool) error
	Delete(id int) error
	SetOpen(id int, isOpen bool) error

	RoleAndLocation(userID int) (role string, locationID *int, err error)
}
