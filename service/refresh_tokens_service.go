package service

import (
	"context"
	"fmt"
	"log"
	"time"
	"users-microservice/config"
	"users-microservice/dto"
	"users-microservice/models"
	"users-microservice/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshTokenService struct {
	RefreshTokenRepository repository.RefreshTokenRepository
	UserService            *UserService
	config                 config.Config
}

func NewRefreshTokenService(RefreshTokenRepo repository.RefreshTokenRepository, config *config.Config) *RefreshTokenService {
	return &RefreshTokenService{
		RefreshTokenRepository: RefreshTokenRepo,
		config:                 *config,
	}
}

func (service *RefreshTokenService) CreateRefreshTokenService(ctx context.Context, refreshToken *dto.RefreshTokenCreateDTO) error {
	timeNow := time.Now()
	Id, errId := uuid.NewRandom()
	if errId != nil {
		return fmt.Errorf("error: error generating token_id: %w", errId)
	}
	user, findUserErr := service.UserService.FindUserByIDService(ctx, refreshToken.UserId)
	refreshTokenModel := models.RefreshToken{
		ID:             Id.String(),
		UserId:         refreshToken.UserId,
		Jti: 			refreshToken.Jti,
		IssuedAt:       timeNow,
		Expires:        refreshToken.ExpiresAt,
		Revoked:        false,
		SessionVersion: user.SessionVersion,
	}
	if findUserErr != nil {
		return ErrUserNotFound
	}

	err := service.RefreshTokenRepository.CreateRefreshToken(ctx, &refreshTokenModel)
	if err != nil {
		return fmt.Errorf("internal server error: %s", err)
	}
	return nil
}

func (service *RefreshTokenService) RefreshAccessToken(ctx context.Context, jwtString string) (string, error) {
	// 1. Obtener usuario
	// 2. Obtener el id del usuario
	// 3. Obtener los claims del jwt
	claims, verifyClaim := extractClaims(jwtString, service.config)
	if !verifyClaim {
		return uuid.Nil.String(), fmt.Errorf("token validation failed")
	}
	// 1. Sacar el user ID
	userId, claimErr := claims.GetSubject()
	if claimErr != nil {
		return uuid.Nil.String(), fmt.Errorf("claims err")
	}
	// Obtener el jti
	jti, ok := claims["jti"].(string)
	if !ok {
	return uuid.Nil.String(), fmt.Errorf("missing jti claim")
	}
	// Obtener el usuario
	user, userFindErr := service.UserService.userService.FindUserByID(ctx, userId)
	if userFindErr != nil {
		return uuid.Nil.String(), fmt.Errorf("find user err")
	}
	//	Buscar el refresh token que este vinculado al user id
	refreshTokenUser, findTokenErr := service.RefreshTokenRepository.FindRefreshTokenByID(ctx, jti)
	if findTokenErr != nil {
		return uuid.Nil.String(), fmt.Errorf("token validation failed")
	}
	// Obtener la fecha de expiracion
	tokenExpiryTime, tokenExpiryTimeErr := claims.GetExpirationTime()
	if tokenExpiryTimeErr != nil {
		return uuid.Nil.String(), fmt.Errorf("token with no due date")
	}
	// Verificar que el token no este expirado
	if refreshTokenUser.SessionVersion < user.SessionVersion {
		return uuid.Nil.String(), fmt.Errorf("invalid token")
	}

	if refreshTokenUser.Revoked {
	revokeErr := service.RefreshTokenRepository.RevokeAllTokenFromUser(ctx, refreshTokenUser.UserId)
	if revokeErr != nil {
		return uuid.Nil.String(), fmt.Errorf("error revoking all tokens: %w", revokeErr)
	}
	return uuid.Nil.String(), fmt.Errorf("token reuse detected — all tokens revoked")
}

	if time.Now().After(tokenExpiryTime.Time) {
		return uuid.Nil.String(), fmt.Errorf("token expired")
	} else {
		// Crear nuevo JWT
		newJwt, createJwtErr := createJwtToken(user, service)
		if createJwtErr != nil {
			return uuid.Nil.String(), fmt.Errorf("new token err")
		}
		refreshTokenPlain, uuidErr := uuid.NewRandom()
		if uuidErr != nil {
		}
		// Crear el DTO del nuevo refresh token
		newRefreshToken := dto.RefreshTokenCreateDTO{
			UserId:     userId,
			Jti: refreshTokenPlain.String(),
			ExpiresAt:  time.Now().Add(service.config.REFRESH_TOKEN_CONFIG.EXPIRY_TIME),
		}
		// Crear el nuevo refresh token
		service.CreateRefreshTokenService(ctx, &newRefreshToken)
		// Revoke the old refresh token
		revokeTokenErr := service.RefreshTokenRepository.RevokeToken(ctx, refreshTokenUser.Jti)
		if revokeTokenErr != nil {
			return uuid.Nil.String(), fmt.Errorf("revoke token failed")
		}
		return newJwt, nil
	}
	//	~ Hacer un script de limpieza de todos los token revokados (Se hace en el main. Talvez haga una funcion en el
	// 	  repo de refresh token para mantener el orden en mi codigo)
}

func extractClaims(tokenStr string, config config.Config) (jwt.MapClaims, bool) {
	hmacSecretString := config.JWT_SECRET_KEY
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})
	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func createJwtToken(user *models.User, service *RefreshTokenService) (string, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return	"", nil
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"name":  user.Name,
			"email": user.Email,
			"sub":   user.UserId,
			"jti":   tokenId,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
			"iss":   "users-microservice",
			"iat":   time.Now().Unix(),
			"aud":   "contacts-service",
		})

	jwtToken, signErr := token.SignedString([]byte(service.config.JWT_SECRET_KEY))
	if signErr != nil {
		return "Error al firmar el token", signErr
	}
	return jwtToken, nil
}

// TODO: Dividir por funciones el metodo de refresh token, minimo 3
