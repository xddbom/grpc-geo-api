package weather_test    // ?

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/domain/weather"
)

func TestCoordinates_Validate(t *testing.T) {
    tests := []struct {
        name    string
        input   weather.Coordinates
        wantErr bool
    }{
        {
            name:    "valid coordinates",
            input:   weather.Coordinates{Latitude: 50.45, Longitude: 30.52},
            wantErr: false,
        },
        {
            name:    "invalid latitude too big",
            input:   weather.Coordinates{Latitude: 100, Longitude: 0},
            wantErr: true,
        },
        {
            name:    "invalid latitude too small",
            input:   weather.Coordinates{Latitude: -100, Longitude: 0},
            wantErr: true,
        },
        {
            name:    "invalid longitude too big",
            input:   weather.Coordinates{Latitude: 0, Longitude: 200},
            wantErr: true,
        },
        {
            name:    "invalid longitude too small",
            input:   weather.Coordinates{Latitude: 0, Longitude: -200},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.input.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
