package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var LocalMongo *mongo.Client
var TestCollection *mongo.Collection

func ConnectMongo() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	LocalMongo = client
	TestCollection = client.Database("localdb").Collection("test_collection")
	log.Println("Connected to MongoDB")
}

func InsertToMongo(data map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := TestCollection.InsertOne(ctx, data)
	return err
}

func ReadFromMongo(id string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result map[string]interface{}
	err := TestCollection.FindOne(ctx, map[string]interface{}{"_id": id}).Decode(&result)
	return result, err
}
