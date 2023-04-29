package handler

import (
	"bwastartup/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users})
}

func (h *userHandler) NewAccount(c *gin.Context) {
	c.HTML(http.StatusOK, "user_new.html", nil)
}

func (h *userHandler) Create(c *gin.Context) {
	var input user.FormCreateUserInput
	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err.Error()
		c.HTML(http.StatusOK, "user_new.html", input)
		return
	}

	registerInput := user.RegisterUserInput{}
	registerInput.Name = input.Name
	registerInput.Occupation = input.Occupation
	registerInput.Email = input.Email
	registerInput.Password = input.Password

	_, err = h.userService.RegisterUser(registerInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

// show edit user page
func (h *userHandler) Edit(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	payload := user.FormUpdateUserInput{}

	user, err := h.userService.GetUserById(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", nil)
		return
	}

	payload.Id = user.Id
	payload.Name = user.Name
	payload.Email = user.Email
	payload.Occupation = user.Occupation
	payload.Error = nil

	c.HTML(http.StatusOK, "user_edit.html", payload)
}

func (h *userHandler) Update(c *gin.Context) {
	var form user.FormUpdateUserInput

	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	errBind := c.ShouldBind(&form)
	if errBind != nil {
		form.Error = errBind
		c.HTML(http.StatusOK, "user_edit.html", form)
		return
	}

	form.Id = id

	_, err := h.userService.UpdateUser(form)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/users")
}
