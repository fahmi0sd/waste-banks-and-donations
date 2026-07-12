package calculator

import (
	"github.com/fahmi0sd/waste-banks-and-donations/service/calculator"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) ActivePrice(categoryID int) (categoryName string, pricePerUnit float64, found bool, err error) {
	type row struct {
		Name         string
		PricePerUnit float64
	}
	var res row
	dbErr := r.db.Table("waste_category AS c").
		Select("c.name, p.price_per_unit").
		Joins("JOIN waste_price p ON p.category_id = c.id AND p.effective_until IS NULL").
		Where("c.id = ? AND c.is_active = true", categoryID).
		Take(&res).Error

	if dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			return "", 0, false, nil
		}
		return "", 0, false, dbErr
	}
	return res.Name, res.PricePerUnit, true, nil
}

var _ calculator.Repository = (*GormRepository)(nil)
