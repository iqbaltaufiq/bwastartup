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
