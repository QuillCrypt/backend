# QuillCrypt Backend

QuillCrypt is a secure, end-to-end encrypted messaging platform for developers. This repository contains the backend service, built with Go and designed using Hexagonal Architecture (Ports and Adapters) for high maintainability and scalability.

> **Note**: This project is currently under active development. Developers are advised NOT TO USE THIS BACKEND until versioned releases start as breaking changes may appear. Many features listed below are planned and will be implemented in upcoming releases.

## 🚀 Planned Features

- **End-to-End Encryption Support**: Robust key management and secure message routing protocols.
- **Real-time Communication**: High-performance WebSocket-based messaging for instant delivery.
- **Hexagonal Architecture**: Deeply decoupled business logic for maximum flexibility and testability.
- **Scalable Infrastructure**: Integration with PostgreSQL for persistent storage and Redis for session management.
- **Fully Containerized**: Standardized Docker environment for consistent development and deployment.

## 🛡️ Cryptography Notice

This software is designed to facilitate encrypted communications. Please be aware that the export, import, and use of certain cryptographic software may be restricted in some jurisdictions. Before using or distributing this software, please verify the local laws and regulations regarding the use of cryptography.

## 🛠 Tech Stack

- **Language**: [Go](https://go.dev/) (v1.26+)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Cache/Session**: [Redis](https://redis.io/)
- **Logging**: [Uber-Go Zap](https://github.com/uber-go/zap)
- **Configuration**: [envconfig](https://github.com/kelseyhightower/envconfig) & [godotenv](https://github.com/joho/godotenv)
- **Containerization**: [Docker](https://www.docker.com/) & Docker Compose

## 📂 Project Structure

```text
├── cmd/                   # Application entry points (e.g., the main server)
├── internal/              # Private application code
│   ├── api/               # HTTP delivery layer
│   │   ├── handler/       # Request handlers (controllers)
│   │   ├── middleware/    # HTTP middleware (auth, logging, etc.)
│   │   └── router/        # Route definitions
│   ├── config/            # Configuration and environment management
│   ├── core/              # Business logic (Hexagonal Core)
│   │   ├── domain/        # Domain entities and models
│   │   ├── port/          # Interface definitions (repository and service ports)
│   │   └── service/       # Business logic implementations
│   ├── repository/        # Data persistence implementations
│   │   ├── postgres/      # PostgreSQL repository implementation
│   │   └── redis/         # Redis repository implementation (sessions, cache)
│   └── websocket/         # Real-time communication hub and client logic
└── pkg/                   # Publicly shared utilities (e.g., logger)
```

## 🚦 Getting Started

### Prerequisites

- Go 1.26 or higher
- Docker and Docker Compose

### Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/quillcrypt-backend.git
   cd quillcrypt-backend
   ```

2. **Configure Environment**:
   Copy the example environment file and adjust the values as needed:
   ```bash
   cp .env.example .env
   ```

3. **Start Infrastructure**:
   Use Docker Compose to spin up the database and Redis:
   ```bash
   docker-compose up -d
   ```

4. **Run the Server**:
   ```bash
   go run cmd/server/main.go
   ```

## 📖 API Documentation

*(Planned - Swagger/OpenAPI documentation will be integrated as development progresses)*

## 🔐 Security

QuillCrypt is being built with a security-first mindset. The backend is designed to handle message routing and key exchange while ensuring that only the intended recipients can decrypt messages via client-side or hardware-backed keys.

## 📄 License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)** - see the [LICENSE](LICENSE) file for details.
