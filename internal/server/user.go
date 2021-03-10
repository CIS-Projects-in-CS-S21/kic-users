package server

import (
	"context"

	"go.uber.org/zap"

	pbcommon "github.com/kic/users/pkg/proto/common"
	pbusers "github.com/kic/users/pkg/proto/users"
)

type UsersService struct {
	pbusers.UnimplementedUsersServer

	logger *zap.SugaredLogger
}

func (s *UsersService) GetJWTToken(ctx context.Context, req *pbusers.GetJWTTokenRequest) (*pbusers.GetJWTTokenResponse, error) {
	return nil, nil
}

func (s *UsersService) AddUser(ctx context.Context, req *pbusers.AddUserRequest) (*pbusers.AddUserResponse, error) {
	resp := &pbusers.AddUserResponse{
		Success: true,
		CreatedUser: &pbcommon.User{
			UserID:   123,
			UserName: "test",
			Email:    "test@test.com",
		},
		Errors: []pbusers.AddUserError{},
	}
	return resp, nil
}

func (s *UsersService) GetUserByUsername(context.Context, *pbusers.GetUserByUsernameRequest) (*pbusers.GetUserByUsernameResponse, error) {
	return nil, nil
}

func (s *UsersService) GetUserByID(context.Context, *pbusers.GetUserByIDRequest) (*pbusers.GetUserByIDResponse, error) {
	return nil, nil
}

func (s *UsersService) GetUserNameByID(context.Context, *pbusers.GetUserNameByIDRequest) (*pbusers.GetUserNameByIDResponse, error) {
	return nil, nil
}

func (s *UsersService) DeleteUserByID(context.Context, *pbusers.DeleteUserByIDRequest) (*pbusers.DeleteUserByIDResponse, error) {
	return nil, nil
}

func (s *UsersService) UpdateUserInfo(context.Context, *pbusers.UpdateUserInfoRequest) (*pbusers.UpdateUserInfoResponse, error) {
	return nil, nil
}

