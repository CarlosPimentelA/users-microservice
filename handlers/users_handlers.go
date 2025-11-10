package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"users-microservice/dto"
	"users-microservice/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	Service   *service.UserService
	Validator *validator.Validate
}

func NewUserHandler(userService *service.UserService, validator *validator.Validate, refreshTokenService *service.RefreshTokenService) *UserHandler {
	return &UserHandler{
		Service:   userService,
		Validator: validator,
	}
}

func SetupRoutes(g *gin.Engine, userHandler *UserHandler, refreshTokenHandler *RefreshTokenHandler) {
	userPath := g.Group("/users")
	{
		userPath.POST("", userHandler.HandleCreateUser)
		userPath.POST("/login", userHandler.HandleLoginUser)
	}
	refreshTokenPath := g.Group("refresh")
	{
		refreshTokenPath.POST("", refreshTokenHandler.HandleRefreshToken)
	}
}

func (handler *UserHandler) HandleCreateUser(gc *gin.Context) {
	ctx, cancel := context.WithTimeout(gc.Request.Context(), 10*time.Second)
	defer cancel()
	newUser := new(dto.UserDTO)
	err := gc.BindJSON(newUser)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}
	validationErr := handler.Validator.Struct(newUser)
	if validationErr != nil {
		if validateError, ok := validationErr.(validator.ValidationErrors); ok {
			fmt.Println("Errors: ", validateError)
			var errorMessage []string
			for _, er := range validateError {
				message := translateFieldErr(er)
				errorMessage = append(errorMessage, message)
			}
			gc.JSON(http.StatusBadRequest, gin.H{
				"status":        http.StatusBadRequest,
				"error":         "VALIDATION_FAILED",
				"error_details": errorMessage,
			})
			return
		}
		return
	}
	userDTO, serviceErr := handler.Service.CreateUserService(ctx, newUser)
	if serviceErr != nil {
		statusCode, errorMessage := MapErrorToHttp(serviceErr)
		gc.JSON(statusCode, gin.H{
			"status":  statusCode,
			"error":   http.StatusText(statusCode),
			"message": errorMessage,
		})
		return
	}
	gc.JSON(http.StatusCreated, userDTO)
}

func (handler *UserHandler) HandleLoginUser(gc *gin.Context) {
	ctx, cancel := context.WithTimeout(gc.Request.Context(), 10*time.Second)
	defer cancel()
	newLogin := new(dto.LoginRequestDTO)
	err := gc.BindJSON(newLogin)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}
	validationErr := handler.Validator.Struct(newLogin)
	if validationErr != nil {
		if validateError, ok := validationErr.(validator.ValidationErrors); ok {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status":        http.StatusBadRequest,
				"error":         "INVALID_CREDENTIALS",
				"error_details": validateError,
			})
			return
		}
		return
	}
	jwt, authErr := handler.Service.AuthenticationService(ctx, newLogin.Email, newLogin.Password)
	if authErr != nil {
		gc.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   http.StatusText(http.StatusBadRequest),
			"message": authErr,
		})
		return
	}
	gc.JSON(http.StatusCreated, jwt)
}

func MapErrorToHttp(err error) (int, string) {
	if errors.Is(err, service.ErrUserNotFound) {
		return http.StatusNotFound, "The requested resource was not found."
	}
	if errors.Is(err, service.ErrUpdateFailed) {
		return http.StatusConflict, "The update failed due to a current data conflict."
	}
	if errors.Is(err, service.ErrEmailConflict) {
		return http.StatusConflict, "This email is already registered"
	}
	if errors.Is(err, service.ErrInvalidCredencials) {
		return http.StatusBadRequest, "Invalid credentials"
	}
	return http.StatusInternalServerError, "An unexpected error occurred on the server."
}

func translateFieldErr(fe validator.FieldError) string {
	fieldName := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("El campo %s es obligatorio.", fieldName)
	case "email":
		return fmt.Sprintf("El campo %s debe ser una dirección de email válida.", fieldName)
	case "min":
		return fmt.Sprintf("El campo %s debe tener un mínimo de %s caracteres.", fieldName, fe.Param())
	case "max":
		return fmt.Sprintf("El campo %s no debe exceder los %s caracteres.", fieldName, fe.Param())
	default:
		return fmt.Sprintf("El campo %s falló la validación.", fieldName)
	}
}
