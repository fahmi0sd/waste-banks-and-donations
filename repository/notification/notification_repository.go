package notification

import (
	"github.com/fahmi0sd/waste-banks-and-donations/service/notification"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

var _ notification.Repository = (*GormRepository)(nil)

func (r *GormRepository) UserContact(userID int) (name string, email string, err error) {
	type row struct {
		Name  string
		Email string
	}
	var res row
	dbErr := r.db.Table("users").
		Select("name, email").
		Where("id = ?", userID).
		Take(&res).Error
	if dbErr != nil {
		return "", "", dbErr
	}
	return res.Name, res.Email, nil
}

func (r *GormRepository) Log(n notification.NotificationLog) error {
	return r.db.Create(&n).Error
}

func (r *GormRepository) ListByUser(userID int) ([]notification.NotificationLog, error) {
	var list []notification.NotificationLog
	err := r.db.Where("user_id = ?", userID).Order("sent_at DESC").Find(&list).Error
	return list, err
}
