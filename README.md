# Flood Monitor

A Go application for monitoring severe weather conditions and sending real-time alerts. Currently monitors precipitation and thunderstorms using weather data APIs, with plans for modular provider architecture to support specialized flood monitoring data sources.

## Current Architecture

### Core Components
- **Weather Monitoring**: Fetches hourly precipitation and weather code data  
- **Alert Logic**: Triggers on ≥50mm precipitation in 3 hours OR thunderstorm detection (weather codes 95-99)
- **Notification System**: HTTP POST alerts via ntfy
- **Scheduler**: Configurable interval-based monitoring (default: 1 hour)

### Data Flow
1. Load configuration from environment variables
2. Start periodic monitoring ticker  
3. For each check:
   - Query weather API for 3-hour forecast
   - Calculate precipitation totals and check for thunderstorms
   - Send alert if thresholds exceeded

### Current Implementation
- **Weather Provider**: Open-Meteo API integration (`fetchForecast()` in `cmd/weather-alert/main.go:87`)
- **Alert Thresholds**: Hardcoded 50mm precipitation + thunderstorm codes 95-99
- **Single-file Architecture**: All logic in `main.go` (144 lines)

The application is designed to run as a long-running process in Kubernetes but can also run standalone.

## Configuration

The application is configured via environment variables:

| Variable        | Description                                              |
|-----------------|----------------------------------------------------------|
| `LATITUDE`      | Latitude of the location to monitor.                     |
| `LONGITUDE`     | Longitude of the location to monitor.                    |
| `NFTY_URL`      | Full URL of the ntfy topic to publish messages to.       |
| `CHECK_INTERVAL`| Optional interval (Go duration) between checks. Default `1h`. |

## Running locally

```
LATITUDE=52.52 LONGITUDE=13.41 NFTY_URL=https://ntfy.sh/weather-alert \
    go run ./cmd/weather-alert
```

## Kubernetes

A simple deployment could look like this:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: flood-monitor
spec:
  replicas: 1
  selector:
    matchLabels:
      app: flood-monitor
  template:
    metadata:
      labels:
        app: flood-monitor
    spec:
      containers:
      - name: flood-monitor
        image: yourrepo/flood-monitor:latest
        env:
        - name: LATITUDE
          value: "52.52"
        - name: LONGITUDE
          value: "13.41"
        - name: NFTY_URL
          value: "https://ntfy.sh/weather-alert"
```

Build the Docker image with the provided `Dockerfile` and deploy using the
manifest above or integrate into your existing setup.

## Development Roadmap

### Next Steps
1. **Simple Weather Provider Interface** - Basic abstraction to swap data sources
2. **Unit Testing** - Test core logic with mock weather data  
3. **Configuration Flexibility** - Customizable alert thresholds

### Testing Approach
- Mock weather data for unit testing core alert logic
- Simple test scenarios: heavy rain, thunderstorms, normal weather
