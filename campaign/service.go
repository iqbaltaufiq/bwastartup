package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userId int) ([]Campaign, error)
	GetCampaign(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(campaignId GetCampaignDetailInput, payload CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// fetch campaign from database
// if client sent user_id,
// then fetch all campaigns that belong to that user_id
// if not, then fetch all campaigns from database
func (s *service) GetCampaigns(userId int) ([]Campaign, error) {
	if userId != 0 {
		campaigns, err := s.repository.FindByUserId(userId)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s *service) GetCampaign(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindById(input.Id)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserId = input.User.Id

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.Id)
	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *service) UpdateCampaign(campaignId GetCampaignDetailInput, payload CreateCampaignInput) (Campaign, error) {
	campaign, errFind := s.repository.FindById(campaignId.Id)
	if errFind != nil {
		return campaign, errFind
	}

	if payload.User.Id != campaign.UserId {
		return Campaign{}, errors.New("Unauthorized")
	}

	campaign.Name = payload.Name
	campaign.ShortDescription = payload.ShortDescription
	campaign.Description = payload.Description
	campaign.GoalAmount = payload.GoalAmount
	campaign.Perks = payload.Perks

	updatedCampaign, errUpdate := s.repository.Update(campaign)
	if errUpdate != nil {
		return updatedCampaign, errUpdate
	}

	return updatedCampaign, nil
}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	campaign, errFind := s.repository.FindById(input.CampaignId)
	if errFind != nil {
		return CampaignImage{}, errFind
	}

	if campaign.UserId != input.User.Id {
		return CampaignImage{}, errors.New("Unauthorized")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1

		_, errMark := s.repository.MarkImageAsNonPrimary(input.CampaignId)
		if errMark != nil {
			return CampaignImage{}, errMark
		}
	}

	image := CampaignImage{}
	image.CampaignId = input.CampaignId
	image.IsPrimary = isPrimary
	image.FileName = fileLocation

	newImage, errCreate := s.repository.CreateImage(image)
	if errCreate != nil {
		return newImage, errCreate
	}

	return newImage, nil
}
