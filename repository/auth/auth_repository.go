package auth

import (
	"github.com/fahmi0sd/waste-banks-and-donations/service/auth"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

var _ auth.Repository = (*GormRepository)(nil)

type userRow struct {
	ID           int
	Name         string
	Email        string
	Phone        string
	PasswordHash string `gorm:"column:password_hash"`
	Role         string
}

func (r *GormRepository) Create(name, email, phone, passwordHash string) (auth.UserRecord, error) {
	values := map[string]interface{}{
		"name":          name,
		"email":         email,
		"password_hash": passwordHash,
		"role":          "user",
	}
	if phone != "" {
		values["phone"] = phone
	}

	if err := r.db.Table("users").Create(values).Error; err != nil {
		return auth.UserRecord{}, err
	}

	var res userRow
	if err := r.db.Table("users").Where("email = ?", email).Take(&res).Error; err != nil {
		return auth.UserRecord{}, err
	}

	return auth.UserRecord{
		ID: res.ID, Name: res.Name, Email: res.Email, Phone: res.Phone,
		PasswordHash: res.PasswordHash, Role: res.Role,
	}, nil
}

func (r *GormRepository) FindByEmail(email string) (auth.UserRecord, bool, error) {
	var res userRow
	err := r.db.Table("users").Where("email = ?", email).Take(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return auth.UserRecord{}, false, nil
		}
		return auth.UserRecord{}, false, err
	}
	return auth.UserRecord{
		ID: res.ID, Name: res.Name, Email: res.Email, Phone: res.Phone,
		PasswordHash: res.PasswordHash, Role: res.Role,
	}, true, nil
}

func (r *GormRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Table("users").Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
