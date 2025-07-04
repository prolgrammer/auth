package middleware

import (
	"auth/pkg/logger"
	"github.com/gin-gonic/gin"
)

type middleware struct {
	logger  logger.Logger
	manager SessionService
}

type Middleware interface {
	Authenticate(c *gin.Context)
	HandleErrors(c *gin.Context)
}

func NewMiddleware(manager SessionService, logger logger.Logger) Middleware {
	return &middleware{logger: logger, manager: manager}
}
