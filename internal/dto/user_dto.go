package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" gorm:"not null"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" gorm:"not null"`
}

type GetUserResponse struct {
	Email string
}

type UpdateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

type UpdateUserResponse struct {
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
