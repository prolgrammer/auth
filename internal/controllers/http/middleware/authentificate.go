package middleware

import (
	"auth/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *middleware) Authenticate(c *gin.Context) {
	token := c.GetHeader("Authorization")

	claims, err := m.manager.ParseToken(token)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, usecases.ErrNotAValidAccessToken)
		return
	}

	c.Set("account_id", claims.AccountId)
}
