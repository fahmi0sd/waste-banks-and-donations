package calculator

type Repository interface {
	ActivePrice(categoryID int) (categoryName string, pricePerUnit float64, found bool, err error)
}
