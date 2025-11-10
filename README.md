# Users Microservice

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![MongoDB](https://img.shields.io/badge/MongoDB-Database-47A248?logo=mongodb)
![Gin](https://img.shields.io/badge/Gin-Web_Framework-00ADD8)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

A simple and secure **user authentication microservice** built in **Go** using **Gin** and **MongoDB**.  
It implements user registration, login, and secure JWT-based authentication with **refresh token rotation**.

---

## Features

- 🔐 **User registration and login**
- 🪪 **JWT-based authentication**
- 🔄 **Refresh token rotation**
- 🧂 **Password hashing with bcrypt**
- 💾 **MongoDB persistence**
- 🧱 **Clean layered architecture (repository / service / handler)**
- ⚙️ **Configuration via `.env` file**
- 🧪 **Testable via Postman**

---

## Architecture Overview

```
users-microservice/
│
├── config/                  # Environment and app configuration
├── db/                      # MongoDB connection setup
├── dto/                     # Data Transfer Objects (DTOs)
├── handlers/                # Gin HTTP handlers
├── models/                  # Database models (User, RefreshToken)
├── repository/              # MongoDB repositories
├── service/                 # Business logic (UserService, RefreshTokenService)
├── main.go                  # Entry point
└── .env                     # Environment variables
```

---

## Environment Variables

The `.env` file contains runtime configuration:

```env
MONGO_URI=mongodb://localhost:27017
DB_NAME=users_db
DB_COLLECTION_USERS=users
DB_COLLECTION_REFRESH_TOKENS=refresh_tokens
JWT_SECRET=supersecretkey
```

---

## API Endpoints

| Method | Endpoint         | Description               | Body Example |
|:-------|:-----------------|:--------------------------|:--------------|
| `POST` | `/users/register` | Register a new user        | `{ "name": "John", "last_name": "Doe", "email": "john@doe.com", "password": "123456" }` |
| `POST` | `/users/login`    | Authenticate user & get JWT | `{ "email": "john@doe.com", "password": "12345678" }` |
| `POST` | `/token/refresh`  | Request a new access token  | `{ "jwt": "<token>" }` |

> The refresh token is **managed server-side**, not stored on the client.

---

## Authentication Flow

1. **User registers** → data is hashed and stored.
2. **User logs in** → a JWT + refresh token is generated.
3. **Access token** expires → client requests `/token/refresh`.
4. **Refresh token** rotation occurs; old tokens are revoked.

---

## Technologies

- [Go](https://go.dev/)
- [Gin](https://github.com/gin-gonic/gin)
- [MongoDB](https://www.mongodb.com/)
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [uuid](https://pkg.go.dev/github.com/google/uuid)
- [jwt](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)
- [validator](https://github.com/go-playground/validator)

---

## Running the Project

### 1. Clone the repo
```bash
git clone https://github.com/yourusername/users-microservice.git
cd users-microservice
```

### 2. Configure environment
```bash
cp .env.example .env
```

### 3. Run MongoDB
Make sure MongoDB is running locally or remotely.

### 4. Start the server
```bash
go run .
```

Server runs by default on:  
 `http://localhost:8080`

---

## Project Structure Example

```
service/
├── users_service.go
├── refresh_tokens_service.go
repository/
├── user_repository.go
├── refresh_token_repository.go
handlers/
├── users_handler.go
├── refresh_token_handler.go
```

---

## Example Request

### Login (Postman)

**POST** `http://localhost:8080/users/login`
```json
{
  "email": "john@doe.com",
  "password": "123456"
}
```

**Response**
```json
{
  "user_id": "c2a2d460-7e1a-4b4f-a9ef-1a41b72fa1a9",
  "name": "John",
  "email": "john@doe.com",
  "jwt": "<access_token>"
}
```

---

## License

This project is licensed under the [MIT License](LICENSE).

```
MIT License

Copyright (c) 2025

Permission is hereby granted, free of charge, to any person obtaining a copy.
```

---

## Future Improvements

- Add **logout endpoint** (invalidate refresh tokens)
- Add **unit tests**
- Add **API test**
- Enhance error management

---

Developed by **Carlos Pimentel**  
