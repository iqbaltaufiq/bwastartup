package handler

import (
	"bwastartup/user"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type sessionHandler struct {
	userService user.Service
}

func NewSessionHandler(userService user.Service) *sessionHandler {
	return &sessionHandler{userService}
}

func (h *sessionHandler) New(c *gin.Context) {
	c.HTML(http.StatusOK, "session_new.html", nil)
}

func (h *sessionHandler) Create(c *gin.Context) {
	var input user.UserLoginInput

	bindErr := c.ShouldBind(&input)
	if bindErr != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	user, loginErr := h.userService.Login(input)
	if loginErr != nil || user.Role != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.Id)
	session.Set("user_name", user.Name)
	session.Save()

	c.Redirect(http.StatusFound, "/users")
}

func (h *sessionHandler) Destroy(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/login")
}
