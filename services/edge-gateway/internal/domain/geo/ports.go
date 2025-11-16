package geo

import (
	"context"
)

type GeoProvider interface {
    Search(ctx context.Context, query string) (*GeoPoint, error)
    Reverse(ctx context.Context, coords Coordinates) (*GeoPoint, error)
}

type GeoService interface {
    ResolveLocation(ctx context.Context, query string) (*GeoPoint, error)
    ResolveCoordinates(ctx context.Context, coords Coordinates) (*GeoPoint, error)
}

