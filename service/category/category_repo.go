package category

type Repository interface {
	Create(c Category) (Category, error)
	List() ([]Category, error)
	GetByID(id int) (Category, bool, error)
	Update(id int, name *string, parentCategoryID *int, unit *string, isActive *bool) error
	Delete(id int) error

	RoleOf(userID int) (string, error)
}
