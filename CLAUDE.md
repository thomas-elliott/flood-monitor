# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running
```bash
# Run locally with environment variables
LATITUDE=52.52 LONGITUDE=13.41 NFTY_URL=https://ntfy.sh/weather-alert go run ./cmd/weather-alert

# Build the binary
go build -o weather-alert ./cmd/weather-alert

# Build Docker image
docker build -t flood-monitor .
```

### Go Commands
```bash
# Format code
go fmt ./...

# Vet code for potential issues
go vet ./...

# Download dependencies
go mod download

# Update dependencies
go mod tidy
```

## Code Architecture

This is a single-purpose Go application that monitors weather conditions and sends alerts via ntfy. The entire application logic is contained in a single file: `cmd/weather-alert/main.go`.

### Key Components

- **Weather Monitoring**: Uses Open-Meteo API to fetch hourly precipitation and weather code data
- **Alert System**: Sends HTTP POST messages to configurable ntfy endpoints
- **Configuration**: Environment variable-based configuration (LATITUDE, LONGITUDE, NFTY_URL, CHECK_INTERVAL)

### Core Logic Flow
1. Load configuration from environment variables
2. Start ticker-based monitoring loop (default 1 hour interval)
3. For each check:
   - Fetch 3-hour weather forecast from Open-Meteo API
   - Calculate total precipitation in next 3 hours
   - Check for thunderstorm weather codes (95-99)
   - Send alert if precipitation ≥50mm OR thunderstorm detected

### Data Structures
- `config`: Holds latitude, longitude, check interval, and ntfy URL
- `forecastResponse`: Maps to Open-Meteo API JSON response structure

### Environment Variables
- `LATITUDE`: Location latitude for weather monitoring
- `LONGITUDE`: Location longitude for weather monitoring  
- `NFTY_URL`: Full URL of ntfy topic for alert messages
- `CHECK_INTERVAL`: Optional Go duration string (default "1h")

## Deployment

The application is designed to run as a single replica in Kubernetes with the provided Dockerfile. It's a long-running process that checks weather conditions periodically and sends alerts when thresholds are exceeded.