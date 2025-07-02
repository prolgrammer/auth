package usecases

import (
	"auth/internal/controllers/requests"
	"auth/internal/controllers/responses"
	"auth/internal/repositories"
	"context"
	"errors"
	"fmt"
	"net/http"
)

type getTokensUseCase struct {
	userRepo       GetTokensUserRepository
	sessionRepo    GetTokensSessionRepository
	hashProvider   GetTokensHashService
	cookieService  GetTokensCookieService
	sessionManager GetTokensSessionService
}

type GetTokensUseCase interface {
	GetTokens(context context.Context, writer http.ResponseWriter, request requests.GetTokenRequest, ip, userAgent string) (responses.Session, error)
}

func NewGetTokensUseCase(
	userRepo GetTokensUserRepository,
	sessionRepo GetTokensSessionRepository,
	hashProvider GetTokensHashService,
	cookieService GetTokensCookieService,
	sessionManager GetTokensSessionService,
) GetTokensUseCase {
	return &getTokensUseCase{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		hashProvider:   hashProvider,
		sessionManager: sessionManager,
		cookieService:  cookieService,
	}
}

func (uc *getTokensUseCase) GetTokens(context context.Context, writer http.ResponseWriter, request requests.GetTokenRequest, ip, userAgent string) (responses.Session, error) {
	userId := request.UserId
	user, err := uc.userRepo.SelectByUserId(context, userId)
	if err != nil {
		if errors.Is(err, repositories.ErrEntityNotFound) {
			return responses.Session{}, fmt.Errorf("failed to find user: %w", ErrEntityNotFound)
		}
		return responses.Session{}, fmt.Errorf("failed to find user: %w", err)
	}

	err = uc.sessionRepo.DeleteByUserId(context, user.Id)
	if err != nil {
		return responses.Session{}, fmt.Errorf("failed to delete session: %w", err)
	}

	session, err := uc.sessionManager.CreateSession(user)
	if err != nil {
		return responses.Session{}, fmt.Errorf("%w: couldn't create session", err)
	}

	session.IP = ip
	session.UserAgent = userAgent

	hashedRefreshToken, err := uc.hashProvider.GenerateHash(session.RefreshToken)
	if err != nil {
		return responses.Session{}, fmt.Errorf("%w: failed to hash refresh token", err)
	}

	refreshToken := session.RefreshToken
	session.RefreshToken = string(hashedRefreshToken)

	err = uc.sessionRepo.Insert(context, session)
	if err != nil {
		return responses.Session{}, fmt.Errorf("%w: failed to insert session", err)
	}

	uc.cookieService.Set(
		writer,
		"access_token",
		session.AccessToken,
		session.AccessExpiresAt,
	)

	refreshSessionResponse := responses.NewSession(session.AccessToken, refreshToken, session.AccessExpiresAt.Unix())

	return refreshSessionResponse, nil
}
