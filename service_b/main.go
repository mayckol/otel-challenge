package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mayckol/otel-challenge/service-b/handler"
	"github.com/mayckol/otel-challenge/service-b/http_client"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	serviceBClient := http_client.NewViaCepClient(
		os.Getenv("SERVICE_B_URL"),
		true,
	)

	weatherClient := http_client.NewWeatherClientClient(
		os.Getenv("WEATHER_API_URL"),
		os.Getenv("WEATHER_API_KEY"),
		true,
	)

	weatherHandler := handler.NewWeatherHandler(serviceBClient, weatherClient)

	mux := http.NewServeMux()
	mux.Handle("/service-b", otelhttp.NewHandler(http.HandlerFunc(weatherHandler.Weather), "ServiceA"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	fmt.Printf("Service B running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
