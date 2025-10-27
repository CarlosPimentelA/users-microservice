package repository

import (
	"context"
	"errors"
	"users-microservice/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUser(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, email string, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, email string) error
	UpdateField(ctx context.Context, email string, newValue interface{}, field string) (*models.User, error)
}

type mongoUserRepository struct {
	collection *mongo.Collection
}

// UpdateField implements UserRepository.

var ErrUserNotFound = errors.New("user not found")

func NewMongoUserRepository(client *mongo.Client, dbName string, collectionName string) UserRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &mongoUserRepository{
		collection: collection,
	}
}

// CreateUser implements UserRepository.
func (repo *mongoUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := repo.collection.InsertOne(ctx, user)
	return err
}

// DeleteUser implements UserRepository.
func (repo *mongoUserRepository) DeleteUser(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	_, err := repo.collection.DeleteOne(ctx, filter)
	return err
}

// FindUser implements UserRepository.
func (repo *mongoUserRepository) FindUser(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := repo.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser implements UserRepository.
func (repo *mongoUserRepository) UpdateUser(ctx context.Context, email string, user *models.User) (*models.User, error) {
	var userUpdated models.User
	filter := bson.M{"email": email}
	replacement := models.User{
		Name:         user.Name,
		Email:        user.Email,
		LastName:     user.LastName,
		PasswordHash: user.PasswordHash,
		UserId:       user.UserId}
		config := options.FindOneAndReplace().SetReturnDocument(options.After)
	mongoErr := repo.collection.FindOneAndReplace(ctx, filter, replacement, config).Decode(&userUpdated)
	if mongoErr != nil {
		return  nil, mongoErr
	}
	return &userUpdated, nil
}

func (repo *mongoUserRepository) UpdateField( ctx context.Context, email string, newValue interface{}, field string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			field: newValue,
		},
	}
	config := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := repo.collection.FindOneAndUpdate(ctx, filter, update, config).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}