package database

import (
	"context"
	common "github.com/kic/users/pkg/proto/common"
	"go.uber.org/zap"
)

type MockRepository struct {
	db *map[int]*common.User

	logger *zap.SugaredLogger

	idCounter int
}


func NewMockRepository(db *map[int]*common.User, logger *zap.SugaredLogger) *MockRepository {
	return &MockRepository{
		db:     db,
		logger: logger,
	}
}

func searchDBByUsername(db *map[int]*common.User, username string) (int, *common.User) {
	for key, value := range *db {
		if value.UserName == username {
			return key, value
		}
	}
	return -1, nil
}

func (s *MockRepository) checkIfUsernameAvailable(username string) bool {
	_, user := searchDBByUsername(s.db, username)

	if user == nil {
		return true
	}

	if user.UserName == username {
		return false
	}
	return true
}

func searchDBByEmail(db *map[int]*common.User, email string) (int, *common.User) {
	for key, value := range *db {
		if value.Email == email {
			return key, value
		}
	}
	return -1, nil
}

func (s *MockRepository) checkIfEmailAvailable(email string) bool {
	_, user := searchDBByEmail(s.db, email)

	if user == nil {
		return true
	}

	if user.Email == email {
		return false
	}
	return false
}

func (s *MockRepository) AddUser(user *common.User) (int, error) {
	ok := true

	if !s.checkIfUsernameAvailable(user.UserName) {
		ok = false
	}

	if !s.checkIfEmailAvailable(user.Email) {
		ok = false
	}

	if ok {
		database := *s.db
		database[s.idCounter] = user
		s.idCounter++
		return s.idCounter - 1, nil
	}

	return -1, nil
}

func (s *MockRepository) GetUser (user *common.User) (*common.User, error) {
	userNameQuery := user.UserName

	_, foundUser := searchDBByUsername(s.db, userNameQuery)

	return foundUser, nil
}

func (s *MockRepository) GetUserByID(id int64) (*common.User, error) {

	foundUser := (*s.db)[int(id)]

	return foundUser, nil
}

func (s *MockRepository) DeleteUserByID(userID int64) error {
	delete(*s.db, int(userID))
	return nil
}

func (s *MockRepository) UpdateUserInfo(ctx context.Context, user *common.User) error {
	ok := true

	if user.Email != "" { // update Email if it's been changed
		if !s.checkIfEmailAvailable(user.Email) {
			s.logger.Debug("Email not available")
			s.logger.Debugf("Result of s.checkIfEmailAvailable(%v): %v", user.Email, s.checkIfEmailAvailable(user.Email))
			ok = false
		}
		s.logger.Debugf("Current ok (in email case): %v", ok)

		if ok {
			(*s.db)[int(user.UserID)].Email = user.Email
		}

	}

	if user.UserName != "" { // update Username if it's been changed
		if !s.checkIfUsernameAvailable(user.UserName) {
			s.logger.Debug("Username not available")
			ok = false
		}
		s.logger.Debugf("Current ok (in username case): %v", ok)

		if ok {
			(*s.db)[int(user.UserID)].UserName = user.UserName
		}

	}

	if (user.Birthday.Day != 0 && user.Birthday.Month != 0 && user.Birthday.Year != 0) { // update Birthday if it's been changed
		(*s.db)[int(user.UserID)].Birthday = user.Birthday
	}

	if user.City != "" { // update Password if it's been changed
		(*s.db)[int(user.UserID)].City = user.City
	}

	if user.Bio != "" { // update Bio if it's been changed
		(*s.db)[int(user.UserID)].Bio = user.Bio
	}

	return nil
}

