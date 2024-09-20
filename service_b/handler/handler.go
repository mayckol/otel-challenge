package handler

import (
	"encoding/json"
	"github.com/mayckol/otel-challenge/service-b/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/mayckol/otel-challenge/service-b/http_client"
)

type WeatherHandler struct {
	ViaCepClient  http_client.ViaCepClientInterface
	WeatherClient http_client.WeatherClientInterface
}

func NewWeatherHandler(viaCepClient http_client.ViaCepClientInterface, weatherClient http_client.WeatherClientInterface) *WeatherHandler {
	return &WeatherHandler{
		ViaCepClient:  viaCepClient,
		WeatherClient: weatherClient,
	}
}

func (h *WeatherHandler) Weather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("method not allowed")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	zipcode := r.URL.Query().Get("zipcode")
	serializedZipCode := utils.ZipCode(zipcode)
	if !serializedZipCode.IsValid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if _, err := w.Write([]byte("invalid zipcode")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	addressDetails, err := h.ViaCepClient.AddressDetails(zipcode)
	if err != nil {
		slices := strings.Split(err.Error(), "@")
		if len(slices) != 2 {
			if _, err := w.Write([]byte("error making request")); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		message, statusCode := slices[0], slices[1]
		code, _ := strconv.Atoi(statusCode)
		w.WriteHeader(code)
		w.Write([]byte(message))
		return
	}

	if addressDetails == nil {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("can not find zipcode")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	weatherAPIResponse, err := h.WeatherClient.WeatherDetails(addressDetails.Localidade)
	if err != nil {
		slices := strings.Split(err.Error(), "@")
		if len(slices) != 2 {
			if _, err := w.Write([]byte("error making request")); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		message, statusCode := slices[0], slices[1]
		code, _ := strconv.Atoi(statusCode)
		w.WriteHeader(code)
		w.Write([]byte(message))
		return
	}

	if weatherAPIResponse == nil {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("can not find zipcode")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	tempK := utils.RoundToDecimal(weatherAPIResponse.Current.TempC+273.15, 1)
	weatherResponse := http_client.WeatherResponse{
		City:  addressDetails.Localidade,
		TempC: utils.RoundToDecimal(weatherAPIResponse.Current.TempC, 1),
		TempF: utils.RoundToDecimal(weatherAPIResponse.Current.TempF, 1),
		TempK: tempK,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(weatherResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
