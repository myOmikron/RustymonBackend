package models

import (
	"time"
)

type Session struct {
	Common
	UserID     uint
	SessionKey string
	ValidUntil time.Time
}
