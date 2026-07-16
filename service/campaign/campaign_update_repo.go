package campaign

type UpdateRepository interface {
	CreateUpdate(update CampaignUpdate) (CampaignUpdate, error)
	ListUpdate(campaignID int) ([]CampaignUpdate, error)
}
