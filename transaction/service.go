package transaction

import (
	"bwastartup/campaign"
	"bwastartup/payment"
	"errors"
)

type Service interface {
	GetTransactionsByCampaignId(input GetTxByCampaignIdInput) ([]Transaction, error)
	GetTransactionsByUserId(userId int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
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

func (s *service) GetTransactionsByUserId(userId int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserId(userId)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.CampaignId = input.CampaignId
	transaction.Amount = input.Amount
	transaction.UserId = input.User.Id
	transaction.Status = "pending"

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{
		Id:     newTransaction.Id,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}