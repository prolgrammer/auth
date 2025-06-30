package middleware

import (
	"auth/internal/usecases"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *middleware) HandleErrors(c *gin.Context) {
	if len(c.Errors) > 0 {
		err := c.Errors.Last()

		if errors.Is(err, usecases.ErrWrongPassword) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, usecases.ErrEntityNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, "Internal Server Error")
	}
}
