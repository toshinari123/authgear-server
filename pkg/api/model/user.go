package model

import (
	"time"
)

type User struct {
	Meta
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	IsAnonymous   bool       `json:"is_anonymous"`
	IsVerified    bool       `json:"is_verified"`
	IsDisabled    bool       `json:"is_disabled"`
	DisableReason *string    `json:"disable_reason"`
}
