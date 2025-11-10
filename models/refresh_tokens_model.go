package models

import (
	"time"
)

type RefreshToken struct {
	ID             string    `bson:"_id"`
	UserId         string    `bson:"user_id"`
	Jti     	   string    `bson:"jti"`
	IssuedAt       time.Time `bson:"created_at"`
	Expires        time.Time `bson:"expiry_time"`
	Revoked        bool      `bson:"revoked"`
	SessionVersion int       `json:"-" validate:"required" bson:"session"`
}
