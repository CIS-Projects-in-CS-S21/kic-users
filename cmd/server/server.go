package main

import (
	"fmt"
	"github.com/kic/users/pkg/database"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
	"os"
	"os/signal"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"

	"github.com/kic/users/internal/server"
	"github.com/kic/users/pkg/logging"
	pbusers "github.com/kic/users/pkg/proto/users"
)

func main() {
	IsProduction := os.Getenv("PRODUCTION") != ""
	dbPass := os.Getenv("DB_PASS")
	var logger *zap.SugaredLogger
	var dbConnString string
	if IsProduction {
		logger = logging.CreateLogger(zapcore.InfoLevel)
		dbConnString = fmt.Sprintf("root:%v@tcp(mysql.kic.svc.cluster.local:3306)/kic_users_prod?charset=utf8mb4&parseTime=True&loc=Local", dbPass)
	} else {
		logger = logging.CreateLogger(zapcore.DebugLevel)
		dbConnString = fmt.Sprintf("root:%v@tcp(mysql.kic.svc.cluster.local:3306)/kic_users_test?charset=utf8mb4&parseTime=True&loc=Local", dbPass)
	}

	ListenAddress := ":" + os.Getenv("PORT")

	listener, err := net.Listen("tcp", ListenAddress)
	if err != nil {
		logger.Fatalf("Unable to listen on %v: %v", ListenAddress, err)
	}

	grpcServer := grpc.NewServer()

	db, err := gorm.Open(mysql.Open(dbConnString), &gorm.Config{})

	if err != nil {
		logger.Fatalf("Unable connect to db %v", err)
	}

	err = db.AutoMigrate(&database.UserModel{})

	if err != nil {
		logger.Fatalf("Unable migrate tables to db %v", err)
	}

	repo := database.NewSQLRepository(db, logger)

	serv := server.NewUsersService(repo, logger)

	pbusers.RegisterUsersServer(grpcServer, serv)
	authv3.RegisterAuthorizationServer(grpcServer, serv)

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
