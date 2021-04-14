package database

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

type MockRepository struct {
	db map[uint]*UserModel

	logger *zap.SugaredLogger
	idCounter uint
}


func NewMockRepository(db map[uint]*UserModel, logger *zap.SugaredLogger) *MockRepository {
	return &MockRepository{
		db:        db,
		logger:    logger,
		idCounter: uint(len(db)),
	}
}

func (m* MockRepository) AddUser(ctx context.Context, user *UserModel) (int64, error) {
	for _, val := range m.db {
		if val.Username == user.Username || val.Email == user.Email {
			return -1, errors.New("username or email taken")
		}
	}
	user.ID = m.idCounter
	m.db[m.idCounter] = user
	m.idCounter++
	return int64(user.ID), nil
}

func (m* MockRepository) GetUser(ctx context.Context, user *UserModel) (*UserModel, error) {
	for _, val := range m.db {
		if val.Username == user.Username || val.Email == user.Email {
			return val, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m* MockRepository) GetUserByID(ctx context.Context, id int64) (*UserModel, error) {
	if val, ok := m.db[uint(id)]; ok {
		return val, nil
	}
	return nil, errors.New("user not found")
}

func (m* MockRepository) DeleteUserByID(ctx context.Context, id int64) error {
	if _, ok := m.db[uint(id)]; ok {
		delete(m.db, uint(id))
		return nil
	}
	return errors.New("user not found")
}

func (m* MockRepository) UpdateUserInfo(ctx context.Context, user *UserModel) error {
	if _, ok := m.db[user.ID]; !ok {
		return errors.New("update user not found")
	}
	m.db[user.ID] = user
	return nil
}


