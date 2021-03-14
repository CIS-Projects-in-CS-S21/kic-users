package database

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SQLRepository struct {
	db *gorm.DB

	logger *zap.SugaredLogger
}

func NewSQLRepository(db *gorm.DB, logger *zap.SugaredLogger) *SQLRepository {
	return &SQLRepository{
		db:     db,
		logger: logger,
	}
}

func (s *SQLRepository) checkIfUsernameAvailable(username string) bool {
	var user UserModel
	s.db.Where(&UserModel{Username: username}).First(&user)
	if user.Username == username {
		s.logger.Debugf("Username not available: %v", username)
		return false
	}
	return true
}

func (s *SQLRepository) checkIfEmailAvailable(email string) bool {
	var user UserModel
	s.db.Where(&UserModel{Email: email}).First(&user)
	if user.Email == email {
		s.logger.Debugf("Email not available: %v", email)
		return false
	}
	return true
}

func (s *SQLRepository) AddUser(ctx context.Context, user *UserModel) (int64, error) {
	ok := true

	if !s.checkIfUsernameAvailable(user.Username) {
		s.logger.Debug("Username not available")
		ok = false
	}

	s.logger.Debugf("Current ok: %v", ok)

	if !s.checkIfEmailAvailable(user.Email) {
		s.logger.Debug("Email not available")
		s.logger.Debugf("Result of s.checkIfEmailAvailable(%v): %v", user.Email, s.checkIfEmailAvailable(user.Email))
		ok = false
	}

	s.logger.Debugf("Current ok: %v", ok)

	if ok {
		s.db.Create(user)
		return int64(user.ID), nil
	}

	s.logger.Debugf("Did not insert record %v", user)

	return -1, nil
}

func (s *SQLRepository) GetUser(ctx context.Context, user *UserModel) (*UserModel, error) {
	toReturn := &UserModel{}
	transaction := s.db.Where(user).Find(&toReturn)

	return toReturn, transaction.Error
}

func (s *SQLRepository) GetUserByID(ctx context.Context, id int64) (*UserModel, error) {
	toReturn := &UserModel{}
	transaction := s.db.First(&toReturn, id)

	return toReturn, transaction.Error
}

func (s *SQLRepository) DeleteUserByID(ctx context.Context, userID int64) error {
	transaction := s.db.Delete(&UserModel{}, userID)
	return transaction.Error
}

func (s *SQLRepository) UpdateUserInfo(context.Context, *UserModel) error {
	return nil
}
