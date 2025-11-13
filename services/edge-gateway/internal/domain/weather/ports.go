package weather

import (
	"context"
)

type WeatherService interface {
    GetWeatherByCoordinates(ctx context.Context, coords Coordinates) (*Weather, error)
    GetWeatherByCity(ctx context.Context, city string) (*Weather, error)
}

type WeatherProvider interface {
    FetchWeatherByCoordinates(ctx context.Context, coords Coordinates) (*Weather, error)
    FetchWeatherByCity(ctx context.Context, city string) (*Weather, error) 
}