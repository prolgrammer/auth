package usecases

import (
	"auth/internal/entities"
	"context"
)

type (
	SignInUserRepository interface {
		SelectByEmail(context.Context, entities.Email) (entities.User, error)
	}

	SignInSessionRepository interface {
		Insert(context.Context, entities.Session) error
	}

	SignInHashService interface {
		GenerateHash(stringToHash string) ([]byte, error)
		CompareStringAndHash(string, string) bool
	}

	SignInSessionService interface {
		CreateSession(account entities.User) (entities.Session, error)
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
)
