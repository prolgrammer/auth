package usecases

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"auth/internal/controllers/requests"
	"auth/internal/entities"
	"auth/internal/repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	mockSignInUserRepo       *MockSignInUserRepository
	mockSignInSessionRepo    *MockSignInSessionRepository
	mockSignInHashService    *MockSignInHashService
	mockSignInSessionService *MockSignInSessionService
	mockSignInCookieService  *MockSignInCookieService
)

func initSignInMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSignInUserRepo = NewMockSignInUserRepository(ctrl)
	mockSignInSessionRepo = NewMockSignInSessionRepository(ctrl)
	mockSignInHashService = NewMockSignInHashService(ctrl)
	mockSignInSessionService = NewMockSignInSessionService(ctrl)
	mockSignInCookieService = NewMockSignInCookieService(ctrl)
}

func TestSignInUseCase_SignIn_Success(t *testing.T) {
	ctx := context.Background()
	initSignInMocks(t)

	request := &requests.SignIn{
		Email:    "test@mail.ru",
		Password: "password123",
	}
	userAgent := "test-agent"
	ip := "127.0.0.1"
	writer := http.ResponseWriter(nil)

	user := entities.User{
		Id:       "user-id",
		Email:    entities.Email("test@mail.ru"),
		Password: entities.Password("hashed-password"),
	}

	expectedAccessToken := "new-access-token"
	expectedRefreshToken := "new-refresh-token"

	session := entities.Session{
		AccessToken:     expectedAccessToken,
		RefreshToken:    expectedRefreshToken,
		AccessExpiresAt: time.Now().Add(1 * time.Hour),
		UserId:          "user-id",
		IP:              ip,
		UserAgent:       userAgent,
	}

	mockSignInUserRepo.EXPECT().SelectByEmail(ctx, entities.Email("test@mail.ru")).Return(user, nil)
	mockSignInHashService.EXPECT().CompareStringAndHash("password123", string(user.Password)).Return(true)
	mockSignInSessionRepo.EXPECT().DeleteByUserId(ctx, "user-id").Return(nil)
	mockSignInSessionService.EXPECT().CreateSession(user).Return(session, nil)
	mockSignInHashService.EXPECT().GenerateHash("new-refresh-token").Return([]byte("hashed-refresh-token"), nil)
	mockSignInSessionRepo.EXPECT().Insert(ctx, gomock.AssignableToTypeOf(entities.Session{})).Return(nil)
	mockSignInCookieService.EXPECT().Set(writer, "access_token", expectedAccessToken, session.AccessExpiresAt)

	useCase := NewSignInUseCase(
		mockSignInUserRepo,
		mockSignInSessionRepo,
		mockSignInHashService,
		mockSignInSessionService,
		mockSignInCookieService)

	response, err := useCase.SignIn(ctx, writer, request, userAgent, ip)

	assert.NoError(t, err)
	assert.Equal(t, "user-id", response.Id)
	assert.Equal(t, expectedAccessToken, response.Session.AccessToken)
	assert.Equal(t, expectedRefreshToken, response.Session.RefreshToken)
}

func TestSignInUseCase_SignIn_UserNotFound(t *testing.T) {
	ctx := context.Background()
	initSignInMocks(t)

	request := &requests.SignIn{
		Email:    "nobody@mail.ru",
		Password: "password123",
	}

	mockSignInUserRepo.EXPECT().SelectByEmail(ctx, entities.Email("nobody@mail.ru")).Return(entities.User{}, repositories.ErrEntityNotFound)

	useCase := NewSignInUseCase(
		mockSignInUserRepo,
		mockSignInSessionRepo,
		mockSignInHashService,
		mockSignInSessionService,
		mockSignInCookieService)

	response, err := useCase.SignIn(ctx, nil, request, "", "")

	assert.Error(t, err)
	assert.Empty(t, response.Id)
	assert.ErrorIs(t, err, ErrEntityNotFound)
}

func TestSignInUseCase_SignIn_WrongPassword(t *testing.T) {
	ctx := context.Background()
	initSignInMocks(t)

	request := &requests.SignIn{
		Email:    "test@mail.ru",
		Password: "wrong-password",
	}
	user := entities.User{
		Id:       "user-id",
		Email:    "test@mail.ru",
		Password: "hashed-password",
	}

	mockSignInUserRepo.EXPECT().SelectByEmail(ctx, entities.Email("test@mail.ru")).Return(user, nil)
	mockSignInHashService.EXPECT().CompareStringAndHash("wrong-password", string(user.Password)).Return(false)

	useCase := NewSignInUseCase(
		mockSignInUserRepo,
		mockSignInSessionRepo,
		mockSignInHashService,
		mockSignInSessionService,
		mockSignInCookieService)

	response, err := useCase.SignIn(ctx, nil, request, "", "")

	assert.Error(t, err)
	assert.Empty(t, response.Id)
	assert.ErrorIs(t, err, ErrWrongPassword)
}

func TestSignInUseCase_SignIn_DeleteSessionError(t *testing.T) {
	ctx := context.Background()
	initSignInMocks(t)

	request := &requests.SignIn{
		Email:    "test@mail.ru",
		Password: "password123",
	}
	user := entities.User{
		Id:       "user-id",
		Email:    "test@mail.ru",
		Password: "hashed-password",
	}

	mockSignInUserRepo.EXPECT().SelectByEmail(ctx, entities.Email("test@mail.ru")).Return(user, nil)
	mockSignInHashService.EXPECT().CompareStringAndHash("password123", string(user.Password)).Return(true)
	mockSignInSessionRepo.EXPECT().DeleteByUserId(ctx, "user-id").Return(fmt.Errorf("database error"))

	useCase := NewSignInUseCase(
		mockSignInUserRepo,
		mockSignInSessionRepo,
		mockSignInHashService,
		mockSignInSessionService,
		mockSignInCookieService)

	response, err := useCase.SignIn(ctx, nil, request, "", "")

	assert.Error(t, err)
	assert.Empty(t, response.Id)
	assert.Contains(t, err.Error(), "failed to delete session")
}

func TestSignInUseCase_SignIn_CreateSessionError(t *testing.T) {
	ctx := context.Background()
	initSignInMocks(t)

	request := &requests.SignIn{
		Email:    "test@mail.ru",
		Password: "password123",
	}
	user := entities.User{
		Id:       "user-id",
		Email:    "test@mail.ru",
		Password: entities.Password("hashed-password"),
	}

	mockSignInUserRepo.EXPECT().SelectByEmail(ctx, entities.Email("test@mail.ru")).Return(user, nil)
	mockSignInHashService.EXPECT().CompareStringAndHash("password123", string(user.Password)).Return(true)
	mockSignInSessionRepo.EXPECT().DeleteByUserId(ctx, "user-id").Return(nil)
	mockSignInSessionService.EXPECT().CreateSession(user).Return(entities.Session{}, fmt.Errorf("session error"))

	useCase := NewSignInUseCase(
		mockSignInUserRepo,
		mockSignInSessionRepo,
		mockSignInHashService,
		mockSignInSessionService,
		mockSignInCookieService)

	response, err := useCase.SignIn(ctx, nil, request, "", "")

	assert.Error(t, err)
	assert.Empty(t, response.Id)
	assert.Contains(t, err.Error(), "couldn't create session")
}
