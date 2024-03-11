package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	amodels "server/internal/modules/auth/models"
	"server/internal/modules/chat/models"
)

func (c *ChatController) Register(ctx *gin.Context) {
	req := models.RegisterRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	u := amodels.User{
		Username: req.Username,
		Password: req.Password,
	}
	err = c.Auth.Register(u)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, models.RegisterResponse{Message: "Registration sucsess"})

}
func (c *ChatController) Login(ctx *gin.Context) {
	req := models.LoginRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	u := amodels.User{
		Username: req.Username,
		Password: req.Password,
	}
	token, id, username, err := c.Auth.Login(u)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, models.LoginResponse{
		Token:    fmt.Sprintf("%s", token),
		ID:       id,
		Username: username,
	})
}
