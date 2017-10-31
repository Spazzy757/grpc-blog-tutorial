package main

import (
	"log"
	"net"
	"strings"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/Spazzy757/grpcpoc/user"
)

const (
	port = ":50051"
)

// server is used to implement customer.CustomerServer.
type server struct {
	savedUsers []*pb.UserRequest
}


// CreateUser creates a new Customer
func (s *server) CreateUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	s.savedUsers = append(s.savedUsers, in)
	return &pb.UserResponse{Id: in.Id, Success: true}, nil
}

// GetUser returns all customers by given filter
func (s *server) GetUsers(filter *pb.UserFilter, stream pb.User_GetUsersServer) error {
	for _, user := range s.savedUsers {
		if filter.Keyword != "" {
			if !strings.Contains(user.Name, filter.Keyword) {
				continue
			}
		}
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterUserServer(s, &server{})
	s.Serve(lis)
}