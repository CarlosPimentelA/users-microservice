package repository

import (
	"context"
	"users-microservice/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	FindRefreshTokenByHash(ctx context.Context, hashToken string) (*models.RefreshToken, error)
	FindRefreshTokenByID(ctx context.Context, tokenID string) (*models.RefreshToken, error)
	RevokeToken(ctx context.Context, tokenId string) error
	RevokeAllTokenFromUser(ctx context.Context, userId string) error
}

type mongoRefreshTokenRepository struct {
	collection *mongo.Collection
}

// FindRefreshTokenByID implements RefreshTokenRepository.
func (m *mongoRefreshTokenRepository) FindRefreshTokenByID(ctx context.Context, jti string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	filter := bson.M{"jti": jti}
	err := m.collection.FindOne(ctx, filter).Decode(&refreshToken)
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// RevokeAllUserToken implements RefreshTokenRepository.
func (m *mongoRefreshTokenRepository) RevokeAllTokenFromUser(ctx context.Context, userId string) error {
	filter := bson.M{"user_id": userId}
	update := bson.M{
		"$set": bson.M{
			"revoked": true,
		},
	}
	_, err := m.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrUserNotFound
		}
		return err
	}
	return nil
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
	var refreshToken models.RefreshToken
	filter := bson.M{"token": hashToken}
	err := m.collection.FindOne(ctx, filter).Decode(&refreshToken)
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// RevokeToken implements RefreshTokenRepository.
func (m *mongoRefreshTokenRepository) RevokeToken(ctx context.Context, tokenId string) error {
	var refreshToken models.RefreshToken
	filter := bson.M{"jti": tokenId}
	update := bson.M{
		"$set": bson.M{
			"revoked": true,
		},
	}
	config := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := m.collection.FindOneAndUpdate(ctx, filter, update, config).Decode(&refreshToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
