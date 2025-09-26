# BM Staff - Go Microservice

A Go microservice built with Clean Architecture principles, featuring User management with Oracle database integration.

## 🏗️ Architecture

This project follows Clean Architecture with the following layers:

- **Domain Layer**: Business entities, value objects, and domain services
- **Use Cases Layer**: Application business rules and orchestration
- **Interface Adapters Layer**: HTTP handlers, repositories, and external interfaces
- **Infrastructure Layer**: Database connections, logging, and configuration

## 🚀 Features

- **Clean Architecture** implementation
- **Oracle Database** integration
- **RESTful API** with Gin framework
- **Structured logging** with Zap
- **Error handling** with standardized error codes
- **Request validation** with Go Playground Validator
- **Dependency injection** with Google Wire
- **Configuration management** with Viper
- **Database migrations** support

## 📁 Project Structure

```
bm-staff/
├── cmd/                           # Application entrypoints
│   └── api/                      # HTTP API server
├── internal/                     # Private application code
│   ├── domain/                   # Business entities and rules
│   ├── usecases/                # Application business rules
│   ├── interfaces/              # Interface adapters
│   ├── infrastructure/          # Frameworks and drivers
│   └── di/                      # Dependency injection
├── pkg/                         # Public packages
├── configs/                     # Configuration files
├── migrations/                  # Database migrations
└── docs/                        # Documentation
```

## 🛠️ Prerequisites

- Go 1.21 or higher
- Oracle Database 12c or higher
- Access to Oracle database at `192.168.7.248:1521`

## 📦 Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd bm-staff
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up Oracle database:
```bash
# Run the migration script
sqlplus HXT_RS/HXT_RS@192.168.7.248:1521/orclstb @migrations/001_create_users.sql
```

## 🚀 Running the Application

1. Start the application:
```bash
go run cmd/api/main.go
```

2. The API will be available at `http://localhost:8080`

## 📚 API Endpoints

### Users

- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `GET /api/v1/users` - List users (with pagination)

### Health Check

- `GET /health` - Health check endpoint

## 📝 Example API Usage

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  }'
```

### Get User
```bash
curl http://localhost:8080/api/v1/users/{user-id}
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe_updated",
    "email": "john.updated@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user-id}
```

## ⚙️ Configuration

Configuration is managed through YAML files in the `configs/` directory:

- `development.yaml` - Development environment settings
- Database connection settings
- Server configuration
- Logging configuration

## 🧪 Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## 📊 Database Schema

### BMSF Table Prefix Convention
All application tables use the **BMSF_** prefix to identify BM Staff Framework tables:
- **BMSF_USER** - User management table
- **BMSF_ORDER** - Order management table (future)
- **BMSF_PRODUCT** - Product management table (future)

This convention helps:
- Avoid conflicts with existing database tables
- Identify system-specific tables easily
- Apply security policies by table prefix
- Organize database objects by system

### BMSF_USER Table
- `id` (VARCHAR2(36)) - Primary key (UUID)
- `username` (VARCHAR2(50)) - Unique username
- `email` (VARCHAR2(255)) - Unique email address
- `first_name` (VARCHAR2(100)) - First name
- `last_name` (VARCHAR2(100)) - Last name
- `phone` (VARCHAR2(20)) - Phone number
- `status` (VARCHAR2(20)) - User status (ACTIVE, INACTIVE, PENDING, BLOCKED)
- `created_at` (TIMESTAMP) - Creation timestamp
- `updated_at` (TIMESTAMP) - Last update timestamp
- `created_by` (VARCHAR2(36)) - Creator user ID
- `updated_by` (VARCHAR2(36)) - Last updater user ID
- `deleted_at` (TIMESTAMP) - Soft delete timestamp
- `version` (NUMBER(10)) - Optimistic locking version
- `tenant_id` (VARCHAR2(36)) - Multi-tenant identifier

## 🔧 Development

### Code Generation
```bash
# Generate mocks
go generate ./...

# Generate Wire dependencies
wire ./internal/di
```

### Linting
```bash
# Run linter
golangci-lint run

# Format code
go fmt ./...
goimports -w .
```

## 📋 Error Codes

The application uses standardized error codes:

- **1xxx** - System Errors
- **2xxx** - Validation Errors
- **3xxx** - Authentication/Authorization
- **4xxx** - Business Logic
- **5xxx** - External Dependencies

## 🚀 Deployment

### Docker
```bash
# Build image
docker build -t bm-staff .

# Run container
docker run -p 8080:8080 bm-staff
```

### Environment Variables
- `DATABASE_HOST` - Database host
- `DATABASE_PORT` - Database port
- `DATABASE_USERNAME` - Database username
- `DATABASE_PASSWORD` - Database password
- `DATABASE_SERVICE_NAME` - Database service name

## 📄 License

This project is licensed under the MIT License.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📞 Support

For support, please contact the development team or create an issue in the repository.
