# Bulb Talk Server

This is the server component of the Bulb Talk application, a chat platform built with Go.

## Project Structure

The project follows SOLID principles and clean architecture:

```
.
├── cmd/                  # Application entry points
│   ├── migration/        # Migrate database
│   └── talk-server/      # Main server application 
├── internal/             # Internal packages
│   ├── db/               # Database connections
│   ├── handler/          # HTTP handlers
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic layer
│   └── utils/            # Utility functions
├── pkg/                  # Public packages
│   ├── authenticator/    # Authentication utilities
│   └── log/              # Logging utilities
└── test/                 # Test utilities
```

## Architecture

The application follows a clean architecture with the following layers:

1. **Repository Layer**: Handles data access and persistence
   - Interfaces defined in `internal/repository/interfaces.go`
   - Implementations in `internal/repository/{postgres,redis}/`

2. **Service Layer**: Contains business logic
   - Interfaces defined in `internal/service/interfaces.go`
   - Implementations in `internal/service/`

3. **Handler Layer**: Manages HTTP requests and responses
   - Implementations in `internal/handler/`

## SOLID Principles

The application follows SOLID principles:

- **Single Responsibility Principle**: Each class has a single responsibility
- **Open/Closed Principle**: Classes are open for extension but closed for modification
- **Liskov Substitution Principle**: Implementations can be substituted for their interfaces
- **Interface Segregation Principle**: Specific interfaces for different concerns
- **Dependency Inversion Principle**: High-level modules depend on abstractions

## Running the Application

1. Install dependencies:
   ```
   go mod download
   ```

2. Set up environment variables:
   ```
   cp cmd/talk-server/.env.example cmd/talk-server/.env
   # Edit .env file with your configuration
   ```

3. Run the application:
   ```
   cd cmd/talk-server
   go run main.go
   ```

## API Endpoints

### Public Endpoints

- `POST /signup`: Register a new user
- `POST /login`: Log in a user
- `POST /authenticate`: Request an authentication number
- `POST /checkauth`: Check an authentication number
- `GET /chat`: WebSocket endpoint for chat
- `GET /messages`: Get messages for a room

### Authorized Endpoints

All authorized endpoints require a JWT token in the Authorization header.

#### Friend Management
- `POST /auth/getfriends`: Get friend list
- `POST /auth/addfriends`: Add a friend by phone number
- `POST /auth/blockfriend`: Block a friend
- `POST /auth/unblockfriend`: Unblock a friend

#### Room Management
- `POST /auth/rooms`: Get room list
- `POST /auth/createrooms`: Create a new room
- `POST /auth/adduser`: Add a user to a room
- `POST /auth/removeuser`: Remove a user from a room

## Testing

Run tests with:
```
go test ./...
``` 