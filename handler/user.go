package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput
	errBind := c.ShouldBindJSON(&input)
	if errBind != nil {
		errors := helper.FormatValidationError(errBind)
		errMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed.", http.StatusUnprocessableEntity, "error", errMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, errRegister := h.userService.RegisterUser(input)
	if errRegister != nil {
		response := helper.APIResponse("Register account failed.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, errGen := h.authService.GenerateToken(newUser.Id)
	if errGen != nil {
		response := helper.APIResponse("Failed generating token", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account created successfully.", http.StatusOK, "OK", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.UserLoginInput
	errBind := c.ShouldBindJSON(&input)
	if errBind != nil {
		errors := helper.FormatValidationError(errBind)
		errMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, errLogin := h.userService.Login(input)
	if errLogin != nil {
		errMessage := gin.H{"errors": errLogin.Error()}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, errGen := h.authService.GenerateToken(loggedInUser.Id)
	if errGen != nil {
		response := helper.APIResponse("Login failed.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)
	response := helper.APIResponse("Login success", http.StatusOK, "OK", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput
	if errBind := c.ShouldBindJSON(&input); errBind != nil {
		errors := helper.FormatValidationError(errBind)
		errMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Error at checking email", http.StatusUnprocessableEntity, "error", errMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isAvailable, errCheck := h.userService.IsEmailAvailable(input)
	if errCheck != nil {
		errMessage := gin.H{"errors": errCheck.Error()}
		response := helper.APIResponse("Error at checking email", http.StatusUnprocessableEntity, "error", errMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isAvailable,
	}

	metaMessage := "Email has been registered"
	if isAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, errGet := c.FormFile("avatar")
	if errGet != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to get image from form.", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// id should be taken form JWT token
	dummyId := 1
	path := fmt.Sprintf("images/%d-%s", dummyId, file.Filename)

	errUpload := c.SaveUploadedFile(file, path)
	if errUpload != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload image to server.", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, errSave := h.userService.SaveAvatar(dummyId, path)
	if errSave != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to update data in database.", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Success uploaded file to server", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
