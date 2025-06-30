package entities

import "time"

type Session struct {
	AccountId       string
	AccessToken     string
	RefreshToken    string
	AccessExpiresAt time.Time
	ExpiresAt       time.Time
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
