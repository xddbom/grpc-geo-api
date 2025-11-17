workspace {

    model {
        user = person "User"
        gateway = softwareSystem "Edge Gateway"
        weather = softwareSystem "Weather Provider (OpenWeather)"
        geo = softwareSystem "Geo Provider (Nominatim)"

        user -> gateway "requests data"
        gateway -> weather "fetches weather"
        gateway -> geo "fetches coordinates"
    }

    views {
        systemContext gateway "c1" {
            include *
            autolayout
        }

        container gateway "c2" {
            include *
            autolayout
        }

        theme default
    }
}
