package usecases

import (
	"context"
	"fmt"
	"testing"
	"time"

	"auth/internal/entities"
	"auth/internal/repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	mockGetUserRepo *MockGetUserUserRepository
)

func initGetUserMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockGetUserRepo = NewMockGetUserUserRepository(ctrl)
}

func TestGetUserUseCase_GetUserUseCase_Success(t *testing.T) {
	ctx := context.Background()
	initGetUserMocks(t)

	userId := "user-id"
	now := time.Now()
	user := entities.User{
		Id:               userId,
		Email:            "test@mail.ru",
		RegistrationDate: now,
	}

	mockGetUserRepo.EXPECT().SelectByUserId(ctx, userId).Return(user, nil)

	useCase := NewGetUserUseCase(mockGetUserRepo)

	result, err := useCase.GetUserUseCase(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, userId, result.Id)
	assert.Equal(t, "test@mail.ru", result.Email)
	assert.Equal(t, now, result.RegistrationDate)
}

func TestGetUserUseCase_GetUserUseCase_UserNotFound(t *testing.T) {
	ctx := context.Background()
	initGetUserMocks(t)

	userId := "nonexistent-user"

	mockGetUserRepo.EXPECT().SelectByUserId(ctx, userId).Return(entities.User{}, repositories.ErrEntityNotFound)

	useCase := NewGetUserUseCase(mockGetUserRepo)

	result, err := useCase.GetUserUseCase(ctx, userId)

	assert.Error(t, err)
	assert.Empty(t, result.Id)
	assert.ErrorIs(t, err, ErrEntityNotFound)
}

func TestGetUserUseCase_GetUserUseCase_RepositoryError(t *testing.T) {
	ctx := context.Background()
	initGetUserMocks(t)

	userId := "user-id"

	mockGetUserRepo.EXPECT().SelectByUserId(ctx, userId).Return(entities.User{}, fmt.Errorf("db error"))

	useCase := NewGetUserUseCase(mockGetUserRepo)

	result, err := useCase.GetUserUseCase(ctx, userId)

	assert.Error(t, err)
	assert.Empty(t, result.Id)
	assert.Contains(t, err.Error(), "failed to find account")
}
