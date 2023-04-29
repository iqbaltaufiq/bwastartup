package campaign

import (
	"bwastartup/user"
	"time"

	"github.com/leekchan/accounting"
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

func (c *Campaign) GoalAmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(c.GoalAmount)
}

type CampaignImage struct {
	Id         int
	CampaignId int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
