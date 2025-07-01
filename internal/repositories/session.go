package repositories

import (
	"auth/internal/entities"
	"context"
)

type SessionRepository interface {
	Insert(context context.Context, session entities.Session) error
	SelectByRefreshToken(context context.Context, refreshToken string) (entities.Session, error)
	Update(context context.Context, session entities.Session) error
	DeleteByUserId(context context.Context, userId string) error
}

type sessionRepository struct {
	insertSessionCommand        InsertSessionCommand
	selectByRefreshTokenCommand SelectByRefreshTokenCommand
	updateCommand               UpdateSessionCommand
	deleteByUserIdCommand       DeleteByUserIdCommand
}

func NewSessionRepository(
	insertSessionCommand InsertSessionCommand,
	selectByRefreshTokenCommand SelectByRefreshTokenCommand,
	updateCommand UpdateSessionCommand,
	deleteByUserIdCommand DeleteByUserIdCommand,
) SessionRepository {

	return &sessionRepository{
		insertSessionCommand:        insertSessionCommand,
		selectByRefreshTokenCommand: selectByRefreshTokenCommand,
		updateCommand:               updateCommand,
		deleteByUserIdCommand:       deleteByUserIdCommand,
	}
}

func (s *sessionRepository) Insert(context context.Context, session entities.Session) error {
	return s.insertSessionCommand.Execute(context, session)
}

func (s *sessionRepository) SelectByRefreshToken(context context.Context, refreshToken string) (entities.Session, error) {
	return s.selectByRefreshTokenCommand.Execute(context, refreshToken)
}

func (s *sessionRepository) Update(context context.Context, session entities.Session) error {
	return s.updateCommand.Execute(context, session)
}

func (s *sessionRepository) DeleteByUserId(context context.Context, userId string) error {
	return s.deleteByUserIdCommand.Execute(context, userId)
}
