package dto

import "time"

type RefreshTokenResponse struct {
	RefreshToken string
	ExpiresAt    time.Time
}