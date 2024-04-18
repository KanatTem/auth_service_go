package auth

import (
	"auth_service/internal/domain/models"
	"auth_service/internal/lib/jwt"
	"auth_service/internal/lib/logger"
	"auth_service/internal/storage"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserStore interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (userId int64, err error)
	GetUser(ctx context.Context, email string) (user models.User, err error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appId int) (app models.App, err error)
}

type Auth struct {
	log         *slog.Logger
	userStore   UserStore
	appProvider AppProvider
	tokenTTL    time.Duration
}

func New(
	log *slog.Logger,
	userStore UserStore,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userStore:   userStore,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL, // Время жизни возвращаемых токенов
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (userId int64, err error) {

	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("password", pass),
	)

	//хэшим пароль
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		log.Error("Failed to hash password", logger.Err(err))
	}

	userID, err := a.userStore.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("User already exists", logger.Err(err))
		} else {
			log.Error("Failed to save user", logger.Err(err))
		}
		return 0, fmt.Errorf("%s: %w", op, err)

	}
	return userID, nil
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (token string, err error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("try to login user")

	user, err := a.userStore.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("User not found", logger.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("Failed to get user", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("Invalid password", logger.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.GetApp(ctx, appID)

	if err != nil {
		a.log.Error("Failed to get app", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Successfully logged in")

	token, err = jwt.NewToken(user, app, a.tokenTTL)

	if err != nil {
		a.log.Error("Failed to create token", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil

}
