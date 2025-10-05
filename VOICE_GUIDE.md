# 🎙️ Руководство по выбору голоса

## 🔊 Доступные голоса Yandex SpeechKit

### Мужские голоса

| Голос | Описание | Рекомендуется для |
|-------|----------|-------------------|
| `ermil` | Эмоциональный, выразительный | Экскурсии, рассказы ⭐ **ТЕКУЩИЙ** |
| `filipp` | Нейтральный, спокойный | Информационные материалы |
| `omazh` | Глубокий, серьезный | Исторические рассказы |
| `zahar` | Спокойный, приятный | Медитации, спокойные туры |

### Женские голоса

| Голос | Описание | Рекомендуется для |
|-------|----------|-------------------|
| `alena` | Мягкий, приятный | Экскурсии, аудиокниги |
| `oksana` | Нейтральный | Новости, информация |
| `jane` | Эмоциональный | Художественные рассказы |

### Premium голоса

| Голос | Описание | Рекомендуется для |
|-------|----------|-------------------|
| `marina` | Премиум женский | Профессиональные аудиогиды |
| `alexander` | Премиум мужской | Профессиональные аудиогиды |

---

## 🔧 Как изменить голос

### Способ 1: Через переменную окружения (рекомендуется)

```bash
# Установить голос
export YANDEX_VOICE="ermil"

# Запустить сервер
cd backend
go run cmd/server/main.go
```

**Доступные значения:**
- Мужские: `ermil`, `filipp`, `omazh`, `zahar`
- Женские: `alena`, `oksana`, `jane`
- Premium: `marina`, `alexander`

### Способ 2: Через config.yaml

Отредактируйте `backend/config/config.yaml`:

```yaml
ai:
  yandex_voice: ermil  # Измените на нужный голос
```

### Способ 3: Через .env файл

Создайте/отредактируйте `backend/.env`:

```bash
YANDEX_VOICE=ermil
```

---

## 🎯 Рекомендации по выбору

### Для исторических экскурсий
```bash
export YANDEX_VOICE="omazh"  # Глубокий, серьезный
```

### Для современных туров
```bash
export YANDEX_VOICE="ermil"  # Эмоциональный, живой
```

### Для спокойных прогулок
```bash
export YANDEX_VOICE="zahar"  # Спокойный, приятный
```

### Для женского голоса
```bash
export YANDEX_VOICE="alena"  # Мягкий, приятный
```

---

## 🧪 Тестирование голосов

### Быстрый тест

```bash
# 1. Установить голос
export YANDEX_VOICE="ermil"

# 2. Запустить сервер
cd backend
go run cmd/server/main.go

# 3. Сгенерировать тестовый аудиогид
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 60,
    "custom_pois": [
      {
        "name": "Тестовое место",
        "description": "Это тестовое описание для проверки голоса",
        "latitude": 55.7558,
        "longitude": 37.6173
      }
    ]
  }' \
  -o test_voice.mp3

# 4. Прослушать
afplay test_voice.mp3
```

### Сравнение голосов

```bash
# Создайте скрипт для тестирования всех голосов
for voice in ermil filipp omazh zahar alena oksana jane; do
  echo "Testing voice: $voice"
  export YANDEX_VOICE=$voice
  
  # Перезапустите сервер и сгенерируйте аудио
  curl -X POST http://localhost:8080/api/routes/generate-audio \
    -d '{"start_point": {"lat": 55.7558, "lon": 37.6173}, "duration_minutes": 60, "max_waypoints": 1}' \
    -H "Content-Type: application/json" \
    -o "test_${voice}.mp3"
done
```

---

## 📝 Текущая настройка

**Голос по умолчанию:** `ermil` (эмоциональный мужской)

Установлен в:
- `backend/config/config.yaml` → `yandex_voice: ermil`
- Можно переопределить через `YANDEX_VOICE` env var

---

## 💡 Советы

### 1. Выбор голоса по контенту

- **Исторические места** → `omazh` (глубокий)
- **Современные места** → `ermil` (эмоциональный)
- **Культурные места** → `alena` (мягкий)
- **Архитектура** → `filipp` (нейтральный)

### 2. Тестируйте перед продакшеном

Всегда тестируйте голос на реальном контенте перед использованием.

### 3. Premium голоса

Голоса `marina` и `alexander` могут требовать дополнительной подписки в Yandex Cloud.

---

## 🔗 Полезные ссылки

- [Yandex SpeechKit Docs](https://cloud.yandex.ru/docs/speechkit/tts/voices)
- [Примеры голосов](https://cloud.yandex.ru/docs/speechkit/tts/voices#premium)

---

## ✅ Быстрая смена голоса

```bash
# Остановить сервер (Ctrl+C)

# Установить новый голос
export YANDEX_VOICE="omazh"  # или любой другой

# Запустить сервер снова
cd backend
go run cmd/server/main.go

# Готово! Теперь используется новый голос
```

---

## 🎉 Готово!

Теперь вы можете легко менять голос для аудиогидов!

**Текущий голос:** `ermil` (эмоциональный мужской)
