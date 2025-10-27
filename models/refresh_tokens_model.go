package models

import (
	"time"
)

type RefreshToken struct {
	ID        string
	UserId    string
	TokenHash string
	IssuedAt  time.Time
	Expires time.Time
	Revoked bool
}