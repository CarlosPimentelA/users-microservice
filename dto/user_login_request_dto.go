package dto

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email,min=5,max=40"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}