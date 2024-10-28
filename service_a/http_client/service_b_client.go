package http_client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

type Request struct {
	CEP string `json:"cep"`
}

type ServiceBClient struct {
	Client  *http.Client
	BaseURL string
}

func NewServiceBClient(baseURL string, skipTLSVerification bool) *ServiceBClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerification},
	}

	client := &http.Client{
		Transport: otelhttp.NewTransport(transport),
	}

	return &ServiceBClient{
		Client:  client,
		BaseURL: baseURL,
	}
}

func (v *ServiceBClient) WeatherDetails(ctx context.Context, zipCode string) (*ServiceBResponse, error) {
	url := fmt.Sprintf("%s/service-b?zipcode=%s", v.BaseURL, zipCode)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := v.Client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(err.Error()+"@%d", http.StatusInternalServerError))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("error getting address@%d", http.StatusInternalServerError))
	}

	serviceBResponse := ServiceBResponse{}
	err = json.NewDecoder(resp.Body).Decode(&serviceBResponse)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(err.Error()+"@%d", http.StatusInternalServerError))
	}

	if serviceBResponse.City == "" {
		return nil, errors.New(fmt.Sprintf("can not find city@%d", http.StatusNotFound))
	}
	return &serviceBResponse, nil
}
