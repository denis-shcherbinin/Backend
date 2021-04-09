package entity

import "time"

type Session struct {
	UserAgent    string    `json:"userAgent"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
