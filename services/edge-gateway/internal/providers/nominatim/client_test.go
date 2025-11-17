package nominatim_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/domain/geo"
	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/providers/nominatim"
)

func TestNominatimClient_Search_HappyPath(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Berlin", r.URL.Query().Get("q"))
		assert.Equal(t, "jsonv2", r.URL.Query().Get("format"))
		assert.Equal(t, "1", r.URL.Query().Get("addressdetails"))

		json.NewEncoder(w).Encode([]map[string]any{
			{
				"display_name": "Berlin, Germany",
				"lat":          "52.5200",
				"lon":          "13.4050",
				"address": map[string]any{
					"city":         "Berlin",
					"state":        "Berlin",
					"country":      "Germany",
					"country_code": "de",
				},
			},
		})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := nominatim.NewNominatimClient(server.URL, server.Client())

	res, err := client.Search(context.Background(), "Berlin")
	assert.NoError(t, err)

	assert.Equal(t, "Berlin, Germany", res.Name)
	assert.Equal(t, "Berlin", res.City)
	assert.Equal(t, "Germany", res.Country)
	assert.Equal(t, 52.5200, res.Coordinates.Lat)
	assert.Equal(t, 13.4050, res.Coordinates.Lon)
}

func TestNominatimClient_Search_EmptyResult(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := nominatim.NewNominatimClient(server.URL, server.Client())

	res, err := client.Search(context.Background(), "UnknownPlace")

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no results")
}

func TestNominatimClient_Search_InvalidJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{ invalid json`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := nominatim.NewNominatimClient(server.URL, server.Client())

	res, err := client.Search(context.Background(), "Berlin")

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}


func TestNominatimClient_Search_HTTP500(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`server error`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := nominatim.NewNominatimClient(server.URL, server.Client())

	res, err := client.Search(context.Background(), "Berlin")

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}


func TestNominatimClient_Reverse_HappyPath(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/reverse", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "52.520000", r.URL.Query().Get("lat"))
		assert.Equal(t, "13.405000", r.URL.Query().Get("lon"))

		json.NewEncoder(w).Encode(map[string]any{
			"display_name": "Berlin, Germany",
			"lat":          "52.5200",
			"lon":          "13.4050",
			"address": map[string]any{
				"city":         "Berlin",
				"state":        "Berlin",
				"country":      "Germany",
				"country_code": "de",
			},
		})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := nominatim.NewNominatimClient(server.URL, server.Client())

	coords := geo.Coordinates{Lat: 52.52, Lon: 13.405}
	res, err := client.Reverse(context.Background(), coords)

	assert.NoError(t, err)
	assert.Equal(t, "Berlin", res.City)
	assert.Equal(t, "Germany", res.Country)
}


func TestNominatimClient_Reverse_InvalidJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/reverse", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not a json`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := nominatim.NewNominatimClient(server.URL, server.Client())

	res, err := client.Reverse(context.Background(), geo.Coordinates{Lat: 52.52, Lon: 13.405})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}


func TestNominatimClient_Reverse_HTTP500(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/reverse", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := nominatim.NewNominatimClient(server.URL, server.Client())

	res, err := client.Reverse(context.Background(), geo.Coordinates{Lat: 1, Lon: 1})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}
