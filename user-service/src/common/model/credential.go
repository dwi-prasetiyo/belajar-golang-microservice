package model

import "time"

type Credential struct {
	UserID    string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CredentialWithRole struct {
	UserID    string
	Password  string
	Role      string
}

