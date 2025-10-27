package models

type User struct {
	UserId  string `json:"id" validate:"required" bson:"_id"`
	Name     string `json:"name" validate:"required" bson:"name"`
	LastName string `json:"lastname" validate:"required" bson:"last_name"`
	Email    string `json:"email" validate:"required,email" bson:"email"`
	PasswordHash string `json:"-" validate:"required" bson:"password_hash"`
}