package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService: userService}
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

	formatter := user.FormatUser(newUser, "tokentokentokentoken")

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

	formatter := user.FormatUser(loggedInUser, "tokentokentokentoken")
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
