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

func (c *UserController) Register(ctx *gin.Context) {
	// TODO: create internal/dto/user_entity.go and move registerInput dto there
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Register(input.Email, input.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered"})
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

// TODO: Move login to auth_controller.go (create new file)
func (c *UserController) Login(ctx *gin.Context) {
	// TODO: create internal/dto/user_entity.go and move loginInput dto there
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Login(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *UserController) SetupRoutes(router *gin.Engine) {
	//TODO: group routes under /api/v1
	router.POST("/users", c.Register)
	router.GET("/users/:email", c.GetUser)
	//TODO: move to auth_controller
	router.POST("/login", c.Login)
}
