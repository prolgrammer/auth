package usecases

import (
	"auth/internal/controllers/requests"
	"auth/internal/controllers/responses"
	"auth/internal/entities"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
)

type signUpUseCase struct {
	userRepo       SignUpUserRepository
	sessionRepo    SignUpSessionRepository
	sessionManager SignUpSessionService
	hashService    SignUpHashService
	cookieService  SignUpCookieService
}

type SignUpUseCase interface {
	CreateUser(context context.Context, writer http.ResponseWriter, request requests.SignUp, userAgent, ip string) (responses.SignUp, error)
}

func NewSignUpUseCase(
	userRepo SignUpUserRepository,
	sessionRepo SignUpSessionRepository,
	sessionService SignUpSessionService,
	hashService SignUpHashService,
	cookieService SignInCookieService,
) SignUpUseCase {
	return &signUpUseCase{
		userRepo:       userRepo,
		sessionManager: sessionService,
		sessionRepo:    sessionRepo,
		hashService:    hashService,
		cookieService:  cookieService,
	}
}

func (u *signUpUseCase) CreateUser(context context.Context, writer http.ResponseWriter, request requests.SignUp, userAgent, ip string) (responses.SignUp, error) {
	user := entities.NewUser(request.Email, request.Password)
	err := user.Validate()
	if err != nil {
		return responses.SignUp{}, fmt.Errorf("%w: %w", ErrInvalidEntity, err)
	}

	exists, err := u.userRepo.CheckEmailExists(context, user.Email)
	if err != nil {
		return responses.SignUp{}, fmt.Errorf("%w: failed to check if the email is already taken", err)
	}
	if exists {
		return responses.SignUp{}, fmt.Errorf("%w: email has already been taken", ErrEntityAlreadyExists)
	}

	hashedPassword, err := u.hashService.GenerateHash(request.Password)
	if err != nil {
		return responses.SignUp{}, fmt.Errorf("%w: failed to hash the password", err)
	}

	user.Password = entities.Password(hashedPassword)

	user.Id, err = u.userRepo.Insert(context, user)
	if err != nil {
		return responses.SignUp{}, fmt.Errorf("%w: failed to insert user", err)
	}

	session, err := u.sessionManager.CreateSession(user)
	if err != nil {
		return responses.SignUp{}, fmt.Errorf("%w: failed to create session", err)
	}
	session.IP = ip
	session.UserAgent = userAgent

	hashedRefreshToken, err := u.hashService.GenerateHash(session.RefreshToken)
	if err != nil {
		return responses.SignUp{}, fmt.Errorf("%w: failed to hash refresh token", err)
	}

	encodedRefreshToken := base64.StdEncoding.EncodeToString([]byte(session.RefreshToken))
	session.RefreshToken = string(hashedRefreshToken)

	err = u.sessionRepo.Insert(context, session)
	if err != nil {
		return responses.SignUp{}, fmt.Errorf("%w: failed to insert device", err)
	}

	refreshSessionResponse :=
		responses.NewSession(session.AccessToken, encodedRefreshToken, session.AccessExpiresAt.Unix())

	u.cookieService.Set(
		writer,
		"access_token",
		session.AccessToken,
		session.AccessExpiresAt,
	)

	return responses.NewSignUp(user.Id, refreshSessionResponse), nil
}
