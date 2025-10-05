# 🎧 Умный Аудиогид - Backend

Backend API на Go для генерации умных аудиомаршрутов.

## 🚀 Быстрый старт

### 1. Установка зависимостей

```bash
# Установка Go (если еще не установлен)
brew install go  # macOS
# или скачайте с https://go.dev/dl/

# Установка зависимостей проекта
go mod download
```

### 2. Настройка PostgreSQL

```bash
# Установка PostgreSQL + PostGIS
brew install postgresql postgis  # macOS

# Запуск PostgreSQL
brew services start postgresql

# Создание базы данных
createdb audioguid

# Подключение PostGIS
psql audioguid -c "CREATE EXTENSION IF NOT EXISTS postgis;"
```

### 3. Конфигурация

Скопируйте `.env.example` в `.env` и заполните:

```bash
cp .env.example .env
```

**Важно! Скачайте ключ 2GIS:**
```bash
# Скачать файл 2gis.key по ссылке:
# https://disk.yandex.com/d/w7oHITw-8OTL7Q
# И поместить в ./keys/2gis.key
```

Отредактируйте `.env`:

```env
# 2GIS API (уже настроено)
GAPIS_API_KEY=1bfdf79c-0a09-4ea6-882e-afc09421c8d5
GAPIS_APP_ID=ru.2gishackathon.app06.02
GAPIS_KEY_FILE=./keys/2gis.key

# Yandex Cloud (получить на https://console.cloud.yandex.ru/)
YANDEX_API_KEY=...
YANDEX_FOLDER_ID=...
YANDEX_VOICE=alena

# Yandex Search API (опционально, https://yandex.ru/dev/xml/)
YANDEX_SEARCH_USER=...
YANDEX_SEARCH_KEY=...

# Database
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/audioguid?sslmode=disable

# Server
PORT=8080
```

### 4. Импорт тестовых данных

```bash
# Запустить скрипт импорта POI
cd ../scripts
python3 import_sample_pois.py
```

### 5. Запуск сервера

```bash
go run cmd/server/main.go
```

Сервер запустится на `http://localhost:8080`

---

## 📡 API Endpoints

### Health Check
```bash
GET /api/health
```

### Генерация маршрута
```bash
POST /api/routes/generate
Content-Type: application/json

{
  "start_point": {
    "lat": 55.7558,
    "lon": 37.6173
  },
  "duration_minutes": 60,
  "epochs": ["soviet"],
  "interests": ["architecture", "history"],
  "max_waypoints": 5
}
```

**Ответ:**
```json
{
  "route_id": "uuid",
  "name": "Советская Москва: Архитектура",
  "total_distance": 4200,
  "estimated_duration": 60,
  "waypoints": [
    {
      "id": "uuid",
      "name": "ВДНХ",
      "coordinates": {"lat": 55.8304, "lon": 37.6325},
      "epoch": "soviet",
      "category": "architecture",
      "order": 1,
      "content": {
        "text": "Рассказ о месте...",
        "audio_url": "/api/audio/uuid",
        "duration_seconds": 180
      }
    }
  ]
}
```

### Получить детали маршрута
```bash
GET /api/routes/:route_id
```

### Список мест интереса
```bash
GET /api/pois?epoch=soviet&category=architecture
```

### Получить аудио для одной точки
```bash
GET /api/audio/:waypoint_id
```

**Ответ:** MP3 файл

### Получить полный аудиогид маршрута
```bash
GET /api/routes/:route_id/audio
```

**Ответ:** MP3 файл со всеми точками маршрута, объединенными в один аудиогид

**Возможные статусы:**
- `200 OK` - аудиогид готов, возвращается MP3 файл
- `206 Partial Content` - часть аудио еще генерируется
- `404 Not Found` - аудио еще не сгенерировано

**Пример:**
```bash
# Скачать полный аудиогид
curl http://localhost:8080/api/routes/550e8400-e29b-41d4-a716-446655440000/audio -o route_guide.mp3

# Воспроизвести
afplay route_guide.mp3  # macOS
mpg123 route_guide.mp3  # Linux
```

---

## 🏗 Структура проекта

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Точка входа
├── internal/
│   ├── api/                  # HTTP handlers
│   │   ├── routes.go
│   │   └── route_handler.go
│   ├── services/             # Бизнес-логика
│   │   ├── gis_service.go    # 2GIS API
│   │   ├── poi_service.go    # Места интереса
│   │   ├── content_service.go # OpenAI
│   │   └── tts_service.go    # ElevenLabs
│   ├── models/               # Модели данных
│   │   └── models.go
│   └── database/             # БД
│       └── postgres.go
├── config/
│   └── config.yaml
├── keys/
│   └── 2gis.key             # Скачать!
├── audio/                   # Сгенерированные аудио
├── go.mod
└── .env
```

---

## 🧪 Тестирование

### Тест генерации маршрута

```bash
curl -X POST http://localhost:8080/api/routes/generate \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 60,
    "epochs": ["soviet"],
    "interests": ["architecture"]
  }'
```

### Тест получения аудио

```bash
# Получить route_id из предыдущего запроса
curl http://localhost:8080/api/routes/{route_id}

# Скачать аудио
curl http://localhost:8080/api/audio/{waypoint_id} -o audio.mp3
```

---

## 🔧 Разработка

### Добавление нового POI

```go
poi := models.POI{
    Name:        "ВДНХ",
    Description: "Выставка достижений народного хозяйства",
    Latitude:    55.8304,
    Longitude:   37.6325,
    Epoch:       "soviet",
    Category:    "architecture",
    Importance:  10,
    YearBuilt:   1939,
    Architect:   "Вячеслав Олтаржевский",
    Style:       "Сталинский ампир",
}

db.Create(&poi)
```

### Добавление нового эндпоинта

1. Добавьте handler в `internal/api/route_handler.go`
2. Зарегистрируйте в `internal/api/routes.go`

---

## 📝 TODO для MVP

- [x] Базовая структура проекта
- [x] Модели данных
- [x] API endpoints
- [x] Интеграция с 2GIS API
- [x] Интеграция с OpenAI
- [x] Интеграция с ElevenLabs
- [ ] Импорт тестовых POI
- [ ] Тестирование всех эндпоинтов
- [ ] Docker compose для dev окружения
- [ ] Frontend интеграция

---

## 🐛 Troubleshooting

### PostgreSQL не запускается
```bash
brew services restart postgresql
```

### Ошибка "database does not exist"
```bash
createdb audioguid
psql audioguid -c "CREATE EXTENSION postgis;"
```

### Ошибка "2GIS API key invalid"
Проверьте, что файл `keys/2gis.key` скачан и находится в правильном месте.

### Ошибка "YandexGPT API error"
Проверьте, что API ключ правильный и сервисный аккаунт имеет роли:
- `ai.languageModels.user` (для YandexGPT)
- `ai.speechkit-tts.user` (для SpeechKit)

### Ошибка "Yandex Search API error"
Проверьте user ID и ключ на https://yandex.ru/dev/xml/ (Search API опциональный)

---

## 📚 Документация API

- [2GIS API Docs](https://docs.2gis.com/ru)
- [YandexGPT Docs](https://cloud.yandex.ru/docs/yandexgpt/)
- [Yandex SpeechKit Docs](https://cloud.yandex.ru/docs/speechkit/)
- [Yandex Search API Docs](https://yandex.ru/dev/xml/doc/)

---

**Happy coding! 🚀**

