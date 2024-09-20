package http_client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
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
	transport := &http.Transport{}
	if skipTLSVerification {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &ServiceBClient{
		Client:  &http.Client{Transport: transport},
		BaseURL: baseURL,
	}
}

func (v *ServiceBClient) WeatherDetails(zipCode string) (*ServiceBResponse, error) {
	url := fmt.Sprintf("%s/service-b?zipcode=%s", v.BaseURL, zipCode)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(url), nil)
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
