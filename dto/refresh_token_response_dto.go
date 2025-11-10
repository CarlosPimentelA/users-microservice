package dto

import "time"

type RefreshTokenResponse struct {
	ExpiresAt    time.Time
}
