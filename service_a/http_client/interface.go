package http_client

import "context"

type ServiceBClientInterface interface {
	WeatherDetails(ctx context.Context, zipCode string) (*ServiceBResponse, error)
}
