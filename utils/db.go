package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(mongoDbURI string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDbURI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	// ping the database
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}

func GetCollection(db *mongo.Client, collectionName string) *mongo.Collection {
	// In our system mongo db uses "dev" as the namespace for all environment's collection so hard coding this value here
	collection := db.Database("dev").Collection(collectionName)
	return collection
}
