package geo

import (
	"context"
)

type geoService struct {
    provider GeoProvider
}

func (s *geoService) ResolveLocation(ctx context.Context, query string) (*GeoPoint, error) {
    return s.provider.Search(ctx, query)
}

func (s *geoService) ResolveCoordinates(ctx context.Context, coords Coordinates) (*GeoPoint, error) {
    return s.provider.Reverse(ctx, coords)
}
