package pkg

import (
	"auth/config"
	"auth/internal/entities"
	"time"
)

type sessionService struct {
	config  config.TokenConfiguration
	access  TokenService
	refresh TokenService
}

type SessionService interface {
	CreateSession(account entities.User) (entities.Session, error)
	ParseToken(token string) (entities.AccessTokenClaims, error)
}

func NewSessionService(
	config config.TokenConfiguration,
	access TokenService,
	refresh TokenService) SessionService {
	return &sessionService{config: config, access: access, refresh: refresh}
}

func (t *sessionService) CreateSession(account entities.User) (entities.Session, error) {
	accessExpiresAt := time.Now().Add(t.config.AccessTokenDuration)
	refreshExpiresAt := time.Now().Add(t.config.RefreshTokenDuration)

	accessClaims := entities.NewClaims(account.Id, accessExpiresAt)
	access, err := t.access.CreateAccessToken(accessClaims)
	if err != nil {
		return entities.Session{}, err
	}

	refresh, err := t.refresh.CreateRefreshToken()
	if err != nil {
		return entities.Session{}, err
	}

	return entities.Session{
		UserId:          account.Id,
		AccessToken:     access,
		AccessExpiresAt: accessExpiresAt,
		RefreshToken:    refresh,
		ExpiresAt:       refreshExpiresAt,
	}, nil
}

func (t *sessionService) ParseToken(token string) (entities.AccessTokenClaims, error) {
	claims, err := t.access.ParseAccessToken(token)
	if err != nil {
		return entities.AccessTokenClaims{}, err
	}

	return claims, nil
}
