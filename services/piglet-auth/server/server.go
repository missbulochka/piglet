package auth

import (
	"context"
	"google.golang.org/grpc"
	authv1 "piglet/internal/proto/grpc/auth/service"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	panic("Login: implement me")
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
	panic("Register: implement me")
}

func (s *serverAPI) UpdateUser(
	ctx context.Context,
	req *authv1.UpdateUserRequest,
) (*authv1.UpdateUserResponse, error) {
	panic("Update: implement me")
}
