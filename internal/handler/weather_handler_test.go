package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-app/internal/handler"
	"weather-app/internal/models"
	"weather-app/internal/service"
)

type mockWeatherService struct {
	GetWeatherFunc    func(ctx context.Context, city string) (*models.WeatherData, error)
	UpdateWeatherFunc func(ctx context.Context, city string) error
}

func (m *mockWeatherService) GetWeather(ctx context.Context, city string) (*models.WeatherData, error) {
	return m.GetWeatherFunc(ctx, city)
}

func (m *mockWeatherService) UpdateWeather(ctx context.Context, city string) error {
	return m.UpdateWeatherFunc(ctx, city)
}

func TestHandleWeather_GetWeather_Success(t *testing.T) {
	mockService := &mockWeatherService{
		GetWeatherFunc: func(ctx context.Context, city string) (*models.WeatherData, error) {
			return &models.WeatherData{
				City:        "TestCity",
				Temp:        25.0,
				Description: "Sunny",
			}, nil
		},
	}

	h := handler.NewWeatherHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/weather?city=TestCity", nil)
	rec := httptest.NewRecorder()

	h.HandleWeather(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var weather models.WeatherData
	if err := json.NewDecoder(rec.Body).Decode(&weather); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if weather.City != "TestCity" {
		t.Errorf("expected city %q, got %q", "TestCity", weather.City)
	}
}

func TestHandleWeather_GetWeather_CityNotFound(t *testing.T) {
	mockService := &mockWeatherService{
		GetWeatherFunc: func(ctx context.Context, city string) (*models.WeatherData, error) {
			return nil, service.ErrWeatherNotFound
		},
	}

	h := handler.NewWeatherHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/weather?city=UnknownCity", nil)
	rec := httptest.NewRecorder()

	h.HandleWeather(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandleWeather_PutWeather_Success(t *testing.T) {
	mockService := &mockWeatherService{
		UpdateWeatherFunc: func(ctx context.Context, city string) error {
			return nil
		},
	}

	h := handler.NewWeatherHandler(mockService)
	req := httptest.NewRequest(http.MethodPut, "/weather?city=TestCity", nil)
	rec := httptest.NewRecorder()

	h.HandleWeather(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expectedBody := "Weather data updated successfully"
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestHandleWeather_PutWeather_Failure(t *testing.T) {
	mockService := &mockWeatherService{
		UpdateWeatherFunc: func(ctx context.Context, city string) error {
			return service.ErrUpdateFailed
		},
	}

	h := handler.NewWeatherHandler(mockService)
	req := httptest.NewRequest(http.MethodPut, "/weather?city=TestCity", nil)
	rec := httptest.NewRecorder()

	h.HandleWeather(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
