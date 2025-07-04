package repositories

import (
	"auth/internal/entities"
	"context"
)

type (
	InsertUserCommand interface {
		Execute(context context.Context, user entities.User) (string, error)
	}
	SelectUserByIdCommand interface {
		Execute(context context.Context, id string) (entities.User, error)
	}
	SelectUserByEmailCommand interface {
		Execute(context context.Context, email entities.Email) (entities.User, error)
	}
)

type (
	InsertSessionCommand interface {
		Execute(ctx context.Context, session entities.Session) error
	}

	SelectByUserIdCommand interface {
		Execute(ctx context.Context, refreshToken string) (entities.Session, error)
	}

	DeleteByUserIdCommand interface {
		Execute(ctx context.Context, userId string) error
	}

	UpdateSessionCommand interface {
		Execute(ctx context.Context, session entities.Session) error
	}
)
