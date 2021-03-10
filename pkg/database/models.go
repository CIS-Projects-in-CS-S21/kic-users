package database

import (
	pbcommon "github.com/kic/users/pkg/proto/common"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Email           string
	Username 		string
	Password 		string
	Birthday        *pbcommon.Date
	City            string
}

func NewUserModel(username, email, password, city string, birthday *pbcommon.Date) *UserModel {
	return &UserModel{
		Email:    email,
		Username: username,
		Password: password,
		Birthday: birthday,
		City:     city,
	}
}