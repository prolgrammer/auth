package usecases

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	mockLogoutSessionRepo   *MockLogoutSessionRepository
	mockLogoutCookieService *MockLogoutCookieService
)

func initLogoutMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockLogoutSessionRepo = NewMockLogoutSessionRepository(ctrl)
	mockLogoutCookieService = NewMockLogoutCookieService(ctrl)
}

func TestLogoutUseCase_Logout_Success(t *testing.T) {
	ctx := context.Background()
	initLogoutMocks(t)

	userId := "user-id"
	writer := http.ResponseWriter(nil)

	mockLogoutCookieService.EXPECT().Clear(writer, "access_token")
	mockLogoutSessionRepo.EXPECT().DeleteByUserId(ctx, userId).Return(nil)

	useCase := NewLogoutUseCase(
		mockLogoutSessionRepo,
		mockLogoutCookieService)

	err := useCase.Logout(ctx, writer, userId)

	assert.NoError(t, err)
}

func TestLogoutUseCase_Logout_DeleteError(t *testing.T) {
	ctx := context.Background()
	initLogoutMocks(t)

	userId := "user-id"
	writer := http.ResponseWriter(nil)

	mockLogoutCookieService.EXPECT().Clear(writer, "access_token")
	mockLogoutSessionRepo.EXPECT().DeleteByUserId(ctx, userId).Return(fmt.Errorf("database error"))

	useCase := NewLogoutUseCase(
		mockLogoutSessionRepo,
		mockLogoutCookieService)

	err := useCase.Logout(ctx, writer, userId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}
