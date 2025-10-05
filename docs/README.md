# 📚 Документация Audio Guide API

Добро пожаловать в документацию проекта Audio Guide!

---

## 🚀 Быстрый старт

### Для разработчиков
- **[Шпаргалка](guides/ШПАРГАЛКА.md)** - Все основные команды на одной странице
- **[Быстрый старт аудиогида](guides/QUICK_START_AUDIO.md)** - Начните за 3 минуты

### Для пользователей API
- **[Swagger UI](http://localhost:8080/swagger/index.html)** - Интерактивная документация API (требуется запущенный сервер)

---

## 📖 Структура документации

### 📡 API Documentation (`api/`)

Документация по REST API эндпоинтам:

- **[Асинхронный API аудиогида](api/AUDIO_GUIDE_API.md)**
  - Создание маршрута
  - Получение аудио после генерации
  - Проверка статуса
  - Примеры интеграции (JS, Python, Swift)

- **[Синхронный API аудиогида](api/SYNC_AUDIO_API.md)** ⚡
  - Один запрос → готовый MP3
  - Простая интеграция
  - Примеры использования

### 📘 Guides (`guides/`)

Практические руководства:

- **[Шпаргалка](guides/ШПАРГАЛКА.md)**
  - Команды для запуска
  - Примеры запросов
  - Troubleshooting

- **[Быстрый старт аудиогида](guides/QUICK_START_AUDIO.md)**
  - Установка
  - Первый запрос
  - Получение аудио

### 📚 Reference (`reference/`)

Справочная информация:

- **[Полная документация](reference/ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md)**
  - Архитектура системы
  - Модели данных
  - Потоки данных
  - База данных

- **[Changelog аудиогида](reference/CHANGELOG_AUDIO_GUIDE.md)**
  - История изменений
  - Технические детали
  - Roadmap

- **[Сводка изменений](reference/SUMMARY.md)**
  - Краткая сводка функций
  - Статистика изменений

- **[Сводка синхронного API](reference/SYNC_AUDIO_SUMMARY.md)**
  - Краткое описание синхронного API

### 📊 Diagrams (`diagrams/`)

Визуальные диаграммы:

- **[Диаграммы](diagrams/ДИАГРАММЫ.md)**
  - Архитектура системы
  - Потоки данных
  - Схема базы данных
  - Последовательности операций

---

## 🎯 Быстрая навигация

### По функциональности

#### Создание маршрута
- [Асинхронный способ](api/AUDIO_GUIDE_API.md#эндпоинт) - `POST /api/routes/generate`
- [Синхронный способ](api/SYNC_AUDIO_API.md#эндпоинт) - `POST /api/routes/generate-audio` ⚡

#### Получение аудио
- [Аудио одной точки](api/AUDIO_GUIDE_API.md#получить-аудио-для-одной-точки) - `GET /api/audio/:waypoint_id`
- [Аудио всего маршрута](api/AUDIO_GUIDE_API.md#получить-полный-аудиогид-маршрута) - `GET /api/routes/:route_id/audio`

#### Работа с POI
- [Список мест](api/AUDIO_GUIDE_API.md) - `GET /api/pois`
- [Детали места](api/AUDIO_GUIDE_API.md) - `GET /api/pois/:poi_id`

### По типу задачи

#### Я хочу начать работу
→ [Шпаргалка](guides/ШПАРГАЛКА.md)  
→ [Быстрый старт](guides/QUICK_START_AUDIO.md)

#### Я хочу интегрировать API
→ [Swagger UI](http://localhost:8080/swagger/index.html)  
→ [Синхронный API](api/SYNC_AUDIO_API.md)  
→ [Асинхронный API](api/AUDIO_GUIDE_API.md)

#### Я хочу понять архитектуру
→ [Полная документация](reference/ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md)  
→ [Диаграммы](diagrams/ДИАГРАММЫ.md)

#### Я столкнулся с проблемой
→ [Troubleshooting в Шпаргалке](guides/ШПАРГАЛКА.md#troubleshooting)  
→ [Troubleshooting в API docs](api/AUDIO_GUIDE_API.md#troubleshooting)

---

## 🔧 Swagger / OpenAPI

### Интерактивная документация

После запуска сервера доступна по адресу:

```
http://localhost:8080/swagger/index.html
```

### Возможности Swagger UI

- ✅ Просмотр всех эндпоинтов
- ✅ Описание параметров и моделей
- ✅ Примеры запросов/ответов
- ✅ Тестирование API прямо в браузере
- ✅ Скачивание OpenAPI спецификации

### Экспорт спецификации

OpenAPI спецификация доступна в форматах:

- **JSON:** `http://localhost:8080/swagger/doc.json`
- **YAML:** См. `backend/docs/swagger.yaml`

---

## 📊 Сравнение API методов

| Характеристика | Асинхронный | Синхронный |
|----------------|-------------|------------|
| **Эндпоинт** | `/routes/generate` | `/routes/generate-audio` |
| **Запросов** | 2-3 | 1 |
| **Время ответа** | Мгновенно | 2-5 минут |
| **Результат** | route_id | MP3 файл |
| **Сложность** | Средняя | Простая |
| **Использование** | Сложные сценарии | Простые сценарии |

---

## 🎓 Примеры использования

### Синхронный API (рекомендуется для начала)

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 60,
    "epochs": ["soviet"],
    "interests": ["architecture"],
    "max_waypoints": 3
  }' \
  -o my_guide.mp3
```

### Асинхронный API

```bash
# 1. Создать маршрут
ROUTE_ID=$(curl -X POST http://localhost:8080/api/routes/generate \
  -H "Content-Type: application/json" \
  -d '{"start_point": {"lat": 55.7558, "lon": 37.6173}, "duration_minutes": 60}' \
  | jq -r '.route_id')

# 2. Подождать 2-3 минуты
sleep 180

# 3. Скачать аудио
curl http://localhost:8080/api/routes/$ROUTE_ID/audio -o guide.mp3
```

---

## 🛠️ Инструменты

### Тестовые скрипты

В корне проекта:

- `./generate_audio_route.sh` - Синхронная генерация аудиогида
- `./test_route.sh` - Полный тест асинхронного API
- `./test_audio_guide.sh <route_id>` - Тест скачивания аудио

### Backend скрипты

В папке `backend/`:

- `./run.sh` - Запуск сервера
- `go run cmd/server/main.go` - Запуск в dev режиме

---

## 📦 Технологии

- **Backend:** Go 1.21+, Gin Framework
- **Database:** PostgreSQL 16 + PostGIS
- **API Documentation:** Swagger/OpenAPI 3.0
- **External APIs:** 2GIS, YandexGPT, Yandex TTS

---

## 🔗 Полезные ссылки

### Внешние ресурсы
- [2GIS API Docs](https://docs.2gis.com/)
- [YandexGPT Docs](https://cloud.yandex.ru/docs/yandexgpt/)
- [Yandex SpeechKit Docs](https://cloud.yandex.ru/docs/speechkit/)

### Внутренние ресурсы
- [Backend README](../backend/README.md)
- [Основной README](../README.md)

---

## 📞 Поддержка

При возникновении вопросов:

1. Проверьте [Шпаргалку](guides/ШПАРГАЛКА.md)
2. Изучите [Swagger UI](http://localhost:8080/swagger/index.html)
3. Посмотрите [Troubleshooting](api/AUDIO_GUIDE_API.md#troubleshooting)

---

## 🎉 Быстрые команды

```bash
# Запустить сервер
cd backend && go run cmd/server/main.go

# Открыть Swagger UI (после запуска сервера)
open http://localhost:8080/swagger/index.html

# Сгенерировать аудиогид
./generate_audio_route.sh

# Проверить здоровье API
curl http://localhost:8080/api/health
```

---

<div align="center">
  <b>Создано с ❤️ для хакатона 2GIS</b>
  <br>
  <sub>Умные аудиомаршруты для незабываемых прогулок</sub>
</div>
