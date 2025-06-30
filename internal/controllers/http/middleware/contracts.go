package middleware

import "auth/internal/entities"

type (
	SessionService interface {
		ParseToken(string) (entities.AccessTokenClaims, error)
	}
)
