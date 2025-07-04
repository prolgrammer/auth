package usecases

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
	"time"

	"auth/internal/controllers/requests"
	"auth/internal/entities"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	mockSignUpUserRepo       *MockSignUpUserRepository
	mockSignUpSessionRepo    *MockSignUpSessionRepository
	mockSignUpHashService    *MockSignUpHashService
	mockSignUpSessionService *MockSignUpSessionService
	mockSignUpCookieService  *MockSignUpCookieService
)

func initSignUpMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSignUpUserRepo = NewMockSignUpUserRepository(ctrl)
	mockSignUpSessionRepo = NewMockSignUpSessionRepository(ctrl)
	mockSignUpHashService = NewMockSignUpHashService(ctrl)
	mockSignUpSessionService = NewMockSignUpSessionService(ctrl)
	mockSignUpCookieService = NewMockSignUpCookieService(ctrl)
}

func TestSignUpUseCase_CreateUser_Success(t *testing.T) {
	ctx := context.Background()
	initSignUpMocks(t)

	request := requests.SignUp{
		Email:    "test@mail.ru",
		Password: "password123",
	}
	userAgent := "test-agent"
	ip := "127.0.0.1"

	session := entities.Session{
		AccessToken:     "new-access-token",
		RefreshToken:    "new-refresh-token",
		AccessExpiresAt: time.Now().Add(1 * time.Hour),
		UserId:          "new-user-id",
		IP:              ip,
		UserAgent:       userAgent,
	}

	expectedUserId := "new-user-id"
	expectedAccessToken := "new-access-token"
	expectedRefreshToken := "new-refresh-token"
	writer := http.ResponseWriter(nil)

	mockSignUpUserRepo.EXPECT().CheckEmailExists(ctx, entities.Email(request.Email)).Return(false, nil)
	mockSignUpHashService.EXPECT().GenerateHash(request.Password).Return([]byte("hashedpassword"), nil)
	mockSignUpUserRepo.EXPECT().Insert(ctx, gomock.AssignableToTypeOf(entities.User{})).Return(expectedUserId, nil)
	mockSignUpSessionService.EXPECT().CreateSession(gomock.AssignableToTypeOf(entities.User{})).Return(session, nil)
	mockSignUpHashService.EXPECT().GenerateHash(session.RefreshToken).Return([]byte("hashed-refresh-token"), nil)
	mockSignUpSessionRepo.EXPECT().Insert(ctx, gomock.AssignableToTypeOf(entities.Session{})).Return(nil)
	mockSignUpCookieService.EXPECT().Set(writer, "access_token", expectedAccessToken, session.AccessExpiresAt)

	useCase := NewSignUpUseCase(
		mockSignUpUserRepo,
		mockSignUpSessionRepo,
		mockSignUpSessionService,
		mockSignUpHashService,
		mockSignUpCookieService)

	response, err := useCase.CreateUser(ctx, writer, request, userAgent, ip)

	assert.NoError(t, err)
	assert.Equal(t, expectedAccessToken, response.Session.AccessToken)
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(expectedRefreshToken)), response.Session.RefreshToken)
}

func TestSignUpUseCase_CreateUser_InvalidEntity(t *testing.T) {
	ctx := context.Background()
	initSignUpMocks(t)

	request := requests.SignUp{
		Email:    "invalid-email",
		Password: "password123",
	}

	useCase := signUpUseCase{}

	response, err := useCase.CreateUser(ctx, nil, request, "", "")

	assert.Error(t, err)
	assert.Empty(t, response.Id)
	assert.ErrorIs(t, err, ErrInvalidEntity)
}

func TestSignUpUseCase_CreateUser_EmailExists(t *testing.T) {
	ctx := context.Background()
	initSignUpMocks(t)

	request := requests.SignUp{
		Email:    "exists@mail.ru",
		Password: "password123",
	}

	useCase := NewSignUpUseCase(
		mockSignUpUserRepo,
		mockSignUpSessionRepo,
		mockSignUpSessionService,
		mockSignUpHashService,
		mockSignUpCookieService)

	mockSignUpUserRepo.EXPECT().CheckEmailExists(ctx, entities.Email("exists@mail.ru")).Return(true, nil)

	response, err := useCase.CreateUser(ctx, nil, request, "", "")

	assert.Error(t, err)
	assert.Empty(t, response.Id)
	assert.ErrorIs(t, err, ErrEntityAlreadyExists)
}

func TestSignUpUseCase_CreateUser_HashError(t *testing.T) {
	ctx := context.Background()
	initSignUpMocks(t)

	request := requests.SignUp{
		Email:    "test@mail.ru",
		Password: "password123",
	}

	mockSignUpUserRepo.EXPECT().CheckEmailExists(ctx, entities.Email("test@mail.ru")).Return(false, nil)
	mockSignUpHashService.EXPECT().GenerateHash("password123").Return(nil, fmt.Errorf("hash error"))

	useCase := NewSignUpUseCase(
		mockSignUpUserRepo,
		mockSignUpSessionRepo,
		mockSignUpSessionService,
		mockSignUpHashService,
		mockSignUpCookieService)

	response, err := useCase.CreateUser(ctx, nil, request, "", "")

	assert.Error(t, err)
	assert.Empty(t, response.Id)
	assert.Contains(t, err.Error(), "failed to hash the password")
}

func TestSignUpUseCase_CreateUser_InsertError(t *testing.T) {
	ctx := context.Background()
	initSignUpMocks(t)

	useCase := signUpUseCase{
		userRepo:    mockSignUpUserRepo,
		hashService: mockSignUpHashService,
	}

	request := requests.SignUp{
		Email:    "test@mail.ru",
		Password: "password123",
	}

	mockSignUpUserRepo.EXPECT().CheckEmailExists(ctx, entities.Email("test@mail.ru")).Return(false, nil)
	mockSignUpHashService.EXPECT().GenerateHash("password123").Return([]byte("hashedpassword"), nil)
	mockSignUpUserRepo.EXPECT().Insert(ctx, gomock.Any()).Return("", fmt.Errorf("insert error"))

	response, err := useCase.CreateUser(ctx, nil, request, "", "")

	assert.Error(t, err)
	assert.Empty(t, response.Id)
	assert.Contains(t, err.Error(), "failed to insert user")
}
