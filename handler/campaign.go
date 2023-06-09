package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, errGet := h.service.GetCampaigns(userId)
	if errGet != nil {
		response := helper.APIResponse("Failed to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success fetching campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	errBind := c.ShouldBindUri(&input)
	if errBind != nil {
		response := helper.APIResponse("Failed binding URI.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, errGet := h.service.GetCampaign(input)
	if errGet != nil {
		response := helper.APIResponse("Failed getting campaign.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success getting campaign", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput
	errBind := c.ShouldBindJSON(&input)
	if errBind != nil {
		errors := helper.FormatValidationError(errBind)
		errMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed parsing body request", http.StatusUnprocessableEntity, "error", errMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, errCreate := h.service.CreateCampaign(input)
	if errCreate != nil {
		response := helper.APIResponse("Failed creating new campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success creating new campaign", http.StatusCreated, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusCreated, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput
	errBindURI := c.ShouldBindUri(&input)
	if errBindURI != nil {
		response := helper.APIResponse("Failed to parse URI", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var payload campaign.CreateCampaignInput
	errBindJSON := c.ShouldBindJSON(&payload)
	if errBindJSON != nil {
		response := helper.APIResponse("Failed to parse json", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user := c.MustGet("currentUser").(user.User)
	payload.User = user

	updatedCampaign, errUpdate := h.service.UpdateCampaign(input, payload)
	if errUpdate != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success updating a campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput
	errBind := c.ShouldBind(&input)
	if errBind != nil {
		errVal := helper.FormatValidationError(errBind)
		data := gin.H{
			"is_uploaded": false,
			"errors":      errVal,
		}

		response := helper.APIResponse("Failed binding request", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	file, errFile := c.FormFile("file")
	if errFile != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("File not found", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user := c.MustGet("currentUser").(user.User)
	path := fmt.Sprintf("images/%d-%s", user.Id, file.Filename)
	input.User = user

	errSave := c.SaveUploadedFile(file, path)
	if errSave != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to save file", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, errSave = h.service.SaveCampaignImage(input, path)
	if errSave != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to save file", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Image has been uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
