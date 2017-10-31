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
)

const (
	address = "127.0.0.1:50051"
)

// createCustomer calls the RPC method CreateCustomer of CustomerServer
func createCustomer(client pb.UserClient, user *pb.UserRequest) {
	resp, err := client.CreateUser(context.Background(), user)
	if err != nil {
		log.Fatalf("Could not create Customer: %v", err)
	}
	if resp.Success {
		log.Printf("A new Customer has been added with id: %d", resp.Id)
	}
}

// getCustomers calls the RPC method GetCustomers of CustomerServer
func getCustomers(client pb.UserClient, filter *pb.UserFilter) {
	// calling the streaming API
	stream, err := client.GetUsers(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}
		log.Printf("Customer: %v", customer)
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

	var id int32
	fmt.Print("Enter ID:")
	_, err = fmt.Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Users Name:")
	name, _ := reader.ReadString('\n')
	fmt.Print("Enter Users Email:")
	email, _ := reader.ReadString('\n')
	fmt.Print("Enter Users Password:")
	password, _ := reader.ReadString('\n')

	customer := &pb.UserRequest{
		Id:    id,
		Name:  name,
		Email: email,
		Password: password,
	}
	// Create a new customer
	createCustomer(client, customer)
	//// Filter with an empty Keyword
	filter := &pb.UserFilter{Keyword: "Brendan"}
	getCustomers(client, filter)
}
