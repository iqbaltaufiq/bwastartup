package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserId(userId int) ([]Campaign, error)
	FindById(campaignId int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

// fetch all campaigns in database
func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign

	errFind := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if errFind != nil {
		return campaigns, errFind
	}

	return campaigns, nil
}

// fetch all campaigns that is created by specific user
func (r *repository) FindByUserId(userId int) ([]Campaign, error) {
	var campaigns []Campaign

	errFind := r.db.Where("user_id = ?", userId).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if errFind != nil {
		return campaigns, errFind
	}

	return campaigns, nil
}

// fetch one campaign with detailed information
// also preload the campaign images
// and the person who created this campaign
func (r *repository) FindById(campaignId int) (Campaign, error) {
	var campaign Campaign

	err := r.db.Preload("User").Preload("CampaignImages").First(&campaign, campaignId).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

// create new campaign to database
func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
