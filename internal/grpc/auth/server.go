package auth

import (
	"context"
	"errors"
	ssov1 "github.com/KanatTem/sso_proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct { //имлементация api протобаффа
	ssov1.UnimplementedAuthServer
	auth Auth // логика регистрации/auth
}

type Auth interface { // интерфейс логики регистрации/auth
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userIds int64, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login( //имплеметация апи протобафа
	ctx context.Context,
	in *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {

	// валидация -> вызов логики auth.Login

	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "missing email")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "missing password")
	}
	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing appID")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{Token: token}, nil

}

func (s *serverAPI) Register(
	ctx context.Context,
	in *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "missing email")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "missing password")
	}

	userId, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &ssov1.RegisterResponse{UserId: userId}, nil

}
