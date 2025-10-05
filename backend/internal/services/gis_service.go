package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GISService работает с 2GIS API
type GISService struct {
	apiKey  string
	appID   string
	baseURL string
	client  *http.Client
}

func NewGISService() *GISService {
	return &GISService{
		apiKey:  os.Getenv("GAPIS_API_KEY"),
		appID:   os.Getenv("GAPIS_APP_ID"),
		baseURL: "https://catalog.api.2gis.com",
		client:  &http.Client{},
	}
}

// FindPlacesNearby ищет места рядом с точкой
func (s *GISService) FindPlacesNearby(lat, lon, radius float64, query string) ([]Place, error) {
	url := fmt.Sprintf("%s/3.0/items?key=%s&point=%f,%f&radius=%f&q=%s",
		s.baseURL, s.apiKey, lon, lat, radius, query)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call 2GIS API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("2GIS API error: %s", string(body))
	}

	var result PlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Result.Items, nil
}

// BuildRoute строит маршрут между точками
func (s *GISService) BuildRoute(points []RoutePoint, routeType string) (*RouteResult, error) {
	// Формируем строку точек
	pointsStr := ""
	for i, p := range points {
		if i > 0 {
			pointsStr += ";"
		}
		pointsStr += fmt.Sprintf("%f,%f", p.Lon, p.Lat)
	}

	url := fmt.Sprintf("%s/routing/7.0.0/global?key=%s&type=%s&points=%s&output=json",
		s.baseURL, s.apiKey, routeType, pointsStr)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call routing API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("routing API error: %s", string(body))
	}

	var result RouteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode route response: %w", err)
	}

	if len(result.Result) == 0 {
		return nil, fmt.Errorf("no route found")
	}

	return &result.Result[0], nil
}

// Структуры для 2GIS API

type Place struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address struct {
		Name string `json:"name"`
	} `json:"address"`
	Point struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"point"`
}

type PlacesResponse struct {
	Meta struct {
		Code int `json:"code"`
	} `json:"meta"`
	Result struct {
		Items []Place `json:"items"`
		Total int     `json:"total"`
	} `json:"result"`
}

type RoutePoint struct {
	Lat float64
	Lon float64
}

type RouteResult struct {
	TotalDuration int     `json:"total_duration"` // секунды
	TotalDistance int     `json:"total_distance"` // метры
	Type          string  `json:"type"`
	Geometry      [][]float64 `json:"geometry"`
}

type RouteResponse struct {
	Result []RouteResult `json:"result"`
}

