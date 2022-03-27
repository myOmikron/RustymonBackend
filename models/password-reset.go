package models

import (
	"github.com/myOmikron/echotools/utilitymodels"
	"time"
)

type PasswordReset struct {
	ID         uint               `gorm:"primarykey" json:"id"`
	UserID     uint               `json:"user_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User       utilitymodels.User `json:"-"`
	Token      string             `json:"token" gorm:"not null"`
	ValidUntil time.Time          `json:"valid_until" gorm:"not null"`
}
