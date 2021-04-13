package database

import (
	"go.uber.org/zap"
	common "github.com/kic/users/pkg/proto/common"
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

func (s *MockRepository) GetUser (user *UserModel) (*UserModel, error) {



	return nil, nil
}

func (s *MockRepository) GetUserByID(id int64) (*UserModel, error) {


	return nil, nil
}

func (s *MockRepository) DeleteUserByID(userID int64) error {
	return nil
}

