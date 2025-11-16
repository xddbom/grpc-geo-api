package nominatim

type searchResult struct {
    DisplayName string  `json:"display_name"`
    Lat         string  `json:"lat"`
    Lon         string  `json:"lon"`
    Address     struct {
        City        string `json:"city"`
        State       string `json:"state"`
        Country     string `json:"country"`
        CountryCode string `json:"country_code"`
    } `json:"address"`
}
