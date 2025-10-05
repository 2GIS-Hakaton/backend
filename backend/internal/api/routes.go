package api

import (
	"github.com/dimmy-kor/audioguid/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes регистрирует все API эндпоинты
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// Инициализация сервисов
	gisService := services.NewGISService()
	poiService := services.NewPOIService(db)
	contentService := services.NewContentService()
	ttsService := services.NewTTSService()

	// Хендлеры
	routeHandler := NewRouteHandler(gisService, poiService, contentService, ttsService, db)

	// API группа
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Audioguid API is running",
			})
		})

		// Маршруты
		routes := api.Group("/routes")
		{
			routes.POST("/generate", routeHandler.GenerateRoute)
			routes.POST("/generate-audio", routeHandler.GenerateRouteWithAudio)
			routes.GET("/:route_id", routeHandler.GetRoute)
			routes.GET("/:route_id/audio", routeHandler.GetRouteAudio)
		}

		// Места интереса
		pois := api.Group("/pois")
		{
			pois.GET("", routeHandler.GetPOIs)
			pois.GET("/:poi_id", routeHandler.GetPOI)
		}

		// Аудио
		api.GET("/audio/:waypoint_id", routeHandler.GetAudio)
	}
}
