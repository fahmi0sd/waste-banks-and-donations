package category

import (
	"github.com/fahmi0sd/waste-banks-and-donations/service/category"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

var _ category.Repository = (*GormRepository)(nil)

type categoryRow struct {
	ID               int
	Name             string
	ParentCategoryID *int `gorm:"column:parent_category_id"`
	Unit             string
	IsActive         bool `gorm:"column:is_active"`
}

func toCategory(r categoryRow) category.Category {
	return category.Category{
		ID: r.ID, Name: r.Name, ParentCategoryID: r.ParentCategoryID,
		Unit: r.Unit, IsActive: r.IsActive,
	}
}

const selectCols = "id, name, parent_category_id, unit, is_active"

func (r *GormRepository) Create(c category.Category) (category.Category, error) {
	values := map[string]interface{}{
		"name":      c.Name,
		"unit":      c.Unit,
		"is_active": c.IsActive,
	}
	if c.ParentCategoryID != nil {
		values["parent_category_id"] = *c.ParentCategoryID
	}

	if err := r.db.Table("waste_category").Create(values).Error; err != nil {
		return category.Category{}, err
	}

	var res categoryRow
	if err := r.db.Table("waste_category").Select(selectCols).
		Where("name = ?", c.Name).Order("id DESC").Take(&res).Error; err != nil {
		return category.Category{}, err
	}
	return toCategory(res), nil
}

func (r *GormRepository) List() ([]category.Category, error) {
	var rows []categoryRow
	if err := r.db.Table("waste_category").Select(selectCols).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]category.Category, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCategory(row))
	}
	return list, nil
}

func (r *GormRepository) GetByID(id int) (category.Category, bool, error) {
	var res categoryRow
	err := r.db.Table("waste_category").Select(selectCols).Where("id = ?", id).Take(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return category.Category{}, false, nil
		}
		return category.Category{}, false, err
	}
	return toCategory(res), true, nil
}

func (r *GormRepository) Update(id int, name *string, parentCategoryID *int, unit *string, isActive *bool) error {
	values := map[string]interface{}{}
	if name != nil {
		values["name"] = *name
	}
	if parentCategoryID != nil {
		values["parent_category_id"] = *parentCategoryID
	}
	if unit != nil {
		values["unit"] = *unit
	}
	if isActive != nil {
		values["is_active"] = *isActive
	}
	if len(values) == 0 {
		return nil
	}
	return r.db.Table("waste_category").Where("id = ?", id).Updates(values).Error
}

func (r *GormRepository) Delete(id int) error {
	return r.db.Table("waste_category").Where("id = ?", id).Delete(nil).Error
}

func (r *GormRepository) RoleOf(userID int) (string, error) {
	var role string
	err := r.db.Table("users").Select("role").Where("id = ?", userID).Take(&role).Error
	return role, err
}
