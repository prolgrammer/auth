package entities

import "time"

type Session struct {
	Id              int
	UserId          string
	AccessToken     string
	RefreshToken    string
	UserAgent       string
	IP              string
	AccessExpiresAt time.Time
	ExpiresAt       time.Time
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
