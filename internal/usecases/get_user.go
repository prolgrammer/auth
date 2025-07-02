package usecases

import (
	"auth/internal/controllers/responses"
	"auth/internal/repositories"
	"context"
	"errors"
	"fmt"
)

type getUserUseCase struct {
	userRepo GetUserUserRepository
}

type GetUserUseCase interface {
	GetUserUseCase(context context.Context, userId string) (responses.User, error)
}

func NewGetUserUseCase(userRepo GetUserUserRepository) GetUserUseCase {
	return &getUserUseCase{userRepo: userRepo}
}

func (g getUserUseCase) GetUserUseCase(context context.Context, userId string) (responses.User, error) {
	user, err := g.userRepo.SelectByUserId(context, userId)
	if err != nil {
		if errors.Is(err, repositories.ErrEntityNotFound) {
			return responses.User{}, fmt.Errorf("%w account not found", ErrEntityNotFound)
		}
		return responses.User{}, fmt.Errorf("%w failed to find account", err)
	}

	userResponse := responses.User{
		Id:               user.Id,
		Email:            string(user.Email),
		RegistrationDate: user.RegistrationDate,
	}

	return userResponse, nil
}
