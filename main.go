package main

import (
	"log"
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

	// 1. Configuración y conexión
	config, err := config.LoadConfig()
	client, mongoErr := db.Db_connection()
	if err != nil || mongoErr != nil {
		log.Fatalf("Error cargando configuración o conectando a la DB: %v %v", err, mongoErr)
	}

	// 2. Crear repositorios
	userRepo := repository.NewMongoUserRepository(client, config.DB_NAME, config.DB_COLLECTION_USERS)
	refreshTokenRepo := repository.NewRefreshTokenRepository(client, config.DB_NAME, config.DB_COLLECTION_REFRESH_TOKENS)

	// 3. Crear servicios con dependencias circulares
	var userService *service.UserService
	var refreshTokenService *service.RefreshTokenService

	// Inicializar el refreshTokenService sin UserService todavía
	refreshTokenService = service.NewRefreshTokenService(refreshTokenRepo, config)

	// Inicializar el userService con el refreshTokenService
	userService = service.NewUserService(userRepo, refreshTokenService, config)

	// Ahora conectar ambos correctamente
	refreshTokenService.UserService = userService

	// 4. Handlers y validación
	validate := validator.New()
	userHandler := handlers.NewUserHandler(userService, validate, refreshTokenService)
	refreshTokenHandler := handlers.NewRefreshTokenHandler(refreshTokenService, validate)

	// 5. Rutas
	handlers.SetupRoutes(router, userHandler, refreshTokenHandler)

	// 6. Ejecutar servidor
	router.Run()
}