# Temperature by Zip Code Service

## Overview

This project consists of two services:

- **Service A**: Receives a ZIP code via POST, validates it, and forwards it to Service B.
- **Service B**: Receives a ZIP code, retrieves the city from ViaCEP, fetches the current weather from WeatherAPI, and returns temperature data.

Both services implement OpenTelemetry (OTEL) tracing and Zipkin for distributed tracing.

## Prerequisites

- Docker
- Docker Compose
- A WeatherAPI key (sign up at [WeatherAPI](https://www.weatherapi.com/) to get one)