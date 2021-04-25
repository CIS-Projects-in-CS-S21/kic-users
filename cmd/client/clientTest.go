package main

import (
	"context"
	"fmt"
	pbcommon "github.com/kic/users/pkg/proto/common"
	pbusers "github.com/kic/users/pkg/proto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func shouldAddUser(client pbusers.UsersClient) int64 {
	in := &pbusers.AddUserRequest{
		Email:           "qdn@gmail.com",
		DesiredUsername: "qdnovinger",
		DesiredPassword: "password134234",
		Birthday: &pbcommon.Date{
			Year:  1998,
			Month: 8,
			Day:   21,
		},
		City: "Scranton",
	}

	addRes, err := client.AddUser(context.Background(), in)

	if err != nil {
		log.Fatalf("fail to add user: %v", err)
	}

	_, err = client.AddUser(context.Background(), in)

	if err == nil {
		log.Fatalf("should fail to add duplicate user")
	}

	log.Printf("shouldAddUser Success")

	return addRes.CreatedUser.UserID
}

func shouldGetToken(client pbusers.UsersClient) string {
	tokRes, err := client.GetJWTToken(context.Background(), &pbusers.GetJWTTokenRequest{
		Username: "qdnovinger",
		Password: "password134234",
	})

	if err != nil {
		log.Fatalf("fail to get token: %v", err)
	}

	log.Printf("shouldGetToken Success")

	return tokRes.Token
}

func shouldGetUser(ctx context.Context, client pbusers.UsersClient) {
	usernameRes, err := client.GetUserByUsername(ctx, &pbusers.GetUserByUsernameRequest{Username: "qdnovinger"})

	if err != nil {
		log.Fatalf("fail to get user by name: %v", err)
	}

	if usernameRes.User.Email != "qdn@gmail.com" || usernameRes.User.UserName != "qdnovinger" {
		log.Fatalf("got incorrect user with username")
	}

	log.Printf("shouldGetUser Success")
}

func shouldGetUserByID(ctx context.Context, client pbusers.UsersClient, uid int64) {
	unameByIDRes, err := client.GetUserNameByID(ctx, &pbusers.GetUserNameByIDRequest{UserID: uid})

	if err != nil {
		log.Fatalf("fail to get username by id: %v", err)
	}

	if unameByIDRes.Username != "qdnovinger" {
		log.Fatalf("Incorrect response from GetUserNameByID: %v", unameByIDRes.Username)
	}

	log.Printf("shouldGetUserByID Success")
}

func shouldUpdateUser(ctx context.Context, client pbusers.UsersClient, uid int64) {
	updateReq := &pbusers.UpdateUserInfoRequest{
		UserID:          uid,
		Email:           "",
		DesiredUsername: "hot_mama_RAWR_XD",
		DesiredPassword: "",
		Birthday:        nil,
		City:            "Philadelphia",
		Bio:             "Hey guys I am Ryan",
		Triggers:        "stuff",
		IsPrivate:       "1",
	}

	_, err := client.UpdateUserInfo(ctx, updateReq)

	if err != nil {
		log.Fatalf("fail to update user: %v", err)
	}

	usernameRes, err := client.GetUserByUsername(ctx, &pbusers.GetUserByUsernameRequest{Username: "hot_mama_RAWR_XD"})

	if err != nil || usernameRes.User.UserName != "hot_mama_RAWR_XD" ||
		usernameRes.User.UserID != uid ||
		usernameRes.User.IsPrivate != "1" ||
		usernameRes.User.Triggers != "stuff" {
		log.Fatalf("fail to update user: %v", usernameRes.User)
	}

	log.Printf("shouldUpdateUser Success")
}

func shouldDeleteUser(ctx context.Context, client pbusers.UsersClient, uid int64) {
	deleteRes, err := client.DeleteUserByID(ctx, &pbusers.DeleteUserByIDRequest{UserID: uid})

	if err != nil {
		log.Fatalf("fail to delete user: %v", err)
	}

	if deleteRes.Success != true {
		log.Fatalf("fail to delete user")
	}

	log.Printf("shouldDeleteUser Success")
}

func main() {
	conn, err := grpc.Dial("test.api.keeping-it-casual.com:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pbusers.NewUsersClient(conn)

	uid := shouldAddUser(client)
	token := shouldGetToken(client)

	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", token))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	shouldGetUser(ctx, client)
	shouldGetUserByID(ctx, client, uid)
	shouldUpdateUser(ctx, client, uid)
	shouldDeleteUser(ctx, client, uid)

}
