package transaction

import (
	"bwastartup/campaign"
	"errors"
)

type Service interface {
	GetTransactionsByCampaignId(input GetTxByCampaignIdInput) ([]Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

// Fetch all transactions in a campaign
// This should ONLY be visible to the campaign's owner
func (s *service) GetTransactionsByCampaignId(input GetTxByCampaignIdInput) ([]Transaction, error) {
	campaign, errFind := s.campaignRepository.FindById(input.Id)
	if errFind != nil {
		return []Transaction{}, errFind
	}

	if campaign.UserId != input.User.Id {
		return []Transaction{}, errors.New("you are not the owner of this campaign")
	}

	transactions, errGet := s.repository.GetByCampaignId(input.Id)
	if errGet != nil {
		return transactions, errGet
	}

	return transactions, nil
}
