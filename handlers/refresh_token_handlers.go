package handlers

import (
	"users-microservice/service"

	"github.com/go-playground/validator/v10"
)

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