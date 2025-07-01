package users

import (
	"auth/infrastructure/postgres"
	"auth/infrastructure/postgres/commands"
	"auth/internal/entities"
	"auth/internal/repositories"
	"context"
	sq "github.com/Masterminds/squirrel"
)

type selectUserByIdCommand struct {
	client *postgres.Client
}

func NewSelectUserByIdCommand(client *postgres.Client) repositories.SelectUserByIdCommand {
	return &selectUserByIdCommand{client: client}
}

func (s *selectUserByIdCommand) Execute(context context.Context, id string) (entities.User, error) {
	sql, args, err := s.client.Builder.Select(
		commands.UserIdField,
		commands.UserEmailField,
		commands.UserPasswordField,
		commands.UserCreatedAtField).
		From(commands.UserTable).
		Where(sq.Eq{commands.UserIdField: id}).
		ToSql()
	if err != nil {
		return entities.User{}, err
	}

	return selectUser(context, s.client, sql, args)
}
