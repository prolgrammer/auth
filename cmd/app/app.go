package app

import (
	"auth/config"
	"auth/infrastructure/postgres"
	http2 "auth/internal/controllers/http"
	"auth/internal/controllers/http/middleware"
	"auth/internal/repositories"
	"auth/internal/usecases"
	"auth/pkg"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
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

	postgresClient, err = postgres.NewClient(cfg.PG)
	if err != nil {
		fmt.Printf("couldn't start postgres: %s\n", err.Error())
		return
	}
	err = postgresClient.MigrateUp()
	if err != nil {
		if errors.Is(err, postgres.ErrNoChange) {
			fmt.Printf("postgres has the latest version. nothing to migrate\n")
			return
		}
		fmt.Printf("failed to migrate postgres: %s\n", err.Error())
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
}

func runHTTP(cfg *config.Config) {
	router := gin.Default()
	router.HandleMethodNotAllowed = true

	mw := middleware.NewMiddleware(sessionService)
	http2.InitServiceMiddleware(router)
	http2.NewSignUpController(router, signUpUseCase, mw)
	http2.NewSignInController(router, signInUseCase, mw)
	http2.NewGenerateTokensController(router, generateTokensUseCase, mw)
	http2.NewRefreshSessionController(router, refreshSessionUseCase, mw)

	address := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	fmt.Printf("starting HTTP server on %s\n", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
