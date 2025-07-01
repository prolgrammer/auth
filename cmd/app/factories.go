package app

import (
	"auth/infrastructure/postgres"
	"auth/infrastructure/postgres/commands/sessions"
	"auth/infrastructure/postgres/commands/users"
	"auth/internal/repositories"
)

func CreatePGUserRepo(client *postgres.Client) repositories.UserRepository {
	selectAccountByIdCommand := users.NewSelectUserByIdCommand(client)
	selectAccountByEmailCommand := users.NewSelectUserByEmailCommand(client)
	insertAccountCommand := users.NewInsertUserPGCommand(client)

	return repositories.NewUserRepository(
		selectAccountByIdCommand,
		selectAccountByEmailCommand,
		insertAccountCommand)
}

func CreateSessionRepo(client *postgres.Client) repositories.SessionRepository {
	selectSessionByRefreshTokenCommand := sessions.NewSelectByRefreshTokenCommand(client)
	deleteSessionByUserId := sessions.NewDeleteByUserIdCommand(client)
	insertSessionCommand := sessions.NewInsertSessionCommand(client)
	updateSessionCommand := sessions.NewUpdateSessionPGCommand(client)

	return repositories.NewSessionRepository(
		insertSessionCommand,
		selectSessionByRefreshTokenCommand,
		updateSessionCommand,
		deleteSessionByUserId,
	)
}
