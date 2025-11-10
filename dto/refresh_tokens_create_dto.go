package dto

import "time"

type RefreshTokenCreateDTO struct {
	UserId     string 
	Jti string		   
	ExpiresAt  time.Time 
}
