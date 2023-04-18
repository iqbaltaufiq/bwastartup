package campaign

import (
	"bwastartup/user"
	"time"
)

type Campaign struct {
	Id               int
	UserId           int
	Name             string
	ShortDescription string
	Description      string
	GoalAmount       int
	CurrentAmount    int
	Perks            string
	BackerCount      int
	Slug             string
	CampaignImages   []CampaignImage
	User             user.User
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
