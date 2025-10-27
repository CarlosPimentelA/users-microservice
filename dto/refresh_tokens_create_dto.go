package dto

import "time"

type RefreshTokenCreateDTO struct {
	UserId    string
	TokenPlain string
	ExpiresAt time.Time
	ClientInfo string
}