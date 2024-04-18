package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	authv1 "piglet-auth-service/api/proto/gen"
)

type Auth interface {
	Login(
		ctx context.Context,
		username string,
		password string,
	) (token string, err error)
	RegisterUser(
		ctx context.Context,
		username string,
		email string,
		password string,
	) (userID int64, err error)
	UpdateUser(
		ctx context.Context,
		username string,
		email string,
		oldPassword string,
		password string,
	) (success bool, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	// TODO: setup validation
	if req.GetUsername() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username is empty")
	}
	if req.GetPass() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is empty")
	}

	token, err := s.auth.Login(ctx, req.GetUsername(), req.GetPass())
	if err != nil {
		// TODO ...
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) RegisterUser(
	ctx context.Context,
	req *authv1.RegisterUserRequest,
) (*authv1.RegisterUserResponse, error) {
	// TODO: setup validation
	if req.GetUsername() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username is empty")
	}
	if req.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is empty")
	}
	if req.GetPass() == "" || req.GetConPass() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is empty")
	}
	if req.GetPass() != req.GetConPass() {
		return nil, status.Errorf(codes.InvalidArgument, "password mismatch")
	}

	userID, err := s.auth.RegisterUser(ctx, req.GetUsername(), req.GetEmail(), req.GetPass())
	if err != nil {
		// TODO ...
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	return &authv1.RegisterUserResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) UpdateUser(
	ctx context.Context,
	req *authv1.UpdateUserRequest,
) (*authv1.UpdateUserResponse, error) {
	// TODO: setup validation
	if req.GetUsername() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username is empty")
	}
	if req.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is empty")
	}
	if req.GetOldPass() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "old password is empty")
	}
	if req.GetPass() == "" || req.GetConPass() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is empty")
	}
	if req.GetPass() != req.GetConPass() {
		return nil, status.Errorf(codes.InvalidArgument, "password mismatch")
	}

	success, err := s.auth.UpdateUser(
		ctx,
		req.GetUsername(),
		req.GetEmail(),
		req.GetOldPass(),
		req.GetPass(),
	)
	if err != nil {
		// TODO ...
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	return &authv1.UpdateUserResponse{
		Success: success,
	}, err
}
