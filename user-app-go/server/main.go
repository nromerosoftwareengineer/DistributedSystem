package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	db "user-project-go/dataBase"
	pb "user-project-go/user-app"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	dbService, err := db.NewDBService()
	if err != nil {
		log.Fatalf("Failed to create database service: %v", err)
	}
	defer dbService.Close()
	user := in.GetUser()
	updateUser := db.User{
		UserID:    user.GetUserId(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		City:      user.GetCity(),
		State:     user.GetState(),
		Address1:  user.GetAddress1(),
		Address2:  user.GetAddress2(),
		Zip:       user.GetZip(),
	}

	msg := ""
	msg, err = dbService.UpdateUser(ctx, updateUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}

	return &pb.UpdateUserResponse{Message: msg}, nil
}

func (s *server) ListAllUsers(ctx context.Context, in *pb.ListAllUsersRequest) (*pb.ListAllUserResponse, error) {
	dbService, err := db.NewDBService()

	if err != nil {
		log.Fatalf("Failed to create database service: %v", err)
	}
	defer dbService.Close()

	// Get user from database
	dbUsers, err := dbService.ListUsers(ctx)
	if err != nil {
		log.Fatalf("Failed to retrieve user : %v", err)
	}
	protoUsers := make([]*pb.User, len(dbUsers))
	for _, dbUser := range dbUsers {
		protoUser := &pb.User{
			UserId:    dbUser.UserID,
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			City:      dbUser.City,
			State:     dbUser.State,
			Address1:  dbUser.Address1,
			Address2:  dbUser.Address2,
			Zip:       dbUser.Zip,
		}
		protoUsers = append(protoUsers, protoUser)
	}
	// Then return in your ListUsers response:
	return &pb.ListAllUserResponse{
		Users: protoUsers,
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	dbService, err := db.NewDBService()
	if err != nil {
		log.Fatalf("Failed to create database service: %v", err)
	}
	defer dbService.Close()
	err = dbService.DeleteUser(ctx, in.GetUserId())
	if err != nil {
		log.Printf("Error deleting user: %v", err)
	}

	msg := "Deleted userid: " + in.UserId + " successfully"

	return &pb.DeleteUserResponse{Message: msg}, nil
}

func (s *server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	dbService, err := db.NewDBService()
	if err != nil {
		log.Fatalf("Failed to create database service: %v", err)
	}
	defer dbService.Close()

	// Get user from database
	dbUser, err := dbService.GetUser(ctx, in.UserId)
	if err != nil {
		log.Fatalf("Failed to retrieve user : %v", err)
	}

	// Map the database User to protobuf User
	protoUser := &pb.User{
		UserId:    dbUser.UserID,    // Use the actual values from dbUser
		FirstName: dbUser.FirstName, // instead of hardcoded values
		LastName:  dbUser.LastName,
		City:      dbUser.City,
		State:     dbUser.State,
		Address1:  dbUser.Address1,
		Address2:  dbUser.Address2,
		Zip:       dbUser.Zip,
	}

	return &pb.GetUserResponse{
		User: protoUser,
	}, nil

}

func (s *server) PutUser(ctx context.Context, in *pb.UserPutRequest) (*pb.UserPutResponse, error) {
	dbService, err := db.NewDBService()
	if err != nil {
		log.Fatalf("Failed to create database service: %v", err)
	}
	defer dbService.Close()
	user := in.GetUser()
	newUser := db.User{
		UserID:    user.GetUserId(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		City:      user.GetCity(),
		State:     user.GetState(),
		Address1:  user.GetAddress1(),
		Address2:  user.GetAddress2(),
		Zip:       user.GetZip(),
	}

	msg := ""
	msg, err = dbService.CreateUser(ctx, newUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}

	return &pb.UserPutResponse{Message: msg}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
