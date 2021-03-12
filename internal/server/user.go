package server

import (
	"context"
	pbcommon "github.com/kic/users/pkg/proto/common"
	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/kic/users/pkg/database"
	pbusers "github.com/kic/users/pkg/proto/users"
)

type UsersService struct {
	pbusers.UnimplementedUsersServer

	db database.Repository
	keyset jwk.Set

	logger *zap.SugaredLogger
}

func NewUsersService(db database.Repository, logger *zap.SugaredLogger) *UsersService {
	secretKey := os.Getenv("SECRET_KEY")
	raw := []byte(secretKey)

	jkey, _ := jwk.New(raw)

	keyset := jwk.NewSet()
	keyset.Add(jkey)

	return &UsersService{
		db:                       db,
		keyset: 				  keyset,
		logger:                   logger,
	}
}

func (s *UsersService) GetJWTToken(ctx context.Context, req *pbusers.GetJWTTokenRequest) (*pbusers.GetJWTTokenResponse, error) {
	valid, err := s.ValidateUser(req.Username, req.Password)

	if err != nil {
		return nil, err
	}

	if !valid {
		return &pbusers.GetJWTTokenResponse{
			Token: "",
			Error: pbusers.GetJWTTokenResponse_INVALID_PASSWORD,
		}, nil
	}

	userData, err := s.db.GetUser(context.TODO(), &database.UserModel{Username: req.Username})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not access user")
	}

	token, err := s.GenerateJWT(int64(userData.ID))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not generate token")
	}

	resp := &pbusers.GetJWTTokenResponse{
		Token: token,
		Error: -1,
	}

	return resp, nil
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

func (s *UsersService) GetUserByID(ctx context.Context, req *pbusers.GetUserByIDRequest) (*pbusers.GetUserByIDResponse, error) {
	usr, err := s.db.GetUserByID(context.TODO(), req.GetUserID())

	if err != nil {
		return &pbusers.GetUserByIDResponse{
			Success: false,
			User:    nil,
			Errors:  nil,
		}, err
	}

	resp := &pbusers.GetUserByIDResponse{
		Success: true,
		User:    &pbcommon.User{
			UserID:   int64(usr.ID),
			UserName: usr.Username,
			Email:    usr.Email,
		},
		Errors:  nil,
	}

	return resp, nil
}

func (s *UsersService) GetUserNameByID(ctx context.Context, req *pbusers.GetUserNameByIDRequest) (*pbusers.GetUserNameByIDResponse, error) {
	usr, err := s.db.GetUserByID(context.TODO(), req.GetUserID())

	if err != nil {
		return &pbusers.GetUserNameByIDResponse{
			Username: "",
		}, err
	}

	resp := &pbusers.GetUserNameByIDResponse{
		Username: usr.Username,
	}

	return resp, nil
}

func (s *UsersService) DeleteUserByID(ctx context.Context, req *pbusers.DeleteUserByIDRequest) (*pbusers.DeleteUserByIDResponse, error) {
	headers, ok :=  metadata.FromIncomingContext(ctx)

	if !ok {
		s.logger.Debugf("Failed to get headers from incoming call in DeleteUserByID")
		return nil, status.Errorf(codes.Unauthenticated, "Send token along with request")
	}

	token := headers[authHeader][0]

	tok, err := s.DecodeJWT(token)

	if err != nil {
		s.logger.Debugf("Failed to decode token")
		return nil, status.Errorf(codes.Internal, "Failed to decode token")
	}

	strID, _ := tok.Get("uid")
	tokID, err := strconv.Atoi(strID.(string))

	if int64(tokID) != req.UserID || err != nil {
		return &pbusers.DeleteUserByIDResponse{
			Success: false,
		}, status.Errorf(codes.Unauthenticated, "Cannot delete another user's account")
	}

	err = s.db.DeleteUserByID(context.TODO(), req.UserID)

	if err != nil {
		s.logger.Debugf("Failed to get headers from incoming call in DeleteUserByID")
	}

	return &pbusers.DeleteUserByIDResponse{
		Success: true,
	}, nil
}

func (s *UsersService) UpdateUserInfo(context.Context, *pbusers.UpdateUserInfoRequest) (*pbusers.UpdateUserInfoResponse, error) {
	return nil, nil
}

