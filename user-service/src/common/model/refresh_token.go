package model

import "time"

type RefreshToken struct {
	Token     string
	UserID    string
	CreatedAt  time.Time
	UpdatedAt time.Time
}


type RefreshTokenWithRole struct {
	Token     string
	UserID    string
	Role      string
}