package repository

import (
	"context"
	"users-microservice/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	FindRefreshTokenByHash(ctx context.Context, hashToken string) (*models.RefreshToken, error)
	RevokeToken(ctx context.Context, hashToken string) error
}

type mongoRefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository(client *mongo.Client, dbName string, collectionName string) RefreshTokenRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &mongoRefreshTokenRepository{
		collection: collection,
	}
}

// CreateRefreshToken implements RefreshTokenRepository.
func (m *mongoRefreshTokenRepository) CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	_, err := m.collection.InsertOne(ctx, refreshToken)
	return err
}

// FindRefreshTokenByHash implements RefreshTokenRepository.
func (m *mongoRefreshTokenRepository) FindRefreshTokenByHash(ctx context.Context, hashToken string) (*models.RefreshToken, error) {
	panic("unimplemented")
}

// RevokeToken implements RefreshTokenRepository.
func (m *mongoRefreshTokenRepository) RevokeToken(ctx context.Context, hashToken string) error {
	panic("unimplemented")
}

