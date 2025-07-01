package sessions

import (
	"auth/infrastructure/postgres"
	"auth/infrastructure/postgres/commands"
	"auth/internal/entities"
	"auth/internal/repositories"
	"context"
)

type insertSessionCommand struct {
	client *postgres.Client
}

func NewInsertSessionCommand(client *postgres.Client) repositories.InsertSessionCommand {
	return &insertSessionCommand{client: client}
}

func (c *insertSessionCommand) Execute(ctx context.Context, session entities.Session) error {
	sql, args, err := c.client.Builder.
		Insert(commands.SessionTable).
		Columns(
			commands.SessionUserIdField,
			commands.SessionRefreshTokenHash,
			commands.SessionUserAgentField,
			commands.SessionIPField,
			commands.SessionExpiresAtField,
		).
		Values(
			session.UserId,
			session.RefreshToken,
			session.UserAgent,
			session.IP,
			session.ExpiresAt,
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = c.client.Pool.Exec(ctx, sql, args...)
	return err
}
