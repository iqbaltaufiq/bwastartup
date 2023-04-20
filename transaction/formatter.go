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
