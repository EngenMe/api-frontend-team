package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/EngenMe/api-frontend-team/internal/dto"
	"github.com/EngenMe/api-frontend-team/internal/middleware"
	"github.com/EngenMe/api-frontend-team/internal/model"
	"github.com/EngenMe/api-frontend-team/internal/repository"
	"github.com/EngenMe/api-frontend-team/internal/service"
	"github.com/EngenMe/api-frontend-team/pkg/db"
	"github.com/EngenMe/api-frontend-team/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func init() {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatal("Error loading .env file for tests")
	}
}

func runTestServer() *httptest.Server {
	dbConn = db.InitDB()

	userRepo := repository.NewUserRepository(dbConn)
	tokenRepo := repository.NewTokenRepo(dbConn)
	authService := service.NewAuthService(userRepo, tokenRepo)
	userService := service.NewUserService(userRepo)
	authController := NewAuthController(authService)
	UserController := NewUserController(userService)

	router := gin.Default()

	apiV1Auth := router.Group("/api/v1/auth")
	apiV1User := router.Group("/api/v1/user")

	apiV1User.Use(middleware.AuthenticationMiddleware())

	authController.SetupAuthRoutes((apiV1Auth))
	UserController.SetupUserRoutes(apiV1User)

	return httptest.NewServer(router)
}

func tearDown(userId string) {
	dbConn.Unscoped().Delete(&model.User{}, "id = ?", userId)
	dbConn.Unscoped().Delete(&model.Token{}, "user_id = ?", userId)

}

func Test_post_api_integration_test_register(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	t.Run("it should return 201 when user is registered", func(t *testing.T) {
		user := &dto.RegisterRequest{
			Email:    "test1@exmple.com",
			Password: "password",
			Name:     "testuser",
		}
		body, _ := json.Marshal(user)

		resp, _ := http.Post(ts.URL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(body))
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status code 201, got %d", resp.StatusCode)
		}
		var responseBody dto.AuthUserResponse
		json.NewDecoder(resp.Body).Decode(&responseBody)
		var createdUser model.User
		dbConn.First(&createdUser, "email = ?", user.Email)
		assert.Equal(t, createdUser.Email, responseBody.User.Email)
		assert.Equal(t, user.Email, responseBody.User.Email)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, responseBody.Tokens.Access)
		assert.NotEmpty(t, responseBody.Tokens.Refresh)
		tearDown(responseBody.User.Id)

	})

	t.Run("it should return 500 when user already exists", func(t *testing.T) {
		user := &model.User{
			Email:    "test5@exmple.com",
			Password: "password",
			Name:     "testuser",
		}
		dbConn.Create(&user)
		var createdUser model.User
		dbConn.First(&createdUser, "email = ?", user.Email)
		body, _ := json.Marshal(user)
		resp, _ := http.Post(ts.URL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(body))
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		tearDown(strconv.FormatUint(uint64(createdUser.ID), 10))
	})

	t.Run("it should return 400 when email is invalid", func(t *testing.T) {
		user := &dto.RegisterRequest{
			Email:    "test1example.com",
			Password: "password",
			Name:     "testuser",
		}
		body, _ := json.Marshal(user)

		resp, _ := http.Post(ts.URL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(body))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.EqualError(t, errors.New("invalid email format"), "invalid email format")
	})

	t.Run("it should return 400 when user is not valid", func(t *testing.T) {
		user := &dto.RegisterRequest{
			Email:    "",
			Password: "",
			Name:     ""}
		body, _ := json.Marshal(user)
		resp, _ := http.Post(ts.URL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(body))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.EqualError(t, errors.New("Bad Request"), "Bad Request")
	})

}

func Test_post_api_integration_test_login(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	user := &model.User{
		Email:    "test@example.com",
		Password: "123456",
		Name:     "testuser",
	}
	hPassword, _ := utils.HashPassword(user.Password)
	user.Password = hPassword
	dbConn.Create(&user)
	var createdUser model.User
	dbConn.First(&createdUser, "email = ?", user.Email)
	t.Run("it should return 200 when user is logged in", func(t *testing.T) {
		user := &dto.LoginRequest{
			Email:    "test@example.com",
			Password: "123456",
		}
		body, _ := json.Marshal(user)
		resp, _ := http.Post(ts.URL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(body))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		tearDown(strconv.FormatUint(uint64(createdUser.ID), 10))
	})

	t.Run("it should return 401 when user is not found", func(t *testing.T) {
		user := &dto.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(user)
		resp, _ := http.Post(ts.URL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(body))
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.EqualError(t, errors.New("Invalid credentials"), "Invalid credentials")
		tearDown(strconv.FormatUint(uint64(createdUser.ID), 10))
	})

}

func Test_post_api_integration_test_refresh(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	user := &model.User{
		Email:    "test@example.com",
		Password: "123456",
		Name:     "testuser",
	}
	hPassword, _ := utils.HashPassword(user.Password)
	user.Password = hPassword
	dbConn.Create(&user)
	var createdUser model.User
	dbConn.First(&createdUser, "email = ?", user.Email)
	t.Run("it should return 200 when refresh token is valid", func(t *testing.T) {
		user := &dto.LoginRequest{
			Email:    "test@example.com",
			Password: "123456",
		}
		body, _ := json.Marshal(user)
		resp, _ := http.Post(ts.URL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(body))

		var responseBody dto.AuthUserResponse
		json.NewDecoder(resp.Body).Decode(&responseBody)
		refreshToken := responseBody.Tokens.Refresh.Token

		refreshBody, _ := json.Marshal(&RefreshTokenRequest{
			RefreshToken: refreshToken,
		})
		resp, _ = http.Post(ts.URL+"/api/v1/auth/refresh", "application/json", bytes.NewBuffer(refreshBody))

		bodyBytes, _ := io.ReadAll(resp.Body)

		var respBody *dto.RefreshTokenResonse
		err := json.Unmarshal(bodyBytes, &respBody)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, respBody.Refresh)
		assert.NotEmpty(t, respBody.Access)
		assert.NotEqual(t, respBody.Refresh, refreshToken)
		assert.NotEqual(t, respBody.Access, refreshToken)

		tearDown(strconv.FormatUint(uint64(createdUser.ID), 10))

	})
}
