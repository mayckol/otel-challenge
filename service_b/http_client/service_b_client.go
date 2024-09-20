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

type ViaCepClient struct {
	Client  *http.Client
	BaseURL string
}

func NewViaCepClient(baseURL string, skipTLSVerification bool) *ViaCepClient {
	transport := &http.Transport{}
	if skipTLSVerification {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &ViaCepClient{
		Client:  &http.Client{Transport: transport},
		BaseURL: baseURL,
	}
}

func (v *ViaCepClient) AddressDetails(zipCode string) (*ViaCepResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/json", v.BaseURL, zipCode), nil)
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

	viaCepResponse := ViaCepResponse{}
	err = json.NewDecoder(resp.Body).Decode(&viaCepResponse)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(err.Error()+"@%d", http.StatusInternalServerError))
	}

	if viaCepResponse.Localidade == "" {
		return nil, errors.New(fmt.Sprintf("can not find zipcode@%d", http.StatusNotFound))
	}
	return &viaCepResponse, nil
}
