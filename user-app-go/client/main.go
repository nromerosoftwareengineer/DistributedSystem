package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	pb "user-project-go/user-app"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:50051: %v", err)
	}
	defer conn.Close()

	us := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//user := &pb.User{
	//	UserId:    "1021",
	//	FirstName: "Liza",
	//	LastName:  "wick",
	//	City:      "Lufkin",
	//	State:     "CA",
	//	Address1:  "123 npc St",
	//	Address2:  "Apt 22",
	//	Zip:       "92345",
	//}
	//
	//r, err := us.PutUser(ctx, &pb.UserPutRequest{
	//	User: user,
	//})
	//log.Printf("Response from gRPC server's: %s", r.GetMessage())

	updatedUser := pb.User{
		UserId:    "1001",
		FirstName: "wizard",
		LastName:  "nick",
		City:      "Fortworth",
		State:     "TX",
		Address1:  "3012",
		Address2:  "Apt 22",
		Zip:       "76106",
	}

	r, err := us.UpdateUser(ctx, &pb.UpdateUserRequest{User: &updatedUser})

	log.Printf("Response from gRPC server's UpdateUser service: %s", r.GetMessage())

	user, err := us.GetUser(ctx, &pb.GetUserRequest{UserId: "1001"})

	log.Printf("Response from gRPC server's: %s", user.GetUser().FirstName)

	//msg, err := us.DeleteUser(ctx, &pb.DeleteUserRequest{UserId: "1021"})

	//log.Printf("Response from gRPC server's DeleteUser service: %s", msg.GetMessage())

	userAllList, err := us.ListAllUsers(context.Background(), &pb.ListAllUsersRequest{})

	for _, user := range userAllList.Users {
		log.Println("User:", user)
	}

	if err != nil {
		log.Fatalf("error calling function SayHello: %v", err)
	}

	log.Printf("Response from gRPC server's: %s", user.GetUser().FirstName)
}
