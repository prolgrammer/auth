package users

import (
	"auth/infrastructure/postgres"
	"auth/infrastructure/postgres/commands"
	"auth/internal/entities"
	"auth/internal/repositories"
	"context"
	sq "github.com/Masterminds/squirrel"
)

type selectUserByEmailCommand struct {
	client *postgres.Client
}

func NewSelectUserByEmailCommand(client *postgres.Client) repositories.SelectUserByEmailCommand {
	return &selectUserByEmailCommand{client: client}
}

func (s *selectUserByEmailCommand) Execute(context context.Context, email entities.Email) (entities.User, error) {
	sql, args, err := s.client.Builder.
		Select(
			commands.UserIdField,
			commands.UserEmailField,
			commands.UserPasswordField,
			commands.UserCreatedAtField,
		).
		From(commands.UserTable).
		Where(sq.Eq{commands.UserEmailField: email}).
		ToSql()
	if err != nil {
		return entities.User{}, err
	}

	return selectUser(context, s.client, sql, args)
}
