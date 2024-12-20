package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"weather-app/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWeatherRepository struct {
	mock.Mock
}

func (m *MockWeatherRepository) GetWeather(ctx context.Context, city string) (*models.WeatherData, error) {
	args := m.Called(ctx, city)
	if data, ok := args.Get(0).(*models.WeatherData); ok {
		return data, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWeatherRepository) UpdateWeather(ctx context.Context, weather *models.WeatherData) error {
	args := m.Called(ctx, weather)
	return args.Error(0)
}

func TestWeatherService_UpdateWeather(t *testing.T) {
	mockRepo := new(MockWeatherRepository)
	service := NewWeatherService(mockRepo, "http://mockapi.com", "mockapikey")

	ctx := context.Background()
	city := "TestCity"

	mockRepo.On("UpdateWeather", mock.Anything, mock.MatchedBy(func(data *models.WeatherData) bool {
		return data.City == "TestCity" && data.Description == "Cloudy" && data.Temp == 20
	})).Return(nil).Once()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"main": { "temp": 20 },
			"weather": [{ "description": "Cloudy" }],
			"name": "TestCity"
		}`)
	}))
	defer server.Close()

	service.baseURL = server.URL

	err := service.UpdateWeather(ctx, city)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
