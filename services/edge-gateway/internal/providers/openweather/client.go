package openweather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
    "strconv"

	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/domain/weather"
	"go.uber.org/zap"
)

type OpenWeatherProvider struct {
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

var GeoBaseURL = "https://api.openweathermap.org/geo/1.0/direct"
var OwmBaseURL = "https://api.openweathermap.org/data/3.0/onecall"

func New(key string, c *http.Client, logger *zap.Logger) *OpenWeatherProvider {
	return &OpenWeatherProvider{
		apiKey:     key,
		httpClient: c,
		logger:     logger,
	}
}

func (p *OpenWeatherProvider) FetchWeatherByCity(ctx context.Context, city string) (*weather.Weather, error) {
	p.logger.Debug("OpenWeather: fetch by city", zap.String("city", city))

	base, err := url.Parse(GeoBaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid GeoAPI base URL: %w", err)
	}

	params := url.Values{}
	params.Set("q", city)
	params.Set("limit", "1")
	params.Set("appid", p.apiKey)
	base.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
	if err != nil {
		p.logger.Error("failed to build city lookup request", zap.Error(err))
		return nil, err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.Error("city lookup failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("geo API error: %d %s", resp.StatusCode, body)
	}

	var results []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("city not found")
	}

	coords := weather.Coordinates{
		Latitude:  results[0].Lat,
		Longitude: results[0].Lon,
	}

	return p.FetchWeatherByCoordinates(ctx, coords)
}

func (p *OpenWeatherProvider) FetchWeatherByCoordinates(ctx context.Context, coords weather.Coordinates) (*weather.Weather, error) {
	if err := coords.Validate(); err != nil {
		return nil, err
	}

	p.logger.Debug("requesting weather from OpenWeather",
		zap.Float64("lat", coords.Latitude),
		zap.Float64("lon", coords.Longitude),
	)

	base, err := url.Parse(OwmBaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid OWM base URL: %w", err)
	}

	params := url.Values{}
    params.Set("lat", strconv.FormatFloat(coords.Latitude, 'f', 6, 64))
    params.Set("lon", strconv.FormatFloat(coords.Longitude, 'f', 6, 64))
	params.Set("appid", p.apiKey)
	params.Set("units", "metric")
	base.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
	if err != nil {
		p.logger.Error("failed to build request", zap.Error(err))
		return nil, err
	}

	start := time.Now()
	resp, err := p.httpClient.Do(req)
	duration := time.Since(start)

	p.logger.Debug("OWM request completed",
		zap.Duration("latency", duration),
	)

	if err != nil {
		p.logger.Error("http request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		p.logger.Warn("unexpected OpenWeather response",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)),
		)
		return nil, fmt.Errorf("OWM error %d: %s", resp.StatusCode, string(body))
	}

	var ow OneCallResponse
	if err := json.NewDecoder(resp.Body).Decode(&ow); err != nil {
		p.logger.Error("failed to decode OpenWeather JSON", zap.Error(err))
		return nil, err
	}

	return mapToDomain(&ow, coords), nil
}