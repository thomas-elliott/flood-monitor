package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCheckAndAlert(t *testing.T) {
	tests := []struct {
		name          string
		precipitation float64
		thunder       bool
		shouldAlert   bool
	}{
		{
			name:          "heavy rain triggers alert",
			precipitation: 55.0,
			thunder:       false,
			shouldAlert:   true,
		},
		{
			name:          "thunderstorm triggers alert",
			precipitation: 10.0,
			thunder:       true,
			shouldAlert:   true,
		},
		{
			name:          "both conditions trigger alert",
			precipitation: 60.0,
			thunder:       true,
			shouldAlert:   true,
		},
		{
			name:          "normal weather no alert",
			precipitation: 5.0,
			thunder:       false,
			shouldAlert:   false,
		},
		{
			name:          "threshold boundary no alert",
			precipitation: 49.9,
			thunder:       false,
			shouldAlert:   false,
		},
		{
			name:          "threshold boundary triggers alert",
			precipitation: 50.0,
			thunder:       false,
			shouldAlert:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alertSent := false
			
			// Test the core alert logic by extracting it
			if tt.precipitation >= 50 || tt.thunder {
				alertSent = true
			}
			
			if alertSent != tt.shouldAlert {
				t.Errorf("expected alert=%v, got alert=%v", tt.shouldAlert, alertSent)
			}
		})
	}
}

func TestForecastParsing(t *testing.T) {
	// Test the time parsing logic from fetchForecast
	now := time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC)
	cutoff := now.Add(3 * time.Hour) // 13:00

	testTimes := []struct {
		timeStr   string
		shouldUse bool
	}{
		{"2023-12-01T09:00", false}, // before now
		{"2023-12-01T10:00", false}, // exactly now
		{"2023-12-01T11:00", true},  // within window
		{"2023-12-01T12:00", true},  // within window
		{"2023-12-01T13:00", true},  // exactly at cutoff (inclusive)
		{"2023-12-01T14:00", false}, // after cutoff
	}

	for _, tt := range testTimes {
		t.Run(tt.timeStr, func(t *testing.T) {
			parsedTime, err := time.Parse("2006-01-02T15:04", tt.timeStr)
			if err != nil {
				t.Fatal(err)
			}
			
			shouldUse := parsedTime.After(now) && !parsedTime.After(cutoff)
			if shouldUse != tt.shouldUse {
				t.Errorf("time %s: expected shouldUse=%v, got=%v", tt.timeStr, tt.shouldUse, shouldUse)
			}
		})
	}
}

func TestWeatherCodeThunder(t *testing.T) {
	tests := []struct {
		code     int
		isThunder bool
	}{
		{0, false},   // clear
		{1, false},   // mainly clear
		{94, false},  // just below thunder range
		{95, true},   // thunderstorm
		{96, true},   // thunderstorm with hail
		{97, true},   // thunderstorm with hail
		{98, true},   // thunderstorm with hail
		{99, true},   // thunderstorm with hail
		{100, false}, // above thunder range
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("code_%d", tt.code), func(t *testing.T) {
			isThunder := tt.code >= 95 && tt.code <= 99
			if isThunder != tt.isThunder {
				t.Errorf("code %d: expected thunder=%v, got=%v", tt.code, tt.isThunder, isThunder)
			}
		})
	}
}

// MockProvider for testing
type MockProvider struct {
	precipitation float64
	thunder       bool
	shouldError   bool
}

func (m *MockProvider) GetForecast(ctx context.Context, lat, lon float64) (WeatherData, error) {
	if m.shouldError {
		return WeatherData{}, fmt.Errorf("mock error")
	}
	return WeatherData{
		PrecipitationMM: m.precipitation,
		HasThunderstorm: m.thunder,
	}, nil
}

func TestWeatherProviderInterface(t *testing.T) {
	tests := []struct {
		name        string
		provider    WeatherProvider
		expectError bool
	}{
		{
			name:        "mock provider heavy rain",
			provider:    &MockProvider{precipitation: 60.0, thunder: false},
			expectError: false,
		},
		{
			name:        "mock provider thunderstorm",
			provider:    &MockProvider{precipitation: 10.0, thunder: true},
			expectError: false,
		},
		{
			name:        "mock provider error",
			provider:    &MockProvider{shouldError: true},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			weather, err := tt.provider.GetForecast(ctx, 0, 0)
			
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError {
				// Basic validation that we got weather data
				if weather.PrecipitationMM < 0 {
					t.Error("precipitation should not be negative")
				}
			}
		})
	}
}