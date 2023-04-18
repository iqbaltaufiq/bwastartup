package campaign

import "strings"

type CampaignFormatter struct {
	Id               int    `json:"id"`
	UserId           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

type CampaignDetailUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignDetailImagesFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

type CampaignDetailFormatter struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	ImageURL         string   `json:"image_url"`
	GoalAmount       int      `json:"goal_amount"`
	CurrentAmount    int      `json:"current_amount"`
	UserId           int      `json:"user_id"`
	Slug             string   `json:"slug"`
	Perks            []string `json:"perks"`
	User             CampaignDetailUserFormatter
	Images           []CampaignDetailImagesFormatter
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.Id = campaign.Id
	campaignFormatter.UserId = campaign.UserId
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	campaignsFormatter := []CampaignFormatter{}

	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}

// format response body for detailed campaign response
func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	c := CampaignDetailFormatter{}

	c.Id = campaign.Id
	c.UserId = campaign.UserId
	c.Name = campaign.Name
	c.ShortDescription = campaign.ShortDescription
	c.CurrentAmount = campaign.CurrentAmount
	c.GoalAmount = campaign.GoalAmount
	c.UserId = campaign.UserId
	c.Slug = campaign.Slug
	c.ImageURL = ""

	perks := strings.Split(campaign.Perks, ",")
	c.Perks = perks

	if len(campaign.CampaignImages) > 0 {
		c.ImageURL = campaign.CampaignImages[0].FileName
	}

	u := CampaignDetailUserFormatter{}

	u.Name = campaign.User.Name
	u.ImageURL = campaign.User.AvatarFileName

	images := []CampaignDetailImagesFormatter{}

	for _, image := range campaign.CampaignImages {
		i := CampaignDetailImagesFormatter{}

		i.ImageURL = image.FileName
		isPrimary := false

		if image.IsPrimary == 1 {
			isPrimary = true
		}

		i.IsPrimary = isPrimary
		images = append(images, i)
	}

	c.User = u
	c.Images = images

	return c
}
