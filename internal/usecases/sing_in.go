package usecases

type singInUseCase struct {
	userRepo       SignInUserRepository
	sessionRepo    SignInSessionRepository
	hashProvider   SignInHashService
	sessionManager SignInSessionService
}
