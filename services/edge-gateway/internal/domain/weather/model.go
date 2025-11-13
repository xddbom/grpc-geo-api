package weather

import (
	"fmt"
	"errors"
)

type Coordinates struct {
	Latitude float64 	
	Longitude float64 	
}

func (c Coordinates) Validate() error {
	if c.Latitude < -90 || c.Latitude > 90 {
		return fmt.Errorf("invalid latitude: %f (must be between -90 and 90)", c.Latitude)
	}
	if c.Longitude < -180 || c.Longitude > 180 {
        return errors.New("longitude out of range")
    }
    return nil
}

type Temperature struct {
    Actual    float64
    FeelsLike float64
}

type Wind struct {
    Speed float64
    Deg   int
}

type Weather struct {
    Coordinates Coordinates
    Temp        Temperature
    Wind        Wind
    Humidity    float64
    Condition   string
}
