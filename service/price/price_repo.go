package price

type Repository interface {
	CategoryExists(categoryID int) (bool, error)
	CloseActivePrice(categoryID int) error
	InsertPrice(categoryID int, pricePerUnit float64, setBy int) (Price, error)
	ActivePrice(categoryID int) (Price, bool, error)
	History(categoryID int) ([]Price, error)

	RoleOf(userID int) (string, error)
}
