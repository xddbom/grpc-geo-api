package geo

import (
	"fmt"
)

type Coordinates struct {
    Lat float64
    Lon float64
}

func (c Coordinates) Validate() error {
    if c.Lat < -90 || c.Lat > 90 {
        return fmt.Errorf("latitude out of range")
    }
    if c.Lon < -180 || c.Lon > 180 {
        return fmt.Errorf("longitude out of range")
    }
    return nil
}

type GeoPoint struct {
    Name        string
    City        string
    State       string
    Country     string
    CountryCode string
    Coordinates Coordinates
}
