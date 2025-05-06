package controller

import (
	"net/http"

	"github.com/EngenMe/api-frontend-team/internal/dto"
	"github.com/EngenMe/api-frontend-team/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{service}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var userDTO dto.LoginRequest

	if err := ctx.ShouldBindJSON(&userDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Login(&userDTO)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *AuthController) Register(ctx *gin.Context) {
	var userDTO dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&userDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Register(&userDTO); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func (c *AuthController) SetupAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", c.Register)
	router.POST("/login", c.Login)
}
