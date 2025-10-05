package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// POI (Point of Interest) - место интереса
type POI struct {
	ID           uuid.UUID               `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string                  `gorm:"not null" json:"name"`
	Description  string                  `gorm:"type:text" json:"description"`
	Latitude     float64                 `gorm:"not null" json:"latitude"`
	Longitude    float64                 `gorm:"not null" json:"longitude"`
	Epoch        string                  `gorm:"index" json:"epoch"`          // medieval, imperial, soviet, modern
	Category     string                  `gorm:"index" json:"category"`       // architecture, history, culture, religion, art
	Importance   int                     `gorm:"default:5" json:"importance"` // 1-10
	YearBuilt    int                     `json:"year_built,omitempty"`
	Architect    string                  `json:"architect,omitempty"`
	Style        string                  `json:"style,omitempty"`
	Photos       pq.StringArray          `gorm:"type:text[]" json:"photos" swaggertype:"array,string"`
	WikipediaURL string                  `json:"wikipedia_url,omitempty"`
	Metadata     *map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty" swaggerignore:"true"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

// TableName overrides the default table name (optional, pois is default)
// func (POI) TableName() string {
// 	return "pois"
// }

// Route - маршрут
type Route struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name              string
	Description       string         `gorm:"type:text"`
	TotalDistance     float64        // метры
	EstimatedDuration int            // минуты
	Waypoints         []Waypoint     `gorm:"foreignKey:RouteID"`
	Epochs            pq.StringArray `gorm:"type:text[]"`
	Categories        pq.StringArray `gorm:"type:text[]"`
	CreatedAt         time.Time
}

// Waypoint - точка на маршруте
type Waypoint struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RouteID   uuid.UUID `gorm:"type:uuid;not null;index"`
	POIID     uuid.UUID `gorm:"type:uuid;not null"`
	POI       POI       `gorm:"foreignKey:POIID"`
	Order     int       `gorm:"not null"`
	Content   *Content  `gorm:"foreignKey:WaypointID"`
	CreatedAt time.Time
}

// Content - контент для точки (текст + аудио)
type Content struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WaypointID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Text       string    `gorm:"type:text;not null"`
	AudioURL   string
	AudioPath  string
	Duration   int            // секунды
	Photos     pq.StringArray `gorm:"type:text[]"`
	Generated  bool           `gorm:"default:false"`
	CreatedAt  time.Time
}

// RouteRequest - запрос на создание маршрута
type RouteRequest struct {
	StartPoint      Point    `json:"start_point" binding:"required"`
	DurationMinutes int      `json:"duration_minutes" binding:"required,min=15,max=180"`
	Epochs          []string `json:"epochs"`
	Interests       []string `json:"interests"`
	MaxWaypoints    int      `json:"max_waypoints"`
	// Конкретные места для маршрута (опционально)
	POIIDs     []string    `json:"poi_ids,omitempty"`     // UUID мест из базы
	CustomPOIs []CustomPOI `json:"custom_pois,omitempty"` // Свои места с координатами
}

// CustomPOI - пользовательское место для маршрута
type CustomPOI struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Epoch       string  `json:"epoch,omitempty"`
	Category    string  `json:"category,omitempty"`
}

// Point - географическая точка
type Point struct {
	Lat float64 `json:"lat" binding:"required"`
	Lon float64 `json:"lon" binding:"required"`
}

// RouteResponse - ответ с маршрутом
type RouteResponse struct {
	RouteID           string            `json:"route_id"`
	Name              string            `json:"name"`
	TotalDistance     float64           `json:"total_distance"`
	EstimatedDuration int               `json:"estimated_duration"`
	Waypoints         []WaypointDetails `json:"waypoints"`
}

// WaypointDetails - детали точки маршрута
type WaypointDetails struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Coordinates Point           `json:"coordinates"`
	Epoch       string          `json:"epoch"`
	Category    string          `json:"category"`
	Order       int             `json:"order"`
	Content     *ContentDetails `json:"content,omitempty"`
}

// ContentDetails - детали контента
type ContentDetails struct {
	Text            string   `json:"text"`
	AudioURL        string   `json:"audio_url,omitempty"`
	DurationSeconds int      `json:"duration_seconds"`
	Photos          []string `json:"photos"`
}
