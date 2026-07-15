package user

import (
	"time"

	"github.com/fahmi0sd/waste-banks-and-donations/service/user"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

var _ user.Repository = (*GormRepository)(nil)

type userRow struct {
	ID         int
	Name       string
	Email      string
	Phone      string
	Role       string
	LocationID *int      `gorm:"column:location_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func toProfile(r userRow) user.Profile {
	return user.Profile{
		ID: r.ID, Name: r.Name, Email: r.Email, Phone: r.Phone,
		Role: r.Role, LocationID: r.LocationID,
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
	}
}

const selectCols = "id, name, email, COALESCE(phone, '') AS phone, role, location_id, created_at"

func (r *GormRepository) FindByID(id int) (user.Profile, bool, error) {
	var res userRow
	err := r.db.Table("users").Select(selectCols).Where("id = ?", id).Take(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user.Profile{}, false, nil
		}
		return user.Profile{}, false, err
	}
	return toProfile(res), true, nil
}

func (r *GormRepository) UpdateProfile(id int, name, phone *string, locationID *int) error {
	values := map[string]interface{}{"updated_at": time.Now()}
	if name != nil {
		values["name"] = *name
	}
	if phone != nil {
		values["phone"] = *phone
	}
	if locationID != nil {
		values["location_id"] = *locationID
	}
	return r.db.Table("users").Where("id = ?", id).Updates(values).Error
}
