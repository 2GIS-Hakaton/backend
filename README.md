# 🎧 Умный Аудиогид - Audio Guide API

> Система генерации персонализированных аудиомаршрутов с использованием 2GIS API, YandexGPT и Yandex TTS

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![2GIS](https://img.shields.io/badge/2GIS-API-00C853?style=flat)](https://docs.2gis.com/)
[![Yandex](https://img.shields.io/badge/Yandex-Cloud-FF0000?style=flat)](https://cloud.yandex.ru/)

## 🚀 Быстрый Старт

### За 3 минуты

```bash
# 1. Установить переменные
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
export YANDEX_API_KEY="your-key"
export YANDEX_FOLDER_ID="your-folder"

# 2. Запустить сервер
cd backend
go run cmd/server/main.go

# 3. Сгенерировать аудиогид (в другом терминале)
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -d '{"start_point": {"lat": 55.7558, "lon": 37.6173}, "duration_minutes": 90, "max_waypoints": 3}' \
  -H "Content-Type: application/json" \
  -o route_guide.mp3

# 4. Воспроизвести
afplay route_guide.mp3
```

**📖 Подробная инструкция:** [GETTING_STARTED.md](GETTING_STARTED.md)

---

## 📚 Документация

### 🎯 Начните здесь
- **[GETTING_STARTED.md](GETTING_STARTED.md)** - 🚀 Полное руководство по запуску
- **[FINAL_SETUP.md](FINAL_SETUP.md)** - ✅ Финальная настройка и исправления
- **[Swagger UI](http://localhost:8080/swagger/index.html)** - 🔧 Интерактивная API документация

### 📖 Детальная документация
- **[docs/README.md](docs/README.md)** - Центр документации
- **[docs/api/](docs/api/)** - API документация (синхронный и асинхронный)
- **[docs/guides/](docs/guides/)** - Руководства и шпаргалки
- **[docs/reference/](docs/reference/)** - Справочники и архитектура
- **[docs/diagrams/](docs/diagrams/)** - Визуальные диаграммы

---

## 🌟 Возможности

### ✅ Что уже работает

- **🗺️ Генерация маршрутов** - создание персонализированных пешеходных маршрутов
- **📍 12 POI в базе** - исторические места Москвы (советская эпоха, средневековье)
- **🔍 Умный поиск** - фильтрация по эпохам и категориям
- **📊 PostgreSQL** - надежное хранение данных
- **🌐 REST API** - простой и понятный интерфейс

### ⚡ С Yandex API ключами

- **✍️ Генерация текста** - увлекательные исторические рассказы через YandexGPT
- **🔊 Синтез аудио** - озвучивание текста через Yandex TTS
- **🔎 Поиск информации** - дополнительные факты через Yandex Search

---

## 📡 API Endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/health` | Проверка работоспособности |
| `GET` | `/api/pois` | Список всех мест интереса |
| `GET` | `/api/pois?epoch=soviet` | Фильтрация по эпохе |
| `GET` | `/api/pois/:id` | Детали конкретного места |
| `POST` | `/api/routes/generate` | Создание маршрута (асинхронно) |
| `POST` | `/api/routes/generate-audio` | **⚡ Создать маршрут и сразу получить MP3** |
| `GET` | `/api/routes/:id` | Получение маршрута с контентом |
| `GET` | `/api/routes/:id/audio` | Полный аудиогид маршрута (MP3) |
| `GET` | `/api/audio/:waypoint_id` | Аудио одной точки (MP3) |

---

## 🎯 Примеры Использования

### Создать маршрут "Советская архитектура"

```bash
curl -X POST http://localhost:8080/api/routes/generate \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 60,
    "epochs": ["soviet"],
    "interests": ["architecture"],
    "max_waypoints": 3
  }'
```

**Ответ:**
```json
{
  "route_id": "ea7a7092-8fa1-44c2-a656-55f38d4822e3",
  "name": "Советская Москва: Архитектура",
  "estimated_duration": 60,
  "waypoints": [
    {
      "id": "waypoint-1",
      "name": "Красная площадь",
      "coordinates": {"lat": 55.7539, "lon": 37.6208},
      "epoch": "medieval",
      "order": 1
    },
    {
      "id": "waypoint-2",
      "name": "Кремль",
      "coordinates": {"lat": 55.752, "lon": 37.6175},
      "epoch": "medieval",
      "order": 2
    }
  ]
}
```

### Получить маршрут с контентом

```bash
curl http://localhost:8080/api/routes/ea7a7092-8fa1-44c2-a656-55f38d4822e3
```

### Сгенерировать аудиогид сразу (синхронно) ⚡

```bash
# Один запрос - сразу получаете MP3!
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 60,
    "epochs": ["soviet"],
    "interests": ["architecture"],
    "max_waypoints": 3
  }' \
  -o my_audio_guide.mp3

# Воспроизвести
afplay my_audio_guide.mp3  # macOS
```

### Или скачать аудиогид существующего маршрута

```bash
# Скачать весь аудиогид одним файлом
curl http://localhost:8080/api/routes/ea7a7092-8fa1-44c2-a656-55f38d4822e3/audio -o route_guide.mp3

# Воспроизвести
afplay route_guide.mp3  # macOS
```

### Скачать аудио одной точки

```bash
curl http://localhost:8080/api/audio/waypoint-1 -o audio.mp3
```

---

## 🏗️ Архитектура

```
┌─────────────┐
│   Client    │ HTTP Requests
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────┐
│   GIN Web Server (:8080)        │
│  ┌──────────────────────────┐   │
│  │   Route Handler          │   │
│  └────────┬─────────────────┘   │
└───────────┼─────────────────────┘
            │
    ┌───────┴───────┐
    │               │
    ▼               ▼
┌─────────┐   ┌─────────────┐
│Services │   │  Database   │
│         │   │             │
│POI      │   │ PostgreSQL  │
│GIS      │   │             │
│Content  │   │ - pois      │
│TTS      │   │ - routes    │
└────┬────┘   │ - waypoints │
     │        │ - contents  │
     │        └─────────────┘
     ▼
┌────────────────┐
│ External APIs  │
│ - 2GIS         │
│ - YandexGPT    │
│ - Yandex TTS   │
└────────────────┘
```

**Подробные диаграммы:** См. [ДИАГРАММЫ.md](ДИАГРАММЫ.md)

---

## 📊 База Данных

### Таблицы

- **`pois`** - Места интереса (12 записей)
- **`routes`** - Созданные маршруты
- **`waypoints`** - Точки маршрутов
- **`contents`** - Сгенерированный контент (текст + аудио)

### Связи

```
pois (1) ──< waypoints (N)
routes (1) ──< waypoints (N)
waypoints (1) ──< contents (1)
```

**Схема БД:** См. [ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md - База данных](ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md#-база-данных)

---

## 🔧 Технологии

| Компонент | Технология | Описание |
|-----------|-----------|----------|
| **Backend** | Go 1.21+ | Основной язык |
| **Web Framework** | Gin | HTTP роутер |
| **ORM** | GORM | PostgreSQL ORM |
| **Database** | PostgreSQL 16 | Хранение данных |
| **Maps & Routing** | 2GIS API | Карты и маршруты |
| **Text Generation** | YandexGPT | AI генерация текста |
| **Audio** | Yandex TTS | Text-to-Speech |
| **Search** | Yandex Search | Дополнительная информация |

---

## 📁 Структура Проекта

```
2gis_hackaton/
├── backend/                    # 🎯 Backend приложение
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # Точка входа
│   ├── internal/
│   │   ├── api/               # HTTP handlers
│   │   ├── services/          # Бизнес-логика
│   │   ├── models/            # Модели данных
│   │   └── database/          # БД подключение
│   ├── config/
│   │   └── config.yaml        # Конфигурация
│   ├── audio/                 # Сгенерированные MP3
│   └── go.mod                 # Зависимости
│
├── scripts/                   # 🔧 Вспомогательные скрипты
│   ├── create_tables.sql      # SQL схема
│   └── import_sample_pois.py  # Импорт POI
│
├── ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md     # 📖 Детальная документация
├── ДИАГРАММЫ.md               # 🎨 Визуальные диаграммы
├── ШПАРГАЛКА.md               # ⚡ Быстрая справка
└── README.md                  # 📄 Этот файл
```

---

## 🎓 Как Это Работает

### 1. Пользователь создает запрос

```json
POST /api/routes/generate
{
  "start_point": {"lat": 55.7558, "lon": 37.6173},
  "duration_minutes": 60,
  "epochs": ["soviet"],
  "interests": ["architecture"]
}
```

### 2. Система находит POI

- Вычисляет радиус поиска: `duration × 83 м/мин`
- Ищет места в БД по фильтрам (epoch, category)
- Вычисляет расстояние (формула Haversine)
- Сортирует по важности

### 3. Создает маршрут

- Сохраняет маршрут в БД
- Создает waypoints для каждого POI
- Возвращает route_id клиенту

### 4. Генерирует контент (асинхронно)

Для каждого waypoint:
1. **YandexGPT** генерирует текст (2-3 минуты чтения)
2. **Yandex TTS** создает аудио из текста
3. Сохраняет MP3 в `./audio/`
4. Записывает в таблицу `contents`

### 5. Клиент получает результат

```json
GET /api/routes/{route_id}
{
  "waypoints": [{
    "name": "ВДНХ",
    "content": {
      "text": "Перед вами знаменитая ВДНХ...",
      "audio_url": "/api/audio/waypoint-id",
      "duration_seconds": 180
    }
  }]
}
```

**Полный поток:** См. [ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md - Поток данных](ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md#-полный-поток-данных)

---

## 🔐 Настройка

### Минимальная конфигурация (без Yandex)

```bash
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
```

**Работает:**
- ✅ Генерация маршрутов
- ✅ Поиск POI
- ✅ API endpoints

**Не работает:**
- ❌ Генерация текста
- ❌ Создание аудио

### Полная конфигурация (с Yandex)

```bash
# База данных
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"

# Yandex Cloud
export YANDEX_API_KEY="AQVN..."
export YANDEX_FOLDER_ID="b1g..."
export YANDEX_VOICE="alena"
```

**Как получить Yandex ключи:**
1. Регистрация: https://console.cloud.yandex.ru/
2. Создать Service Account
3. Добавить роли: `ai.languageModels.user`, `ai.speechkit-tts.user`
4. Создать API Key
5. Скопировать Folder ID

---

## 📍 Доступные POI (12)

| Название | Эпоха | Категория | Координаты |
|----------|-------|-----------|-----------|
| ВДНХ | soviet | architecture | 55.8304, 37.6325 |
| Павильон Космос | soviet | history | 55.8283, 37.6308 |
| Фонтан Дружба народов | soviet | art | 55.8277, 37.6319 |
| Рабочий и колхозница | soviet | art | 55.8311, 37.6278 |
| МГУ | soviet | architecture | 55.7033, 37.5297 |
| Останкинская башня | soviet | architecture | 55.8194, 37.6119 |
| Парк Горького | soviet | culture | 55.7304, 37.6012 |
| Гостиница Украина | soviet | architecture | 55.7526, 37.5676 |
| Третьяковская галерея | imperial | art | 55.7415, 37.6206 |
| Красная площадь | medieval | history | 55.7539, 37.6208 |
| Храм Василия Блаженного | medieval | religion | 55.7525, 37.6231 |
| Кремль | medieval | history | 55.7520, 37.6175 |

---

## 🧪 Тестирование

### Автоматический тест

```bash
# Синхронная генерация (один запрос → MP3)
./generate_audio_route.sh

# Полный тест с генерацией маршрута (асинхронно)
./test_route.sh

# Быстрый тест скачивания аудиогида
./test_audio_guide.sh <route_id>
```

### Ручное тестирование

```bash
# 1. Health check
curl http://localhost:8080/api/health

# 2. Список POI
curl http://localhost:8080/api/pois | python3 -m json.tool

# 3. Создать маршрут
curl -X POST http://localhost:8080/api/routes/generate \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 90,
    "epochs": ["soviet", "medieval"],
    "interests": ["architecture", "history"]
  }' | python3 -m json.tool
```

---

## 🐛 Troubleshooting

### Сервер не запускается

```bash
# Проверить порт
lsof -ti:8080 | xargs kill -9

# Проверить БД
PGPASSWORD=changeme psql -h 51.250.86.178 -p 30101 -U nike -d audioguid -c "SELECT 1;"
```

### "No POIs found"

```bash
# Увеличить duration_minutes или расширить фильтры
{
  "duration_minutes": 120,
  "epochs": ["soviet", "medieval", "imperial"],
  "interests": ["architecture", "history", "art"]
}
```

### Контент не генерируется

- Проверьте наличие Yandex API ключей
- Подождите 1-2 минуты (генерация асинхронная)
- Проверьте логи сервера

**Больше решений:** См. [ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md - Troubleshooting](ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md#-troubleshooting)

---

## 📈 Roadmap

- [ ] Добавить больше POI (100+)
- [ ] Импорт данных из ЦОДД
- [ ] Оптимизация маршрутов по времени
- [ ] Учет погоды и дорожных работ
- [ ] WebSocket для real-time обновлений
- [ ] Frontend приложение
- [ ] Mobile приложение (iOS/Android)
- [ ] Голосовое управление
- [ ] Рекомендации на основе истории

---

## 👥 Команда

Проект создан для хакатона 2GIS

---

## 📄 Лицензия

MIT License - свободное использование

---

## 🔗 Полезные Ссылки

- [2GIS API Docs](https://docs.2gis.com/)
- [YandexGPT Docs](https://cloud.yandex.ru/docs/yandexgpt/)
- [Yandex SpeechKit Docs](https://cloud.yandex.ru/docs/speechkit/)
- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)

---

## 💬 Поддержка

Если у вас возникли вопросы:
1. Прочитайте [ПОЛНУЮ_ДОКУМЕНТАЦИЮ.md](ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md)
2. Посмотрите [ШПАРГАЛКУ.md](ШПАРГАЛКА.md)
3. Изучите [ДИАГРАММЫ.md](ДИАГРАММЫ.md)

---

<div align="center">
  <b>Создано с ❤️ для хакатона 2GIS</b>
  <br>
  <sub>Умные аудиомаршруты для незабываемых прогулок</sub>
</div>
