package main

import (
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kic/users/internal/server"
	"github.com/kic/users/pkg/logging"
	pbusers "github.com/kic/users/pkg/proto/users"
)

func main() {
	IsProduction := os.Getenv("PRODUCTION") != ""
	var logger *zap.SugaredLogger
	if IsProduction {
		logger = logging.CreateLogger(zapcore.InfoLevel)
	} else {
		logger = logging.CreateLogger(zapcore.DebugLevel)
	}

	ListenAddress := ":" + os.Getenv("PORT")

	listener, err := net.Listen("tcp", ListenAddress)
	if err != nil {
		logger.Fatalf("Unable to listen on %v: %v", ListenAddress, err)
	}

	grpcServer := grpc.NewServer()

	pbusers.RegisterUsersServer(grpcServer, &server.UsersService{})


	go func() {
		defer listener.Close()
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatalf("Failed to serve: %v", err)
		}
	}()


	defer grpcServer.Stop()

	// the server is listening in a goroutine so hang until we get an interrupt signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
}