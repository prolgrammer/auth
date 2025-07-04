package http

import (
	"auth/internal/controllers"
	"auth/internal/controllers/http/middleware"
	"auth/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type webhookController struct {
	logger logger.Logger
}

func NewWebhookController(
	handler *gin.Engine,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	wc := &webhookController{
		logger: logger,
	}

	handler.POST("/auth/webhook", wc.HandleWebhook, middleware.HandleErrors)
}

func (wc *webhookController) HandleWebhook(c *gin.Context) {
	var message map[string]interface{}
	if err := c.ShouldBindJSON(&message); err != nil {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	wc.logger.Info().Msgf("Message: %+v\n", message)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "webhook received",
	})
}
