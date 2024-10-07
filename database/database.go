package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB(uri string, databaseName string, collectionName string) (*mongo.Client, *mongo.Collection, error) {
	fmt.Println("Connecting to database..")
	clientOptions := options.Client().ApplyURI(uri) // Update the URI as needed
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	collection := client.Database(databaseName).Collection(collectionName) // Change database and collection names as needed

	return client, collection, nil
}
