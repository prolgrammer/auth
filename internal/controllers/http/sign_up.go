package http

import (
	"auth/internal/controllers"
	"auth/internal/controllers/http/middleware"
	"auth/internal/controllers/requests"
	"auth/internal/usecases"
	"auth/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type signupController struct {
	logger logger.Logger
	user   usecases.SignUpUseCase
}

func NewSignUpController(
	handler *gin.Engine,
	user usecases.SignUpUseCase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	u := &signupController{
		logger: logger,
		user:   user,
	}

	handler.POST("/auth/signup", u.SignUp, middleware.HandleErrors)
}

// SignUp godoc
// @Summary      регистрация нового пользователя
// @Description  регистрация нового пользователя
// @Accept       json
// @Produce      json
// @Param request body requests.SignUp true "структура запрос"
// @Success      200  {object}  responses.SignUp
// @Failure 400 {object} string "некорректный формат запроса"
// @Failure 409 {object} string "пользователь уже существует"
// @Failure 500 {object} string "внутренняя ошибка сервера"
// @Router       /auth/signup [post]
func (u *signupController) SignUp(c *gin.Context) {
	var userRequest requests.SignUp

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	response, err := u.user.CreateUser(c, c.Writer, userRequest, userAgent, ip)

	if err != nil {
		middleware.AddGinError(c, errors.Wrap(err, "failed to create account"))
		return
	}

	c.JSON(http.StatusOK, response)
}
