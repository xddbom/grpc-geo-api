package grpc

import (
    "context"

    pb "github.com/xddbom/grpc-geo-api/api/gen/go/weather_service/v1"
    "github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/domain/weather"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type WeatherServer struct {
    pb.UnimplementedWeatherServiceServer
    svc weather.WeatherService
}

func NewWeatherServer(svc weather.WeatherService) *WeatherServer {
    return &WeatherServer{svc: svc}
}

func (s *WeatherServer) GetWeatherByCoordinates(ctx context.Context, req *pb.GetWeatherByCoordinatesRequest) (*pb.GetWeatherResponse, error) {
    coords := weather.Coordinates{
        Latitude:  req.Lat,
        Longitude: req.Lon,
    }

    w, err := s.svc.GetWeatherByCoordinates(ctx, coords)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return mapToProto(w), nil
}

func (s *WeatherServer) GetWeatherByCity(ctx context.Context, req *pb.GetWeatherByCityRequest) (*pb.GetWeatherResponse, error) {
    w, err := s.svc.GetWeatherByCity(ctx, req.City)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return mapToProto(w), nil
}

func mapToProto(w *weather.Weather) *pb.GetWeatherResponse {
    return &pb.GetWeatherResponse{
        Lat:       w.Coordinates.Latitude,
        Lon:       w.Coordinates.Longitude,

        Actual:    w.Temp.Actual,
        FeelsLike: w.Temp.FeelsLike,

        WindSpeed: w.Wind.Speed,
        WindDeg:   int32(w.Wind.Deg),

        Humidity:  w.Humidity,
        Condition: w.Condition,
    }
}
