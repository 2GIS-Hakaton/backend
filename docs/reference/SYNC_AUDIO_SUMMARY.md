# ⚡ Синхронная генерация аудиогида - Сводка

## ✨ Что добавлено

Новая ручка для **синхронной** генерации аудиогида - отправляете параметры маршрута, получаете готовый MP3 файл одним запросом!

---

## 🎯 Новый эндпоинт

```
POST /api/routes/generate-audio
```

**Вход:** JSON с параметрами маршрута  
**Выход:** MP3 файл с аудиогидом

---

## 🚀 Как использовать

### Самый простой способ

```bash
./generate_audio_route.sh
```

### С помощью curl

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

afplay my_guide.mp3
```

---

## 📊 Сравнение

| Характеристика | Асинхронный API | Синхронный API |
|----------------|-----------------|----------------|
| **Запросов** | 2-3 | 1 |
| **Ожидание** | Нужно проверять статус | Автоматическое |
| **Время ответа** | Мгновенно (route_id) | 2-5 минут (MP3) |
| **Сложность** | Средняя | Простая |
| **Использование** | Сложные сценарии | Простые сценарии |

---

## 📁 Новые файлы

1. **`generate_audio_route.sh`** - Скрипт для тестирования
2. **`SYNC_AUDIO_API.md`** - Полная документация
3. **`SYNC_AUDIO_SUMMARY.md`** - Эта сводка

## 📝 Измененные файлы

1. **`backend/internal/api/route_handler.go`** - Добавлена функция `GenerateRouteWithAudio()`
2. **`backend/internal/api/routes.go`** - Зарегистрирован роут `/routes/generate-audio`
3. **`README.md`** - Обновлена таблица API и примеры

---

## ⏱️ Время выполнения

- **1 точка:** ~30-60 секунд
- **3 точки:** ~2-3 минуты
- **5 точек:** ~3-5 минут

---

## ✅ Преимущества

- ✅ **Один запрос** вместо нескольких
- ✅ **Автоматическая генерация** всего контента
- ✅ **Готовый MP3** на выходе
- ✅ **Простота использования**
- ✅ **Обратная совместимость** - старые ручки работают

---

## 🎓 Пример workflow

```bash
# 1. Запустить сервер
cd backend
export DATABASE_URL="postgresql://nike:<password>@51.250.86.178:<port>/audioguid?sslmode=disable"
export YANDEX_API_KEY="your-key"
export YANDEX_FOLDER_ID="your-folder"
go run cmd/server/main.go

# 2. В другом терминале - сгенерировать аудиогид
./generate_audio_route.sh

# 3. Готово! Файл route_audio_*.mp3 создан
```

---

## 📚 Документация

- **Полная документация:** [SYNC_AUDIO_API.md](SYNC_AUDIO_API.md)
- **Асинхронный API:** [AUDIO_GUIDE_API.md](AUDIO_GUIDE_API.md)
- **Быстрый старт:** [QUICK_START_AUDIO.md](QUICK_START_AUDIO.md)
- **Основной README:** [README.md](README.md)

---

## 🎉 Готово!

Теперь можно получить аудиогид **одним запросом**:

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -d '{"start_point": {"lat": 55.7558, "lon": 37.6173}, "duration_minutes": 60, "epochs": ["soviet"], "interests": ["architecture"]}' \
  -o guide.mp3
```

**Просто, быстро, удобно!** 🚀
