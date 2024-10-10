package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() error {
	DB_URL := os.Getenv("DB_URL")

	if DB_URL == "" {
		return fmt.Errorf("DB_URL variable is not set in .env")
	}

	clientOptions := options.Client().ApplyURI(DB_URL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to ping MongoDB: %v", err)
	}

	DB = client.Database("sai")

	fmt.Println("Connected to MongoDB successfully!")
	return nil
}
