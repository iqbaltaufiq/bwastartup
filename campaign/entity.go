package campaign

import "time"

type Campaign struct {
	Id               int
	UserId           int
	Name             string
	ShortDescription string
	Description      string
	GoalAmount       string
	CurrentAmount    string
	Perks            string
	BackerCount      int
	Slug             string
	CampaignImages   []CampaignImage
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CampaignImage struct {
	Id         int
	CampaignId int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
