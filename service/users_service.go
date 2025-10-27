package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"users-microservice/config"
	"users-microservice/dto"
	"users-microservice/models"
	"users-microservice/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailConflict = errors.New("email is already registered")
var ErrUserNotFound = errors.New("user not found")
var ErrUpdateFailed = errors.New("update failed")
var ErrInvalidCredencials = errors.New("invalid credencials")
var ErrInternalServer = errors.New("internal server error")
type UserService struct {
	userService repository.UserRepository
	config config.Config
}

type FieldUpdateFunc func(context.Context, string, string) (*models.User, error)

func NewUserService(userRepo repository.UserRepository, config config.Config) *UserService {
	return &UserService{
		userService: userRepo,
		config: config,
	}
}

func (service *UserService) CreateUserService(ctx context.Context, userDTO *dto.UserDTO) (*dto.UserDTO, error) {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("error: hashing the password: %w", err)		
	}
	userId, err := uuid.NewRandom()
	if err != nil{
		return nil, fmt.Errorf("error: error generating user_id: %w",err)
	}
	
	user := models.User{
		Name: userDTO.Name,
		LastName: userDTO.LastName,
		Email: userDTO.Email,
		UserId: userId.String(),
		PasswordHash: string(passwordHashed),
	}

	duplicateEmail, findErr := service.userService.FindUser(ctx, userDTO.Email)
	if findErr != nil && !errors.Is(findErr, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("error: email duplicate validation: %w", findErr)
	}
	if duplicateEmail != nil {
		return nil, ErrEmailConflict
	}
	repoError := service.userService.CreateUser(ctx, &user)
	if repoError != nil {
		return nil, fmt.Errorf("error: register failed in db: %w", repoError)
	}
	return mapModelToDTO(&user), nil
}

func (service *UserService) FindUserService(ctx context.Context, email string) (*dto.UserDTO, error) {
		user, err := service.userService.FindUser(ctx, email)
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		if err != nil {
			return nil, fmt.Errorf("error: db error: %w", err)
		}
		
		
		return mapModelToDTO(user), nil
}

func (service *UserService) UpdateUserService(ctx context.Context, email string, user *models.User) (*dto.UserDTO, error) {
	user, err := service.userService.UpdateUser(ctx, email, user)
	if err != nil {
		return nil, fmt.Errorf("error: update error: %w", err)
	}
	
	return mapModelToDTO(user), nil
}

func (service *UserService) DeleteUserService(ctx context.Context, email string) error {
	err := service.userService.DeleteUser(ctx, email)
	if err != nil {
		return fmt.Errorf("error: error deleting the user with this email: %s: %w", email, err)
	}
		return nil
}

func (service *UserService) UpdateFieldService(ctx context.Context, email string, newValue string, fieldFunc FieldUpdateFunc, errMessage string) (*dto.UserDTO, error) {
	userModified, err := fieldFunc(ctx, email, newValue)
	if err != nil {
		return nil, fmt.Errorf("error: error modifing the %s: %w", errMessage, err)
	}
		return mapModelToDTO(userModified), nil
}

func (service *UserService) AuthenticationService(ctx context.Context, email string, password string) (*dto.AuthResponse, error) {
	user, err := service.userService.FindUser(ctx, email)
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		if !errors.Is(err, ErrUserNotFound) {
			return nil, ErrInternalServer
		}
	credencialErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if credencialErr != nil {
		return nil, ErrInvalidCredencials
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "name": user.Name, 
		"email": user.Email,
		"sub": user.UserId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iss": "users-microservice",
		"iat": time.Now().Unix(),
		"aud": "contacts-service",
	})

	tokenString, signErr := token.SignedString(service.config.JWT_SECRET_KEY)
	if signErr != nil {
		return nil, signErr
	}
	refreshToken, uuidErr := uuid.NewRandom()
	/*
	1. Guardar el refresh token en la db con el servicio de los refresh token.
	2. Hacer que los usuarios puedan tener varios dispositivos logeados de manera independiente. 
		- Ejemplo: Si hago log out en mi telefono que mi cuenta en la computadora no haga log out tambien.
	3. Mejorar los mensajes de error de la parte de los refresh token(repository, service, handlers [cuando haya]) y revisar los mensajes
	de error que dan los usuarios. 
	4. Tener los mensajes de error centralizados en la config. 
		~ Tener los mensajes en un campo que sea error, y que hay esten todas las posibles repuestas de error.
			- Tener todos los errores divididos en los errores de usuario y los de token.
			- Tener errores generales para ambos.
		~ Utilizar esta config en toda la app.
	*/
	if uuidErr != nil {
		return nil, uuidErr
	}
	response := dto.AuthResponse{
		UserId: user.UserId,
		Name: user.Name,
		Email: email,
		JWT: tokenString,
	}
	return &response, nil
}

func mapModelToDTO(model *models.User) *dto.UserDTO {
	var userDTO = dto.UserDTO{
		Name: model.Name,
		LastName: model.LastName,
		Email: model.Email,
		Password: model.PasswordHash,
	}
	return &userDTO
}