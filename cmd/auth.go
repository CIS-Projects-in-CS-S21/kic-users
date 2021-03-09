package main

import (
	"context"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/zap"
)

type AuthService struct {
	logger *zap.SugaredLogger
}

type AuthServer struct {

}

type User struct {
	username string
	password string
}

type DataBase interface {
	CheckUser(user *User) (bool, error)
}

func (s* AuthServer) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	c := authv3.CheckResponse{}
	return &c, nil
}

func(user *User) CheckUser() (bool, error) {
	return true,nil
}


