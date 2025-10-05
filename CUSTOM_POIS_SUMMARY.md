# ✨ Новая функция: Свои места для аудиогида

## 🎯 Что добавлено

Теперь API поддерживает **3 способа** создания маршрутов:

1. **Автопоиск** - по критериям (как раньше)
2. **Конкретные POI ID** - выбрать места из базы
3. **Свои места** - указать координаты и описания ⭐ НОВОЕ

---

## 🚀 Быстрый пример

### Свои места

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 90,
    "custom_pois": [
      {
        "name": "Моя любимая кофейня",
        "description": "Уютное место с лучшим кофе",
        "latitude": 55.7558,
        "longitude": 37.6173
      },
      {
        "name": "Парк, где я гуляю",
        "description": "Мое любимое место для пробежек",
        "latitude": 55.7304,
        "longitude": 37.6012
      }
    ]
  }' \
  -o my_personal_guide.mp3
```

### Конкретные места из базы

```bash
# Сначала получить ID
curl http://localhost:8080/api/pois | jq '.[] | {id, name}'

# Затем использовать
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 90,
    "poi_ids": ["uuid-1", "uuid-2", "uuid-3"]
  }' \
  -o route_guide.mp3
```

---

## 📝 Изменения в API

### RouteRequest (обновлен)

Добавлены новые поля:

```json
{
  "start_point": {"lat": 55.7558, "lon": 37.6173},
  "duration_minutes": 90,
  
  // НОВОЕ: Конкретные POI ID
  "poi_ids": ["uuid-1", "uuid-2"],
  
  // НОВОЕ: Свои места
  "custom_pois": [
    {
      "name": "Название",
      "description": "Описание",
      "latitude": 55.7558,
      "longitude": 37.6173,
      "epoch": "modern",      // опционально
      "category": "culture"   // опционально
    }
  ]
}
```

### Приоритет

1. Если указаны `poi_ids` → используются они
2. Иначе если указаны `custom_pois` → используются они
3. Иначе → автоматический поиск по критериям

---

## 🔧 Технические детали

### Новые модели

```go
type RouteRequest struct {
    // ... существующие поля ...
    POIIDs     []string    `json:"poi_ids,omitempty"`
    CustomPOIs []CustomPOI `json:"custom_pois,omitempty"`
}

type CustomPOI struct {
    Name        string  `json:"name" binding:"required"`
    Description string  `json:"description"`
    Latitude    float64 `json:"latitude" binding:"required"`
    Longitude   float64 `json:"longitude" binding:"required"`
    Epoch       string  `json:"epoch,omitempty"`
    Category    string  `json:"category,omitempty"`
}
```

### Новые функции

- `getPOIsByIDs(ids []string)` - получение POI по ID
- `createCustomPOIs(customPOIs []CustomPOI)` - создание новых POI

### Обновленные эндпоинты

- `POST /api/routes/generate` - поддерживает все 3 режима
- `POST /api/routes/generate-audio` - поддерживает все 3 режима

---

## 📚 Документация

**Полное руководство:** [CUSTOM_POIS_GUIDE.md](CUSTOM_POIS_GUIDE.md)

**Содержит:**
- Детальные примеры
- Форматы данных
- Советы по использованию
- Сравнение способов

---

## ✅ Готово к использованию!

Попробуйте создать свой персональный аудиогид прямо сейчас:

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 60,
    "custom_pois": [
      {
        "name": "Ваше место",
        "description": "Расскажите о нем",
        "latitude": 55.7558,
        "longitude": 37.6173
      }
    ]
  }' \
  -o my_guide.mp3
```
