# User Management Service

A gRPC-based user management service with PostgreSQL backend. This service provides CRUD operations for user data management with support for pagination.

## Features

- gRPC API for user management
- PostgreSQL database integration
- Basic CRUD operations (Create, Read, Update, Delete)
- Pagination support for user listing
- Environment variable configuration
- Structured error handling

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 14 or higher
- Protocol Buffers compiler (protoc)
- Docker (optional)

## Installation

1. Clone the repository:
```bash
git clone [your-repository-url]
cd [repository-name]
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
```bash
# Create the PostgreSQL database
createdb noelromero

# Set up environment variables
export DATABASE_URL="postgresql://noelromero:your_password@localhost:5432/noelromero"
```

4. Create the users table:
```sql
CREATE TABLE users (
    userid VARCHAR(50) PRIMARY KEY NOT NULL,
    firstName VARCHAR(100) NOT NULL,
    lastName VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(50) NOT NULL,
    address1 VARCHAR(255) NOT NULL,
    address2 VARCHAR(255) NOT NULL,
    zip VARCHAR(20) NOT NULL
);
```

## Configuration

The service can be configured using environment variables:

```bash
# Database configuration
DATABASE_URL=postgresql://username:password@localhost:5432/dbname

# Server configuration (optional)
GRPC_PORT=50051
```

## Usage

1. Start the server:
```bash
go run cmd/server/main.go
```

2. Run the client:
```bash
go run cmd/client/main.go
```

### Example API Calls

Creating a new user:
```go
client := pb.NewUserServiceClient(conn)
response, err := client.CreateUser(ctx, &pb.UserPutRequest{
    User: &pb.User{
        UserId: "1001",
        FirstName: "John",
        LastName: "Doe",
        City: "San Francisco",
        State: "CA",
        Address1: "123 Main St",
        Address2: "Apt 4B",
        Zip: "94105",
    },
})
```

Retrieving a user:
```go
response, err := client.GetUser(ctx, &pb.GetUserRequest{
    UserId: "1001",
})
```

## Project Structure

```
.
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── client/
│       └── main.go
├── internal/
│   ├── db/
│   │   └── service.go
│   └── server/
│       └── server.go
├── proto/
│   └── user-app/
│       ├── helloworld.proto
│       ├── helloworld.pb.go
│       └── helloworld.pb.go
├── go.mod
├── go.sum
└── README.md
```

## Development

To regenerate the Protocol Buffers code:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/user/user.proto
```

## Testing

Run the tests:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, please open an issue in the GitHub repository.
