package usecases

import (
	"auth/internal/entities"
	"context"
	"net/http"
	"time"
)

type (
	SignInUserRepository interface {
		SelectByEmail(context.Context, entities.Email) (entities.User, error)
	}

	SignInSessionRepository interface {
		Insert(context.Context, entities.Session) error
		DeleteByUserId(context.Context, string) error
	}

	SignInHashService interface {
		GenerateHash(stringToHash string) ([]byte, error)
		CompareStringAndHash(string, string) bool
	}

	SignInSessionService interface {
		CreateSession(account entities.User) (entities.Session, error)
	}

	SignInCookieService interface {
		Set(w http.ResponseWriter, name, value string, expires time.Time)
	}

	SignUpUserRepository interface {
		CheckEmailExists(context.Context, entities.Email) (bool, error)
		Insert(context.Context, entities.User) (string, error)
	}

	SignUpSessionRepository interface {
		Insert(context.Context, entities.Session) error
	}

	SignUpSessionService interface {
		CreateSession(user entities.User) (entities.Session, error)
	}

	SignUpHashService interface {
		GenerateHash(stringToHash string) ([]byte, error)
	}

	SignUpCookieService interface {
		Set(w http.ResponseWriter, name, value string, expires time.Time)
	}
)

type (
	GenerateTokensUserRepository interface {
		SelectByUserId(context.Context, string) (entities.User, error)
	}

	GenerateTokensSessionRepository interface {
		Insert(context.Context, entities.Session) error
		DeleteByUserId(context.Context, string) error
	}

	GenerateTokensHashService interface {
		GenerateHash(stringToHash string) ([]byte, error)
	}

	GenerateTokensCookieService interface {
		Set(w http.ResponseWriter, name, value string, expires time.Time)
	}

	GenerateTokensSessionService interface {
		CreateSession(account entities.User) (entities.Session, error)
	}
)

type (
	RefreshSessionUserRepository interface {
		SelectByUserId(context.Context, string) (entities.User, error)
	}

	RefreshSessionSessionRepository interface {
		SelectByUserId(context.Context, string) (entities.Session, error)
		DeleteByUserId(context.Context, string) error
		Update(context context.Context, session entities.Session) error
	}

	RefreshSessionSessionService interface {
		ParseToken(token string) (entities.AccessTokenClaims, error)
		CreateSession(account entities.User) (entities.Session, error)
	}

	RefreshSessionCookieService interface {
		Set(w http.ResponseWriter, name, value string, expires time.Time)
	}

	RefreshSessionHashProvider interface {
		GenerateHash(stringToHash string) ([]byte, error)
		CompareStringAndHash(string, string) bool
	}
)

type (
	GetUserUserRepository interface {
		SelectByUserId(context.Context, string) (entities.User, error)
	}
)

type (
	LogoutSessionRepository interface {
		DeleteByUserId(context.Context, string) error
	}
	LogoutCookieService interface {
		Clear(w http.ResponseWriter, name string)
	}
	LogoutSessionService interface {
		ParseToken(token string) (entities.AccessTokenClaims, error)
	}
)
