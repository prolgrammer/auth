package http

import (
	"auth/internal/controllers"
	"auth/internal/controllers/http/middleware"
	"auth/internal/controllers/requests"
	"auth/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getTokensController struct {
	useCase usecases.GetTokensUseCase
}

func NewGetTokensController(
	handler *gin.Engine,
	useCase usecases.GetTokensUseCase,
	middleware middleware.Middleware,
) {
	u := &getTokensController{
		useCase: useCase,
	}

	handler.GET("/token", u.GetTokens, middleware.HandleErrors)
}

func (router *getTokensController) GetTokens(c *gin.Context) {
	var request requests.GetTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()

	response, err := router.useCase.GetTokens(c, c.Writer, request, userAgent, ip)
	if err != nil {
		middleware.AddGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
