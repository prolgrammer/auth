package sessions

import (
	"auth/infrastructure/postgres"
	"auth/infrastructure/postgres/commands"
	"auth/internal/entities"
	"auth/internal/repositories"
	"context"
	"github.com/google/uuid"
)

type selectByRefreshTokenCommand struct {
	client *postgres.Client
}

func NewSelectByUserIdCommand(client *postgres.Client) repositories.SelectByRefreshTokenCommand {
	return &selectByRefreshTokenCommand{client: client}
}

func (c *selectByRefreshTokenCommand) Execute(ctx context.Context, userId string) (entities.Session, error) {

	uuidUser, err := uuid.Parse(userId)
	if err != nil {
		return entities.Session{}, err
	}
	sql, args, err := c.client.Builder.
		Select(
			commands.SessionIdField,
			commands.SessionUserIdField,
			commands.SessionRefreshTokenHash,
			commands.SessionUserAgentField,
			commands.SessionIPField,
			commands.SessionExpiresAtField,
		).
		From(commands.SessionTable).
		Where(commands.SessionUserIdField+" = ?", uuidUser).
		ToSql()
	if err != nil {
		return entities.Session{}, err
	}

	var session entities.Session
	err = c.client.Pool.QueryRow(ctx, sql, args...).Scan(
		&session.Id,
		&session.UserId,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IP,
		&session.ExpiresAt,
	)
	return session, err
}
