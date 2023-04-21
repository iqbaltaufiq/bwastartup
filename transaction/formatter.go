package transaction

import "time"

type CampaignTransactionFormatter struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{}

	formatter.Id = transaction.Id
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	return formatter
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {
	transactionsFormatter := []CampaignTransactionFormatter{}

	if len(transactions) == 0 {
		return transactionsFormatter
	}

	for _, tx := range transactions {
		formatter := FormatCampaignTransaction(tx)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}

type UserTransactionFormatter struct {
	Id        int                       `json:"id"`
	Amount    int                       `json:"amount"`
	Status    string                    `json:"status"`
	CreatedAt time.Time                 `json:"created_at"`
	Campaign  CampaignInUserTxFormatter `json:"campaign"`
}

type CampaignInUserTxFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{}

	formatter.Id = transaction.Id
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	formatCampaign := CampaignInUserTxFormatter{}
	formatCampaign.Name = transaction.Campaign.Name
	formatCampaign.ImageURL = ""

	if len(transaction.Campaign.CampaignImages) > 0 {
		formatCampaign.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	formatter.Campaign = formatCampaign
	return formatter
}

func FormatUserTransactions(transactions []Transaction) []UserTransactionFormatter {
	formatter := []UserTransactionFormatter{}

	if len(transactions) == 0 {
		return formatter
	}

	for _, tx := range transactions {
		f := FormatUserTransaction(tx)
		formatter = append(formatter, f)
	}

	return formatter
}

type TransactionFormatter struct {
	Id         int    `json:"id"`
	CampaignId int    `json:"campaign_id"`
	UserId     int    `json:"user_id"`
	Amount     int    `json:"amount"`
	Code       string `json:"code"`
	Status     string `json:"status"`
	PaymentURL string `json:"payment_url"`
}

func FormatTransaction(transaction Transaction) TransactionFormatter {
	formatter := TransactionFormatter{}
	formatter.Id = transaction.Id
	formatter.CampaignId = transaction.CampaignId
	formatter.UserId = transaction.UserId
	formatter.Amount = transaction.Amount
	formatter.Code = transaction.Code
	formatter.Status = transaction.Status
	formatter.PaymentURL = transaction.PaymentURL

	return formatter
}
