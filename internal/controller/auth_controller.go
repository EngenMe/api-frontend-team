package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/EngenMe/api-frontend-team/internal/dto"
	"github.com/markbates/goth/gothic"

	"github.com/EngenMe/api-frontend-team/internal/service"
	"github.com/EngenMe/api-frontend-team/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service service.AuthService
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{service}
}

// Login godoc
// @Summary Login
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.RefreshTokenResonse
// @Failure 400 {object} error
// @Router /auth/login [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (c *AuthController) Login(ctx *gin.Context) {
	var userDTO dto.LoginRequest

	if err := ctx.ShouldBindJSON(&userDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUserResponse, err := c.service.Login(&userDTO)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": authUserResponse.User, "tokens": authUserResponse.Tokens})
}

// Register godoc
// @Summary Register
// @Description Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "Register request"
// @Success 201 {object} string
// @Failure 400 {object} error
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var userDTO dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&userDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userDTO.Email == " " || userDTO.Name == " " || userDTO.Password == " " {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	authUserResponse, err := c.service.Register(&userDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": authUserResponse.User, "tokens": authUserResponse.Tokens})
}

// RefreshToken godoc
// @Summary Refresh Token
// @Description Refresh JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} dto.RefreshTokenResonse
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Router /auth/refresh-token [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var refreshTokenRequest RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&refreshTokenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := jwt.ParseToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	userIdStr := fmt.Sprintf("%v", claims["user_id"])
	refresh_token, err := c.service.RefreshToken(userIdStr, refreshTokenRequest.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access": &refresh_token.Access, "refresh": &refresh_token.Refresh})
}

func (c *AuthController) BeginAuthHandler(ctx *gin.Context) {
	provider := ctx.Param("provider")
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
	log.Printf("Starting auth for provider: %s", provider)
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (c *AuthController) CallbackHandler(ctx *gin.Context) {
	provider := ctx.Param("provider")
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	authUserResponse, err := c.service.ProviderLogin(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": authUserResponse.User, "tokens": authUserResponse.Tokens})
}

func (c *AuthController) SetupAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", c.Register)
	router.POST("/login", c.Login)
	router.POST("/refresh", c.RefreshToken)
	router.GET("/:provider", c.BeginAuthHandler)
	router.GET("/:provider/callback", c.CallbackHandler)
}
