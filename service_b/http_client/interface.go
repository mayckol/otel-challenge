package http_client

type ViaCepClientInterface interface {
	AddressDetails(zipCode string) (*ViaCepResponse, error)
}

type WeatherClientInterface interface {
	WeatherDetails(locale string) (*WeatherAPIResponse, error)
}
