package dto

type AuthResponse struct {
	UserId       string               `json:"_id" validate:"required"`
	Name         string               `json:"name" validate:"required,min=3,max=15"`
	Email        string               `json:"email" validate:"required,email,min=5,max=40"`
	JWT          string               `json:"token" validate:"required"`
	RefreshToken RefreshTokenResponse `json:"refresh_token" validate:"required"`
}
