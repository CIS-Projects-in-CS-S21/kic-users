package database

import (
	"context"
	pbusers "github.com/kic/users/pkg/proto/users"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SQLRepository struct {
	db gorm.DB

	logger *zap.SugaredLogger
}

func NewSQLRepository(db gorm.DB, logger *zap.SugaredLogger) *SQLRepository {
	return &SQLRepository{
		db:     db,
		logger: logger,
	}
}

func (s *SQLRepository) checkIfUsernameAvailable(username string) bool {
	var user UserModel
	s.db.Where(&UserModel{Username: username}).First(&user)
	if user.Username == username {
		return false
	}
	return true
}

func (s *SQLRepository) checkIfEmailAvailable(email string) bool {
	var user UserModel
	s.db.Where(&UserModel{Email: email}).First(&user)
	if user.Email == email {
		return false
	}
	return false
}

func (s *SQLRepository) AddUser(ctx context.Context, user *UserModel) (int64, []pbusers.AddUserError) {
	var errors []pbusers.AddUserError
	ok := true

	if !s.checkIfUsernameAvailable(user.Username) {
		errors = append(errors, pbusers.AddUserError_DUPLICATE_USERNAME)
		ok = false
	}

	if !s.checkIfEmailAvailable(user.Email) {
		errors = append(errors, pbusers.AddUserError_DUPLICATE_EMAIL)
		ok = false
	}

	if ok {
		s.db.Create(user)
		return int64(user.ID), nil
	}

	return -1, errors
}

func (s *SQLRepository) GetUser(context.Context, *UserModel) (*UserModel, error) {
	return nil, nil
}

func (s *SQLRepository) DeleteUserByID(context.Context, int64) error {
	return nil
}

func (s *SQLRepository) UpdateUserInfo(context.Context, *UserModel) error {
	return nil
}
