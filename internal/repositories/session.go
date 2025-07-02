package repositories

import (
	"auth/internal/entities"
	"context"
)

type SessionRepository interface {
	Insert(context context.Context, session entities.Session) error
	SelectByUserId(context context.Context, userId string) (entities.Session, error)
	Update(context context.Context, session entities.Session) error
	DeleteByUserId(context context.Context, userId string) error
}

type sessionRepository struct {
	insertSessionCommand  InsertSessionCommand
	selectUserIdCommand   SelectByRefreshTokenCommand
	updateCommand         UpdateSessionCommand
	deleteByUserIdCommand DeleteByUserIdCommand
}

func NewSessionRepository(
	insertSessionCommand InsertSessionCommand,
	selectByUserIdCommand SelectByRefreshTokenCommand,
	updateCommand UpdateSessionCommand,
	deleteByUserIdCommand DeleteByUserIdCommand,
) SessionRepository {

	return &sessionRepository{
		insertSessionCommand:  insertSessionCommand,
		selectUserIdCommand:   selectByUserIdCommand,
		updateCommand:         updateCommand,
		deleteByUserIdCommand: deleteByUserIdCommand,
	}
}

func (s *sessionRepository) Insert(context context.Context, session entities.Session) error {
	return s.insertSessionCommand.Execute(context, session)
}

func (s *sessionRepository) SelectByUserId(context context.Context, userId string) (entities.Session, error) {
	return s.selectUserIdCommand.Execute(context, userId)
}

func (s *sessionRepository) Update(context context.Context, session entities.Session) error {
	return s.updateCommand.Execute(context, session)
}

func (s *sessionRepository) DeleteByUserId(context context.Context, userId string) error {
	return s.deleteByUserIdCommand.Execute(context, userId)
}
