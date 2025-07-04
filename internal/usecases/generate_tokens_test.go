package usecases

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"auth/internal/entities"
	"auth/internal/repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	mockGenTokensUserRepo       *MockGenerateTokensUserRepository
	mockGenTokensSessionRepo    *MockGenerateTokensSessionRepository
	mockGenTokensHashService    *MockGenerateTokensHashService
	mockGenTokensSessionService *MockGenerateTokensSessionService
	mockGenTokensCookieService  *MockGenerateTokensCookieService
)

func initGenerateTokensMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockGenTokensUserRepo = NewMockGenerateTokensUserRepository(ctrl)
	mockGenTokensSessionRepo = NewMockGenerateTokensSessionRepository(ctrl)
	mockGenTokensHashService = NewMockGenerateTokensHashService(ctrl)
	mockGenTokensSessionService = NewMockGenerateTokensSessionService(ctrl)
	mockGenTokensCookieService = NewMockGenerateTokensCookieService(ctrl)
}

func TestGenerateTokensUseCase_GenerateTokens_Success(t *testing.T) {
	ctx := context.Background()
	initGenerateTokensMocks(t)

	userId := "user-id"
	ip := "127.0.0.1"
	userAgent := "test-agent"
	writer := http.ResponseWriter(nil)

	user := entities.User{
		Id:    userId,
		Email: "test@mail.ru",
	}
	session := entities.Session{
		AccessToken:     "new-access-token",
		RefreshToken:    "new-refresh-token",
		AccessExpiresAt: time.Now().Add(1 * time.Hour),
		UserId:          userId,
	}

	mockGenTokensUserRepo.EXPECT().SelectByUserId(ctx, userId).Return(user, nil)
	mockGenTokensSessionRepo.EXPECT().DeleteByUserId(ctx, userId).Return(nil)
	mockGenTokensSessionService.EXPECT().CreateSession(user).Return(session, nil)
	mockGenTokensHashService.EXPECT().GenerateHash("new-refresh-token").Return([]byte("hashed-refresh-token"), nil)
	mockGenTokensSessionRepo.EXPECT().Insert(ctx, gomock.AssignableToTypeOf(entities.Session{})).Return(nil)
	mockGenTokensCookieService.EXPECT().Set(writer, "access_token", session.AccessToken, session.AccessExpiresAt)

	useCase := NewGenerateTokensUseCase(
		mockGenTokensUserRepo,
		mockGenTokensSessionRepo,
		mockGenTokensHashService,
		mockGenTokensCookieService,
		mockGenTokensSessionService)

	result, err := useCase.GenerateTokens(ctx, writer, userId, ip, userAgent)

	assert.NoError(t, err)
	assert.Equal(t, "new-access-token", result.AccessToken)
	assert.Equal(t, "new-refresh-token", result.RefreshToken)
}

func TestGenerateTokensUseCase_GenerateTokens_UserNotFound(t *testing.T) {
	ctx := context.Background()
	initGenerateTokensMocks(t)

	userId := "failed-user"

	mockGenTokensUserRepo.EXPECT().SelectByUserId(ctx, userId).Return(entities.User{}, repositories.ErrEntityNotFound)

	useCase := NewGenerateTokensUseCase(
		mockGenTokensUserRepo,
		mockGenTokensSessionRepo,
		mockGenTokensHashService,
		mockGenTokensCookieService,
		mockGenTokensSessionService)

	result, err := useCase.GenerateTokens(ctx, nil, userId, "", "")

	assert.Error(t, err)
	assert.Empty(t, result.AccessToken)
	assert.ErrorIs(t, err, ErrEntityNotFound)
}

func TestGenerateTokensUseCase_GenerateTokens_DeleteSessionError(t *testing.T) {
	ctx := context.Background()
	initGenerateTokensMocks(t)

	userId := "user-id"
	user := entities.User{Id: userId}

	mockGenTokensUserRepo.EXPECT().SelectByUserId(ctx, userId).Return(user, nil)
	mockGenTokensSessionRepo.EXPECT().DeleteByUserId(ctx, userId).Return(fmt.Errorf("database error"))

	useCase := NewGenerateTokensUseCase(
		mockGenTokensUserRepo,
		mockGenTokensSessionRepo,
		mockGenTokensHashService,
		mockGenTokensCookieService,
		mockGenTokensSessionService)

	result, err := useCase.GenerateTokens(ctx, nil, userId, "", "")

	assert.Error(t, err)
	assert.Empty(t, result.AccessToken)
	assert.Contains(t, err.Error(), "failed to delete session")
}

func TestGenerateTokensUseCase_GenerateTokens_CreateSessionError(t *testing.T) {
	ctx := context.Background()
	initGenerateTokensMocks(t)

	userId := "user-id"
	user := entities.User{Id: userId}

	mockGenTokensUserRepo.EXPECT().SelectByUserId(ctx, userId).Return(user, nil)
	mockGenTokensSessionRepo.EXPECT().DeleteByUserId(ctx, userId).Return(nil)
	mockGenTokensSessionService.EXPECT().CreateSession(user).Return(entities.Session{}, fmt.Errorf("session error"))

	useCase := NewGenerateTokensUseCase(
		mockGenTokensUserRepo,
		mockGenTokensSessionRepo,
		mockGenTokensHashService,
		mockGenTokensCookieService,
		mockGenTokensSessionService)

	result, err := useCase.GenerateTokens(ctx, nil, userId, "", "")

	assert.Error(t, err)
	assert.Empty(t, result.AccessToken)
	assert.Contains(t, err.Error(), "couldn't create session")
}

func TestGenerateTokensUseCase_GenerateTokens_HashError(t *testing.T) {
	ctx := context.Background()
	initGenerateTokensMocks(t)

	userId := "user-id"
	user := entities.User{Id: userId}
	session := entities.Session{
		RefreshToken: "new-refresh-token",
	}

	mockGenTokensUserRepo.EXPECT().SelectByUserId(ctx, userId).Return(user, nil)
	mockGenTokensSessionRepo.EXPECT().DeleteByUserId(ctx, userId).Return(nil)
	mockGenTokensSessionService.EXPECT().CreateSession(user).Return(session, nil)
	mockGenTokensHashService.EXPECT().GenerateHash(session.RefreshToken).Return(nil, fmt.Errorf("hash error"))

	useCase := NewGenerateTokensUseCase(
		mockGenTokensUserRepo,
		mockGenTokensSessionRepo,
		mockGenTokensHashService,
		mockGenTokensCookieService,
		mockGenTokensSessionService)

	result, err := useCase.GenerateTokens(ctx, nil, userId, "", "")

	assert.Error(t, err)
	assert.Empty(t, result.AccessToken)
	assert.Contains(t, err.Error(), "failed to hash refresh token")
}
