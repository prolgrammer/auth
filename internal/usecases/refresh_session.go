package usecases

import (
	"auth/internal/controllers/requests"
	"auth/internal/controllers/responses"
	"auth/internal/repositories"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type refreshSessionUseCase struct {
	userRepository    RefreshSessionUserRepository
	sessionRepository RefreshSessionSessionRepository
	sessionService    RefreshSessionSessionService
	cookieService     RefreshSessionCookieService
	hashProvider      RefreshSessionHashProvider
}

type RefreshSessionUseCase interface {
	RefreshSession(context *gin.Context, writer http.ResponseWriter, request requests.RefreshSession, ip, userAgent string) (responses.Session, error)
}

func NewRefreshSessionUseCase(
	userRepository RefreshSessionUserRepository,
	sessionRepository RefreshSessionSessionRepository,
	sessionService RefreshSessionSessionService,
	cookieService RefreshSessionCookieService,
	hashProvider RefreshSessionHashProvider) RefreshSessionUseCase {
	return &refreshSessionUseCase{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		sessionService:    sessionService,
		cookieService:     cookieService,
		hashProvider:      hashProvider,
	}
}

func (r refreshSessionUseCase) RefreshSession(context *gin.Context, writer http.ResponseWriter, request requests.RefreshSession, ip, userAgent string) (responses.Session, error) {
	accessToken := request.AccessToken
	refreshToken := request.RefreshToken

	claims, err := r.sessionService.ParseToken(accessToken)
	if err != nil {
		return responses.Session{}, err
	}

	userId := claims["sub"].(string)
	session, err := r.sessionRepository.SelectByUserId(context, userId)
	if err != nil {
		return responses.Session{}, fmt.Errorf("failed to select session: %s", err)
	}

	if userAgent != session.UserAgent {
		err = r.sessionRepository.DeleteByUserId(context, session.UserId)
		if err != nil {
			return responses.Session{}, fmt.Errorf("failed to delete session: %w", err)
		}
		return responses.Session{}, fmt.Errorf("user-agent mismatch: %w", ErrUnauthorized)
	}

	if session.IP != ip {
		sendWebhook(session.UserId, ip)
	}

	cookieAccessToken, err := context.Cookie("access_token")
	if err != nil {
		return responses.Session{}, fmt.Errorf("failed to get access token from cookie: %w", err)
	}
	if accessToken != cookieAccessToken {
		return responses.Session{}, fmt.Errorf("can not refresh session: %w", ErrNotAValidAccessToken)
	}

	valid := r.hashProvider.CompareStringAndHash(refreshToken, session.RefreshToken)
	if !valid {
		return responses.Session{}, fmt.Errorf("can not refresh session: %w", ErrNotAValidRefreshToken)
	}

	user, err := r.userRepository.SelectByUserId(context, session.UserId)
	if err != nil {
		if errors.Is(err, repositories.ErrEntityNotFound) {
			return responses.Session{}, fmt.Errorf("failed to find user: %w", ErrEntityNotFound)
		}
		return responses.Session{}, fmt.Errorf("failed to find user: %w", err)
	}

	newSession, err := r.sessionService.CreateSession(user)
	if err != nil {
		return responses.Session{}, fmt.Errorf("failed to create session: %w", err)
	}

	newSession.IP = ip
	newSession.UserAgent = userAgent

	hashedRefreshToken, err := r.hashProvider.GenerateHash(newSession.RefreshToken)
	if err != nil {
		return responses.Session{}, fmt.Errorf("failed to hash refresh token: %w", err)
	}
	rawRefreshToken := newSession.RefreshToken
	newSession.RefreshToken = string(hashedRefreshToken)

	err = r.sessionRepository.Update(context, newSession)
	if err != nil {
		return responses.Session{}, fmt.Errorf("failed to save session: %w", err)
	}

	r.cookieService.Set(
		writer,
		"access_token",
		newSession.AccessToken,
		newSession.AccessExpiresAt,
	)

	return responses.NewSession(newSession.AccessToken, rawRefreshToken, newSession.AccessExpiresAt.Unix()), nil
}
