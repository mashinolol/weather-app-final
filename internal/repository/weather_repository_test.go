package repository_test

import (
	"context"
	"testing"
	"time"

	"weather-app/internal/models"
	"weather-app/internal/repository"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestWeatherRepository_GetWeather(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer func() {
		if mt.Client != nil {
			_ = mt.Client.Disconnect(context.Background())
		}
	}()

	mt.Run("success", func(mt *mtest.T) {
		now := time.Now()
		weatherData := models.WeatherData{
			City:        "TestCity",
			Description: "Sunny",
			Temp:        25.0,
			LastUpdated: now,
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "weather.weather", mtest.FirstBatch, bson.D{
			{Key: "city", Value: weatherData.City},
			{Key: "description", Value: weatherData.Description},
			{Key: "temp", Value: weatherData.Temp},
			{Key: "last_updated", Value: weatherData.LastUpdated},
		}))

		repo := repository.NewWeatherRepository(mt.Client, "weatherDB")
		result, err := repo.GetWeather(context.TODO(), "TestCity")

		assert.NoError(t, err)
		assert.Equal(t, weatherData.City, result.City)
		assert.Equal(t, weatherData.Description, result.Description)
		assert.Equal(t, weatherData.Temp, result.Temp)
		assert.WithinDuration(t, weatherData.LastUpdated, result.LastUpdated, time.Second)
	})

	mt.Run("not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "weather.weather", mtest.FirstBatch))

		repo := repository.NewWeatherRepository(mt.Client, "weatherdb")
		result, err := repo.GetWeather(context.TODO(), "NonExistentCity")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "weather data for city NonExistentCity not found")
		assert.Nil(t, result)
	})
}

func TestWeatherRepository_UpdateWeather(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer func() {
		if mt.Client != nil {
			_ = mt.Client.Disconnect(context.Background())
		}
	}()

	mt.Run("success", func(mt *mtest.T) {
		weatherData := models.WeatherData{
			City:        "TestCity",
			Description: "Cloudy",
			Temp:        28.0,
			LastUpdated: time.Now(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := repository.NewWeatherRepository(mt.Client, "weatherdb")
		err := repo.UpdateWeather(context.TODO(), &weatherData)

		assert.NoError(t, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		weatherData := models.WeatherData{
			City:        "TestCity",
			Description: "Rainy",
			Temp:        22.0,
			LastUpdated: time.Now(),
		}
		repo := repository.NewWeatherRepository(mt.Client, "weatherdb")
		err := repo.UpdateWeather(context.TODO(), &weatherData)

		assert.Error(t, err)
	})
}
