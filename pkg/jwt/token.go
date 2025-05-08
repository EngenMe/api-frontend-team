package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token with the given user ID & email and expiration time.

func GenerateToken(userID, email string) (idTokenRes string, idExpRes string, err error) {
	// Create a new JWT token
	idExp := time.Now().Add(time.Minute * 15) // Token expires in 15 minates
	idClaims := TokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(idExp),
		},
	}

	idToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, idClaims).SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err != nil {
		return "", "", err
	}
	return idToken, idExp.String(), nil
}

// parseToken parses the JWT token and returns the claims.
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GenerateRefreshToken(userID string) (idRefreshTokenRes string, refreshExpRes string, err error) {
	rExp := time.Now().Add(72 * time.Hour) // Token expires in 3 days
	idRClaims := RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(rExp),
		},
	}

	idRToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, idRClaims).SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err != nil {
		return "", "", err
	}
	return idRToken, rExp.String(), nil
}
