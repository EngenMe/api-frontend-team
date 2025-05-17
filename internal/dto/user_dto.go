package dto

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthUserResponse struct {
	User   User                `json:"user"`
	Tokens RefreshTokenResonse `json:"tokens"`
}
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" gorm:"not null"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" gorm:"not null"`
}

type GetUserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string
}

type UpdateUserRequest struct {
	Name     string `json:"name"  binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

type UpdateUserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TokenPair struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}
type RefreshTokenResonse struct {
	Access  TokenPair `json:"access"`
	Refresh TokenPair `json:"refresh"`
}
