package database

import (
	"context"
)

// Repository - interface for a data provider that interfaces between the database backend and the grpc server
// enables the repository pattern so that we can swap out the database backend easily
type Repository interface {
	AddUser(context.Context, *UserModel) (int64, error)
	// Provide any info you can to get a user
	GetUser(context.Context, *UserModel) (*UserModel, error)
	GetUserByID(context.Context, int64) (*UserModel, error)
	DeleteUserByID(context.Context, int64) error
	UpdateUserInfo(context.Context, *UserModel) error
}
