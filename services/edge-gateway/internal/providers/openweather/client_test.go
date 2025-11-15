package openweather_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/providers/openweather"
	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/domain/weather"
	"go.uber.org/zap"
)

func TestOpenWeatherProvider_HappyPath(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/geo/1.0/direct", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Kyiv", r.URL.Query().Get("q"))

		json.NewEncoder(w).Encode([]map[string]any{
			{"lat": 50.45, "lon": 30.523},
		})
	})

	mux.HandleFunc("/data/3.0/onecall", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"current": map[string]any{
				"temp":        22.0,
				"feels_like":  21.0,
				"humidity":    40.0,
				"wind_speed":  3.5,
				"wind_deg":    180,
				"weather": []map[string]any{{"description": "clear sky"}},
			},
		})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	openweather.GeoBaseURL = server.URL + "/geo/1.0/direct"
	openweather.OwmBaseURL = server.URL + "/data/3.0/onecall"

	p := openweather.New("test-key", server.Client(), zap.NewNop())

	res, err := p.FetchWeatherByCity(context.Background(), "Kyiv")
	assert.NoError(t, err)

	assert.Equal(t, weather.Coordinates{Latitude: 50.45, Longitude: 30.523}, res.Coordinates)
	assert.Equal(t, 22.0, res.Temp.Actual)
	assert.Equal(t, 21.0, res.Temp.FeelsLike)
	assert.Equal(t, 40.0, res.Humidity)
	assert.Equal(t, 3.5, res.Wind.Speed)
	assert.Equal(t, 180, res.Wind.Deg)
	assert.Equal(t, "clear sky", res.Condition)
}

func TestOpenWeatherProvider_GeoEmptyResult(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/geo/1.0/direct", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	openweather.GeoBaseURL = server.URL + "/geo/1.0/direct"
	openweather.OwmBaseURL = server.URL 

	p := openweather.New("key", server.Client(), zap.NewNop())

	res, err := p.FetchWeatherByCity(context.Background(), "UnknownCity")

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "city not found")
}

func TestOpenWeatherProvider_OneCallInvalidJSON(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/geo/1.0/direct", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"lat":50.45,"lon":30.523}]`))
	})

	mux.HandleFunc("/data/3.0/onecall", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{ invalid json`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	openweather.GeoBaseURL = server.URL + "/geo/1.0/direct"
	openweather.OwmBaseURL = server.URL + "/data/3.0/onecall"

	p := openweather.New("key", server.Client(), zap.NewNop())

	res, err := p.FetchWeatherByCity(context.Background(), "Kyiv")

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestOpenWeatherProvider_GeoAPI500(t *testing.T) {
    mux := http.NewServeMux()

    mux.HandleFunc("/geo/1.0/direct", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(500)
        w.Write([]byte(`internal error`))
    })

    server := httptest.NewServer(mux)
    defer server.Close()

    openweather.GeoBaseURL = server.URL + "/geo/1.0/direct"
    openweather.OwmBaseURL = server.URL

    p := openweather.New("key", server.Client(), zap.NewNop())

    res, err := p.FetchWeatherByCity(context.Background(), "Kyiv")

    assert.Nil(t, res)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "geo API error: 500")
}
