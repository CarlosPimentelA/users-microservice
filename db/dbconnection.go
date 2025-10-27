package db

import (
	"context"
	"fmt"
	"time"
	"users-microservice/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func Db_connection() (*mongo.Client, error) {
	uri, error := config.LoadConfig()
	if error != nil {
		return nil, fmt.Errorf("%s", error.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	optionClient := options.Client().ApplyURI(uri.DB_CONNECTION)
	defer cancel()
	client, err := mongo.Connect(optionClient)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	er := client.Ping(ctx, readpref.Primary())
	if er != nil {
		return nil, fmt.Errorf("%s", er.Error())
	}
	fmt.Println("Successfully pinged MongoDB!")
	defer cancel()
	return client, nil
}
