package server

import (
	"context"
	"fmt"
	pbcommon "github.com/kic/users/pkg/proto/common"
	proto "github.com/kic/users/pkg/proto/users"
	"google.golang.org/grpc/metadata"
	"testing"
)

func Test_ShouldAddUser(t *testing.T) {
	res, err := service.AddUser(context.Background(), &proto.AddUserRequest{
		Email:           "testemail@gmail.com",
		DesiredUsername: "tester",
		DesiredPassword: "tester",
		Birthday: &pbcommon.Date{
			Year:  1990,
			Month: 1,
			Day:   2,
		},
		City: "tester",
	})
	if err != nil {
		t.Errorf("Failed to add unique user to database: %v", err)
	}

	if res.Success != true {
		t.Error("Failed to add unique user to database")
	}
}

func Test_FailShouldAddUser(t *testing.T) {
	res, err := service.AddUser(context.Background(), &proto.AddUserRequest{
		Email:           "qdn@gmail.com",
		DesiredUsername: "qdn123",
		DesiredPassword: "tester",
		Birthday: &pbcommon.Date{
			Year:  1990,
			Month: 1,
			Day:   2,
		},
		City: "tester",
	})

	if err == nil {
		t.Error("Should have received error with duplicate credentials but did not")
	}

	if res.Success == true {
		t.Error("Got success despite duplicate username")
	}
}

func Test_ShouldGetUserByUsername(t *testing.T) {
	resp, err := service.GetUserByUsername(context.Background(), &proto.GetUserByUsernameRequest{Username: "qdn123"})
	if err != nil {
		t.Error("Got an error despite sending a proper user")
	}
	if resp.GetUser().UserName != "qdn123" || resp.GetUser().Email != "qdn@gmail.com" || resp.Success == false {
		t.Error("Did not get full user information with request")
	}
}

func Test_ShouldFailGetUserByUsername(t *testing.T) {
	resp, err := service.GetUserByUsername(context.Background(), &proto.GetUserByUsernameRequest{Username: "fakeperson"})
	if err == nil {
		t.Error("Did not get an error despite sending a improper user")
	}
	if resp.Success != false {
		t.Error("Did not get a false success value despite sending bad user data")
	}
}

func Test_ShouldGetUserByID(t *testing.T) {
	resp, err := service.GetUserByID(context.Background(), &proto.GetUserByIDRequest{UserID: 0})
	if err != nil {
		t.Error("Got an error despite requesting a proper ID")
	}

	if resp.GetUser().UserName != "qdn123" || resp.GetUser().Email != "qdn@gmail.com" || resp.Success == false {
		t.Error("Did not get full user information with ID request")
	}
}

func Test_ShouldFailGetUserByID(t *testing.T) {
	resp, err := service.GetUserByID(context.Background(), &proto.GetUserByIDRequest{UserID: 2000})
	if err == nil {
		t.Error("Got no error despite requesting an improper ID")
	}

	if resp.Success != false {
		t.Error("Did not get a false success value despite sending bad user data")
	}
}

func Test_ShouldGetJWT(t *testing.T) {
	token, err := service.GetJWTToken(context.Background(), &proto.GetJWTTokenRequest{
		Username: "qdn123",
		Password: "password",
	})
	if err != nil {
		t.Errorf("Failed to authenticate user with proper credentials")
	}

	if token.Token == "" {
		t.Errorf("Failed to get JWT with proper credentials")
	}
}

func Test_ShouldFailGetJWT(t *testing.T) {
	token, err := service.GetJWTToken(context.Background(), &proto.GetJWTTokenRequest{
		Username: "fakeuser",
		Password: "fakepassword",
	})
	if err == nil {
		t.Errorf("Did not get an error from JWT request with bad credentials")
	}

	if token != nil {
		t.Errorf("Got a token back with bad credentials")
	}
}

func Test_ShouldDeleteUserByID(t *testing.T) {
	token, _ := service.GetJWTToken(context.Background(), &proto.GetJWTTokenRequest{
		Username: "deleteme",
		Password: "password",
	})
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(authHeader, fmt.Sprintf("Bearer %v", token.Token)),
	)
	resp, err := service.DeleteUserByID(ctx, &proto.DeleteUserByIDRequest{
		UserID: 1,
	})
	if err != nil {
		t.Errorf("Failed to delete user with proper credentials and ID: %v", err)
	}
	if resp.Success != true {
		t.Error("Failed to delete user with proper credentials and ID")
	}
}

func Test_ShouldFailDeleteUserByID(t *testing.T) {
	resp, err := service.DeleteUserByID(context.Background(), &proto.DeleteUserByIDRequest{
		UserID: 1000,
	})
	if err == nil {
		t.Errorf("Did not get an error despite trying to delete invalid user")
	}
	if resp != nil {
		t.Error("Got valid response when deleting user despite bad data")
	}
}

func Test_ShouldUpdateUserInfo(t *testing.T) {
	res, _ := service.AddUser(context.Background(), &proto.AddUserRequest{
		Email:           "updateme@gmail.com",
		DesiredUsername: "updateme",
		DesiredPassword: "updateme",
		Birthday: &pbcommon.Date{
			Year:  1990,
			Month: 1,
			Day:   2,
		},
		City: "tester",
	})

	resp, err := service.UpdateUserInfo(context.Background(), &proto.UpdateUserInfoRequest{
		UserID:          res.CreatedUser.UserID,
		Email:           "iamupdated@gmail.com",
		DesiredUsername: "iamupdated",
		City:            "updated",
		Bio:             "BIObioBIO",
	})
	if err != nil {
		t.Errorf("Got an error updating a user with valid items")
	}
	usr := resp.GetUpdatedUser()
	if resp.Success != true || usr.City != "updated" || usr.Bio != "BIObioBIO" || usr.Email != "iamupdated@gmail.com" {
		t.Errorf("Updated user does not have expected information")
	}
}

func Test_ShouldFailUpdateUserInfo(t *testing.T) {
	resp, err := service.UpdateUserInfo(context.Background(), &proto.UpdateUserInfoRequest{
		UserID:          1000,
		Email:           "iamupdated@gmail.com",
		DesiredUsername: "iamupdated",
		City:            "updated",
		Bio:             "BIObioBIO",
	})
	if err == nil {
		t.Errorf("Did not get an error updating a user with invalid items")
	}
	if resp.Success != false {
		t.Errorf("Did not get failure response despite trying to update invalid user")
	}
}