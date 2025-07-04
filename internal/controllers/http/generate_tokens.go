package http

import (
	"auth/internal/controllers/http/middleware"
	"auth/internal/usecases"
	"auth/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getTokensController struct {
	logger  logger.Logger
	useCase usecases.GenerateTokensUseCase
}

func NewGenerateTokensController(
	handler *gin.Engine,
	useCase usecases.GenerateTokensUseCase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	u := &getTokensController{
		logger:  logger,
		useCase: useCase,
	}

	handler.GET("/auth/token/:user_id", u.GenerateTokens, middleware.HandleErrors)
}

// GenerateTokens godoc
// @Summary      создание токенов
// @Description  создание токенов по id пользователя
// @Accept       json
// @Produce      json
// @Param        user_id path string true "path format"
// @Success      200  {object}  responses.Session
// @Failure 400 {object} string "некорректный формат запроса"
// @Failure 404 {object} string "пользователь не найден"
// @Failure 500 {object} string "внутренняя ошибка сервера"
// @Router       /auth/token/{user_id} [get]
func (router *getTokensController) GenerateTokens(c *gin.Context) {
	userId := c.Param("user_id")

	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()

	response, err := router.useCase.GenerateTokens(c, c.Writer, userId, ip, userAgent)
	if err != nil {
		middleware.AddGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
