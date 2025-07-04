package app

import (
	"auth/config"
	"auth/infrastructure/postgres"
	http2 "auth/internal/controllers/http"
	"auth/internal/controllers/http/middleware"
	"auth/internal/repositories"
	"auth/internal/usecases"
	"auth/pkg"
	"auth/pkg/logger"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	l              logger.Logger
	postgresClient *postgres.Client

	bcryptHashService pkg.HashService
	sessionService    pkg.SessionService
	cookieService     pkg.CookieService

	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository

	signInUseCase         usecases.SignInUseCase
	signUpUseCase         usecases.SignUpUseCase
	generateTokensUseCase usecases.GenerateTokensUseCase
	refreshSessionUseCase usecases.RefreshSessionUseCase
	getUserUseCase        usecases.GetUserUseCase
	logoutUserUseCase     usecases.LogoutUseCase
)

func Run() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	initPackages(cfg)
	initService(cfg)
	initRepository()
	initUseCases()

	defer postgresClient.Close()
	runHTTP(cfg)
}

func initPackages(cfg *config.Config) {
	var err error

	l = logger.NewConsoleLogger(logger.LevelSwitch(cfg.LogLevel))

	l.Info().Msgf("starting postgres client")
	postgresClient, err = postgres.NewClient(cfg.PG, l)
	if err != nil {
		l.Fatal().Msgf("couldn't start postgres: %s", err.Error())
		return
	}
	err = postgresClient.MigrateUp()
	if err != nil {
		if errors.Is(err, postgres.ErrNoChange) {
			l.Info().Msgf("postgres has the latest version. nothing to migrate")
			return
		}
		l.Fatal().Msgf("failed to migrate postgres: %s", err.Error())
	}
}

func initService(cfg *config.Config) {
	bcryptHashService = pkg.NewBcryptHashService()

	accessTokenService := pkg.NewTokenService([]byte(cfg.SecretKey))
	refreshTokenService := pkg.NewTokenService([]byte(cfg.SecretKey))
	sessionService = pkg.NewSessionService(cfg.TokenConfiguration, accessTokenService, refreshTokenService)

	cookieService = pkg.NewCookieService(cfg.Cookie)
}

func initRepository() {
	userRepository = CreatePGUserRepo(postgresClient)
	sessionRepository = CreateSessionRepo(postgresClient)
}

func initUseCases() {
	signUpUseCase = usecases.NewSignUpUseCase(
		userRepository,
		sessionRepository,
		sessionService,
		bcryptHashService,
		cookieService,
	)

	signInUseCase = usecases.NewSignInUseCase(
		userRepository,
		sessionRepository,
		bcryptHashService,
		sessionService,
		cookieService,
	)

	generateTokensUseCase = usecases.NewGenerateTokensUseCase(
		userRepository,
		sessionRepository,
		bcryptHashService,
		cookieService,
		sessionService,
	)

	refreshSessionUseCase = usecases.NewRefreshSessionUseCase(
		userRepository,
		sessionRepository,
		sessionService,
		cookieService,
		bcryptHashService,
	)

	getUserUseCase = usecases.NewGetUserUseCase(userRepository)

	logoutUserUseCase = usecases.NewLogoutUseCase(
		sessionRepository,
		cookieService,
	)
}

func runHTTP(cfg *config.Config) {
	router := gin.Default()
	router.HandleMethodNotAllowed = true

	mw := middleware.NewMiddleware(sessionService, l)
	http2.InitServiceMiddleware(router)
	http2.NewSignUpController(router, signUpUseCase, mw, l)
	http2.NewSignInController(router, signInUseCase, mw, l)
	http2.NewGenerateTokensController(router, generateTokensUseCase, mw, l)
	http2.NewRefreshSessionController(router, refreshSessionUseCase, mw, l)
	http2.NewGetUserController(router, getUserUseCase, mw, l)
	http2.NewLogoutController(router, logoutUserUseCase, mw, l)
	http2.NewWebhookController(router, mw, l)

	address := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	l.Info().Msgf("starting HTTP server on %s", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
