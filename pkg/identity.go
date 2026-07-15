package pkg

import "gorm.io/gorm"

type UserIdentity struct {
	Role       string
	LocationID *int
}

func GetUserIdentity(db *gorm.DB, userID int) (UserIdentity, error) {
	var res UserIdentity
	err := db.Table("users").
		Select("role, location_id").
		Where("id = ?", userID).
		Take(&res).Error
	return res, err
}
