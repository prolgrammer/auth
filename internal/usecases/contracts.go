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
	GetTokensUserRepository interface {
		SelectByUserId(context.Context, string) (entities.User, error)
	}

	GetTokensSessionRepository interface {
		Insert(context.Context, entities.Session) error
		DeleteByUserId(context.Context, string) error
	}

	GetTokensHashService interface {
		GenerateHash(stringToHash string) ([]byte, error)
	}

	GetTokensCookieService interface {
		Set(w http.ResponseWriter, name, value string, expires time.Time)
	}

	GetTokensSessionService interface {
		CreateSession(account entities.User) (entities.Session, error)
	}
)
