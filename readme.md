# API Frontend Team – Auth & User Service

This repository contains a Go-based backend service for user authentication and profile management, containerized using Docker Compose with PostgreSQL.

## 🚀 Features

- User Registration and Login
- JWT-based Authentication with Refresh Tokens
- Authenticated User Profile Management (Get, Update, Delete)
- RESTful API Design with versioning (`/api/v1`)
- Dockerized environment for easy setup

## API Endpoints

### 🔐 Auth Routes (`/api/v1/auth`)

| Method | Endpoint         | Description          |
| ------ | ---------------- | -------------------- |
| POST   | `/register`      | Register new user    |
| POST   | `/login`         | Login with email     |
| POST   | `/refresh-token` | Refresh access token |

#### Example: `LoginRequest`

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Example: `RefreshTokenRequest`

```json
{
  "refresh_token": "your_refresh_token"
}
```

### Auth Routes (`/api/v1/auth`)

| Method | Endpoint  | Description            |
| ------ | --------- | ---------------------- |
| GET    | `/:email` | Get user by email      |
| GET    | `/me`     | Get authenticated user |
| PUT    | `/me`     | Update user info       |
| DELETE | `/me`     | Delete user account    |

## Running the App

```bash
# Clone the repo
git clone https://github.com/EngenMe/api-frontend-team.git
cd api-frontend-team

# Start services
docker-compose up --build

```

## 🔒 Authentication

- All /user routes are protected using middleware that verifies JWT access tokens. Use the /login and /refresh-token endpoints to obtain and renew access.

📬 Response Format

RefreshTokenResponse & login

```json
{
  "tokens": {
    "access": {
      "token": "access_jwt_token",
      "expires": "2025-05-16 10:29:45.171970595 +0000 UTC m=+259275.169498203"
    },
    "refresh": {
      "token": "refresh_jwt_token",
      "expires": "2025-05-19 10:29:45.171970595 +0000 UTC m=+259275.169498203"
    }
  }
}
```
