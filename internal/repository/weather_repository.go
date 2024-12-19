package repository

import (
	"context"
	"errors"
	"fmt"

	"weather-app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// интерфейс
type WeatherRepositoryInterface interface {
	GetWeather(ctx context.Context, city string) (*models.WeatherData, error)
	UpdateWeather(ctx context.Context, weather *models.WeatherData) error
}

type WeatherRepository struct {
	collection *mongo.Collection
}

func NewWeatherRepository(client *mongo.Client, dbName string) *WeatherRepository {
	return &WeatherRepository{
		collection: client.Database(dbName).Collection("weather"),
	}
}

func (r *WeatherRepository) GetWeather(ctx context.Context, city string) (*models.WeatherData, error) {
	filter := bson.M{"city": city}
	var weather models.WeatherData
	err := r.collection.FindOne(ctx, filter).Decode(&weather)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("weather data for city %s not found", city)
		}
		return nil, err
	}
	return &weather, err
}

func (r *WeatherRepository) UpdateWeather(ctx context.Context, weather *models.WeatherData) error {
	filter := bson.M{"city": weather.City}
	update := bson.M{"$set": weather}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
