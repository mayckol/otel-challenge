package http_client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type WeatherClient struct {
	Client  *http.Client
	BaseURL string
	ApiKey  string
}

func NewWeatherClientClient(baseUrl, apiKey string, skipTLSVerification bool) *WeatherClient {
	transport := &http.Transport{}
	if skipTLSVerification {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &WeatherClient{
		Client:  &http.Client{Transport: transport},
		BaseURL: baseUrl,
		ApiKey:  apiKey,
	}
}

func (v *WeatherClient) WeatherDetails(locale string) (*WeatherAPIResponse, error) {
	encodedLocalidade := url.QueryEscape(locale)
	if encodedLocalidade == "" {
		return nil, errors.New(fmt.Sprintf("can not find zipcode@%d", http.StatusNotFound))
	}

	weatherAPIURL := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", v.BaseURL, v.ApiKey, encodedLocalidade)
	req, err := http.NewRequest(http.MethodGet, weatherAPIURL, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(err.Error()+"@%d", http.StatusInternalServerError))
	}

	resp, err := v.Client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(err.Error()+"@%d", http.StatusInternalServerError))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("error getting weather@%d", http.StatusInternalServerError))
	}

	weatherAPIResponse := WeatherAPIResponse{}
	err = json.NewDecoder(resp.Body).Decode(&weatherAPIResponse)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(err.Error()+"@%d", http.StatusInternalServerError))
	}

	return &weatherAPIResponse, nil
}
