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
		CompareStringAndHash(string, string) (bool, error)
	}

	SignInSessionService interface {
		CreateSession(account entities.User) (entities.Session, error)
	}
)
