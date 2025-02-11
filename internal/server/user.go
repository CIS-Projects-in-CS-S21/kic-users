package server

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/kic/users/pkg/database"
	pbcommon "github.com/kic/users/pkg/proto/common"
	pbusers "github.com/kic/users/pkg/proto/users"
)

type UsersService struct {
	pbusers.UnimplementedUsersServer

	db     database.Repository
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
		db:     db,
		keyset: keyset,
		logger: logger,
	}
}

func (s *UsersService) GetJWTToken(ctx context.Context, req *pbusers.GetJWTTokenRequest) (*pbusers.GetJWTTokenResponse, error) {
	s.logger.Debug("Getting JWT token")

	valid, err := s.ValidateUser(req.Username, req.Password)

	if err != nil {
		s.logger.Debugf("User %v is invalid: %v", req.Username, err)
		return nil, err
	}

	if !valid {
		s.logger.Debugf("User %v is invalid", req.Username)
		return nil, status.Errorf(codes.InvalidArgument, "Password incorrect")
	}

	s.logger.Debugf("User %v is valid", req.Username)

	userData, err := s.db.GetUser(context.TODO(), &database.UserModel{Username: req.Username})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not access user")
	}

	token, err := s.GenerateJWT(int64(userData.ID))

	s.logger.Debugf("Generated token: %v", token)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not generate token")
	}

	resp := &pbusers.GetJWTTokenResponse{
		Token: token,
	}

	return resp, nil
}

func (s *UsersService) AddUser(ctx context.Context, req *pbusers.AddUserRequest) (*pbusers.AddUserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.DesiredPassword), bcrypt.DefaultCost)

	if err != nil {
		s.logger.Errorf("Failed to hash password: %v", err)
		return nil, err
	}

	model := database.NewUserModel(
		req.DesiredUsername,
		req.Email,
		string(hashedPassword),
		req.City,
		"",
		req.Birthday,
		req.Triggers,
		req.IsPrivate,
	)

	id, err := s.db.AddUser(context.TODO(), model)

	if id != -1 {
		return &pbusers.AddUserResponse{
			Success: true,
			CreatedUser: &pbcommon.User{
				UserID:   id,
				UserName: model.Username,
				Email:    model.Email,
				Birthday: &pbcommon.Date{
					Year:  int32(model.Birthday.Year()),
					Month: int32(model.Birthday.Month()),
					Day:   int32(model.Birthday.Day()),
				},
				City: model.City,
				Bio:  model.Bio,
				Triggers: model.Triggers,
				IsPrivate: model.Private,
			},
		}, err
	}

	resp := &pbusers.AddUserResponse{
		Success:     false,
		CreatedUser: nil,
	}
	return resp, status.Errorf(codes.AlreadyExists, "User already exists")
}

func (s *UsersService) GetUserByUsername(ctx context.Context, req *pbusers.GetUserByUsernameRequest) (*pbusers.GetUserByUsernameResponse, error) {
	model := &database.UserModel{
		Email:    "",
		Username: req.Username,
		Password: "",
		Birthday: time.Time{},
		City:     "",
	}

	user, err := s.db.GetUser(ctx, model)

	if err != nil || user.Username == "" {
		return &pbusers.GetUserByUsernameResponse{
			Success: false,
			User:    nil,
		}, err
	}

	resp := &pbusers.GetUserByUsernameResponse{
		Success: true,
		User: &pbcommon.User{
			UserID:   int64(user.ID),
			UserName: user.Username,
			Email:    user.Email,
			Birthday: &pbcommon.Date{
				Year:  int32(user.Birthday.Year()),
				Month: int32(user.Birthday.Month()),
				Day:   int32(user.Birthday.Day()),
			},
			City: user.City,
			Bio:  user.Bio,
			Triggers: user.Triggers,
			IsPrivate: user.Private,
		},
	}
	return resp, err
}

func (s *UsersService) GetUserByID(ctx context.Context, req *pbusers.GetUserByIDRequest) (*pbusers.GetUserByIDResponse, error) {
	usr, err := s.db.GetUserByID(context.TODO(), req.GetUserID())

	if err != nil {
		return &pbusers.GetUserByIDResponse{
			Success: false,
			User:    nil,
		}, err
	}

	resp := &pbusers.GetUserByIDResponse{
		Success: true,
		User: &pbcommon.User{
			UserID:   int64(usr.ID),
			UserName: usr.Username,
			Email:    usr.Email,
			Birthday: &pbcommon.Date{
				Year:  int32(usr.Birthday.Year()),
				Month: int32(usr.Birthday.Month()),
				Day:   int32(usr.Birthday.Day()),
			},
			City: usr.City,
			Bio:  usr.Bio,
			Triggers: usr.Triggers,
			IsPrivate: usr.Private,
		},
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
	headers, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		s.logger.Debugf("Failed to get headers from incoming call in DeleteUserByID")
		return nil, status.Errorf(codes.Unauthenticated, "Send token along with request")
	}

	header := headers[authHeader][0]

	s.logger.Debugf("Incoming header: %v", header)

	tokString, err := parseCredentialsFromHeader(header)

	tok, err := s.DecodeJWT(tokString)

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

func (s *UsersService) UpdateUserInfo(ctx context.Context, req *pbusers.UpdateUserInfoRequest) (*pbusers.UpdateUserInfoResponse, error) {
	// creating UpdateUserInfo Response that indicates failure
	// will be returned if any operation fails
	failureResponse := &pbusers.UpdateUserInfoResponse{
		Success:     false,
		UpdatedUser: nil,
	}

	s.logger.Debugf("Starting UpdateUserInfo with req: %v", req)

	var hashedPassword []byte // declaring hashedPassword to potentially be filled in
	var err error             // declaring err variable to hold potential errors

	if req.DesiredPassword != "" { // if password change is requested
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(req.DesiredPassword), bcrypt.DefaultCost) // hash the password
	}

	// if error, log and return and failure
	if err != nil {
		s.logger.Errorf("Failed to hash password: %v", err)
		return failureResponse, status.Errorf(codes.InvalidArgument, "Password cannot be encrypted")
	}

	// create UserModel from updated fields
	model := database.NewUserModel(
		req.DesiredUsername,
		req.Email,
		string(hashedPassword),
		req.City,
		req.Bio,
		req.Birthday,
		req.Triggers,
		req.IsPrivate,
	)

	model.ID = uint(req.UserID)

	s.logger.Debugf("Created new user model: %v", model)

	// attempt to update db with model containing updated information
	err = s.db.UpdateUserInfo(context.TODO(), model)

	// if error, log and return failure
	if err != nil {
		s.logger.Errorf("Failed to Update User Info in database: %v", err)
		return failureResponse, err
	}


	usr, _ := s.db.GetUserByID(context.TODO(), req.GetUserID())

	// creating success response
	resp := &pbusers.UpdateUserInfoResponse{Success: true, UpdatedUser: &pbcommon.User{
		UserID:   int64(usr.ID),
		UserName: usr.Username,
		Email:    usr.Email,
		Birthday: &pbcommon.Date{
			Year:  int32(usr.Birthday.Year()),
			Month: int32(usr.Birthday.Month()),
			Day:   int32(usr.Birthday.Day()),
		},
		City: usr.City,
		Bio:  usr.Bio,
		Triggers: usr.Triggers,
		IsPrivate: usr.Private,
	}}

	s.logger.Debugf("Finished updating info in db, returning: %v", resp)

	// returning success response and nil error
	return resp, nil

}
