package users

import (
	"auth/infrastructure/postgres"
	"auth/infrastructure/postgres/commands"
	"auth/internal/entities"
	"auth/internal/repositories"
	"context"
)

type insertUserPGCommand struct {
	client *postgres.Client
}

func NewInsertUserPGCommand(client *postgres.Client) repositories.InsertUserCommand {
	return &insertUserPGCommand{client: client}
}

func (c *insertUserPGCommand) Execute(context context.Context, user entities.User) (string, error) {
	sql, args, err := c.client.Builder.Insert(commands.UserTable).
		Columns(
			commands.UserEmailField,
			commands.UserPasswordField).
		Values(user.Email, user.Password).
		Suffix("RETURNING " + commands.UserIdField).
		ToSql()
	if err != nil {
		return "", err
	}

	var id string
	err = c.client.Pool.QueryRow(context, sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
