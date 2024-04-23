package suite

import (
	"auth_service/internal/config"
	"context"
	ssov1 "github.com/KanatTem/sso_proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config   // Конфигурация приложения
	AuthClient ssov1.AuthClient // Клиент для взаимодействия с gRPC-сервером Auth
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {

	t.Helper()
	t.Parallel()

	cfg := config.MustLoad()

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	// Адрес gRPC-сервера
	grpcAddress := net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))

	cc, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("failed to create grpc client", err)
	}

	authClient := ssov1.NewAuthClient(cc)

	return ctx, &Suite{t, cfg, authClient}

}
