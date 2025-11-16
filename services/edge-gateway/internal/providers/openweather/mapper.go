package openweather

import (
	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/domain/weather"
)

func mapToDomain(ow *OneCallResponse, coords weather.Coordinates) (*weather.Weather) {
    var condition string
    if len(ow.Current.Weather) > 0 {
        condition = ow.Current.Weather[0].Description
    }

    return &weather.Weather{
        Coordinates: coords,
        Temp: weather.Temperature{
            Actual:    ow.Current.Temp,
            FeelsLike: ow.Current.FeelsLike,
        },
        Wind: weather.Wind{
            Speed: ow.Current.WindSpeed,
            Deg:   ow.Current.WindDeg,
        },
        Humidity:  ow.Current.Humidity,
        Condition: condition,
	}
}