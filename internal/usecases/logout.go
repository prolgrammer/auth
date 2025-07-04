package usecases

import (
	"context"
	"net/http"
)

type logoutUseCase struct {
	sessionRepo   LogoutSessionRepository
	cookieService LogoutCookieService
}

type LogoutUseCase interface {
	Logout(context context.Context, writer http.ResponseWriter, userId string) error
}

func NewLogoutUseCase(
	sessionRepo LogoutSessionRepository,
	cookieService LogoutCookieService,
) LogoutUseCase {
	return &logoutUseCase{
		sessionRepo:   sessionRepo,
		cookieService: cookieService,
	}
}

func (l logoutUseCase) Logout(context context.Context, writer http.ResponseWriter, userId string) error {
	l.cookieService.Clear(writer, "access_token")
	return l.sessionRepo.DeleteByUserId(context, userId)
}
