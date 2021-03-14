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

	res, err := client.AddUser(context.Background(), in)

	fmt.Printf("res: %v\nerr: %v\n", res, err)

	if err != nil {
		log.Fatalf("fail to add user: %v", err)
	}

	res, err = client.AddUser(context.Background(), in)

	fmt.Printf("res: %v\nerr: %v\n", res, err)

	tokRes, err := client.GetJWTToken(context.Background(), &pbusers.GetJWTTokenRequest{
		Username: "qdnovinger",
		Password: "password134234",
	})

	if err != nil {
		log.Fatalf("fail to get token: %v", err)
	}

	fmt.Printf("tokRes: %v\n", tokRes)

	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", tokRes.Token), "x-ext-authz", "allow")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	usernameRes, err := client.GetUserByUsername(ctx, &pbusers.GetUserByUsernameRequest{Username: "qdnovinger"})


	fmt.Printf("res: %v\nerr: %v\n", usernameRes, err)

}