package transaction

import "gorm.io/gorm"

type Repository interface {
	GetByCampaignId(campaignId int) ([]Transaction, error)
	GetByUserId(userId int) ([]Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignId(campaignId int) ([]Transaction, error) {
	var transactions []Transaction
	errFind := r.db.Preload("User").Where("campaign_id = ?", campaignId).Order("id DESC").Find(&transactions).Error
	if errFind != nil {
		return transactions, errFind
	}

	return transactions, nil
}

func (r *repository) GetByUserId(userId int) ([]Transaction, error) {
	var transactions []Transaction

	errFind := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userId).Order("id DESC").Find(&transactions).Error
	if errFind != nil {
		return transactions, errFind
	}

	return transactions, nil
}
