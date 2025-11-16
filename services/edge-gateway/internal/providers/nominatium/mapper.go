package nominatim

import (
	"strconv"

	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/domain/geo"
)

func mapToDomain(r *searchResult) (*geo.GeoPoint, error) {
    lat, err := strconv.ParseFloat(r.Lat, 64)
    if err != nil {
        return nil, err
    }

    lon, err := strconv.ParseFloat(r.Lon, 64)
    if err != nil {
        return nil, err
    }

    return &geo.GeoPoint{
        Name:        r.DisplayName,
        City:        r.Address.City,
        State:       r.Address.State,
        Country:     r.Address.Country,
        CountryCode: r.Address.CountryCode,
        Coordinates: geo.Coordinates{
            Lat:  lat,
            Lon: lon,
        },
    }, nil
}