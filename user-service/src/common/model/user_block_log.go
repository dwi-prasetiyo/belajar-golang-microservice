package model

import "time"

type UserBlockLog struct {
	ID          int        `json:"id"`
	UserID      string     `json:"user_id"`
	Reason      string     `json:"reason"`
	BlockedAt   time.Time  `json:"blocked_at"`
	UnblockedAt *time.Time `json:"unblocked_at"`
	IsActive    bool       `json:"is_active"`
}
