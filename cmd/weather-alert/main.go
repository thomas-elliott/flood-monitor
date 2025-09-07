package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type WeatherData struct {
	PrecipitationMM float64
	HasThunderstorm bool
}

type WeatherProvider interface {
	GetForecast(ctx context.Context, lat, lon float64) (WeatherData, error)
}

type OpenMeteoProvider struct{}

// AucklandProvider for NIWA/MetService data (placeholder)
type AucklandProvider struct {
	// Add fields for API keys, endpoints, etc.
}

type forecastResponse struct {
	Hourly struct {
		Time          []string  `json:"time"`
		Precipitation []float64 `json:"precipitation"`
		Weathercode   []int     `json:"weathercode"`
	} `json:"hourly"`
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	
	provider := &OpenMeteoProvider{}
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		checkAndAlert(cfg, provider)
		<-ticker.C
	}
}

func checkAndAlert(cfg *config, provider WeatherProvider) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	weather, err := provider.GetForecast(ctx, cfg.Latitude, cfg.Longitude)
	if err != nil {
		fmt.Fprintf(os.Stderr, "forecast error: %v\n", err)
		return
	}

	if weather.PrecipitationMM >= 50 || weather.HasThunderstorm {
		message := fmt.Sprintf("Heavy rain (%.1fmm) or thunderstorm expected in next 3h", weather.PrecipitationMM)
		if err := sendAlert(ctx, cfg.NtfyURL, message); err != nil {
			fmt.Fprintf(os.Stderr, "alert error: %v\n", err)
		}
	}
}

type config struct {
	Latitude  float64
	Longitude float64
	Interval  time.Duration
	NtfyURL   string
}

func loadConfig() (*config, error) {
	latStr := os.Getenv("LATITUDE")
	lonStr := os.Getenv("LONGITUDE")
	nftyURL := os.Getenv("NFTY_URL")
	if latStr == "" || lonStr == "" || nftyURL == "" {
		return nil, fmt.Errorf("LATITUDE, LONGITUDE and NFTY_URL must be set")
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return nil, err
	}
	interval := time.Hour
	if iv := os.Getenv("CHECK_INTERVAL"); iv != "" {
		if d, err := time.ParseDuration(iv); err == nil {
			interval = d
		}
	}
	return &config{Latitude: lat, Longitude: lon, Interval: interval, NtfyURL: nftyURL}, nil
}

func (p *OpenMeteoProvider) GetForecast(ctx context.Context, lat, lon float64) (WeatherData, error) {
	precipSum, thunder, err := p.fetchOpenMeteoForecast(ctx, lat, lon)
	return WeatherData{
		PrecipitationMM: precipSum,
		HasThunderstorm: thunder,
	}, err
}

func (p *OpenMeteoProvider) fetchOpenMeteoForecast(ctx context.Context, lat, lon float64) (precipSum float64, thunder bool, err error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&hourly=precipitation,weathercode&forecast_days=1&timezone=UTC", lat, lon)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var fr forecastResponse
	if err = json.Unmarshal(body, &fr); err != nil {
		return
	}
	now := time.Now().UTC()
	cutoff := now.Add(3 * time.Hour)
	for i, tStr := range fr.Hourly.Time {
		t, err2 := time.Parse("2006-01-02T15:04", tStr)
		if err2 != nil {
			continue
		}
		if t.After(now) && !t.After(cutoff) {
			if i < len(fr.Hourly.Precipitation) {
				precipSum += fr.Hourly.Precipitation[i]
			}
			if i < len(fr.Hourly.Weathercode) {
				code := fr.Hourly.Weathercode[i]
				if code >= 95 && code <= 99 {
					thunder = true
				}
			}
		}
	}
	return
}

func sendAlert(ctx context.Context, ntfyURL, message string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ntfyURL, strings.NewReader(message))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("ntfy returned status %s: %s", resp.Status, string(body))
	}
	return nil
}

// GetForecast implements WeatherProvider for Auckland-specific data
func (p *AucklandProvider) GetForecast(ctx context.Context, lat, lon float64) (WeatherData, error) {
	// TODO: Implement NIWA/MetService integration
	// This could be:
	// - NIWA API calls if available
	// - MetService web scraping
	// - Combination of both
	
	return WeatherData{}, fmt.Errorf("AucklandProvider not implemented yet")
}
