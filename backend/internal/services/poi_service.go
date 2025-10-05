package services

import (
	"fmt"
	"math"

	"github.com/dimmy-kor/audioguid/internal/models"
	"gorm.io/gorm"
)

// POIService работает с местами интереса
type POIService struct {
	db *gorm.DB
}

func NewPOIService(db *gorm.DB) *POIService {
	return &POIService{db: db}
}

// FindNearby находит места рядом с точкой
func (s *POIService) FindNearby(
	lat, lon, radiusMeters float64,
	epochs, categories []string,
) ([]models.POI, error) {
	var pois []models.POI

	query := s.db.Model(&models.POI{})

	// Фильтр по эпохе
	if len(epochs) > 0 {
		query = query.Where("epoch IN ?", epochs)
	}

	// Фильтр по категориям
	if len(categories) > 0 {
		query = query.Where("category IN ?", categories)
	}

	// Находим все POI (упрощенно, без PostGIS)
	if err := query.Find(&pois).Error; err != nil {
		return nil, fmt.Errorf("failed to query POIs: %w", err)
	}

	// Фильтруем по расстоянию
	var nearby []models.POI
	for _, poi := range pois {
		distance := calculateDistance(lat, lon, poi.Latitude, poi.Longitude)
		if distance <= radiusMeters {
			nearby = append(nearby, poi)
		}
	}

	// Сортируем по важности
	sortPOIsByImportance(nearby)

	return nearby, nil
}

// GetByID получает POI по ID
func (s *POIService) GetByID(id string) (*models.POI, error) {
	var poi models.POI
	if err := s.db.First(&poi, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &poi, nil
}

// Create создает новый POI
func (s *POIService) Create(poi *models.POI) error {
	return s.db.Create(poi).Error
}

// Вспомогательные функции

// calculateDistance вычисляет расстояние между двумя точками (формула Haversine)
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000.0 // метры

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// sortPOIsByImportance сортирует POI по важности (от большего к меньшему)
func sortPOIsByImportance(pois []models.POI) {
	for i := 0; i < len(pois)-1; i++ {
		for j := i + 1; j < len(pois); j++ {
			if pois[i].Importance < pois[j].Importance {
				pois[i], pois[j] = pois[j], pois[i]
			}
		}
	}
}

