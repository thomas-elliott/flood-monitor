FROM golang:1.20-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
RUN go build -o /bin/weather-alert ./cmd/weather-alert

FROM alpine:3.18
COPY --from=build /bin/weather-alert /usr/local/bin/weather-alert
ENTRYPOINT ["/usr/local/bin/weather-alert"]
