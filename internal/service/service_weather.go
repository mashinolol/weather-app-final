package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weather-app/internal/models"
	"weather-app/internal/repository"
)

type WeatherService struct {
	repo    repository.WeatherRepositoryInterface
	baseURL string
	apiKey  string
}

func NewWeatherService(repo repository.WeatherRepositoryInterface, baseURL, apiKey string) *WeatherService {
	return &WeatherService{
		repo:    repo,
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, city string) (*models.WeatherData, error) {
	return s.repo.GetWeather(ctx, city)
}

func (s *WeatherService) UpdateWeather(ctx context.Context, city string) error {
	url := fmt.Sprintf("%s?appid=%s&q=%s&units=metric", s.baseURL, s.apiKey, city)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch weather data")
	}

	var apiResponse struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
		Name string `json:"name"`
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return err
	}

	weather := &models.WeatherData{
		City:        apiResponse.Name,
		Description: apiResponse.Weather[0].Description,
		Temp:        apiResponse.Main.Temp,
		LastUpdated: time.Now(),
	}
	return s.repo.UpdateWeather(ctx, weather)
}
