package location

import (
	"time"

	"github.com/fahmi0sd/waste-banks-and-donations/service/location"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

var _ location.Repository = (*GormRepository)(nil)

type locationRow struct {
	ID        int
	Name      string
	Address   string
	Latitude  *float64
	Longitude *float64
	IsOpen    bool      `gorm:"column:is_open"`
	OpenTime  string    `gorm:"column:open_time"`
	CloseTime string    `gorm:"column:close_time"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func toLocation(r locationRow) location.Location {
	return location.Location{
		ID: r.ID, Name: r.Name, Address: r.Address,
		Latitude: r.Latitude, Longitude: r.Longitude,
		IsOpen: r.IsOpen, OpenTime: r.OpenTime, CloseTime: r.CloseTime,
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
	}
}

const selectCols = "id, name, COALESCE(address, '') AS address, latitude, longitude, is_open, open_time, close_time, created_at"

func (r *GormRepository) Create(l location.Location) (location.Location, error) {
	values := map[string]interface{}{
		"name":       l.Name,
		"is_open":    l.IsOpen,
		"open_time":  l.OpenTime,
		"close_time": l.CloseTime,
	}
	if l.Address != "" {
		values["address"] = l.Address
	}
	if l.Latitude != nil {
		values["latitude"] = *l.Latitude
	}
	if l.Longitude != nil {
		values["longitude"] = *l.Longitude
	}

	if err := r.db.Table("locations").Create(values).Error; err != nil {
		return location.Location{}, err
	}

	var res locationRow
	if err := r.db.Table("locations").Select(selectCols).
		Where("name = ?", l.Name).Order("id DESC").Take(&res).Error; err != nil {
		return location.Location{}, err
	}
	return toLocation(res), nil
}

func (r *GormRepository) List() ([]location.Location, error) {
	var rows []locationRow
	if err := r.db.Table("locations").Select(selectCols).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]location.Location, 0, len(rows))
	for _, row := range rows {
		list = append(list, toLocation(row))
	}
	return list, nil
}

func (r *GormRepository) GetByID(id int) (location.Location, bool, error) {
	var res locationRow
	err := r.db.Table("locations").Select(selectCols).Where("id = ?", id).Take(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return location.Location{}, false, nil
		}
		return location.Location{}, false, err
	}
	return toLocation(res), true, nil
}

func (r *GormRepository) Update(id int, name, address *string, latitude, longitude *float64, openTime, closeTime *string, isOpen *bool) error {
	values := map[string]interface{}{}
	if name != nil {
		values["name"] = *name
	}
	if address != nil {
		values["address"] = *address
	}
	if latitude != nil {
		values["latitude"] = *latitude
	}
	if longitude != nil {
		values["longitude"] = *longitude
	}
	if openTime != nil {
		values["open_time"] = *openTime
	}
	if closeTime != nil {
		values["close_time"] = *closeTime
	}
	if isOpen != nil {
		values["is_open"] = *isOpen
	}
	if len(values) == 0 {
		return nil
	}
	return r.db.Table("locations").Where("id = ?", id).Updates(values).Error
}

func (r *GormRepository) Delete(id int) error {
	return r.db.Table("locations").Where("id = ?", id).Delete(nil).Error
}

func (r *GormRepository) RoleOf(userID int) (string, error) {
	var role string
	err := r.db.Table("users").Select("role").Where("id = ?", userID).Take(&role).Error
	return role, err
}
