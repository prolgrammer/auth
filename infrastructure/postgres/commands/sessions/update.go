package sessions

import (
	"auth/infrastructure/postgres"
	"auth/infrastructure/postgres/commands"
	"auth/internal/entities"
	"auth/internal/repositories"
	"context"
)

type updateSessionPGCommand struct {
	client *postgres.Client
}

func NewUpdateSessionPGCommand(client *postgres.Client) repositories.UpdateSessionCommand {
	return &updateSessionPGCommand{client: client}
}

func (c *updateSessionPGCommand) Execute(ctx context.Context, session entities.Session) error {
	sql, args, err := c.client.Builder.
		Update(commands.SessionTable).
		Set(commands.SessionRefreshTokenHash, session.RefreshToken).
		Set(commands.SessionUserAgentField, session.UserAgent).
		Set(commands.SessionIPField, session.IP).
		Set(commands.SessionExpiresAtField, session.ExpiresAt).
		Where(commands.SessionUserIdField+" = ?", session.UserId).
		ToSql()
	if err != nil {
		return err
	}

	_, err = c.client.Pool.Exec(ctx, sql, args...)
	return err
}
