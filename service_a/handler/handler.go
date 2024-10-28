package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/mayckol/otel-challenge/service-a/http_client"
	"github.com/mayckol/otel-challenge/service-a/utils"
	"go.opentelemetry.io/otel"
)

type ServiceBHandler struct {
	ServiceBClient http_client.ServiceBClientInterface
}

func NewServiceBHandler(ServiceBClient http_client.ServiceBClientInterface) *ServiceBHandler {
	return &ServiceBHandler{
		ServiceBClient: ServiceBClient,
	}
}

func (h *ServiceBHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("service_a")
	ctx, span := tracer.Start(ctx, "HandleRequest")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var req http_client.Request
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	zipCode := utils.ZipCode(req.CEP)
	if !zipCode.IsValid() {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	ServiceBResponse, err := h.ServiceBClient.WeatherDetails(ctx, string(zipCode))
	if err != nil {
		var statusCode int
		var message string
		parts := splitError(err.Error())
		if len(parts) == 2 {
			message = parts[0]
			statusCode, _ = parseStatusCode(parts[1])
		} else {
			message = "error fetching ServiceB"
			statusCode = http.StatusInternalServerError
		}
		http.Error(w, message, statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ServiceBResponse)
}

func splitError(errMsg string) []string {
	return strings.Split(errMsg, "@")
}

func parseStatusCode(codeStr string) (int, error) {
	return strconv.Atoi(codeStr)
}
