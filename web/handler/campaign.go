package handler

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService     user.Service
}

func NewCampaignHandler(campaignService campaign.Service, userService user.Service) *campaignHandler {
	return &campaignHandler{campaignService, userService}
}

func (h *campaignHandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

func (h *campaignHandler) New(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	input := campaign.FormCreateCampaignInput{}
	input.Users = users

	c.HTML(http.StatusOK, "campaign_new.html", input)
}

func (h *campaignHandler) Create(c *gin.Context) {
	var input campaign.FormCreateCampaignInput

	bindErr := c.ShouldBind(&input)
	if bindErr != nil {
		users, fetchErr := h.userService.GetAllUsers()
		if fetchErr != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
		input.Users = users
		input.Error = bindErr

		c.HTML(http.StatusOK, "campaign_new.html", input)
		return
	}

	user, fetchErr := h.userService.GetUserById(input.UserId)
	if fetchErr != nil {
		c.HTML(http.StatusNotFound, "error.html", nil)
		return
	}

	createCampaignInput := campaign.CreateCampaignInput{}
	createCampaignInput.Name = input.Name
	createCampaignInput.Description = input.Description
	createCampaignInput.GoalAmount = input.GoalAmount
	createCampaignInput.ShortDescription = input.ShortDescription
	createCampaignInput.Perks = input.Perks
	createCampaignInput.User = user

	_, createErr := h.campaignService.CreateCampaign(createCampaignInput)
	if createErr != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) NewImage(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	c.HTML(http.StatusOK, "campaign_image.html", gin.H{"Id": id})
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	// catch file from form
	// save file to local directory
	// insert new entry into table Campaign Image

	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	campaignImageInput := campaign.CreateCampaignImageInput{}

	file, formErr := c.FormFile("image")
	if formErr != nil {
		c.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}

	campaign, fetchErr := h.campaignService.GetCampaign(campaign.GetCampaignDetailInput{Id: id})
	if fetchErr != nil {
		c.HTML(http.StatusNotFound, "error.html", nil)
		return
	}

	user, fetchErr := h.userService.GetUserById(campaign.UserId)
	if fetchErr != nil {
		c.HTML(http.StatusNotFound, "error.html", nil)
		return
	}

	path := fmt.Sprintf("images/%d-%s", campaign.UserId, file.Filename)

	campaignImageInput.CampaignId = campaign.Id
	campaignImageInput.IsPrimary = true
	campaignImageInput.User = user

	uploadErr := c.SaveUploadedFile(file, path)
	if uploadErr != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	_, saveErr := h.campaignService.SaveCampaignImage(campaignImageInput, path)
	if saveErr != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Edit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	campaignData, fetchErr := h.campaignService.GetCampaign(campaign.GetCampaignDetailInput{Id: id})
	if fetchErr != nil {
		c.HTML(http.StatusNotFound, "error.html", nil)
		return
	}

	input := campaign.FormUpdateCampaignInput{
		Id:               campaignData.Id,
		Name:             campaignData.Name,
		ShortDescription: campaignData.ShortDescription,
		Description:      campaignData.Description,
		GoalAmount:       campaignData.GoalAmount,
		Perks:            campaignData.Perks,
	}

	c.HTML(http.StatusOK, "campaign_edit.html", input)
}

func (h *campaignHandler) Update(c *gin.Context) {
	var input campaign.FormUpdateCampaignInput
	id, _ := strconv.Atoi(c.Param("id"))

	input.Id = id

	bindErr := c.ShouldBind(&input)
	if bindErr != nil {
		input.Error = bindErr
		input.Id = id
		c.HTML(http.StatusInternalServerError, "campaign_edit.html", input)
		return
	}

	campaignData, fetchErr := h.campaignService.GetCampaign(campaign.GetCampaignDetailInput{Id: id})
	if fetchErr != nil {
		c.HTML(http.StatusNotFound, "error.html", nil)
		return
	}

	user, fetchErr := h.userService.GetUserById(campaignData.UserId)
	if fetchErr != nil {
		c.HTML(http.StatusNotFound, "error.html", nil)
		return
	}

	payload := campaign.CreateCampaignInput{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		User:             user,
	}

	_, updateErr := h.campaignService.UpdateCampaign(campaign.GetCampaignDetailInput{Id: id}, payload)
	if updateErr != nil {
		c.HTML(http.StatusNotFound, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}
