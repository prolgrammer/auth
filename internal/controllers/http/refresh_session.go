package http

import (
	"auth/internal/controllers"
	"auth/internal/controllers/http/middleware"
	"auth/internal/controllers/requests"
	"auth/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type refreshSessionController struct {
	useCase usecases.RefreshSessionUseCase
}

func NewRefreshSessionController(
	handler *gin.Engine,
	useCase usecases.RefreshSessionUseCase,
	middleware middleware.Middleware) {
	u := &refreshSessionController{
		useCase: useCase,
	}

	handler.POST("/auth/token/update", u.RefreshSession, middleware.HandleErrors)
}

// RefreshSession godoc
// @Summary      обновление сессии
// @Description  возвращает новую пару токенов при отправке старой пары и при условии их валидности
// @Accept       json
// @Produce      json
// @Param request body requests.RefreshSession true "request format"
// @Success      200  {object}  responses.Session
// @Failure 400 {object} string "некорректный формат запроса"
// @Failure 401 {object} string "невалидная пара токенов, либо истекший refresh token"
// @Failure 500 {object} string "внутренняя ошибка сервера"
// @Router       /auth/token/update [post]
func (r *refreshSessionController) RefreshSession(c *gin.Context) {
	request := requests.RefreshSession{}
	if err := c.ShouldBind(&request); err != nil {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	response, err := r.useCase.RefreshSession(c, c.Writer, request, ip, userAgent)
	if err != nil {
		middleware.AddGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
