package api

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dimmy-kor/audioguid/internal/models"
	"github.com/dimmy-kor/audioguid/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RouteHandler struct {
	gisService     *services.GISService
	poiService     *services.POIService
	contentService *services.ContentService
	ttsService     *services.TTSService
	db             *gorm.DB
}

func NewRouteHandler(
	gis *services.GISService,
	poi *services.POIService,
	content *services.ContentService,
	tts *services.TTSService,
	db *gorm.DB,
) *RouteHandler {
	return &RouteHandler{
		gisService:     gis,
		poiService:     poi,
		contentService: content,
		ttsService:     tts,
		db:             db,
	}
}

// GenerateRoute создает новый маршрут
// @Summary      Создать маршрут (асинхронно)
// @Description  Создает маршрут с точками интереса. Аудио генерируется асинхронно в фоне. Поддерживает 3 режима: автопоиск, конкретные POI ID, свои места.
// @Tags         routes
// @Accept       json
// @Produce      json
// @Param        request body models.RouteRequest true "Параметры маршрута (можно указать poi_ids или custom_pois)"
// @Success      200 {object} models.RouteResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /routes/generate [post]
func (h *RouteHandler) GenerateRoute(c *gin.Context) {
	var req models.RouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем точки интереса
	var pois []models.POI
	var err error

	// Вариант 1: Конкретные POI ID из базы
	if len(req.POIIDs) > 0 {
		pois, err = h.getPOIsByIDs(req.POIIDs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if len(req.CustomPOIs) > 0 {
		// Вариант 2: Пользовательские места (создаем новые POI)
		pois, err = h.createCustomPOIs(req.CustomPOIs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Вариант 3: Автоматический поиск по критериям
		radius := float64(req.DurationMinutes * 83)
		if radius > 5000 {
			radius = 5000
		}

		pois, err = h.poiService.FindNearby(
			req.StartPoint.Lat,
			req.StartPoint.Lon,
			radius,
			req.Epochs,
			req.Interests,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find POIs"})
			return
		}

		// Ограничиваем количество точек только для автопоиска
		maxPoints := req.MaxWaypoints
		if maxPoints == 0 || maxPoints > 10 {
			maxPoints = 5
		}
		if len(pois) > maxPoints {
			pois = pois[:maxPoints]
		}
	}

	if len(pois) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No POIs found matching criteria"})
		return
	}

	// Создаем маршрут в БД
	route := &models.Route{
		Name:              generateRouteName(req.Epochs, req.Interests),
		Description:       "Автоматически сгенерированный маршрут",
		TotalDistance:     0,
		EstimatedDuration: req.DurationMinutes,
		Epochs:            req.Epochs,
		Categories:        req.Interests,
	}

	if err := h.db.Create(route).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create route"})
		return
	}

	// Создаем waypoints
	waypoints := make([]models.Waypoint, len(pois))
	for i, poi := range pois {
		waypoint := models.Waypoint{
			RouteID: route.ID,
			POIID:   poi.ID,
			Order:   i + 1,
		}
		waypoints[i] = waypoint
	}

	if err := h.db.Create(&waypoints).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create waypoints"})
		return
	}

	// Загружаем POI данные для каждой waypoint
	for i := range waypoints {
		waypoints[i].POI = pois[i]
	}

	// Генерируем контент асинхронно (в реальности - через очередь)
	go h.generateContentForWaypoints(waypoints)

	// Формируем ответ
	response := h.formatRouteResponse(route, waypoints)

	c.JSON(http.StatusOK, response)
}

// GetRoute получает детали маршрута
// @Summary      Получить маршрут
// @Description  Возвращает детали маршрута с точками и контентом
// @Tags         routes
// @Produce      json
// @Param        route_id path string true "Route ID"
// @Success      200 {object} models.RouteResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /routes/{route_id} [get]
func (h *RouteHandler) GetRoute(c *gin.Context) {
	routeID := c.Param("route_id")

	uid, err := uuid.Parse(routeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
		return
	}

	var route models.Route
	if err := h.db.Preload("Waypoints.POI").Preload("Waypoints.Content").
		First(&route, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	response := h.formatRouteResponse(&route, route.Waypoints)
	c.JSON(http.StatusOK, response)
}

// GetPOIs получает список мест интереса
// @Summary      Список POI
// @Description  Возвращает список всех мест интереса с фильтрацией
// @Tags         pois
// @Produce      json
// @Param        epoch query string false "Фильтр по эпохе (medieval, imperial, soviet, modern)"
// @Param        category query string false "Фильтр по категории (architecture, history, culture, religion, art)"
// @Success      200 {array} models.POI
// @Failure      500 {object} map[string]string
// @Router       /pois [get]
func (h *RouteHandler) GetPOIs(c *gin.Context) {
	var pois []models.POI
	query := h.db.Model(&models.POI{})

	// Фильтры
	if epoch := c.Query("epoch"); epoch != "" {
		query = query.Where("epoch = ?", epoch)
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&pois).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch POIs"})
		return
	}

	c.JSON(http.StatusOK, pois)
}

// GetPOI получает детали места
// @Summary      Получить POI
// @Description  Возвращает детали конкретного места интереса
// @Tags         pois
// @Produce      json
// @Param        poi_id path string true "POI ID"
// @Success      200 {object} models.POI
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /pois/{poi_id} [get]
func (h *RouteHandler) GetPOI(c *gin.Context) {
	poiID := c.Param("poi_id")

	uid, err := uuid.Parse(poiID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid POI ID"})
		return
	}

	var poi models.POI
	if err := h.db.First(&poi, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "POI not found"})
		return
	}

	c.JSON(http.StatusOK, poi)
}

// GenerateRouteWithAudio создает маршрут и сразу генерирует аудио (синхронно)
// @Summary      Создать маршрут и получить MP3 (синхронно)
// @Description  Создает маршрут, генерирует аудио для всех точек и возвращает готовый MP3 файл. Занимает 2-5 минут. Поддерживает 3 режима: автопоиск, конкретные POI ID, свои места с координатами.
// @Tags         routes
// @Accept       json
// @Produce      audio/mpeg
// @Param        request body models.RouteRequest true "Параметры маршрута (можно указать poi_ids или custom_pois)"
// @Success      200 {file} audio/mpeg "MP3 файл с аудиогидом"
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /routes/generate-audio [post]
func (h *RouteHandler) GenerateRouteWithAudio(c *gin.Context) {
	var req models.RouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем точки интереса (используем ту же логику)
	var pois []models.POI
	var err error

	if len(req.POIIDs) > 0 {
		pois, err = h.getPOIsByIDs(req.POIIDs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if len(req.CustomPOIs) > 0 {
		pois, err = h.createCustomPOIs(req.CustomPOIs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		radius := float64(req.DurationMinutes * 83)
		if radius > 5000 {
			radius = 5000
		}

		pois, err = h.poiService.FindNearby(
			req.StartPoint.Lat,
			req.StartPoint.Lon,
			radius,
			req.Epochs,
			req.Interests,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find POIs"})
			return
		}

		maxPoints := req.MaxWaypoints
		if maxPoints == 0 || maxPoints > 10 {
			maxPoints = 5
		}
		if len(pois) > maxPoints {
			pois = pois[:maxPoints]
		}
	}

	if len(pois) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No POIs found matching criteria"})
		return
	}

	// Создаем маршрут в БД
	route := &models.Route{
		Name:              generateRouteName(req.Epochs, req.Interests),
		Description:       "Автоматически сгенерированный маршрут с аудио",
		TotalDistance:     0,
		EstimatedDuration: req.DurationMinutes,
		Epochs:            req.Epochs,
		Categories:        req.Interests,
	}

	if err := h.db.Create(route).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create route"})
		return
	}

	// Создаем waypoints
	waypoints := make([]models.Waypoint, len(pois))
	for i, poi := range pois {
		waypoint := models.Waypoint{
			RouteID: route.ID,
			POIID:   poi.ID,
			Order:   i + 1,
		}
		waypoints[i] = waypoint
	}

	if err := h.db.Create(&waypoints).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create waypoints"})
		return
	}

	// Загружаем POI данные
	for i := range waypoints {
		waypoints[i].POI = pois[i]
	}

	// СИНХРОННАЯ генерация контента
	var audioFiles []string
	for _, waypoint := range waypoints {
		// Генерируем текст
		text, err := h.contentService.GenerateDescription(waypoint.POI)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to generate content for %s: %v", waypoint.POI.Name, err),
			})
			return
		}

		// Генерируем аудио
		audioPath, duration, err := h.ttsService.GenerateAudio(text, waypoint.ID.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to generate audio for %s: %v", waypoint.POI.Name, err),
			})
			return
		}

		// Сохраняем контент
		content := models.Content{
			WaypointID: waypoint.ID,
			Text:       text,
			AudioPath:  audioPath,
			AudioURL:   fmt.Sprintf("/api/audio/%s", waypoint.ID.String()),
			Duration:   duration,
			Photos:     waypoint.POI.Photos,
			Generated:  true,
		}

		if err := h.db.Create(&content).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save content"})
			return
		}

		audioFiles = append(audioFiles, audioPath)
	}

	// Объединяем аудио файлы
	if len(audioFiles) == 1 {
		// Один файл - отдаем напрямую
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=route_%s.mp3", route.ID.String()[:8]))
		c.File(audioFiles[0])
		return
	}

	// Несколько файлов - объединяем
	mergedPath, err := h.mergeAudioFiles(audioFiles, route.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to merge audio files",
			"details": err.Error(),
		})
		return
	}

	// Отдаем объединенный файл
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=route_%s.mp3", route.ID.String()[:8]))
	c.File(mergedPath)
}

// GetAudio получает аудио файл для одной точки
// @Summary      Получить аудио точки
// @Description  Возвращает MP3 файл для конкретной точки маршрута
// @Tags         audio
// @Produce      audio/mpeg
// @Param        waypoint_id path string true "Waypoint ID"
// @Success      200 {file} audio/mpeg "MP3 файл"
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /audio/{waypoint_id} [get]
func (h *RouteHandler) GetAudio(c *gin.Context) {
	waypointID := c.Param("waypoint_id")

	uid, err := uuid.Parse(waypointID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid waypoint ID"})
		return
	}

	var content models.Content
	if err := h.db.First(&content, "waypoint_id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audio not found"})
		return
	}

	if content.AudioPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audio not generated yet"})
		return
	}

	c.File(content.AudioPath)
}

// GetRouteAudio получает объединенный аудиофайл для всего маршрута
// @Summary      Получить аудиогид маршрута
// @Description  Возвращает MP3 файл с аудио для всех точек маршрута
// @Tags         routes
// @Produce      audio/mpeg
// @Param        route_id path string true "Route ID"
// @Success      200 {file} audio/mpeg "MP3 файл"
// @Failure      206 {object} map[string]interface{} "Частично готово"
// @Failure      404 {object} map[string]interface{} "Не готово"
// @Router       /routes/{route_id}/audio [get]
func (h *RouteHandler) GetRouteAudio(c *gin.Context) {
	routeID := c.Param("route_id")

	uid, err := uuid.Parse(routeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
		return
	}

	// Получаем маршрут с waypoints
	var route models.Route
	if err := h.db.Preload("Waypoints").First(&route, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	if len(route.Waypoints) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No waypoints in route"})
		return
	}

	// Получаем все аудио файлы
	var audioFiles []string
	var notReady []string

	for _, wp := range route.Waypoints {
		var content models.Content
		err := h.db.First(&content, "waypoint_id = ?", wp.ID).Error

		if err != nil || content.AudioPath == "" {
			notReady = append(notReady, wp.POI.Name)
			continue
		}

		audioFiles = append(audioFiles, content.AudioPath)
	}

	// Если ни одно аудио не готово
	if len(audioFiles) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "No audio files generated yet",
			"pending": notReady,
			"message": "Audio is being generated. Please try again in 1-2 minutes",
		})
		return
	}

	// Если есть частично готовые
	if len(notReady) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"error":   "Some audio files are not ready yet",
			"ready":   len(audioFiles),
			"pending": notReady,
			"message": fmt.Sprintf("Only %d of %d audio files are ready", len(audioFiles), len(route.Waypoints)),
		})
		return
	}

	// Если только один файл - отдаем его напрямую
	if len(audioFiles) == 1 {
		c.File(audioFiles[0])
		return
	}

	// Объединяем несколько MP3 файлов
	mergedPath, err := h.mergeAudioFiles(audioFiles, routeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to merge audio files",
			"details": err.Error(),
		})
		return
	}

	// Отдаем объединенный файл
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=route_%s.mp3", routeID[:8]))
	c.File(mergedPath)
}

// mergeAudioFiles объединяет несколько MP3 файлов в один
func (h *RouteHandler) mergeAudioFiles(files []string, routeID string) (string, error) {
	// Создаем путь для объединенного файла
	mergedPath := fmt.Sprintf("./audio/route_%s_full.mp3", routeID)

	// Создаем выходной файл
	outFile, err := os.Create(mergedPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Копируем все файлы последовательно
	for _, filePath := range files {
		inFile, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to open %s: %w", filePath, err)
		}

		_, err = io.Copy(outFile, inFile)
		inFile.Close()

		if err != nil {
			return "", fmt.Errorf("failed to copy %s: %w", filePath, err)
		}
	}

	return mergedPath, nil
}

// Вспомогательные функции

func (h *RouteHandler) generateContentForWaypoints(waypoints []models.Waypoint) {
	for _, waypoint := range waypoints {
		// Генерируем текст
		text, err := h.contentService.GenerateDescription(waypoint.POI)
		if err != nil {
			continue
		}

		// Генерируем аудио
		audioPath, duration, err := h.ttsService.GenerateAudio(text, waypoint.ID.String())
		if err != nil {
			continue
		}

		// Сохраняем контент
		content := models.Content{
			WaypointID: waypoint.ID,
			Text:       text,
			AudioPath:  audioPath,
			AudioURL:   fmt.Sprintf("/api/audio/%s", waypoint.ID.String()),
			Duration:   duration,
			Photos:     waypoint.POI.Photos,
			Generated:  true,
		}

		h.db.Create(&content)
	}
}

func (h *RouteHandler) formatRouteResponse(route *models.Route, waypoints []models.Waypoint) models.RouteResponse {
	waypointDetails := make([]models.WaypointDetails, len(waypoints))

	for i, wp := range waypoints {
		detail := models.WaypointDetails{
			ID:          wp.ID.String(),
			Name:        wp.POI.Name,
			Description: wp.POI.Description,
			Coordinates: models.Point{
				Lat: wp.POI.Latitude,
				Lon: wp.POI.Longitude,
			},
			Epoch:    wp.POI.Epoch,
			Category: wp.POI.Category,
			Order:    wp.Order,
		}

		if wp.Content != nil {
			detail.Content = &models.ContentDetails{
				Text:            wp.Content.Text,
				AudioURL:        wp.Content.AudioURL,
				DurationSeconds: wp.Content.Duration,
				Photos:          wp.Content.Photos,
			}
		}

		waypointDetails[i] = detail
	}

	return models.RouteResponse{
		RouteID:           route.ID.String(),
		Name:              route.Name,
		TotalDistance:     route.TotalDistance,
		EstimatedDuration: route.EstimatedDuration,
		Waypoints:         waypointDetails,
	}
}

func generateRouteName(epochs []string, interests []string) string {
	name := "Экскурсия"

	if len(epochs) > 0 {
		epochNames := map[string]string{
			"medieval": "Средневековая",
			"imperial": "Императорская",
			"soviet":   "Советская",
			"modern":   "Современная",
		}
		if epochName, ok := epochNames[epochs[0]]; ok {
			name = epochName + " Москва"
		}
	}

	if len(interests) > 0 {
		interestNames := map[string]string{
			"architecture": "Архитектура",
			"history":      "История",
			"culture":      "Культура",
			"religion":     "Религия",
			"art":          "Искусство",
		}
		if interestName, ok := interestNames[interests[0]]; ok {
			name += ": " + interestName
		}
	}

	return name
}

// getPOIsByIDs получает POI по их ID
func (h *RouteHandler) getPOIsByIDs(ids []string) ([]models.POI, error) {
	var pois []models.POI
	
	// Конвертируем строки в UUID
	uuids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uid, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("invalid POI ID: %s", id)
		}
		uuids = append(uuids, uid)
	}
	
	// Получаем POI из базы в том же порядке
	for _, uid := range uuids {
		var poi models.POI
		if err := h.db.First(&poi, "id = ?", uid).Error; err != nil {
			return nil, fmt.Errorf("POI not found: %s", uid)
		}
		pois = append(pois, poi)
	}
	
	return pois, nil
}

// createCustomPOIs создает временные POI из пользовательских данных
func (h *RouteHandler) createCustomPOIs(customPOIs []models.CustomPOI) ([]models.POI, error) {
	pois := make([]models.POI, 0, len(customPOIs))
	
	for _, custom := range customPOIs {
		poi := models.POI{
			Name:        custom.Name,
			Description: custom.Description,
			Latitude:    custom.Latitude,
			Longitude:   custom.Longitude,
			Epoch:       custom.Epoch,
			Category:    custom.Category,
			Importance:  5, // Средняя важность
		}
		
		// Сохраняем в базу
		if err := h.db.Create(&poi).Error; err != nil {
			return nil, fmt.Errorf("failed to create custom POI: %w", err)
		}
		
		pois = append(pois, poi)
	}
	
	return pois, nil
}
