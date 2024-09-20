package http_client

type ServiceBClientInterface interface {
	WeatherDetails(zipCode string) (*ServiceBResponse, error)
}
