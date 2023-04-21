package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetTxByCampaignIdInput

	errBind := c.ShouldBindUri(&input)
	if errBind != nil {
		errVal := helper.FormatValidationError(errBind)
		data := gin.H{"errors": errVal}
		response := helper.APIResponse("Failed binding request", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	input.User = c.MustGet("currentUser").(user.User)

	transactions, errGet := h.service.GetTransactionsByCampaignId(input)
	if errGet != nil {
		response := helper.APIResponse("Failed fetching transactions", http.StatusNotFound, "error", errGet.Error())
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Transactions fetched successfully.", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	user := c.MustGet("currentUser").(user.User)

	transactions, err := h.service.GetTransactionsByUserId(user.Id)
	if err != nil {
		response := helper.APIResponse("Failed fetching transactions", http.StatusNotFound, "error", err.Error())
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Success fetching transactions", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	errBind := c.ShouldBindJSON(&input)
	if errBind != nil {
		errors := helper.FormatValidationError(errBind)
		errMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed parsing body request", http.StatusUnprocessableEntity, "error", errMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	input.User = c.MustGet("currentUser").(user.User)

	newTransaction, errCreate := h.service.CreateTransaction(input)
	if errCreate != nil {
		errMessage := gin.H{"errors": errCreate.Error()}

		response := helper.APIResponse("Failed creating transaction", http.StatusBadRequest, "error", errMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("A new transaction has been created", http.StatusCreated, "success", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusCreated, response)
}
