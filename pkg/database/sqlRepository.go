package database

import (
	"context"
	"errors"
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

	return -1, errors.New("username or email taken")
}

func (s *SQLRepository) GetUser(ctx context.Context, user *UserModel) (*UserModel, error) {
	toReturn := &UserModel{}
	transaction := s.db.Where("username = ?", user.Username).First(&toReturn)

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

func (s *SQLRepository) UpdateUserInfo(ctx context.Context, user *UserModel) error {
	ok := true
	var tx *gorm.DB // declaring response variable DB, which will be returned form s.db.Update()

	if user.Email != "" { // update Email if it's been changed
		if !s.checkIfEmailAvailable(user.Email) {
			s.logger.Debug("Email not available")
			s.logger.Debugf("Result of s.checkIfEmailAvailable(%v): %v", user.Email, s.checkIfEmailAvailable(user.Email))
			ok = false
		}
		s.logger.Debugf("Current ok (in email case): %v", ok)

		if ok {
			tx = s.db.Model(&UserModel{}).Where("id = ?", user.ID).Update("Email", user.Email)
			if tx.Error != nil { // return error if there is one
				return tx.Error
			}
		}

	}

	if user.Username != "" { // update Username if it's been changed
		if !s.checkIfUsernameAvailable(user.Username) {
			s.logger.Debug("Username not available")
			ok = false
		}
		s.logger.Debugf("Current ok (in username case): %v", ok)

		if ok {
			tx = s.db.Model(&UserModel{}).Where("id = ?", user.ID).Update("Username", user.Username)
			if tx.Error != nil { // return error if there is one
				return tx.Error
			}
		}

	}

	if user.Password != "" { // update Password if it's been changed
		tx = s.db.Model(&UserModel{}).Where("id = ?", user.ID).Update("Password", user.Password)
		if tx.Error != nil { // return error if there is one
			return tx.Error
		}
	}

	if user.Password != "" { // update Password if it's been changed
		tx = s.db.Model(&UserModel{}).Where("id = ?", user.ID).Update("Password", user.Password)
		if tx.Error != nil { // return error if there is one
			return tx.Error
		}
	}

	if !user.Birthday.IsZero() { // update Birthday if it's been changed
		tx = s.db.Model(&UserModel{}).Where("id = ?", user.ID).Update("Birthday", user.Birthday)
		if tx.Error != nil { // return error if there is one
			return tx.Error
		}
	}

	if user.City != "" { // update Password if it's been changed
		tx = s.db.Model(&UserModel{}).Where("id = ?", user.ID).Update("City", user.City)
		if tx.Error != nil { // return error if there is one
			return tx.Error
		}
	}

	if user.Bio != "" {
		tx = s.db.Model(&UserModel{}).Where("id = ?", user.ID).Update("Bio", user.Bio)
		if tx.Error != nil { // return error if there is one
			return tx.Error
		}
	}

	if user.Triggers != "" {
		tx = s.db.Model(&UserModel{}).Where("id = ?", user.ID).Update("Triggers", user.Triggers)
		if tx.Error != nil { // return error if there is one
			return tx.Error
		}
	}
	if user.Private != "" {
		tx = s.db.Model(&UserModel{}).Where("id = ?", user.ID).Update("Private", user.Private)
		if tx.Error != nil { // return error if there is one
			return tx.Error
		}
	}

	return nil
}
