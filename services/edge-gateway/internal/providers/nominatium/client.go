package nominatim

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/domain/geo"
)

type NominatimClient struct {
    http    *http.Client
    baseURL string
}

func NewNominatimClient(baseURL string, httpClient *http.Client) *NominatimClient {
    return &NominatimClient{
        http:    httpClient,
        baseURL: baseURL,
    }
}

func (c *NominatimClient) Search(ctx context.Context, query string) (*geo.GeoPoint, error) {
    base, _ := url.Parse(c.baseURL + "/search")
    params := url.Values{}
    params.Set("q", query)
    params.Set("format", "jsonv2")
    params.Set("addressdetails", "1")

    base.RawQuery = params.Encode()

    req, _ := http.NewRequestWithContext(ctx, "GET", base.String(), nil)
    req.Header.Set("User-Agent", "geo-api-demo")

    resp, err := c.http.Do(req)
    if err != nil {
        return nil, fmt.Errorf("nominatim search error: %w", err)
    }
    defer resp.Body.Close()

    var results []searchResult
    if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
        return nil, err
    }

    if len(results) == 0 {
        return nil, fmt.Errorf("no results")
    }

    return mapToDomain(&results[0])
}

func (c *NominatimClient) Reverse(ctx context.Context, coords geo.Coordinates) (*geo.GeoPoint, error) {
    base, _ := url.Parse("https://nominatim.openstreetmap.org/reverse")
    params := url.Values{}
    params.Set("lat", fmt.Sprintf("%f", coords.Lat))
    params.Set("lon", fmt.Sprintf("%f", coords.Lon))
    params.Set("format", "jsonv2")
    params.Set("addressdetails", "1")

    base.RawQuery = params.Encode()

    req, _ := http.NewRequestWithContext(ctx, "GET", base.String(), nil)
    req.Header.Set("User-Agent", "geo-api-demo")

    resp, err := c.http.Do(req)
    if err != nil {
        return nil, fmt.Errorf("nominatim reverse error: %w", err)
    }
    defer resp.Body.Close()

    var result searchResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return mapToDomain(&result)
}

