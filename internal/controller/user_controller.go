package controller

import (
	"net/http"

	"github.com/EngenMe/api-frontend-team/internal/dto"

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

	ctx.JSON(
		http.StatusOK, dto.GetUserResponse{
			Email: user.Email,
		},
	)
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.service.GetUserById(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(
		http.StatusOK, dto.GetUserResponse{
			Email: user.Email,
		},
	)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var userDTO dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&userDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.service.UpdateUser(userIdStr, userDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(
		http.StatusOK, dto.UpdateUserResponse{
			Email: user.Email,
		},
	)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err := c.service.DeleteUser(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *UserController) SetupUserRoutes(router *gin.RouterGroup) {
	router.GET("/:email", c.GetUser)
	router.GET("/me", c.GetProfile)
	router.PUT("/me", c.UpdateUser)
	router.DELETE("/me", c.DeleteUser)
}
