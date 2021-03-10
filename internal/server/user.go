package server

import (
	"context"
	pbcommon "github.com/kic/users/pkg/proto/common"
	"golang.org/x/crypto/bcrypt"

	"go.uber.org/zap"

	"github.com/kic/users/pkg/database"
	pbusers "github.com/kic/users/pkg/proto/users"
)

type UsersService struct {
	pbusers.UnimplementedUsersServer

	db database.Repository

	logger *zap.SugaredLogger
}

func NewUsersService(db database.Repository, logger *zap.SugaredLogger) *UsersService {
	return &UsersService{
		db:                       db,
		logger:                   logger,
	}
}

func (s *UsersService) GetJWTToken(ctx context.Context, req *pbusers.GetJWTTokenRequest) (*pbusers.GetJWTTokenResponse, error) {
	return nil, nil
}

func (s *UsersService) AddUser(ctx context.Context, req *pbusers.AddUserRequest) (*pbusers.AddUserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.DesiredPassword), bcrypt.DefaultCost)

	if err != nil {
		s.logger.Errorf("Failed to hash password: %v", err)
		return nil, err
	}

	model := &database.UserModel{
		Email:    req.Email,
		Username: req.DesiredUsername,
		Password: string(hashedPassword),
		Birthday: req.Birthday,
		City:     req.City,
	}

	id, insertErrors := s.db.AddUser(context.TODO(), model)

	if id == -1 {
		return &pbusers.AddUserResponse{
			Success:     true,
			CreatedUser: &pbcommon.User{
				UserID:   id,
				UserName: model.Username,
				Email:    model.Email,
			},
			Errors:      nil,
		}, nil
	}

	resp := &pbusers.AddUserResponse{
		Success: false,
		CreatedUser: nil,
		Errors: insertErrors,
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

