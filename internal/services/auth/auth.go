package auth

import (
	"auth_service/internal/domain/models"
	"golang.org/x/net/context"
	"log/slog"
	"time"
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

func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (userID int64, err error) {

	const op = "Auth.RegisterNewUser"

}
