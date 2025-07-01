package usecases

import (
	"auth/internal/controllers/requests"
	"auth/internal/controllers/responses"
	"auth/internal/entities"
	"auth/internal/repositories"
	"context"
	"errors"
	"fmt"
	"net/http"
)

type signInUseCase struct {
	userRepo       SignInUserRepository
	sessionRepo    SignInSessionRepository
	hashProvider   SignInHashService
	cookieService  SignInCookieService
	sessionManager SignInSessionService
}

type SignInUseCase interface {
	SignIn(context context.Context, writer http.ResponseWriter, request *requests.SignIn, userAgent, ip string) (responses.SignIn, error)
}

func NewSignInUseCase(
	userRepo SignInUserRepository,
	sessionRepo SignInSessionRepository,
	hashProvider SignInHashService,
	sessionManager SignInSessionService,
	cookieService SignInCookieService,
) SignInUseCase {
	return &signInUseCase{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		hashProvider:   hashProvider,
		sessionManager: sessionManager,
		cookieService:  cookieService,
	}
}

func (u *signInUseCase) SignIn(context context.Context, writer http.ResponseWriter, request *requests.SignIn, userAgent, ip string) (responses.SignIn, error) {
	email := request.Email
	user, err := u.userRepo.SelectByEmail(context, entities.Email(email))
	if err != nil {
		if errors.Is(err, repositories.ErrEntityNotFound) {
			return responses.SignIn{}, fmt.Errorf("failed to find user: %w", ErrEntityNotFound)
		}
		return responses.SignIn{}, fmt.Errorf("failed to find user: %w", err)
	}

	match := u.hashProvider.CompareStringAndHash(request.Password, string(user.Password))
	if !match {
		return responses.SignIn{}, fmt.Errorf("failed to compare password: %w", ErrWrongPassword)
	}

	err = u.sessionRepo.DeleteByUserId(context, user.Id)
	if err != nil {
		return responses.SignIn{}, fmt.Errorf("failed to delete session: %w", err)
	}

	session, err := u.sessionManager.CreateSession(user)
	if err != nil {
		return responses.SignIn{}, fmt.Errorf("%w: couldn't create session", err)
	}

	session.IP = ip
	session.UserAgent = userAgent

	hashedRefreshToken, err := u.hashProvider.GenerateHash(session.RefreshToken)
	if err != nil {
		return responses.SignIn{}, fmt.Errorf("%w: failed to hash refresh token", err)
	}

	refreshToken := session.RefreshToken
	session.RefreshToken = string(hashedRefreshToken)

	err = u.sessionRepo.Insert(context, session)
	if err != nil {
		return responses.SignIn{}, fmt.Errorf("%w: failed to insert session", err)
	}

	u.cookieService.Set(
		writer,
		"access_token",
		session.AccessToken,
		session.AccessExpiresAt,
	)

	refreshSessionResponse := responses.NewSession(session.AccessToken, refreshToken, session.AccessExpiresAt.Unix())

	return responses.NewSignIn(user.Id, refreshSessionResponse), nil
}
