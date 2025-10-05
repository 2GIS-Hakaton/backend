# 🎧 Быстрый старт - Аудиогид API

## За 3 минуты

### 1️⃣ Запустить сервер

```bash
cd backend
export DATABASE_URL="postgresql://nike:<password>@51.250.86.178:<port>/audioguid?sslmode=disable"
go run cmd/server/main.go
```

### 2️⃣ Создать маршрут

```bash
curl -X POST http://localhost:8080/api/routes/generate \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 60,
    "epochs": ["soviet"],
    "interests": ["architecture"],
    "max_waypoints": 3
  }' | python3 -m json.tool
```

**Скопируйте `route_id` из ответа!**

### 3️⃣ Скачать аудиогид

```bash
# Замените <route_id> на ваш ID
curl http://localhost:8080/api/routes/<route_id>/audio -o my_guide.mp3

# Воспроизвести
afplay my_guide.mp3  # macOS
```

---

## Автоматический тест

```bash
# Полный тест
./test_route.sh

# Скачать аудиогид (используйте route_id из вывода)
./test_audio_guide.sh <route_id>
```

---

## Статусы ответа

### ✅ 200 OK - Готово
Возвращает MP3 файл

### ⏳ 206 Partial Content - Частично готово
```json
{
  "error": "Some audio files are not ready yet",
  "ready": 2,
  "pending": ["Красная площадь", "ГУМ"],
  "message": "Only 2 of 4 audio files are ready"
}
```
**→ Подождите 1-2 минуты и попробуйте снова**

### ⚠️ 404 Not Found - Не готово
```json
{
  "error": "No audio files generated yet",
  "pending": ["Красная площадь", "ГУМ", "Собор"],
  "message": "Audio is being generated. Please try again in 1-2 minutes"
}
```
**→ Подождите 1-2 минуты и попробуйте снова**

---

## Полная документация

📖 См. [AUDIO_GUIDE_API.md](AUDIO_GUIDE_API.md)

---

## Troubleshooting

### Аудио не готово?
```bash
# Проверьте статус маршрута
curl http://localhost:8080/api/routes/<route_id> | python3 -m json.tool

# Если content == null, подождите 1-2 минуты
```

### Сервер не отвечает?
```bash
# Проверьте, что сервер запущен
curl http://localhost:8080/api/health
```

### Ошибка 404?
```bash
# Проверьте правильность route_id
# Убедитесь, что маршрут существует
curl http://localhost:8080/api/routes/<route_id>
```

---

## Готово! 🎉

Теперь вы можете получать полные аудиогиды одним запросом!
