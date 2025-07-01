package repositories

import (
	"auth/internal/entities"
	"context"
	"errors"
)

type userRepo struct {
	selectUserByIdCommand    SelectUserByIdCommand
	selectUserByEmailCommand SelectUserByEmailCommand
	insertUserCommand        InsertUserCommand
}

type UserRepository interface {
	Insert(context context.Context, user entities.User) (string, error)
	SelectById(context context.Context, id string) (entities.User, error)
	SelectByEmail(context context.Context, email entities.Email) (entities.User, error)
	CheckEmailExists(context context.Context, email entities.Email) (bool, error)
}

func NewUserRepository(
	selectUserByIdCommand SelectUserByIdCommand,
	selectUserByEmailCommand SelectUserByEmailCommand,
	insertUserCommand InsertUserCommand) UserRepository {
	return &userRepo{
		selectUserByIdCommand:    selectUserByIdCommand,
		selectUserByEmailCommand: selectUserByEmailCommand,
		insertUserCommand:        insertUserCommand,
	}
}

func (u *userRepo) Insert(context context.Context, user entities.User) (id string, err error) {
	return u.insertUserCommand.Execute(context, user)
}

func (u *userRepo) SelectById(context context.Context, id string) (entities.User, error) {
	return u.selectUserByIdCommand.Execute(context, id)
}

func (u *userRepo) SelectByEmail(context context.Context, email entities.Email) (entities.User, error) {
	return u.selectUserByEmailCommand.Execute(context, email)
}

func (u *userRepo) CheckEmailExists(context context.Context, email entities.Email) (bool, error) {
	_, err := u.SelectByEmail(context, email)

	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
