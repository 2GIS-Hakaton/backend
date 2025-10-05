# 🎧 Синхронная генерация аудиогида

## Новая функциональность

Добавлена ручка для **синхронной** генерации аудиогида - отправляете параметры маршрута, получаете готовый MP3 файл.

---

## 🚀 Эндпоинт

```
POST /api/routes/generate-audio
```

### Запрос

**Content-Type:** `application/json`

```json
{
  "start_point": {
    "lat": 55.7558,
    "lon": 37.6173
  },
  "duration_minutes": 60,
  "epochs": ["soviet", "medieval"],
  "interests": ["architecture", "history"],
  "max_waypoints": 3
}
```

### Ответ

**Content-Type:** `audio/mpeg`

Возвращает MP3 файл с аудиогидом для всего маршрута.

---

## ⚡ Преимущества

| Старый способ | Новый способ |
|---------------|--------------|
| 1. POST /routes/generate → route_id | 1. POST /routes/generate-audio → MP3 |
| 2. Ждать 2-3 минуты | |
| 3. GET /routes/:id/audio → MP3 | |
| **3 шага, асинхронно** | **1 шаг, синхронно** |

---

## 📝 Примеры использования

### Базовый пример

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
  -o my_audio_guide.mp3

# Воспроизвести
afplay my_audio_guide.mp3
```

### С использованием скрипта

```bash
# Простой запуск (параметры по умолчанию)
./generate_audio_route.sh

# С кастомными параметрами
./generate_audio_route.sh 55.7558 37.6173 90 '["soviet","medieval"]' '["architecture","history"]' 4
```

### С параметрами

```bash
# Советская архитектура, 60 минут, 3 точки
./generate_audio_route.sh 55.7558 37.6173 60 '["soviet"]' '["architecture"]' 3

# Средневековая история, 90 минут, 5 точек
./generate_audio_route.sh 55.7539 37.6208 90 '["medieval"]' '["history"]' 5
```

---

## ⏱️ Время выполнения

Запрос **синхронный** - сервер ждет генерации всех аудио и только потом возвращает результат:

- **1 точка** = ~30-60 секунд
- **3 точки** = ~2-3 минуты
- **5 точек** = ~3-5 минут

**Важно:** Не закрывайте соединение во время генерации!

---

## 🔧 Как это работает

```
1. Клиент отправляет параметры маршрута
   ↓
2. Сервер находит POI по критериям
   ↓
3. Для каждой точки:
   - Генерирует текст (YandexGPT)
   - Создает аудио (Yandex TTS)
   - Сохраняет в БД
   ↓
4. Объединяет все MP3 в один файл
   ↓
5. Возвращает готовый аудиогид
```

**Все происходит синхронно в одном запросе!**

---

## 📊 Сравнение с асинхронным API

### Асинхронный (существующий)

```bash
# Шаг 1: Создать маршрут
curl -X POST /api/routes/generate -d '{...}' 
# → {"route_id": "xxx"}

# Шаг 2: Подождать 2-3 минуты
sleep 180

# Шаг 3: Скачать аудио
curl /api/routes/xxx/audio -o guide.mp3
```

**Плюсы:**
- ✅ Быстрый ответ (сразу route_id)
- ✅ Можно проверять статус
- ✅ Подходит для фоновой обработки

**Минусы:**
- ❌ Нужно делать несколько запросов
- ❌ Нужно ждать и проверять готовность
- ❌ Сложнее для простых сценариев

### Синхронный (новый)

```bash
# Один запрос
curl -X POST /api/routes/generate-audio -d '{...}' -o guide.mp3
```

**Плюсы:**
- ✅ Один запрос
- ✅ Сразу готовый результат
- ✅ Проще для простых сценариев
- ✅ Не нужно хранить route_id

**Минусы:**
- ❌ Долгий ответ (2-5 минут)
- ❌ Нельзя проверить прогресс
- ❌ Может timeout на медленных соединениях

---

## 🎯 Когда использовать

### Используйте синхронный API если:
- Нужен простой workflow
- Готовы ждать результата
- Не нужно сохранять маршрут
- Делаете разовый запрос

### Используйте асинхронный API если:
- Нужен быстрый ответ
- Хотите показать прогресс
- Нужно сохранить маршрут для повторного использования
- Делаете массовую генерацию

---

## ⚙️ Требования

### Обязательные переменные окружения

```bash
# База данных
export DATABASE_URL="postgresql://user:pass@host:30101/db?sslmode=disable"

# Yandex Cloud (для генерации аудио)
export YANDEX_API_KEY="AQVN..."
export YANDEX_FOLDER_ID="b1g..."
```

**Без Yandex API ключей аудио не сгенерируется!**

---

## 🐛 Обработка ошибок

### 400 Bad Request
```json
{
  "error": "Invalid request parameters"
}
```
**Решение:** Проверьте формат JSON

### 404 Not Found
```json
{
  "error": "No POIs found matching criteria"
}
```
**Решение:** Увеличьте `duration_minutes` или измените фильтры

### 500 Internal Server Error
```json
{
  "error": "Failed to generate content for ВДНХ: API key not set"
}
```
**Решение:** Установите `YANDEX_API_KEY` и `YANDEX_FOLDER_ID`

### Timeout
Если запрос слишком долгий (>5 минут), может произойти timeout.

**Решение:** 
- Уменьшите `max_waypoints`
- Используйте асинхронный API

---

## 💡 Советы

### 1. Оптимальное количество точек
```json
{
  "max_waypoints": 3  // Оптимально: 2-3 минуты
}
```

### 2. Увеличение timeout в curl
```bash
curl --max-time 300 -X POST ...  # 5 минут
```

### 3. Показ прогресса
```bash
curl --progress-bar -X POST ... -o guide.mp3
```

### 4. Проверка перед запросом
```bash
# Проверить сервер
curl http://localhost:8080/api/health

# Проверить ключи
echo $YANDEX_API_KEY
```

---

## 📚 Полная документация

- **Асинхронный API:** [AUDIO_GUIDE_API.md](AUDIO_GUIDE_API.md)
- **Быстрый старт:** [QUICK_START_AUDIO.md](QUICK_START_AUDIO.md)
- **Backend README:** [backend/README.md](backend/README.md)

---

## 🧪 Тестирование

### Автоматический тест

```bash
./generate_audio_route.sh
```

### Ручной тест

```bash
# 1. Запустить сервер
cd backend
export DATABASE_URL="..."
export YANDEX_API_KEY="..."
export YANDEX_FOLDER_ID="..."
go run cmd/server/main.go

# 2. В другом терминале
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 60,
    "epochs": ["soviet"],
    "interests": ["architecture"],
    "max_waypoints": 2
  }' \
  -o test_guide.mp3

# 3. Проверить
ls -lh test_guide.mp3
afplay test_guide.mp3
```

---

## 🎉 Итог

Теперь можно получить аудиогид **одним запросом**:

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -d '{"start_point": {"lat": 55.7558, "lon": 37.6173}, "duration_minutes": 60, "epochs": ["soviet"], "interests": ["architecture"]}' \
  -o guide.mp3
```

**Просто, быстро, удобно!** 🚀
