package transaction

import "bwastartup/user"

type GetTxByCampaignIdInput struct {
	Id   int `uri:"id" binding:"required"`
	User user.User
}
