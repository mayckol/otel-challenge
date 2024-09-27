package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mayckol/otel-challenge/service-a/handler"
	"github.com/mayckol/otel-challenge/service-a/http_client"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	//shutdown := initTracer()
	//defer shutdown()

	serviceBClient := http_client.NewServiceBClient(
		os.Getenv("SERVICE_B_URL"),
		true,
	)

	ServiceBHandler := handler.NewServiceBHandler(serviceBClient)

	mux := http.NewServeMux()
	mux.Handle("/service-a", otelhttp.NewHandler(http.HandlerFunc(ServiceBHandler.Handle), "ServiceA"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	fmt.Printf("Service A running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
