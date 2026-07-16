package campaign

import (
	campaignService "github.com/fahmi0sd/waste-banks-and-donations/service/campaign"
)

func (r *GormRepository) CreateUpdate(update campaignService.CampaignUpdate) (campaignService.CampaignUpdate, error) {
	err := r.db.
		Create(&update).Error
	if err != nil {
		return campaignService.CampaignUpdate{}, err
	}
	return update, nil
}

func (r *GormRepository) ListUpdate(campaignID int) ([]campaignService.CampaignUpdate, error) {
	var updates []campaignService.CampaignUpdate

	err := r.db.
		Where("campaign_id = ?", campaignID).
		Order("published_at DESC").
		Find(&updates).Error
	if err != nil {
		return nil, err
	}
	return updates, nil
}
