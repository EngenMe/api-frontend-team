# Wikidocify API Server

This is a Go-based API server for wikidocify, providing user authentication and profile management, containerized with Docker Compose and PostgreSQL.

## 🚀 Features
- User registration and login with JWT-based authentication
- Profile management (get, update, delete)
- RESTful API with versioning (`/api/v1`)

## Getting Started

### Prerequisites
- Docker and Docker Compose

### Running the Project
1. Clone the repository:
```bash
git clone https://github.com/EngenMe/api-frontend-team.git
cd api-frontend-team
```
2. Build and start the services:
```bash
docker-compose up --build
```
3. Access the Swagger UI to test the API at:
    - [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## 🔒 API Endpoints
- **Auth Routes** (`/api/v1/auth`): Register (`/register`), Login (`/login`), Refresh Token (`/refresh-token`)
- **User Routes** (`/api/v1/user`): Get user by email (`/:email`), Get/update/delete authenticated user (`/me`)

All `/user` routes require JWT authentication. Use `/login` or `/refresh-token` to obtain tokens.