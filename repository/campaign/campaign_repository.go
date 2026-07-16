package campaign

import (
	"errors"

	campaignService "github.com/fahmi0sd/waste-banks-and-donations/service/campaign"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

var _ campaignService.Repository = (*GormRepository)(nil)

func (r *GormRepository) List() ([]campaignService.Campaign, error) {
	var campaigns []campaignService.Campaign

	err := r.db.Where("status <> ?", campaignService.StatusClosed).Order("created_at DESC").Find(&campaigns).Error
	if err != nil {
		return nil, err
	}

	return campaigns, nil
}

func (r *GormRepository) Get(id int) (campaignService.Campaign, error) {
	var campaign campaignService.Campaign

	err := r.db.
		Where("id = ?", id).
		First(&campaign).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campaignService.Campaign{}, errors.New("campaign tidak ditemukan")
		}
		return campaignService.Campaign{}, err
	}
	return campaign, nil
}

func (r *GormRepository) Create(c campaignService.Campaign) (campaignService.Campaign, error) {
	if err := r.db.Create(&c).Error; err != nil {
		return campaignService.Campaign{}, err
	}
	return c, nil
}

func (r *GormRepository) Update(c campaignService.Campaign) (campaignService.Campaign, error) {
	result := r.db.
		Model(&campaignService.Campaign{}).
		Where("id = ?", c.ID).
		Updates(map[string]interface{}{
			"title":          c.Title,
			"description":    c.Description,
			"organizer":      c.Organizer,
			"target_amount":  c.TargetAmount,
			"status":         c.Status,
			"current_amount": c.CurrentAmount,
		})
	if result.Error != nil {
		return campaignService.Campaign{}, result.Error
	}

	if result.RowsAffected == 0 {
		return campaignService.Campaign{}, errors.New("campaign tidak ditemukan")
	}
	return c, nil
}

func (r *GormRepository) Delete(id int) error {
	result := r.db.
		Model(&campaignService.Campaign{}).
		Where("id = ?", id).
		Update("status", campaignService.StatusClosed)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("campaign tidak ditemukan")
	}
	return nil
}
