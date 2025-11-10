package dto

type UserDTO struct {
	Name     string `json:"name" validate:"required,min=3,max=15"`
	LastName string `json:"lastname" validate:"required,min=4,max=15"`
	Email    string `json:"email" validate:"required,email,min=5,max=40"`
	Password string `json:"password" validate:"required,min=8"`
}
