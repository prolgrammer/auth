package usecases

import (
	"auth/internal/controllers/requests"
	"auth/internal/entities"
	"auth/internal/repositories"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	mockRefreshUserRepo       *MockRefreshSessionUserRepository
	mockRefreshSessionRepo    *MockRefreshSessionSessionRepository
	mockRefreshSessionService *MockRefreshSessionSessionService
	mockRefreshCookieService  *MockRefreshSessionCookieService
	mockRefreshHashProvider   *MockRefreshSessionHashProvider
)

func initRefreshMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRefreshUserRepo = NewMockRefreshSessionUserRepository(ctrl)
	mockRefreshSessionRepo = NewMockRefreshSessionSessionRepository(ctrl)
	mockRefreshSessionService = NewMockRefreshSessionSessionService(ctrl)
	mockRefreshCookieService = NewMockRefreshSessionCookieService(ctrl)
	mockRefreshHashProvider = NewMockRefreshSessionHashProvider(ctrl)
}

func TestRefreshSessionUseCase_RefreshSession_Success(t *testing.T) {
	ctx := &gin.Context{}
	initRefreshMocks(t)

	request := requests.RefreshSession{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}
	ip := "127.0.0.1"
	userAgent := "test-agent"
	writer := http.ResponseWriter(nil)

	claims := map[string]interface{}{"sub": "user-id"}
	oldSession := entities.Session{
		UserId:       "user-id",
		UserAgent:    userAgent,
		IP:           "127.0.0.1",
		RefreshToken: "hashed-refresh-token",
	}
	user := entities.User{
		Id:    "user-id",
		Email: "test@mail.ru",
	}
	newSession := entities.Session{
		AccessToken:     "new-access-token",
		RefreshToken:    "new-refresh-token",
		AccessExpiresAt: time.Now().Add(1 * time.Hour),
		UserId:          "user-id",
	}

	mockRefreshSessionService.EXPECT().ParseToken("access-token").Return(claims, nil)
	mockRefreshSessionRepo.EXPECT().SelectByUserId(ctx, "user-id").Return(oldSession, nil)
	mockRefreshHashProvider.EXPECT().CompareStringAndHash("refresh-token", oldSession.RefreshToken).Return(true)
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	ctx.Request.AddCookie(&http.Cookie{Name: "access_token", Value: "access-token"})
	mockRefreshUserRepo.EXPECT().SelectByUserId(ctx, "user-id").Return(user, nil)
	mockRefreshSessionService.EXPECT().CreateSession(user).Return(newSession, nil)
	mockRefreshHashProvider.EXPECT().GenerateHash("new-refresh-token").Return([]byte("hashed-new-refresh"), nil)
	mockRefreshSessionRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
	mockRefreshCookieService.EXPECT().Set(writer, "access_token", newSession.AccessToken, newSession.AccessExpiresAt)

	useCase := NewRefreshSessionUseCase(
		mockRefreshUserRepo,
		mockRefreshSessionRepo,
		mockRefreshSessionService,
		mockRefreshCookieService,
		mockRefreshHashProvider)

	result, err := useCase.RefreshSession(ctx, writer, request, ip, userAgent)

	assert.NoError(t, err)
	assert.Equal(t, "new-access-token", result.AccessToken)
	assert.Equal(t, "new-refresh-token", result.RefreshToken)
}

func TestRefreshSessionUseCase_RefreshSession_InvalidToken(t *testing.T) {
	ctx := &gin.Context{}
	initRefreshMocks(t)

	request := requests.RefreshSession{
		AccessToken: "invalid-access-token",
	}

	mockRefreshSessionService.EXPECT().ParseToken("invalid-access-token").Return(nil, fmt.Errorf("invalid access token"))

	useCase := NewRefreshSessionUseCase(
		mockRefreshUserRepo,
		mockRefreshSessionRepo,
		mockRefreshSessionService,
		mockRefreshCookieService,
		mockRefreshHashProvider)

	_, err := useCase.RefreshSession(ctx, nil, request, "", "")

	assert.Error(t, err)
}

func TestRefreshSessionUseCase_RefreshSession_UserAgentMismatch(t *testing.T) {
	ctx := &gin.Context{}
	initRefreshMocks(t)

	request := requests.RefreshSession{
		AccessToken: "access-token",
	}
	claims := map[string]interface{}{"sub": "user-id"}
	oldSession := entities.Session{
		UserId:    "user-id",
		UserAgent: "another-agent",
	}

	mockRefreshSessionService.EXPECT().ParseToken("access-token").Return(claims, nil)
	mockRefreshSessionRepo.EXPECT().SelectByUserId(ctx, "user-id").Return(oldSession, nil)
	mockRefreshSessionRepo.EXPECT().DeleteByUserId(ctx, "user-id").Return(nil)
	mockRefreshCookieService.EXPECT().Clear(gomock.Any(), "access_token")

	useCase := NewRefreshSessionUseCase(
		mockRefreshUserRepo,
		mockRefreshSessionRepo,
		mockRefreshSessionService,
		mockRefreshCookieService,
		mockRefreshHashProvider)

	_, err := useCase.RefreshSession(ctx, nil, request, "", "test-agent")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUnauthorized)
}

func TestRefreshSessionUseCase_RefreshSession_InvalidRefreshToken(t *testing.T) {
	ctx := &gin.Context{}
	initRefreshMocks(t)

	request := requests.RefreshSession{
		AccessToken:  "access-token",
		RefreshToken: "invalid-refresh-token",
	}
	claims := map[string]interface{}{"sub": "user-id"}
	oldSession := entities.Session{
		UserId:       "user-id",
		UserAgent:    "test-agent",
		RefreshToken: "hashed-valid-refresh-token",
	}

	mockRefreshSessionService.EXPECT().ParseToken(request.AccessToken).Return(claims, nil)
	mockRefreshSessionRepo.EXPECT().SelectByUserId(ctx, "user-id").Return(oldSession, nil)
	mockRefreshHashProvider.EXPECT().CompareStringAndHash(request.RefreshToken, oldSession.RefreshToken).Return(false)

	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	ctx.Request.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: request.AccessToken,
	})

	useCase := NewRefreshSessionUseCase(
		mockRefreshUserRepo,
		mockRefreshSessionRepo,
		mockRefreshSessionService,
		mockRefreshCookieService,
		mockRefreshHashProvider)

	_, err := useCase.RefreshSession(ctx, nil, request, "", "test-agent")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotAValidRefreshToken)
}

func TestRefreshSessionUseCase_RefreshSession_UserNotFound(t *testing.T) {
	ctx := &gin.Context{}
	initRefreshMocks(t)

	request := requests.RefreshSession{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}
	claims := map[string]interface{}{"sub": "user-id"}
	oldSession := entities.Session{
		UserId:       "user-id",
		UserAgent:    "test-agent",
		RefreshToken: "hashed-valid-refresh-token",
	}

	mockRefreshSessionService.EXPECT().ParseToken("access-token").Return(claims, nil)
	mockRefreshSessionRepo.EXPECT().SelectByUserId(ctx, "user-id").Return(oldSession, nil)
	mockRefreshHashProvider.EXPECT().CompareStringAndHash(request.RefreshToken, "hashed-valid-refresh-token").Return(true)
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	ctx.Request.AddCookie(&http.Cookie{Name: "access_token", Value: request.AccessToken})
	mockRefreshUserRepo.EXPECT().SelectByUserId(ctx, "user-id").Return(entities.User{}, repositories.ErrEntityNotFound)

	useCase := NewRefreshSessionUseCase(
		mockRefreshUserRepo,
		mockRefreshSessionRepo,
		mockRefreshSessionService,
		mockRefreshCookieService,
		mockRefreshHashProvider)

	_, err := useCase.RefreshSession(ctx, nil, request, "", "test-agent")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEntityNotFound)
	assert.Contains(t, err.Error(), "failed to find user")
}
