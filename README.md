# This is part of a blog post for a gRPC server and client

This runs a server that receives a User object via rpc from a client,
use case being communication of user information between microservice containers

## Setup
you will need Go installed (my version being 1.9)

**Start Server**
```bash
go run server/main.go
```

**Start Client**
```bash
go run client/main.go
```

## Usage
This will have a simple command line interface where you add user information,
this will then be transfered to the server where the users information will be stored (no persistence)
Once the server gets all the users information, the client then requests a list of all users
which it then prints out to the terminal
