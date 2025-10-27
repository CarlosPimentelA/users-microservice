package service

import (
	"context"
	"fmt"
	"time"
	"users-microservice/config"
	"users-microservice/dto"
	"users-microservice/models"
	"users-microservice/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)
type RefreshTokenService struct {
	RefreshTokenRepository repository.RefreshTokenRepository
	config config.Config
}

func NewRefreshTokenService(RefreshTokenRepo repository.RefreshTokenRepository, config config.Config) *RefreshTokenService {
	return &RefreshTokenService{
		RefreshTokenRepository: RefreshTokenRepo,
		config: config,
	}
}

func (service *RefreshTokenService) CreateRefreshTokenService(ctx context.Context, refreshToken *dto.RefreshTokenCreateDTO) (string, error) {
	timeNow := time.Now()
	Id, errId := uuid.NewRandom()
	if errId != nil{
		return "", fmt.Errorf("error: error generating token_id: %w",errId)
	}
	tokenHash, hashErr := bcrypt.GenerateFromPassword([]byte(refreshToken.TokenPlain), 12)
	if hashErr != nil {
		return "error", fmt.Errorf("error hashing the token: %s", hashErr)
	}
	refreshTokenModel := models.RefreshToken{
		ID: Id.String(),
		UserId: refreshToken.UserId,
		TokenHash: string(tokenHash),
		IssuedAt: timeNow,
		Expires: timeNow.Add(service.config.REFRESH_TOKEN_CONFIG.EXPIRY_TIME),
		Revoked: false,
	}
	err := service.RefreshTokenRepository.CreateRefreshToken(ctx, &refreshTokenModel)
	if err != nil {
		return "error", fmt.Errorf("internal server error: %s", err)
	}
	return refreshToken.TokenPlain, nil
}