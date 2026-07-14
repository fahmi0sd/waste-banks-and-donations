package queue

import (
	"github.com/fahmi0sd/waste-banks-and-donations/pkg"
	"github.com/fahmi0sd/waste-banks-and-donations/service/queue"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

var _ queue.Repository = (*GormRepository)(nil)

func (r *GormRepository) Create(q queue.Queue) (queue.Queue, error) {
	if err := r.db.Create(&q).Error; err != nil {
		return queue.Queue{}, err
	}
	return q, nil
}

func (r *GormRepository) GetByID(id int) (queue.Queue, error) {
	var q queue.Queue
	if err := r.db.First(&q, "id = ?", id).Error; err != nil {
		return queue.Queue{}, err
	}
	return q, nil
}

func (r *GormRepository) CountTodayByLocation(locationID int) (int, error) {
	var count int64
	err := r.db.Table("queue").
		Where("location_id = ? AND created_at::date = CURRENT_DATE", locationID).
		Count(&count).Error
	return int(count), err
}

func (r *GormRepository) ListWaitingByLocation(locationID int) ([]queue.QueueWithUser, error) {
	var list []queue.QueueWithUser
	err := r.db.Table("queue AS q").
		Select("q.id, q.user_id, u.name AS user_name, q.queue_number, q.preferred_time, q.status, q.created_at").
		Joins("JOIN users u ON u.id = q.user_id").
		Where("q.location_id = ? AND q.status = ?", locationID, queue.StatusWaiting).
		Order("q.created_at ASC").
		Scan(&list).Error
	return list, err
}

func (r *GormRepository) UpdateStatus(id int, status string, verifiedBy int) error {
	return r.db.Table("queue").
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"verified_by": verifiedBy,
		}).Error
}

func (r *GormRepository) LocationInfo(locationID int) (exists bool, isOpen bool, openTime string, closeTime string, err error) {
	type row struct {
		IsOpen    bool
		OpenTime  string
		CloseTime string
	}
	var res row
	dbErr := r.db.Table("locations").
		Select("is_open, open_time, close_time").
		Where("id = ?", locationID).
		Take(&res).Error

	if dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			return false, false, "", "", nil
		}
		return false, false, "", "", dbErr
	}
	return true, res.IsOpen, res.OpenTime, res.CloseTime, nil
}

func (r *GormRepository) AdminLocation(userID int) (role string, locationID *int, err error) {
	identity, dbErr := pkg.GetUserIdentity(r.db, userID)
	if dbErr != nil {
		return "", nil, dbErr
	}
	return identity.Role, identity.LocationID, nil
}
