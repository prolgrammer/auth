package http

import (
	"auth/internal/controllers/http/middleware"
	"auth/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type signInController struct {
	useCase usecases.SignInUseCase
}

func NewSignInController(
	handler *gin.Engine,
	useCase usecases.SignInUseCase,
	middleware middleware.Middleware) {
	u := &signInController{
		useCase: useCase,
	}

	handler.POST("/signin", u.SignIn, middleware.HandleErrors)
}

// SignIn godoc
// @Summary      вход в аккаунт
// @Description  вход в аккаунт с использованием пар логин + пароль или email + пароль для получения токенов
// @Accept       json
// @Produce      json
// @Param request body requests.SignIn true "структура запроса"
// @Success      200  {object}  responses.SignIn
// @Failure 400 {object} string "некорректный формат запроса"
// @Failure 401 {object} string "неправильный пароль"
// @Failure 404 {object} string "пользователь не найден"
// @Failure 500 {object} string "внутренняя ошибка сервера"
// @Router       /signin [post]
func (router *signInController) SignIn(c *gin.Context) {
	var request requests.SignIn
	if err := c.ShouldBind(&request); err != nil {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	response, err := router.useCase.SignIn(c, &request)

	if err != nil {
		middleware.AddGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
	return
}
