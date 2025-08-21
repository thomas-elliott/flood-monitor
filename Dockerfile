FROM golang:1.21-alpine AS build

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download

COPY cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/weather-alert ./cmd/weather-alert

FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=build /bin/weather-alert /usr/local/bin/weather-alert

LABEL org.opencontainers.image.source=https://github.com/thomas-elliott/flood-monitor
LABEL org.opencontainers.image.description="Weather monitoring and flood alert application"
LABEL org.opencontainers.image.licenses=MIT

USER nobody:nobody

ENTRYPOINT ["/usr/local/bin/weather-alert"]
