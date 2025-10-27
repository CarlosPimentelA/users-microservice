Project Setup Guide

To run the application without making code changes, your MongoDB instance must be running and accessible using the URI provided in the DB_CONNECTION variable.
This guide provides the necessary steps to set up and run this project locally. This application is written in Go and uses MongoDB as its database.

Prerequisites
To get started, you must have the following installed on your system:

Go (1.18 or newer preferred)

MongoDB (either a running local instance or a connection URI for a remote service like MongoDB Atlas)

1. Environment Variables (.env file)
The project requires specific configuration settings for database connection and security, which must be stored in an environment file.

You need to create a new file named .env in the root directory of the project and populate it with the following required fields:

Variable	Description
DB_CONNECTION	The full MongoDB connection URI. This is often a string starting with mongodb:// or mongodb+srv://.
DB_NAME	The name of the database to be used within your MongoDB instance (e.g., myAppDB).
DB_COLLECTION_USERS	The collection name used to store user records.
DB_COLLECTION_REFRESH_TOKENS	The collection name used for storing session refresh tokens for authentication.
JWT_SECRET_KEY	A long, random, and secure string used to sign and verify JSON Web Tokens (JWTs). This is critical for security.


3. Installation and Running
Once your .env file is configured and MongoDB is running, you can launch the application:

A. Download Dependencies
Navigate to the project's root directory in your terminal and download all required Go modules:

Bash

go mod tidy
B. Run the Application
You can run the application directly using the Go toolchain:

Bash:
go run .

go build
./<executable_name> # (e.g., ./my-go-app)
The application should now be running and connected to your MongoDB database.
