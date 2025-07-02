package usecases

import (
	"auth/internal/controllers/responses"
	"auth/internal/repositories"
	"context"
	"errors"
	"fmt"
	"net/http"
)

type generateTokensUseCase struct {
	userRepo       GenerateTokensUserRepository
	sessionRepo    GenerateTokensSessionRepository
	hashProvider   GenerateTokensHashService
	cookieService  GenerateTokensCookieService
	sessionManager GenerateTokensSessionService
}

type GenerateTokensUseCase interface {
	GenerateTokens(context context.Context, writer http.ResponseWriter, userId, ip, userAgent string) (responses.Session, error)
}

func NewGenerateTokensUseCase(
	userRepo GenerateTokensUserRepository,
	sessionRepo GenerateTokensSessionRepository,
	hashProvider GenerateTokensHashService,
	cookieService GenerateTokensCookieService,
	sessionManager GenerateTokensSessionService,
) GenerateTokensUseCase {
	return &generateTokensUseCase{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		hashProvider:   hashProvider,
		sessionManager: sessionManager,
		cookieService:  cookieService,
	}
}

func (uc *generateTokensUseCase) GenerateTokens(context context.Context, writer http.ResponseWriter, userId, ip, userAgent string) (responses.Session, error) {
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
