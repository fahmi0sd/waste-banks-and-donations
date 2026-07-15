package price

import (
	"time"

	"github.com/fahmi0sd/waste-banks-and-donations/service/price"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

var _ price.Repository = (*GormRepository)(nil)

type priceRow struct {
	ID             int
	CategoryID     int        `gorm:"column:category_id"`
	PricePerUnit   float64    `gorm:"column:price_per_unit"`
	EffectiveFrom  time.Time  `gorm:"column:effective_from"`
	EffectiveUntil *time.Time `gorm:"column:effective_until"`
	SetBy          *int       `gorm:"column:set_by"`
}

func toPrice(r priceRow) price.Price {
	p := price.Price{
		ID: r.ID, CategoryID: r.CategoryID, PricePerUnit: r.PricePerUnit,
		EffectiveFrom: r.EffectiveFrom.Format(time.RFC3339), SetBy: r.SetBy,
	}
	if r.EffectiveUntil != nil {
		s := r.EffectiveUntil.Format(time.RFC3339)
		p.EffectiveUntil = &s
	}
	return p
}

const selectCols = "id, category_id, price_per_unit, effective_from, effective_until, set_by"

func (r *GormRepository) CategoryExists(categoryID int) (bool, error) {
	var count int64
	err := r.db.Table("waste_category").Where("id = ? AND is_active = true", categoryID).Count(&count).Error
	return count > 0, err
}

func (r *GormRepository) CloseActivePrice(categoryID int) error {
	return r.db.Table("waste_price").
		Where("category_id = ? AND effective_until IS NULL", categoryID).
		Update("effective_until", time.Now()).Error
}

func (r *GormRepository) InsertPrice(categoryID int, pricePerUnit float64, setBy int) (price.Price, error) {
	now := time.Now()
	values := map[string]interface{}{
		"category_id":    categoryID,
		"price_per_unit": pricePerUnit,
		"effective_from": now,
		"set_by":         setBy,
	}
	if err := r.db.Table("waste_price").Create(values).Error; err != nil {
		return price.Price{}, err
	}

	var res priceRow
	if err := r.db.Table("waste_price").Select(selectCols).
		Where("category_id = ? AND effective_until IS NULL", categoryID).
		Order("id DESC").Take(&res).Error; err != nil {
		return price.Price{}, err
	}
	return toPrice(res), nil
}

func (r *GormRepository) ActivePrice(categoryID int) (price.Price, bool, error) {
	var res priceRow
	err := r.db.Table("waste_price").Select(selectCols).
		Where("category_id = ? AND effective_until IS NULL", categoryID).
		Take(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return price.Price{}, false, nil
		}
		return price.Price{}, false, err
	}
	return toPrice(res), true, nil
}

func (r *GormRepository) History(categoryID int) ([]price.Price, error) {
	var rows []priceRow
	err := r.db.Table("waste_price").Select(selectCols).
		Where("category_id = ?", categoryID).
		Order("effective_from DESC").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	list := make([]price.Price, 0, len(rows))
	for _, row := range rows {
		list = append(list, toPrice(row))
	}
	return list, nil
}

func (r *GormRepository) RoleOf(userID int) (string, error) {
	var role string
	err := r.db.Table("users").Select("role").Where("id = ?", userID).Take(&role).Error
	return role, err
}
