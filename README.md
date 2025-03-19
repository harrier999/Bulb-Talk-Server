# Bulb Talk Server

<div align="center">

![Bulb Talk Logo](https://via.placeholder.com/200x200.png?text=Bulb+Talk)

**A modern, real-time chat platform built with Go**

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[Live Demo](https://talk.wasabi-labs.com) • [API Documentation](API_WebSocket_Documentation.md) • [Report Bug](https://github.com/yourusername/bulb-talk/issues)

</div>

## 📋 Table of Contents

- [Overview](#overview)
- [Live Demo](#live-demo)
- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running with Docker](#running-with-docker)
  - [Running Locally](#running-locally)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## 🌟 Overview

Bulb Talk is a modern, real-time chat application built with Go. It provides a robust backend for messaging, friend management, and room-based conversations using WebSockets for real-time communication.

## 💻 Live Demo

Experience Bulb Talk in action by visiting our live demo at [talk.wasabi-labs.com](https://talk.wasabi-labs.com).

Sign up with your phone number, add friends, create chat rooms, and start messaging in real-time!

## ✨ Features

- **User Authentication**: Secure sign-up and login with JWT
- **Friend Management**: Add, block, and unblock friends
- **Real-time Messaging**: Instant messaging using WebSockets
- **Chat Rooms**: Create and manage group conversations
- **Typing Indicators**: See when others are typing
- **Online Presence**: Track when users join and leave rooms
- **RESTful API**: Comprehensive API for client development

## 🏗 Architecture

Bulb Talk is built with clean architecture principles and uses modern technologies:

- **Backend**: Go (Golang) with Gorilla WebSockets and Mux router
- **Database**: PostgreSQL for persistent storage
- **Cache**: Redis for session management and real-time features
- **Authentication**: JWT (JSON Web Tokens) for secure authentication
- **API**: RESTful API with WebSocket support

The system architecture follows SOLID principles:

- **Single Responsibility**: Each component has a single responsibility
- **Open/Closed**: Components are open for extension but closed for modification
- **Liskov Substitution**: Interface implementations are interchangeable
- **Interface Segregation**: Specific interfaces for different concerns
- **Dependency Inversion**: High-level modules depend on abstractions

### Architecture Diagram

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Clients   │────▶│  API Layer  │────▶│  Handlers   │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                                               ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  PostgreSQL │◀───▶│Repositories │◀────│  Services   │
└─────────────┘     └─────────────┘     └─────────────┘
       ▲                   ▲
       │                   │
       │             ┌─────────────┐
       └─────────────│    Redis    │
                     └─────────────┘
```

## 📁 Project Structure

The project follows a clean architecture pattern:

```
.
├── cmd/                  # Application entry points
│   ├── migration/        # Database migration
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

## 🚀 Getting Started

### Prerequisites

To run Bulb Talk Server, you'll need:

- Go (version 1.18 or higher)
- PostgreSQL (version 12 or higher)
- Redis (version 6 or higher)

### Dependencies

The main dependencies include:

- **gorilla/mux**: HTTP router and URL matcher
- **gorilla/websocket**: WebSocket implementation
- **golang-jwt/jwt**: JWT authentication
- **lib/pq**: PostgreSQL driver for Go
- **go-redis/redis**: Redis client for Go
- **google/uuid**: UUID generation
- **joho/godotenv**: Environment variable loading

### Running with Docker

The easiest way to run Bulb Talk is using Docker Compose:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/bulb-talk.git
   cd bulb-talk
   ```

2. Run with Docker Compose:
   ```bash
   docker-compose up
   ```

This will start the server along with PostgreSQL and Redis containers.

### Running Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/bulb-talk.git
   cd bulb-talk
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   cp cmd/talk-server/.env.dist cmd/talk-server/.env
   # Edit .env file with your configuration
   ```

<<<<<<< Updated upstream
3. Run the application:
   ```
   cd cmd/talk-server
=======
4. Run the application:
   ```bash
   cd cmd/talk-server
>>>>>>> Stashed changes
   go run main.go
   ```

The server will start on port 18000 by default.

## 📚 API Endpoints

Bulb Talk provides a comprehensive RESTful API. For detailed documentation, see [API_WebSocket_Documentation.md](API_WebSocket_Documentation.md).

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
- `GET /auth/friends`: Get friend list
- `POST /auth/friends`: Add a friend by phone number
- `PUT /auth/friends/{friendId}/block`: Block a friend
- `PUT /auth/friends/{friendId}/unblock`: Unblock a friend

#### Room Management
- `GET /auth/rooms`: Get room list
- `POST /auth/rooms`: Create a new room
- `POST /auth/rooms/{roomId}/users`: Add a user to a room
- `DELETE /auth/rooms/{roomId}/users/{userId}`: Remove a user from a room

## 🧪 Testing

Run tests with:
```bash
go test ./...
```

Or use the helper script:
```bash
./test/run_tests.sh
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

Built with ❤️ by the Bulb Talk Team 
