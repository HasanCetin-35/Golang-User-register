package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI     = "mongodb://localhost:27017"
	databaseName = "Exercise"
)

func DBinstance() *mongo.Client {
	fmt.Printf("Connecting to MongoDB at %s...\n", mongoURI)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v\n", err)
	}

	// Ping MongoDB to check if the connection is successful
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v\n", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database(databaseName).Collection(collectionName)
}
