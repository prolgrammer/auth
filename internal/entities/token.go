package entities

import "time"

const (
	UserIdClaimName    = "sub"
	ExpiresAtClaimName = "exp"
)

type AccessTokenClaims map[string]any

func (c AccessTokenClaims) AccountId() string { return c[UserIdClaimName].(string) }

func (c AccessTokenClaims) ExpiresAt() time.Time { return time.Unix(c[ExpiresAtClaimName].(int64), 0) }

func NewClaims(accountId string, expiresAt time.Time) AccessTokenClaims {
	return AccessTokenClaims{UserIdClaimName: accountId, ExpiresAtClaimName: expiresAt.Unix()}
}
