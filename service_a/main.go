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

//func initTracer() func() {
//zipkinURL := os.Getenv("ZIPKIN_ENDPOINT")
//exporter, err := zipkin.New(zipkinURL)
//if err != nil {
//	log.Fatalf("Failed to create Zipkin exporter: %v", err)
//}
//
//tp := sdktrace.NewTracerProvider(
//	sdktrace.WithBatcher(exporter),
//	sdktrace.WithResource(resource.NewWithAttributes(
//		semconv.SchemaURL,
//		semconv.ServiceNameKey.String("service-a"),
//	)),
//)
//
//otel.SetTracerProvider(tp)
//
//return func() {
//	if err := tp.Shutdown(context.Background()); err != nil {
//		log.Fatalf("Error shutting down tracer provider: %v", err)
//	}
//}
//}

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
