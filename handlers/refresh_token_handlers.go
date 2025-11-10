package handlers

import (
	"net/http"
	"strings"
	"users-microservice/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
type RefreshTokenRequest struct {
	Token string `validate:"required,jwt"`
}

type RefreshTokenHandler struct {
	Service   *service.RefreshTokenService
	Validator *validator.Validate
}

func NewRefreshTokenHandler(service *service.RefreshTokenService, validator *validator.Validate) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		Service:   service,
		Validator: validator,
	}
}

func (handler *RefreshTokenHandler) HandleRefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization header required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	request := RefreshTokenRequest{Token: tokenString}

	// Validar con el validator
	if err := handler.Validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing token"})
		return
	}

	newAccessToken, err := handler.Service.RefreshAccessToken(ctx, tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
	})
}

// 1. Hacer el endpoint para refrescar automaticamente el jwt del usuario. Handler ‼️
