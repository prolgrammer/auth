package http

import (
	"auth/internal/controllers"
	"auth/internal/controllers/http/middleware"
	"auth/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getUserController struct {
	useCase usecases.GetUserUseCase
}

func NewGetUserController(
	handler *gin.Engine,
	useCase usecases.GetUserUseCase,
	middleware middleware.Middleware,
) {
	g := &getUserController{
		useCase: useCase,
	}

	handler.GET("/auth/user", middleware.Authenticate, g.GetUser, middleware.HandleErrors)
}

// GetUser godoc
// @Summary      запрос на получение пользователя пользователя
// @Description  запрос на получение пользователя с использованием токена, переданного в заголовке "Authorization"
// @Produce      json
// @Param Authorization header string true "access token"
// @Success 200 {object} responses.User
// @Failure 401 {object} string "некорректный access token"
// @Failure 500 {object} string "внутренняя ошибка сервера"
// @Router       /auth/user [get]
func (gc *getUserController) GetUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		middleware.AddGinError(c, controllers.ErrAuthRequired)
		return
	}

	response, err := gc.useCase.GetUserUseCase(c, userId.(string))
	if err != nil {
		middleware.AddGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
