package users

import (
	"auth/infrastructure/postgres"
	"auth/internal/entities"
	"auth/internal/repositories"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

func selectUser(context context.Context, client *postgres.Client, sql string, args []any) (entities.User, error) {
	result := entities.User{}
	row := client.Pool.QueryRow(context, sql, args...)
	err := row.Scan(&result.Id, &result.Email, &result.Password, &result.RegistrationDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.User{}, repositories.ErrEntityNotFound
		}
		return entities.User{}, err
	}
	return result, nil
}
