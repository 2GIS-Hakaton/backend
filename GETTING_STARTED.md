# 🚀 Начало работы с Audio Guide API

> Полное руководство по запуску и использованию проекта

---

## 📋 Содержание

1. [Быстрый старт](#-быстрый-старт)
2. [Требования](#-требования)
3. [Установка](#-установка)
4. [Конфигурация](#-конфигурация)
5. [Запуск](#-запуск)
6. [Использование API](#-использование-api)
7. [Документация](#-документация)
8. [Troubleshooting](#-troubleshooting)

---

## ⚡ Быстрый старт

```bash
# 1. Установить переменные окружения
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
export YANDEX_API_KEY="your-yandex-api-key"
export YANDEX_FOLDER_ID="your-yandex-folder-id"

# 2. Запустить сервер
cd backend
go run cmd/server/main.go

# 3. В другом терминале - сгенерировать аудиогид
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 90,
    "max_waypoints": 3
  }' \
  -o route_guide.mp3

# 4. Воспроизвести
afplay route_guide.mp3
```

---

## 📦 Требования

### Обязательные

- **Go** 1.21+ ([скачать](https://go.dev/dl/))
- **PostgreSQL** 16+ с PostGIS
- **Python** 3.8+ (для скриптов)

### Опциональные (для генерации аудио)

- **Yandex Cloud API ключи** ([получить](https://console.cloud.yandex.ru/))
  - API Key для YandexGPT
  - API Key для Yandex SpeechKit
  - Folder ID

---

## 🔧 Установка

### 1. Клонировать репозиторий

```bash
git clone <repository-url>
cd 2gis_hackaton
```

### 2. Установить зависимости Go

```bash
cd backend
go mod download
```

### 3. Установить зависимости Python

```bash
cd scripts
pip3 install -r requirements.txt
```

### 4. Проверить подключение к БД

```bash
PGPASSWORD=changeme psql -h 51.250.86.178 -p 30101 -U nike -d audioguid -c "SELECT 1;"
```

### 5. Импортировать тестовые POI

```bash
cd scripts
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
python3 import_sample_pois.py
```

Вы должны увидеть:
```
✅ Импорт завершен!
   Добавлено: 12
   Всего POI: 12
```

---

## ⚙️ Конфигурация

### Переменные окружения

Создайте файл `backend/.env`:

```bash
# База данных
DATABASE_URL=postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable

# Yandex Cloud (для генерации аудио)
YANDEX_API_KEY=AQVN...
YANDEX_FOLDER_ID=b1g...
YANDEX_VOICE=alena

# Сервер
PORT=8080
```

### Получение Yandex API ключей

1. Зарегистрируйтесь на https://console.cloud.yandex.ru/
2. Создайте Service Account
3. Добавьте роли:
   - `ai.languageModels.user` (для YandexGPT)
   - `ai.speechkit-tts.user` (для Yandex TTS)
4. Создайте API Key
5. Скопируйте Folder ID из консоли

**Важно:** Без Yandex API ключей аудио не будет генерироваться, но остальной функционал работает.

---

## 🚀 Запуск

### Вариант 1: Через go run (разработка)

```bash
cd backend
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
export YANDEX_API_KEY="your-key"
export YANDEX_FOLDER_ID="your-folder"
go run cmd/server/main.go
```

### Вариант 2: Скомпилировать и запустить

```bash
cd backend
go build -o audioguid cmd/server/main.go
./audioguid
```

### Вариант 3: С .env файлом

```bash
cd backend
# Создайте .env файл с переменными
go run cmd/server/main.go
```

Сервер запустится на порту 8080:
```
🚀 Server starting on port 8080
```

---

## 📡 Использование API

### 1. Проверка работоспособности

```bash
curl http://localhost:8080/api/health
```

Ответ:
```json
{
  "status": "ok",
  "message": "Audioguid API is running"
}
```

### 2. Просмотр доступных POI

```bash
curl http://localhost:8080/api/pois | jq
```

### 3. Генерация аудиогида (синхронно)

**Самый простой способ** - один запрос, сразу получаете MP3:

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

Подождите 2-3 минуты (генерация аудио), затем:

```bash
afplay route_guide.mp3  # macOS
mpg123 route_guide.mp3  # Linux
```

### 4. Генерация маршрута (асинхронно)

Если нужна гибкость:

```bash
# Шаг 1: Создать маршрут
ROUTE_ID=$(curl -X POST http://localhost:8080/api/routes/generate \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 90,
    "max_waypoints": 3
  }' | jq -r '.route_id')

echo "Route ID: $ROUTE_ID"

# Шаг 2: Подождать генерации (2-3 минуты)
sleep 180

# Шаг 3: Скачать аудиогид
curl http://localhost:8080/api/routes/$ROUTE_ID/audio -o guide.mp3
```

---

## 📚 Документация

### Интерактивная документация

После запуска сервера откройте:

```
http://localhost:8080/swagger/index.html
```

**Swagger UI** позволяет:
- Просмотреть все эндпоинты
- Увидеть параметры и модели
- Протестировать API прямо в браузере

### Документация в Markdown

Вся документация структурирована в папке `docs/`:

```
docs/
├── README.md                    # Главный индекс
├── api/                         # API документация
│   ├── AUDIO_GUIDE_API.md      # Асинхронный API
│   └── SYNC_AUDIO_API.md       # Синхронный API
├── guides/                      # Руководства
│   ├── ШПАРГАЛКА.md            # Команды и примеры
│   └── QUICK_START_AUDIO.md    # Быстрый старт
├── reference/                   # Справочники
│   ├── ПОЛНАЯ_ДОКУМЕНТАЦИЯ.md  # Архитектура
│   ├── CHANGELOG_AUDIO_GUIDE.md
│   ├── SUMMARY.md
│   └── SYNC_AUDIO_SUMMARY.md
└── diagrams/                    # Диаграммы
    └── ДИАГРАММЫ.md
```

**Начните с:** [docs/README.md](docs/README.md)

---

## 🐛 Troubleshooting

### Сервер не запускается

**Проблема:** Порт 8080 занят

```bash
# Убить процесс на порту 8080
lsof -ti:8080 | xargs kill -9
```

**Проблема:** Не удается подключиться к БД

```bash
# Проверить подключение
PGPASSWORD=changeme psql -h 51.250.86.178 -p 30101 -U nike -d audioguid -c "SELECT 1;"

# Проверить переменную окружения
echo $DATABASE_URL
```

### Ошибка "No POIs found"

**Причина:** В базе нет POI или они не подходят по критериям

**Решение 1:** Импортировать POI

```bash
cd scripts
export DATABASE_URL="postgresql://nike:changeme@51.250.86.178:30101/audioguid?sslmode=disable"
python3 import_sample_pois.py
```

**Решение 2:** Использовать более широкие критерии

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 120,
    "max_waypoints": 5
  }' \
  -o guide.mp3
```

**Решение 3:** Проверить количество POI

```bash
PGPASSWORD=changeme psql -h 51.250.86.178 -p 30101 -U nike -d audioguid -c "SELECT COUNT(*) FROM pois;"
```

### Ошибка "YANDEX_API_KEY not set"

**Причина:** Не установлены переменные окружения для Yandex Cloud

**Решение:**

```bash
export YANDEX_API_KEY="your-key"
export YANDEX_FOLDER_ID="your-folder"

# Перезапустить сервер
```

Или создайте файл `backend/.env`:

```
YANDEX_API_KEY=AQVN...
YANDEX_FOLDER_ID=b1g...
```

### Аудио не генерируется

**Проверьте:**

1. Установлены ли Yandex API ключи
2. Есть ли права у Service Account
3. Правильный ли Folder ID

```bash
# Проверить переменные
echo $YANDEX_API_KEY
echo $YANDEX_FOLDER_ID

# Проверить логи сервера
# Должны быть сообщения о генерации аудио
```

### Swagger UI не открывается

**Проверьте:**

```bash
# 1. Сервер запущен
curl http://localhost:8080/api/health

# 2. Swagger доступен
curl http://localhost:8080/swagger/doc.json

# 3. Откройте в браузере
open http://localhost:8080/swagger/index.html
```

---

## 🎯 Примеры использования

### Пример 1: Простой маршрут

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 60,
    "max_waypoints": 2
  }' \
  -o simple_guide.mp3
```

### Пример 2: Советская архитектура

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.8304, "lon": 37.6325},
    "duration_minutes": 90,
    "epochs": ["soviet"],
    "interests": ["architecture"],
    "max_waypoints": 4
  }' \
  -o soviet_architecture.mp3
```

### Пример 3: Историческая Москва

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7539, "lon": 37.6208},
    "duration_minutes": 120,
    "epochs": ["medieval", "imperial"],
    "interests": ["history", "religion"],
    "max_waypoints": 5
  }' \
  -o historical_moscow.mp3
```

### Пример 4: Через Swagger UI

1. Откройте http://localhost:8080/swagger/index.html
2. Найдите `POST /api/routes/generate-audio`
3. Нажмите "Try it out"
4. Введите параметры:
   ```json
   {
     "start_point": {"lat": 55.7558, "lon": 37.6173},
     "duration_minutes": 90,
     "max_waypoints": 3
   }
   ```
5. Нажмите "Execute"
6. Скачайте MP3 файл из ответа

---

## 🔗 Полезные ссылки

### Документация проекта
- [Главный README](README.md)
- [Центр документации](docs/README.md)
- [Финальная настройка](FINAL_SETUP.md)

### Swagger UI
- [Интерактивная документация](http://localhost:8080/swagger/index.html)
- [OpenAPI JSON](http://localhost:8080/swagger/doc.json)

### Внешние ресурсы
- [2GIS API](https://docs.2gis.com/)
- [YandexGPT](https://cloud.yandex.ru/docs/yandexgpt/)
- [Yandex SpeechKit](https://cloud.yandex.ru/docs/speechkit/)
- [Go Documentation](https://go.dev/doc/)

---

## 📊 Доступные POI

В базе данных 12 тестовых мест интереса:

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

## ✅ Чеклист запуска

- [ ] Go 1.21+ установлен
- [ ] PostgreSQL доступен (порт 30101)
- [ ] Python 3.8+ установлен
- [ ] Зависимости Go установлены (`go mod download`)
- [ ] Зависимости Python установлены (`pip3 install -r requirements.txt`)
- [ ] POI импортированы (12 записей)
- [ ] Переменные окружения установлены
- [ ] Yandex API ключи получены (опционально)
- [ ] Сервер запущен и отвечает на `/api/health`
- [ ] Swagger UI открывается

---

## 🎉 Готово!

Теперь вы можете:

✅ Генерировать персонализированные аудиомаршруты  
✅ Тестировать API через Swagger UI  
✅ Интегрировать API в свои приложения  
✅ Изучать документацию  

**Следующие шаги:**
1. Откройте [Swagger UI](http://localhost:8080/swagger/index.html)
2. Изучите [документацию](docs/README.md)
3. Сгенерируйте свой первый аудиогид!

---

<div align="center">
  <b>Создано с ❤️ для хакатона 2GIS</b>
  <br>
  <sub>Умные аудиомаршруты для незабываемых прогулок</sub>
</div>
