package handler

import (
	"encoding/json"
	"net/http"
	"weather-app/internal/service"
)

type WeatherHandler struct {
	service *service.WeatherService
}

func NewWeatherHandler(service *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{service: service}
}

func (h *WeatherHandler) HandleWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetWeather(w, r, city)
	case http.MethodPut:
		h.handlePutWeather(w, r, city)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *WeatherHandler) handleGetWeather(w http.ResponseWriter, r *http.Request, city string) {
	ctx := r.Context()
	weather, err := h.service.GetWeather(ctx, city)
	if err != nil {
		http.Error(w, "Weather data not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func (h *WeatherHandler) handlePutWeather(w http.ResponseWriter, r *http.Request, city string) {
	ctx := r.Context()
	err := h.service.UpdateWeather(ctx, city)
	if err != nil {
		http.Error(w, "Failed to update weather data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Weather data updated successfully"))
}
