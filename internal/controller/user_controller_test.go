package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/EngenMe/api-frontend-team/internal/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string
	jwt.RegisteredClaims
}

func createExpiredToken(userID, email string) string {
	// Set expiration time in the past
	expiredAt := time.Now().Add(-1 * time.Minute)

	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET_KEY")
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

func init() {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatal("Error loading .env file for tests")
	}
}

func setup(ts *httptest.Server) dto.AuthUserResponse {

	user := &dto.RegisterRequest{
		Email:    "test00@exmple.com",
		Password: "password",
		Name:     "testuser",
	}
	body, _ := json.Marshal(user)

	resp, _ := http.Post(ts.URL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(body))
	var responseBody dto.AuthUserResponse
	json.NewDecoder(resp.Body).Decode(&responseBody)
	return responseBody

}

func Test_get_profile_api_integration_test(t *testing.T) {
	ts := runTestServer()
	user := setup(ts)
	defer ts.Close()
	t.Cleanup(func() {
		tearDown(user.User.Id)
	})
	t.Run("it should return 200 when user get profile", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/user/me", nil)
		req.Header.Set("Authorization", "Bearer "+user.Tokens.Access.Token)

		client := &http.Client{}
		resp, _ := client.Do(req)
		var responseBody dto.GetUserResponse
		json.NewDecoder(resp.Body).Decode(&responseBody)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, responseBody.Email, user.User.Email)
		assert.Equal(t, responseBody.Name, user.User.Name)
		assert.Equal(t, responseBody.Id, user.User.Id)

	})

	t.Run("it should return 401 when user is not authorized", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/user/me", nil)
		req.Header.Set("Authorization", "Bearer "+"invalid_token")
		client := &http.Client{}
		resp, _ := client.Do(req)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	})
	t.Run("it should return 401 when access token is expired after 15 mins", func(t *testing.T) {
		expToken := createExpiredToken(user.User.Id, user.User.Email)
		user.Tokens.Access.Token = expToken
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/user/me", nil)
		req.Header.Set("Authorization", "Bearer "+user.Tokens.Access.Token)
		client := &http.Client{}
		resp, _ := client.Do(req)
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Response status: %d\nBody: %s", resp.StatusCode, string(body))

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	})
}

func Test_update_profile_api_integration_test(t *testing.T) {
	ts := runTestServer()
	user := setup(ts)
	defer ts.Close()
	t.Cleanup(func() {
		tearDown(user.User.Id)
	})
	t.Run("it should return 200 when user update profile", func(t *testing.T) {
		userUpdate := &dto.UpdateUserRequest{
			Email:    "updateemail@test.com",
			Password: "password",
			Name:     "updateuser",
		}
		body, _ := json.Marshal(userUpdate)
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/user/me", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+user.Tokens.Access.Token)
		client := &http.Client{}
		resp, _ := client.Do(req)
		var responseBody dto.UpdateUserResponse
		json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, responseBody.Email, userUpdate.Email)
		assert.Equal(t, responseBody.Name, userUpdate.Name)
		assert.Equal(t, responseBody.Id, user.User.Id)
	})
	t.Run("it should return 401 when user is not authorized", func(t *testing.T) {
		userUpdate := &dto.UpdateUserRequest{
			Email:    "updateemail@test.com",
			Password: "password",
			Name:     "updateuser",
		}
		body, _ := json.Marshal(userUpdate)
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/user/me", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+"invalid_token")
		client := &http.Client{}
		resp, _ := client.Do(req)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
	t.Run("it should return 400 when user update profile with invalid data", func(t *testing.T) {
		userUpdate := &dto.UpdateUserRequest{
			Email:    "updateemailtest.com",
			Password: "password",
			Name:     "updateuser",
		}
		body, _ := json.Marshal(userUpdate)
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/user/me", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+user.Tokens.Access.Token)
		client := &http.Client{}
		resp, _ := client.Do(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

}

func Test_delete_profile_api_integration_test(t *testing.T) {
	ts := runTestServer()
	user := setup(ts)
	defer ts.Close()
	t.Cleanup(func() {
		tearDown(user.User.Id)
	})
	t.Run("it should return 204 when user delete profile", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/user/me", nil)
		req.Header.Set("Authorization", "Bearer "+user.Tokens.Access.Token)
		client := &http.Client{}
		resp, _ := client.Do(req)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
	t.Run("it should return 401 when user is not authorized", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/user/me", nil)
		req.Header.Set("Authorization", "Bearer "+"invalid_token")
		client := &http.Client{}
		resp, _ := client.Do(req)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
