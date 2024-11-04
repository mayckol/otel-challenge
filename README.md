# README

## Overview

This application consists of two services, **Service A** and **Service B**, designed to provide current weather information based on a Brazilian ZIP code (CEP). The services are built using Go and implement distributed tracing with OpenTelemetry (OTEL) and Zipkin.

- **Service A**: Receives a CEP via a POST request, validates it, and forwards it to Service B if valid.
- **Service B**: Receives the CEP, retrieves city information from the ViaCEP API, fetches current weather data from the WeatherAPI, and returns the temperature in Celsius, Fahrenheit, and Kelvin, along with the city name.

Both services are containerized using Docker and orchestrated with Docker Compose, which also includes the OTEL Collector and Zipkin for tracing.

## Project Structure
```
├── docker-compose.yml
├── otel-collector-config.yml
├── service_a
│   ├── Dockerfile
│   ├── main.go
│   ├── handler
│   │   └── handler.go
│   ├── http_client
│   │   └── service_b_client.go
│   └── utils
│       └── zipcode.go
└── service_b
    ├── Dockerfile
    ├── main.go
    ├── handler
    │   └── handler.go
    ├── http_client
    │   ├── via_cep_client.go
    │   └── weather_client.go
    └── utils
    └── numbers.go
```

## Objectives

- **Input Validation**: Service A ensures that the input CEP is an 8-digit string.
- **Data Retrieval**: Service B fetches city information using the ViaCEP API and weather data using the WeatherAPI.
- **Temperature Conversion**: Converts temperatures to Celsius, Fahrenheit, and Kelvin.
- **Distributed Tracing**: Implements OpenTelemetry and Zipkin for tracing across services.

## Prerequisites

- **Docker** and **Docker Compose** installed.
- **WeatherAPI Key**: Obtain an API key from [WeatherAPI](https://www.weatherapi.com/) and replace `your_weatherapi_key_here` in the environment variables.

## How to Start the Application

1. **Clone the Repository**:

```bash
git clone https://github.com/mayckol/otel-challenge.git
cd otel-challenge
```
2. **Update WeatherAPI Key:**
   Replace your_weatherapi_key_here in the docker-compose.yml with your actual WeatherAPI key.
3. **Build and Start the Services:**:
```bash
docker-compose up --build
```
This command builds the Docker images and starts all services, including Service A, Service B, the OTEL Collector, and Zipkin.
4. **Service A:**
   - Receives a POST request with JSON payload { "cep": "29902555" }.
   - Validates the CEP to ensure it is an 8-digit string.
   - Forwards the CEP to Service B if valid.
   - Returns a 422 Unprocessable Entity error if invalid.
- **Service B:**
   - Receives the CEP from Service A.
   - Fetches city information from the ViaCEP API.
   - Retrieves current weather data for the city from the WeatherAPI.
   - Converts temperature to Celsius, Fahrenheit, and Kelvin.
   - Returns a JSON response with city and temperature data.
- **Tracing:**
   - Both services are instrumented with OpenTelemetry.
   - Traces are collected by the OTEL Collector and exported to Zipkin.

## Testing the Application
### Valid Request Example:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"cep": "01001000"}' http://localhost:8081/service-a
```
Expected Response:
```
{
  "city": "São Paulo",
  "temp_c": 25.5,
  "temp_f": 77.9,
  "temp_k": 298.7
}
```

Invalid CEP Format:
- Response: 422 Unprocessable Entity
- Message: invalid zipcode

## Viewing Traces in Zipkin
1. Access Zipkin UI:
Open http://localhost:9411 in your web browser.

2. Search for Traces:
Click on "Run Query" to see the latest traces.

3. Analyze Traces:
- View the distributed trace spanning Service A and Service B.
- Analyze the spans to understand the request flow and performance.

### Shutdown
To stop the services, press Ctrl+C in the terminal where Docker Compose is running, or run:
```bash
docker-compose down
```