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

type DataBase interface {
}

func (s* AuthServer) Check(ctx context.Context, req *CheckRequest) (*CheckResponse, error) {

}


