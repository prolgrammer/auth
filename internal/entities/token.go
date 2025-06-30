package entities

import "time"

const (
	AccountIdClaimName = "sub"
	ExpiresAtClaimName = "exp"
)

type AccessTokenClaims map[string]any

func (c AccessTokenClaims) AccountId() string { return c[AccountIdClaimName].(string) }

func (c AccessTokenClaims) ExpiresAt() time.Time { return time.Unix(c[ExpiresAtClaimName].(int64), 0) }

func NewClaims(accountId string, expiresAt time.Time) AccessTokenClaims {
	return AccessTokenClaims{AccountIdClaimName: accountId, ExpiresAtClaimName: expiresAt.Unix()}
}
