package config

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	BaseURL      string
	APIKey       string
	MongoURI     string
	DatabaseName string
}

func NewConfig() *Config {
	return &Config{
		BaseURL:      getEnv("BASE_URL", ""),
		APIKey:       getEnv("API_KEY", ""),
		MongoURI:     getEnv("MONGO_URI", ""),
		DatabaseName: "weatherdb",
	}
}

func ConnectMongo(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongo.Connect(ctx, options.Client().ApplyURI(uri))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
