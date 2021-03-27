package database

import (
	pbcommon "github.com/kic/users/pkg/proto/common"
	"gorm.io/gorm"
	"time"
)

type UserModel struct {
	gorm.Model
	Email           string
	Username 		string
	Password 		string `gorm:"size:255"`
	Birthday        time.Time
	City            string
	Bio				string
}

func NewUserModel(username, email, password, city, bio string, birthday *pbcommon.Date) *UserModel {
	var bday time.Time
	if birthday != nil {
		bday = time.Date(int(birthday.Year), time.Month(birthday.Month), int(birthday.Day), 0,0,0,0, time.Local)
	}
	return &UserModel{
		Email:    email,
		Username: username,
		Password: password,
		Birthday: bday,
		City:     city,
		Bio: 	  bio,
	}
}