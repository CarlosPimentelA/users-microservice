package main

import (
	"users-microservice/config"
	"users-microservice/db"
	"users-microservice/handlers"
	"users-microservice/repository"
	"users-microservice/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	router := gin.Default()
	config, err := config.LoadConfig()
	client, mongoErr := db.Db_connection()
	userRepo := repository.NewMongoUserRepository(client, config.DB_NAME, config.DB_COLLECTION_USERS)
	refreshTokenRepo := repository.NewRefreshTokenRepository(client, config.DB_NAME, config.DB_COLLECTION_REFRESH_TOKENS)
	refreshTokenService := service.NewRefreshTokenService(refreshTokenRepo, *config)
	userService := service.NewUserService(userRepo, *config)
	validate := validator.New()
	userHandler := handlers.NewUserHandler(userService, validate)
	refreshTokenHandler := handlers.NewRefreshTokenHandler(refreshTokenService, validate)
	handlers.SetupRoutes(router, userHandler, refreshTokenHandler)
	if err != nil {

	}
	if mongoErr != nil {

	}
	router.Run()
}