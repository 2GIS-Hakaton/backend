# 📋 Сводка проекта Audio Guide API

## ✅ Что сделано

### 🎯 Основная функциональность

1. **Синхронный API** - `POST /api/routes/generate-audio`
   - Один запрос → готовый MP3 файл
   - Автоматическая генерация контента
   - Объединение аудио всех точек

2. **Асинхронный API** - `POST /api/routes/generate`
   - Создание маршрута
   - Фоновая генерация контента
   - Гибкое управление

3. **Swagger UI** - Интерактивная документация
   - Все эндпоинты задокументированы
   - Тестирование в браузере
   - OpenAPI спецификация

4. **База данных**
   - PostgreSQL с PostGIS
   - 12 тестовых POI
   - 4 таблицы (pois, routes, waypoints, contents)

### 🔧 Исправления

1. **Порт БД** - Обновлен с 5432 на 30101
2. **Скрипт импорта** - Исправлено название таблицы
3. **TTS сервис** - Исправлено кодирование URL параметров
4. **Документация** - Структурирована и объединена

### 📚 Документация

1. **Главные файлы**
   - `README.md` - Обзор проекта
   - `GETTING_STARTED.md` - Полное руководство по запуску
   - `FINAL_SETUP.md` - Финальная настройка

2. **Структурированная документация** (`docs/`)
   - `api/` - API документация
   - `guides/` - Руководства
   - `reference/` - Справочники
   - `diagrams/` - Диаграммы

3. **Удалены дубликаты**
   - DOCUMENTATION_SUMMARY.md
   - QUICK_FIX_PORT.md
   - START_HERE.md
   - SWAGGER_SETUP.md

---

## 📁 Структура проекта

```
2gis_hackaton/
├── README.md                    # Главный обзор
├── GETTING_STARTED.md           # Руководство по запуску
├── FINAL_SETUP.md               # Финальная настройка
├── PROJECT_SUMMARY.md           # Эта сводка
├── .gitignore                   # Обновлен
│
├── backend/                     # Backend приложение
│   ├── cmd/server/main.go      # Точка входа + Swagger
│   ├── internal/
│   │   ├── api/                # HTTP handlers
│   │   ├── services/           # Бизнес-логика
│   │   ├── models/             # Модели данных
│   │   └── database/           # БД
│   ├── docs/                   # Swagger (генерируется)
│   └── config/config.yaml      # Конфигурация
│
├── docs/                        # Документация
│   ├── README.md               # Центр документации
│   ├── api/                    # API docs
│   ├── guides/                 # Руководства
│   ├── reference/              # Справочники
│   └── diagrams/               # Диаграммы
│
├── scripts/                     # Вспомогательные скрипты
│   ├── import_sample_pois.py   # Импорт POI
│   └── create_tables.sql       # SQL схема
│
└── tests/                       # Тестовые скрипты
    ├── generate_audio_route.sh
    ├── test_route.sh
    └── test_audio_guide.sh
```

---

## 🎯 API Эндпоинты

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/health` | Проверка работоспособности |
| `GET` | `/api/pois` | Список мест интереса |
| `GET` | `/api/pois/:id` | Детали места |
| `POST` | `/api/routes/generate` | Создать маршрут (асинхронно) |
| `POST` | `/api/routes/generate-audio` | **Создать + MP3 (синхронно)** ⚡ |
| `GET` | `/api/routes/:id` | Получить маршрут |
| `GET` | `/api/routes/:id/audio` | Получить аудиогид |
| `GET` | `/api/audio/:waypoint_id` | Получить аудио точки |
| `GET` | `/swagger/*` | Swagger UI |

---

## 📊 Статистика

- **Строк кода:** ~3000+
- **API эндпоинтов:** 9
- **POI в базе:** 12
- **Документация:** 15+ файлов
- **Swagger аннотаций:** ~70
- **Тестовых скриптов:** 3

---

## 🔧 Технологии

| Компонент | Технология |
|-----------|-----------|
| Backend | Go 1.21+ |
| Web Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL 16 + PostGIS |
| API Docs | Swagger/OpenAPI 3.0 |
| Maps | 2GIS API |
| Text Generation | YandexGPT |
| Audio | Yandex SpeechKit |

---

## 🚀 Быстрый старт

```bash
# 1. Установить переменные
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
export YANDEX_API_KEY="your-key"
export YANDEX_FOLDER_ID="your-folder"

# 2. Запустить сервер
cd backend && go run cmd/server/main.go

# 3. Сгенерировать аудиогид
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -d '{"start_point": {"lat": 55.7558, "lon": 37.6173}, "duration_minutes": 90, "max_waypoints": 3}' \
  -H "Content-Type: application/json" -o guide.mp3
```

---

## 📚 Ключевые файлы документации

### Для начала работы
1. **[GETTING_STARTED.md](GETTING_STARTED.md)** - Полное руководство
2. **[FINAL_SETUP.md](FINAL_SETUP.md)** - Исправления и настройка

### Для разработчиков
3. **[docs/README.md](docs/README.md)** - Центр документации
4. **[docs/api/SYNC_AUDIO_API.md](docs/api/SYNC_AUDIO_API.md)** - Синхронный API
5. **[docs/api/AUDIO_GUIDE_API.md](docs/api/AUDIO_GUIDE_API.md)** - Асинхронный API

### Для изучения
6. **[docs/reference/ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md](docs/reference/ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md)** - Архитектура
7. **[docs/diagrams/ДИАГРАММЫ.md](docs/diagrams/ДИАГРАММЫ.md)** - Визуальные схемы

---

## ✨ Особенности

### Синхронный API
- ✅ Один запрос → MP3 файл
- ✅ Автоматическая генерация контента
- ✅ Простая интеграция
- ⏱️ Время: 2-5 минут

### Асинхронный API
- ✅ Быстрый ответ (route_id)
- ✅ Проверка статуса
- ✅ Гибкое управление
- ⏱️ Время: мгновенно + фоновая обработка

### Swagger UI
- ✅ Интерактивное тестирование
- ✅ Автогенерация из кода
- ✅ OpenAPI спецификация
- 🔗 http://localhost:8080/swagger/index.html

---

## 🎓 Примеры использования

### cURL

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 90,
    "epochs": ["soviet", "medieval"],
    "interests": ["architecture", "history"],
    "max_waypoints": 3
  }' \
  -o route_guide.mp3
```

### JavaScript

```javascript
const response = await fetch('http://localhost:8080/api/routes/generate-audio', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    start_point: { lat: 55.7558, lon: 37.6173 },
    duration_minutes: 90,
    max_waypoints: 3
  })
});

const blob = await response.blob();
const url = URL.createObjectURL(blob);
const audio = new Audio(url);
audio.play();
```

### Python

```python
import requests

response = requests.post(
    'http://localhost:8080/api/routes/generate-audio',
    json={
        'start_point': {'lat': 55.7558, 'lon': 37.6173},
        'duration_minutes': 90,
        'max_waypoints': 3
    }
)

with open('route_guide.mp3', 'wb') as f:
    f.write(response.content)
```

---

## 🔗 Полезные ссылки

### Проект
- [README.md](README.md) - Главная страница
- [GETTING_STARTED.md](GETTING_STARTED.md) - Руководство
- [Swagger UI](http://localhost:8080/swagger/index.html) - API docs

### Внешние ресурсы
- [2GIS API](https://docs.2gis.com/)
- [YandexGPT](https://cloud.yandex.ru/docs/yandexgpt/)
- [Yandex SpeechKit](https://cloud.yandex.ru/docs/speechkit/)

---

## 🎉 Готово к использованию!

Проект полностью настроен и готов к работе:

- ✅ Swagger UI работает
- ✅ API эндпоинты задокументированы
- ✅ База данных настроена
- ✅ POI импортированы
- ✅ Документация структурирована
- ✅ .gitignore обновлен

**Начните с:** [GETTING_STARTED.md](GETTING_STARTED.md)

---

<div align="center">
  <b>Создано с ❤️ для хакатона 2GIS</b>
  <br>
  <sub>Умные аудиомаршруты для незабываемых прогулок</sub>
</div>
