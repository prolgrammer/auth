package http

import (
	"auth/internal/controllers"
	"auth/internal/controllers/http/middleware"
	"auth/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type logoutController struct {
	useCase usecases.LogoutUseCase
}

func NewLogoutController(
	handler *gin.Engine,
	useCase usecases.LogoutUseCase,
	middleware middleware.Middleware,
) {
	g := &logoutController{
		useCase: useCase,
	}

	handler.POST("/auth/session/logout", middleware.Authenticate, g.Logout, middleware.HandleErrors)
}

// Logout godoc
// @Summary      запрос на закрытие сессий пользователя
// @Description  запрос на закрытие сессий пользователя по их id с использованием токена, переданного в заголовке "Authorization"
// @Produce      json
// @Param Authorization header string true "access token"
// @Success 200 "ok"
// @Failure 401 {object} string "некорректный access token"
// @Failure 500 {object} string "внутренняя ошибка сервера"
// @Router       /auth/session/logout [post]
func (g *logoutController) Logout(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		middleware.AddGinError(ctx, controllers.ErrAuthRequired)
		return
	}

	err := g.useCase.Logout(ctx, ctx.Writer, userId.(string))
	if err != nil {
		middleware.AddGinError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, "successful logout")
}
