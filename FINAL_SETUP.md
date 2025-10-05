# ✅ Финальная настройка - Все исправлено!

## 🔧 Что было исправлено

1. **✅ Порт БД** - Обновлен с 5432 на **30101**
2. **✅ Скрипт импорта** - Исправлено название таблицы `places_of_interest` → `pois`
3. **✅ TTS сервис** - Исправлено кодирование URL параметров (символы `\n` теперь кодируются правильно)
4. **✅ POI импортированы** - 12 тестовых мест в базе
5. **✅ Swagger документация** - Добавлена интерактивная API документация
6. **✅ Структурированная документация** - Все `.md` файлы организованы в `docs/`

---

## 🚀 Запуск (пошагово)

### Шаг 1: Установить переменные окружения

```bash
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
export YANDEX_API_KEY="your-yandex-api-key"
export YANDEX_FOLDER_ID="your-yandex-folder-id"
```

**Важно:** Без Yandex API ключей аудио не сгенерируется!

### Шаг 2: Запустить сервер

```bash
cd backend
go run cmd/server/main.go
```

Вы увидите:
```
🚀 Server starting on port 8080
```

### Шаг 3: Проверить, что все работает

В **другом терминале**:

```bash
# Проверить сервер
curl http://localhost:8080/api/health

# Проверить POI
curl http://localhost:8080/api/pois | jq

# Открыть Swagger UI
open http://localhost:8080/swagger/index.html
```

### Шаг 4: Сгенерировать аудиогид

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

# Подождите 2-3 минуты...

# Воспроизвести
afplay route_guide.mp3
```

---

## 🎯 Или используйте готовый скрипт

```bash
# Установить переменные
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
export YANDEX_API_KEY="your-key"
export YANDEX_FOLDER_ID="your-folder"

# Запустить скрипт
./generate_audio_route.sh
```

---

## 📚 Документация

### Главные ссылки

- **[START_HERE.md](START_HERE.md)** - 🚀 Быстрый старт
- **[docs/README.md](docs/README.md)** - 📚 Центр документации
- **[Swagger UI](http://localhost:8080/swagger/index.html)** - 🔧 Интерактивная API

### API документация

- **[Синхронный API](docs/api/SYNC_AUDIO_API.md)** - Один запрос → MP3
- **[Асинхронный API](docs/api/AUDIO_GUIDE_API.md)** - Гибкая генерация

### Руководства

- **[Шпаргалка](docs/guides/ШПАРГАЛКА.md)** - Все команды
- **[Быстрый старт](docs/guides/QUICK_START_AUDIO.md)** - За 3 минуты

---

## 🐛 Troubleshooting

### Ошибка "invalid control character in URL"

**Решение:** Уже исправлено! TTS сервис теперь правильно кодирует параметры URL.

### Ошибка "No POIs found"

```bash
# Проверить количество POI
PGPASSWORD=changeme psql -h 51.250.86.178 -p 30101 -U nike -d audioguid -c "SELECT COUNT(*) FROM pois;"

# Если 0, импортировать
cd scripts
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
python3 import_sample_pois.py
```

### Ошибка "YANDEX_API_KEY not set"

```bash
# Установить ключи
export YANDEX_API_KEY="your-key"
export YANDEX_FOLDER_ID="your-folder"

# Перезапустить сервер
```

### Сервер не запускается

```bash
# Убить процесс на порту 8080
lsof -ti:8080 | xargs kill -9

# Проверить подключение к БД
PGPASSWORD=changeme psql -h 51.250.86.178 -p 30101 -U nike -d audioguid -c "SELECT 1;"
```

---

## ✨ Что теперь работает

### ✅ API Эндпоинты

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/health` | Проверка работоспособности |
| `GET` | `/api/pois` | Список мест интереса |
| `GET` | `/api/pois/:id` | Детали места |
| `POST` | `/api/routes/generate` | Создать маршрут (асинхронно) |
| `POST` | `/api/routes/generate-audio` | **Создать маршрут + MP3 (синхронно)** ⚡ |
| `GET` | `/api/routes/:id` | Получить маршрут |
| `GET` | `/api/routes/:id/audio` | Получить аудиогид маршрута |
| `GET` | `/api/audio/:waypoint_id` | Получить аудио точки |
| `GET` | `/swagger/*` | Swagger UI |

### ✅ Swagger UI

Интерактивная документация:
```
http://localhost:8080/swagger/index.html
```

Возможности:
- Просмотр всех эндпоинтов
- Описание параметров и моделей
- Тестирование API в браузере
- Экспорт OpenAPI спецификации

### ✅ База данных

- **12 POI** в базе
- **Порт:** 30101
- **Таблицы:** pois, routes, waypoints, contents

### ✅ Документация

Структурированная в `docs/`:
- `api/` - API документация
- `guides/` - Руководства
- `reference/` - Справочники
- `diagrams/` - Диаграммы

---

## 🎓 Примеры использования

### Простой запрос

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -d '{"start_point": {"lat": 55.7558, "lon": 37.6173}, "duration_minutes": 60, "max_waypoints": 2}' \
  -H "Content-Type: application/json" \
  -o guide.mp3
```

### С фильтрами

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 120,
    "epochs": ["soviet", "medieval", "imperial"],
    "interests": ["architecture", "history", "culture", "art"],
    "max_waypoints": 5
  }' \
  -o full_guide.mp3
```

### Через Swagger UI

1. Откройте http://localhost:8080/swagger/index.html
2. Найдите `POST /api/routes/generate-audio`
3. Нажмите "Try it out"
4. Введите параметры
5. Нажмите "Execute"
6. Скачайте MP3 файл

---

## 📊 Статистика проекта

- **API эндпоинтов:** 9
- **POI в базе:** 12
- **Документация:** 15+ файлов
- **Swagger аннотаций:** ~70 строк
- **Тестовых скриптов:** 3

---

## 🎉 Готово!

Все исправлено и готово к работе:

- ✅ Правильный порт БД (30101)
- ✅ POI импортированы
- ✅ TTS сервис исправлен
- ✅ Swagger UI работает
- ✅ Документация структурирована

**Начните с:** [START_HERE.md](START_HERE.md) или откройте [Swagger UI](http://localhost:8080/swagger/index.html)

---

## 💡 Полезные команды

```bash
# Проверить все
curl http://localhost:8080/api/health && \
curl http://localhost:8080/api/pois | jq 'length' && \
echo "✅ Все работает!"

# Регенерировать Swagger
cd backend
~/go/bin/swag init -g cmd/server/main.go --output ./docs

# Пересоздать POI
cd scripts
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
python3 import_sample_pois.py
```

---

<div align="center">
  <b>🚀 Проект готов к использованию!</b>
  <br>
  <sub>Swagger + Структурированная документация + Рабочий API</sub>
</div>
