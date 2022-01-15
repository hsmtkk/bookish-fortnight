package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CosmosDBConnString = "COSMOS_DB_CONN_STRING"

func main() {
	connStr := os.Getenv(CosmosDBConnString)
	if connStr == "" {
		log.Fatalf("%s environment variable must be defined", CosmosDBConnString)
	}

	// connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		log.Fatalf("failed to connect mongo; %s", err)
	}

	// defer disconnect
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("failed to disconnect mongo; %s", err)
		}
	}()

	if err := insertRecord(client); err != nil {
		log.Fatal(err)
	}
	if err := findRecord(client); err != nil {
		log.Fatal(err)
	}
}

func insertRecord(client *mongo.Client) error {
	collection := client.Database("testing").Collection("numbers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	if err != nil {
		return fmt.Errorf("failed to insert record; %w", err)
	}
	id := res.InsertedID
	fmt.Println(id)
	return nil
}

func findRecord(client *mongo.Client) error {
	var result struct {
		Value float64
	}
	collection := client.Database("testing").Collection("numbers")
	filter := bson.D{{"name", "pi"}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := collection.FindOne(ctx, filter).Decode(&result); err == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
	} else if err != nil {
		return fmt.Errorf("failed to find record; %w", err)
	}
	fmt.Println(result)
	return nil
}
