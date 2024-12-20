package main

import (
	"context"
	"log"
	"net/http"

	"weather-app/config"
	"weather-app/internal/handler"
	"weather-app/internal/repository"
	"weather-app/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.NewConfig()

	dbClient, err := config.ConnectMongo(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer dbClient.Disconnect(context.TODO())

	weatherRepo := repository.NewWeatherRepository(dbClient, cfg.DatabaseName)
	weatherService := service.NewWeatherService(weatherRepo, cfg.BaseURL, cfg.APIKey)
	weatherHandler := handler.NewWeatherHandler(weatherService)

	http.HandleFunc("/weather", weatherHandler.HandleWeather)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
