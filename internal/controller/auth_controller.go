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
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
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

	idToken, refreshIdToken, err := c.service.Login(&userDTO)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokens := map[string]string{
		"id_token":      idToken,
		"refresh_token": refreshIdToken,
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokens})
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

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var refreshTokenRequest RefreshTokenRequest
	userIdStr := "3" //get userId from inmemory database "redis"

	if err := ctx.ShouldBindJSON(&refreshTokenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idToken, refreshIdToken, err := c.service.RefreshToken(userIdStr, refreshTokenRequest.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	tokens := map[string]string{
		"id_token":      idToken,
		"refresh_token": refreshIdToken,
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokens})
}

func (c *AuthController) SetupAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", c.Register)
	router.POST("/login", c.Login)
	router.POST("/refresh-token", c.RefreshToken)
}
