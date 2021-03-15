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

func main() {
	conn, err := grpc.Dial("test.api.keeping-it-casual.com:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pbusers.NewUsersClient(conn)

	in := &pbusers.AddUserRequest{
		Email:           "qdn@gmail.com",
		DesiredUsername: "qdnovinger",
		DesiredPassword: "password134234",
		Birthday:        &pbcommon.Date{
			Year:  1998,
			Month: 8,
			Day:   21,
		},
		City:            "Scranton",
	}

	addRes, err := client.AddUser(context.Background(), in)

	fmt.Printf("res: %v\nerr: %v\n", addRes, err)

	if err != nil {
		log.Fatalf("fail to add user: %v", err)
	}

	addRes2, err := client.AddUser(context.Background(), in)

	fmt.Printf("res: %v\nerr: %v\n", addRes2, err)

	tokRes, err := client.GetJWTToken(context.Background(), &pbusers.GetJWTTokenRequest{
		Username: "qdnovinger",
		Password: "password134234",
	})

	if err != nil {
		log.Fatalf("fail to get token: %v", err)
	}

	fmt.Printf("tokRes: %v\n", tokRes)

	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", tokRes.Token))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	usernameRes, err := client.GetUserByUsername(ctx, &pbusers.GetUserByUsernameRequest{Username: "qdnovinger"})


	fmt.Printf("res: %v\nerr: %v\n", usernameRes, err)

	md = metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", tokRes.Token))
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	unameByIDRes, err := client.GetUserNameByID(ctx, &pbusers.GetUserNameByIDRequest{UserID: addRes.CreatedUser.UserID})

	if err != nil {
		log.Fatalf("fail to get username by id: %v", err)
	}

	if unameByIDRes.Username != "qdnovinger" {
		log.Fatalf("Incorrect response from GetUserNameByID: %v", unameByIDRes.Username)
	}

	md = metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", tokRes.Token))
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	userByIDRes, err := client.GetUserByID(ctx,&pbusers.GetUserByIDRequest{UserID: addRes.CreatedUser.UserID} )

	if err != nil {
		log.Fatalf("fail to get user by id: %v", err)
	}

	if userByIDRes.Success != true || userByIDRes.User.UserName != "qdnovinger" || userByIDRes.User.UserID != addRes.CreatedUser.UserID {
		log.Fatalf("Incorrect response from GetUserByID: %v", unameByIDRes.Username)
	}


	// testing UpdateUser() -------------------------------------

	md = metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", tokRes.Token))
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	updateReq := &pbusers.UpdateUserInfoRequest{
		UserID:          addRes.CreatedUser.UserID,
		Email:           "",
		DesiredUsername: "hot_mama_RAWR_XD",
		DesiredPassword: "",
		Birthday:        nil,
		City:            "Philadelphia",
	}

	updateRes, err := client.UpdateUserInfo(ctx, updateReq)

	fmt.Printf("Update res: %v\nerr: %v\n", updateRes, err)

	if err != nil {
		log.Fatalf("fail to update user: %v", err)
	}


	// ---------------------------------------------------------

	// Getting user again

	md = metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", tokRes.Token))
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	usernameRes, err = client.GetUserByUsername(ctx, &pbusers.GetUserByUsernameRequest{Username: "hot_mama_RAWR_XD"})


	fmt.Printf("Get res: %v\nerr: %v\n", usernameRes, err)

	// -------------------------

	md = metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", tokRes.Token))
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	deleteRes, err := client.DeleteUserByID(ctx, &pbusers.DeleteUserByIDRequest{UserID: addRes.CreatedUser.UserID})

	if err != nil {
		log.Fatalf("fail to delete user: %v", err)
	}

	fmt.Printf("deleteRes: %v\n", deleteRes)

}