# Fitness API with GraphQL and Go

A modern, scalable fitness tracking platform built with Go, Encore, and GraphQL. This project provides a robust backend for managing fitness data, including trainee profiles, workout plans, progress tracking, and trainer-trainee interactions.

## 🚀 Technologies

- **Backend**: Go 1.23+
- **API Framework**: Encore
- **GraphQL**: gqlgen
- **Database**: PostgreSQL
- **Containerization**: Docker
- **Authentication**: JWT

## ✨ Features

- **Admin Service**: User management, authentication, and system configuration
- **Trainee Service**: Fitness tracking, workout management, progress monitoring
- **GraphQL API**: Type-safe, self-documenting API with real-time capabilities
- **Database Migrations**: Automated schema management
- **Developer Experience**: Built-in testing, local development, and observability

## 🛠 Prerequisites

- **Go 1.23** or later ([Installation Guide](https://golang.org/doc/install))
- **Encore CLI** (Install with one of the following):
  ```bash
  # macOS
  brew install encoredev/tap/encore
  
  # Linux
  curl -L https://encore.dev/install.sh | bash
  
  # Windows (PowerShell)
  iwr https://encore.dev/install.ps1 | iex
  ```
- **Docker** ([Installation Guide](https://docs.docker.com/get-docker/))
- **Git** ([Installation Guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git))

## 🚀 Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/hadisa/fitness-app-encore.git
   cd fitness-app-encore
   ```
2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Start the development environment**
   Make sure Docker is running, then start the application:
   ```bash
   encore run
   ```

4. **Access the development dashboard**
   Open [http://localhost:9400](http://localhost:9400) in your browser to access Encore's local developer dashboard.

## 🎮 GraphQL Playground

Access the GraphQL Playground at [http://localhost:4000/graphql/playground](http://localhost:4000/graphql/playground) to explore and test the GraphQL API.

## 📚 Project Structure

```
.
├── admin/                  # Admin service
│   ├── migrations/         # Database migrations
│   └── admin.go            # Admin service implementation
├── trainee/                # Trainee service
│   ├── migrations/         # Database migrations
│   └── trainee.go          # Trainee service implementation
├── graphql/                # GraphQL schema and resolvers
│   ├── admin.graphqls      # Admin GraphQL schema
│   ├── admin.resolvers.go  # Admin resolvers
│   ├── trainee.graphqls    # Trainee GraphQL schema
│   └── trainee.resolvers.go # Trainee resolvers
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 🔧 Code Generation

This project uses `gqlgen` for GraphQL code generation. After modifying any `.graphqls` files, run:

```bash
go run github.com/99designs/gqlgen generate
```

## 🌐 API Documentation

### Admin Service
- **Authentication**
- **User Management**
- **System Configuration

### Trainee Service
- **Profile Management**
- **Workout Tracking**
- **Progress Monitoring**
- **Trainer-Trainee Communication**

## 🧪 Running Tests

```bash
# Run all tests
encore test ./...

# Run tests with coverage
encore test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📦 Deployment

This application is designed to be deployed using Encore's cloud platform:

```bash
# Deploy to staging
encore app deploy --env=staging

# Deploy to production
encore app deploy --env=production
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Hadisa Norozi**

- GitHub: [@hadisa](https://github.com/hadisa)
- Email: hadisa.norozi@gmail.com

## 🙏 Acknowledgments

- Built with [Encore](https://encore.dev/)
- GraphQL implementation using [gqlgen](https://gqlgen.com/)
- Database powered by PostgreSQL