package models

import "time"

type WeatherData struct {
	City        string    `bson:"city" json:"city"`
	Description string    `bson:"description" json:"description"`
	Temp        float64   `bson:"temp" json:"temp"`
	LastUpdated time.Time `bson:"last_updated" json:"last_updated"`
}
