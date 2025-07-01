package sessions

import (
	"auth/infrastructure/postgres"
	"auth/infrastructure/postgres/commands"
	"auth/internal/repositories"
	"context"
)

type deleteByUserIdCommand struct {
	client *postgres.Client
}

func NewDeleteByUserIdCommand(client *postgres.Client) repositories.DeleteByUserIdCommand {
	return &deleteByUserIdCommand{client: client}
}

func (c *deleteByUserIdCommand) Execute(ctx context.Context, userId string) error {
	sql, args, err := c.client.Builder.
		Delete(commands.SessionTable).
		Where(commands.SessionUserIdField+" = ?", userId).
		ToSql()
	if err != nil {
		return err
	}

	_, err = c.client.Pool.Exec(ctx, sql, args...)
	return err
}
