package openweather

type OneCallResponse struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`

    Current struct {
        Temp      float64 `json:"temp"`
        FeelsLike float64 `json:"feels_like"`
        Humidity  float64 `json:"humidity"`
        WindSpeed float64 `json:"wind_speed"`
        WindDeg   int     `json:"wind_deg"`
        Weather []struct {
            Description string `json:"description"`
        } `json:"weather"`
    } `json:"current"`
}