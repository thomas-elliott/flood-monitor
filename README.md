# Flood Monitor

This is a small Go application that periodically checks the weather forecast
using the [Open-Meteo](https://open-meteo.com/) API. If heavy rain is expected
(precipitation over 50&nbsp;mm in the next three hours) or a thunderstorm is
forecast, it sends a message to a configured [ntfy](https://ntfy.sh/) endpoint.

The application is intended to run inside a Kubernetes cluster but can also run
stand‑alone.

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
