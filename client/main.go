package main

import (
	pb "github.com/Spazzy757/grpcpoc/user"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"bufio"
	"log"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	address = "127.0.0.1:50051"
)

// createUser calls the RPC method CreateUser of UserServer
func createUser(client pb.UserClient, user *pb.UserRequest) {
	resp, err := client.CreateUser(context.Background(), user)
	if err != nil {
		log.Fatalf("Could not create User: %v", err)
	}
	if resp.Success {
		log.Printf("A new User has been added with id: %d", resp.Id)
	}
}

// getCustomers calls the RPC method GetCustomers of CustomerServer
func getUsers(client pb.UserClient, filter *pb.UserFilter) {
	// calling the streaming API
	stream, err := client.GetUsers(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get users: %v", err)
	}

	// Runs until the io.EOF signal is received from the server
	// Returns all users
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetUsers(_) = _, %v", client, err)
		}
		log.Printf("Users: %v", customer)
	}
}

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Creates a new CustomerClient
	client := pb.NewUserClient(conn)

	for {
		// This is just a simple input of data from the CLI used to get input to send to the server
		// This is for the ID
		var id int32
		fmt.Print("Enter ID:")
		_, err = fmt.Scan(&id)

		if err != nil {
			log.Fatal(err)
		}

		// This is for the strings
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Users Name:")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSuffix(name, "\n")
		fmt.Print("Enter Users Email:")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSuffix(email, "\n")
		fmt.Print("Enter Users Password:")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSuffix(password, "\n")

		customer := &pb.UserRequest{
			Id:    id,
			Name:  name,
			Email: email,
			Password: password,
		}

		// Create a new customer
		createUser(client, customer)

		// Filter with an empty Keyword, this can later be changed to search through the users
		filter := &pb.UserFilter{Keyword: ""}
		getUsers(client, filter)
	}
}
