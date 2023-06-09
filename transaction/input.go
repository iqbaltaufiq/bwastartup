package transaction

import "bwastartup/user"

type GetTxByCampaignIdInput struct {
	Id   int `uri:"id" binding:"required"`
	User user.User
}

type CreateTransactionInput struct {
	CampaignId int `json:"campaign_id" binding:"required"`
	Amount     int `json:"amount" binding:"required"`
	User       user.User
}

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	OrderId           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
}
