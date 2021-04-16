package server

import (
	"context"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	pbusers "github.com/kic/users/pkg/proto/users"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/kic/users/pkg/database"
	"github.com/kic/users/pkg/logging"
)

var service *UsersService

func TestMain(m *testing.M) {
	logger := logging.CreateLogger(zapcore.DebugLevel)
	os.Setenv("SECRET_KEY", "supersecret")

	pass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	seedData := map[uint]*database.UserModel{
		0: {
			Model: gorm.Model{
				ID: 0,
			},
			Email:    "qdn@gmail.com",
			Username: "qdn123",
			Password: string(pass),
			Birthday: time.Time{},
			City:     "",
			Bio:      "",
		},
		1: {
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "deleteme@gmail.com",
			Username: "deleteme",
			Password: string(pass),
			Birthday: time.Time{},
			City:     "",
			Bio:      "",
		},
	}

	repo := database.NewMockRepository(seedData, logger)

	service = NewUsersService(repo, logger)

	os.Exit(m.Run())
}

func Test_ShouldPassAuthentication(t *testing.T) {
	_, err := service.GetJWTToken(context.Background(), &pbusers.GetJWTTokenRequest{
		Username: "qdn123",
		Password: "password",
	})
	if err != nil {
		t.Errorf("Failed to authenticate user with proper credentials")
	}
}

func Test_ShouldPassAuthenticationWithJWT(t *testing.T) {
	token, err := service.GetJWTToken(context.Background(), &pbusers.GetJWTTokenRequest{
		Username: "qdn123",
		Password: "password",
	})
	if err != nil {
		t.Errorf("Failed to authenticate user with proper credentials")
	}
	check, err := service.Check(context.Background(), &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: map[string]string{
						authHeader: "Bearer " + token.GetToken(),
					},
				},
			},
		},
	})
	if err != nil {
		t.Errorf("Failed to authenticate user with proper credentials")
	}
	if check.GetOkResponse() == nil {
		t.Errorf("Did not get OkResponse with proper credentials")
	}
}

func Test_ShouldFailAuthentication(t *testing.T) {
	check, err := service.Check(context.Background(), &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: map[string]string{
						authHeader: "Bearer " + "sdfasdfsdfasdf",
					},
				},
			},
		},
	})
	if err != nil {
		t.Errorf("Failed to authenticate user with proper credentials")
	}
	if check.GetOkResponse() != nil {
		t.Errorf("Got OkResponse with improper credentials")
	}
}
