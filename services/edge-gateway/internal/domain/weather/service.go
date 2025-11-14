package weather

import (
	"context"
)

type weatherService struct {
	provider WeatherProvider
}

func NewWeatherService(provider WeatherProvider) WeatherService {
    return &weatherService{provider: provider}
}


func (s *weatherService) GetWeatherByCity(ctx context.Context, city string) (*Weather, error) {
    return s.provider.FetchWeatherByCity(ctx, city)
}

func (s *weatherService) GetWeatherByCoordinates(ctx context.Context, coords Coordinates) (*Weather, error) {
	if err := coords.Validate(); err != nil {
		return nil, err
	}
    return s.provider.FetchWeatherByCoordinates(ctx, coords)
}



