# Movie Review API

A Go-based REST API for a movie review application built with Gin framework, PostgreSQL, and JWT authentication.

## Features

- User authentication with JWT tokens
- Movie management and listing
- Movie reviews and ratings
- User profiles
- Category filtering
- Search functionality
- Docker support for containerized deployment

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT (golang-jwt)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Containerization**: Docker & Docker Compose

## Project Structure

```
api/
├── config/          # Configuration management
├── models/          # Database models (User, Token, Movie, Review, etc.)
├── repository/      # Database access layer
├── service/         # Business logic layer
├── handler/         # HTTP handlers/controllers
├── middleware/      # Auth middleware and other middleware
├── utils/           # Utility functions
├── sql/             # Database initialization scripts
├── main.go          # Application entry point
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

## Setup & Installation

### Prerequisites
- Go 1.21+
- PostgreSQL 13+ (or use Docker)
- Docker & Docker Compose (optional)

### Local Development

1. Clone the repository and navigate to the api directory:
```bash
cd api
```

2. Copy the environment file:
```bash
cp .env.example .env
```

3. Install dependencies:
```bash
go mod download
```

4. Start PostgreSQL (using Docker):
```bash
docker run -d \
  --name movie_review_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=movie_review \
  -p 5432:5432 \
  postgres:15-alpine
```

5. Run database migrations:
```bash
psql -h localhost -U postgres -d movie_review -f sql/init.sql
```

6. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

### Docker Compose (Recommended)

```bash
docker-compose up -d
```

This will start both PostgreSQL and the API.

## API Endpoints

### Authentication
- `POST /api/auth/signup` - Register new user
- `POST /api/auth/login` - Login and get JWT token
- `POST /api/auth/logout` - Logout and revoke token
- `GET /api/auth/me` - Get current user profile

### Movies
- `GET /api/movies` - List all movies (paginated)
- `GET /api/movies/:id` - Get single movie by ID
- `GET /api/movies/category/:category` - Get movies by category
- `GET /api/movies/search?q=query` - Search movies
- `GET /api/movies/top-rated` - Get top-rated movies
- `GET /api/movies/latest` - Get latest movies
- `POST /api/movies` - Create movie (admin only)
- `PUT /api/movies/:id` - Update movie (admin only)
- `DELETE /api/movies/:id` - Delete movie (admin only)

### Reviews
- `POST /api/movies/:id/reviews` - Create review for a movie
- `GET /api/movies/:id/reviews` - Get all reviews for a movie
- `GET /api/reviews/:id` - Get single review
- `PUT /api/reviews/:id` - Update review (own review only)
- `DELETE /api/reviews/:id` - Delete review (own review only)
- `GET /api/users/:id/reviews` - Get all reviews by user

### Stats
- `GET /api/movies/:id/rating` - Get average rating for a movie
- `GET /api/categories` - Get all available categories

## Example Usage

### Signup
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "secure_password"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "secure_password"
  }'
```

### Get Movies
```bash
curl -X GET "http://localhost:8080/api/movies?page=1&limit=10" \
  -H "Authorization: Bearer <your_jwt_token>"
```

### Create Review
```bash
curl -X POST http://localhost:8080/api/movies/1/reviews \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "rating": 9,
    "comment": "Great movie!"
  }'
```

## Development

### Run tests
```bash
go test ./...
```

### Build binary
```bash
go build -o movie-api
```

### Run with hot reload (requires air)
```bash
go install github.com/cosmtrek/air@latest
air
```

## Configuration

Environment variables (see `.env.example`):
- `DB_HOST` - PostgreSQL host
- `DB_PORT` - PostgreSQL port
- `DB_USER` - PostgreSQL user
- `DB_PASSWORD` - PostgreSQL password
- `DB_NAME` - Database name
- `JWT_SECRET` - JWT secret key (min 32 chars for HS256)
- `JWT_EXPIRATION_MS` - Token expiration time in milliseconds
- `SERVER_PORT` - Server port

## Database Schema

The database includes the following tables:
- `users` - User accounts
- `tokens` - JWT tokens with revocation support
- `movies` - Movie information
- `categories` - Movie categories
- `reviews` - User reviews with ratings
- `movie_categories` - Many-to-many relationship

## Security

- Passwords are hashed with bcrypt
- JWT tokens are signed with HS256
- Tokens can be revoked on logout
- Protected endpoints require valid JWT token
- CORS headers configured for frontend integration

## License

MIT
