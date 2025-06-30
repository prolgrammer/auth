package middleware

import (
	"github.com/gin-gonic/gin"
)

type middleware struct {
	manager SessionService
}

type Middleware interface {
	Authenticate(c *gin.Context)
	HandleErrors(c *gin.Context)
}

func NewMiddleware(manager SessionService) Middleware {
	return &middleware{manager: manager}
}
