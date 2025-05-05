package controller

import (
	"net/http"

	"github.com/EngenMe/api-frontend-team/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service}
}

func (c *UserController) GetUser(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := c.service.GetUserByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// TODO: use user response dto instead of user object
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) SetupUserRoutes(router *gin.RouterGroup) {
	router.GET("/users/:email", c.GetUser)
}
