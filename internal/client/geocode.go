package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const nominatimURL = "https://nominatim.openstreetmap.org/search"

type GeocodingClient struct {
	httpClient *http.Client
}

type GeocodingResult struct {
	Latitude    float64
	Longitude   float64
	DisplayName string
}

type nominatimResponse struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

func NewGeocodingClient() *GeocodingClient {
	return &GeocodingClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *GeocodingClient) Search(city string) (*GeocodingResult, error) {
	params := url.Values{}
	params.Add("q", city)
	params.Add("format", "json")
	params.Add("limit", "1")

	reqURL := nominatimURL + "?" + params.Encode()
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "Horoscope-TUI/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nominatim returned status %d", resp.StatusCode)
	}

	var results []nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("city not found: %s", city)
	}

	var lat, lon float64
	_, err = fmt.Sscanf(results[0].Lat, "%f", &lat)
	if err != nil {
		return nil, fmt.Errorf("parse latitude: %w", err)
	}
	_, err = fmt.Sscanf(results[0].Lon, "%f", &lon)
	if err != nil {
		return nil, fmt.Errorf("parse longitude: %w", err)
	}

	return &GeocodingResult{
		Latitude:    lat,
		Longitude:   lon,
		DisplayName: results[0].DisplayName,
	}, nil
}
