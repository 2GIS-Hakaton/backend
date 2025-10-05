# 📍 Руководство: Свои места для аудиогида

## 🎯 Возможности

Теперь вы можете создавать аудиогиды тремя способами:

1. **Автоматический поиск** - по критериям (эпоха, категория)
2. **Конкретные места из базы** - передать ID мест
3. **Свои места** - передать координаты и описания

---

## 📡 API

### Вариант 1: Автоматический поиск (как раньше)

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 90,
    "epochs": ["soviet"],
    "interests": ["architecture"],
    "max_waypoints": 3
  }' \
  -o route_guide.mp3
```

### Вариант 2: Конкретные места из базы (по ID)

```bash
# Сначала получите список POI и их ID
curl http://localhost:8080/api/pois | jq '.[] | {id, name}'

# Затем используйте нужные ID
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 90,
    "poi_ids": [
      "uuid-красной-площади",
      "uuid-кремля",
      "uuid-вднх"
    ]
  }' \
  -o route_guide.mp3
```

### Вариант 3: Свои места (координаты + описания) ⭐

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 90,
    "custom_pois": [
      {
        "name": "Моя любимая кофейня",
        "description": "Уютное место с лучшим кофе в городе",
        "latitude": 55.7558,
        "longitude": 37.6173,
        "epoch": "modern",
        "category": "culture"
      },
      {
        "name": "Памятник моему деду",
        "description": "Здесь служил мой дедушка",
        "latitude": 55.7539,
        "longitude": 37.6208,
        "epoch": "soviet",
        "category": "history"
      },
      {
        "name": "Место первого свидания",
        "description": "Здесь мы познакомились",
        "latitude": 55.7520,
        "longitude": 37.6175
      }
    ]
  }' \
  -o my_personal_guide.mp3
```

---

## 🎓 Примеры использования

### Пример 1: Экскурсия по любимым местам

```json
{
  "start_point": {"lat": 55.7558, "lon": 37.6173},
  "duration_minutes": 120,
  "custom_pois": [
    {
      "name": "Парк, где я гуляю",
      "description": "Мое любимое место для утренних пробежек",
      "latitude": 55.7304,
      "longitude": 37.6012
    },
    {
      "name": "Кафе 'У Марии'",
      "description": "Здесь подают лучшие круассаны",
      "latitude": 55.7415,
      "longitude": 37.6206
    },
    {
      "name": "Книжный магазин",
      "description": "Провожу здесь каждую субботу",
      "latitude": 55.7526,
      "longitude": 37.5676
    }
  ]
}
```

### Пример 2: Семейная история

```json
{
  "start_point": {"lat": 55.7558, "lon": 37.6173},
  "duration_minutes": 90,
  "custom_pois": [
    {
      "name": "Дом бабушки",
      "description": "Здесь выросла моя бабушка. Дом 1950-х годов постройки",
      "latitude": 55.8304,
      "longitude": 37.6325,
      "epoch": "soviet",
      "category": "history"
    },
    {
      "name": "Школа деда",
      "description": "Дедушка учился здесь в 1960-х",
      "latitude": 55.8283,
      "longitude": 37.6308,
      "epoch": "soviet",
      "category": "culture"
    },
    {
      "name": "Место встречи родителей",
      "description": "Здесь познакомились мои родители",
      "latitude": 55.8277,
      "longitude": 37.6319,
      "epoch": "soviet",
      "category": "culture"
    }
  ]
}
```

### Пример 3: Туристический маршрут

```json
{
  "start_point": {"lat": 55.7539, "lon": 37.6208},
  "duration_minutes": 120,
  "custom_pois": [
    {
      "name": "Красная площадь",
      "description": "Главная площадь России, сердце Москвы",
      "latitude": 55.7539,
      "longitude": 37.6208,
      "epoch": "medieval",
      "category": "history"
    },
    {
      "name": "ГУМ",
      "description": "Главный универсальный магазин, построен в 1893 году",
      "latitude": 55.7551,
      "longitude": 37.6211,
      "epoch": "imperial",
      "category": "architecture"
    },
    {
      "name": "Мавзолей Ленина",
      "description": "Усыпальница Владимира Ленина",
      "latitude": 55.7536,
      "longitude": 37.6198,
      "epoch": "soviet",
      "category": "history"
    }
  ]
}
```

### Пример 4: Комбинированный (ID + свои места)

```json
{
  "start_point": {"lat": 55.7558, "lon": 37.6173},
  "duration_minutes": 90,
  "poi_ids": [
    "uuid-красной-площади",
    "uuid-кремля"
  ],
  "custom_pois": [
    {
      "name": "Моя любимая смотровая площадка",
      "description": "Отсюда открывается лучший вид на Кремль",
      "latitude": 55.7520,
      "longitude": 37.6180
    }
  ]
}
```

**Примечание:** Если указаны и `poi_ids`, и `custom_pois`, приоритет у `poi_ids`.

---

## 📝 Формат данных

### CustomPOI

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `name` | string | ✅ Да | Название места |
| `description` | string | ❌ Нет | Описание (используется для генерации рассказа) |
| `latitude` | float64 | ✅ Да | Широта |
| `longitude` | float64 | ✅ Да | Долгота |
| `epoch` | string | ❌ Нет | Эпоха (medieval, imperial, soviet, modern) |
| `category` | string | ❌ Нет | Категория (architecture, history, culture, religion, art) |

---

## 🔧 Как это работает

### Для конкретных POI ID

1. API получает список ID
2. Ищет эти места в базе данных
3. Генерирует аудио в том порядке, в котором переданы ID
4. Возвращает MP3 файл

### Для своих мест

1. API получает список мест с координатами
2. **Создает новые POI в базе данных**
3. Генерирует рассказ на основе названия и описания
4. Создает аудио
5. Возвращает MP3 файл

**Важно:** Свои места сохраняются в базу и получают UUID. Вы можете использовать их позже!

---

## 💡 Советы

### 1. Хорошее описание = лучший рассказ

```json
// ❌ Плохо
{
  "name": "Место",
  "description": "Красивое"
}

// ✅ Хорошо
{
  "name": "Парк Победы",
  "description": "Мемориальный комплекс, открытый в 1995 году в честь 50-летия Победы в Великой Отечественной войне. Здесь установлен монумент Победы высотой 141.8 метра."
}
```

### 2. Указывайте эпоху и категорию

Это помогает AI генерировать более точный и интересный рассказ:

```json
{
  "name": "Дом-музей Пушкина",
  "description": "Здесь жил великий поэт",
  "latitude": 55.7415,
  "longitude": 37.6206,
  "epoch": "imperial",      // Помогает AI понять контекст
  "category": "culture"     // Определяет стиль рассказа
}
```

### 3. Логичный порядок

Места обрабатываются в том порядке, в котором вы их передали:

```json
{
  "custom_pois": [
    {"name": "Начало маршрута", ...},
    {"name": "Середина маршрута", ...},
    {"name": "Конец маршрута", ...}
  ]
}
```

---

## 🎯 Примеры использования через Swagger

1. Откройте http://localhost:8080/swagger/index.html
2. Найдите `POST /api/routes/generate-audio`
3. Нажмите "Try it out"
4. Вставьте JSON с вашими местами
5. Нажмите "Execute"
6. Скачайте MP3

---

## 🔍 Получение ID существующих мест

### Все места

```bash
curl http://localhost:8080/api/pois | jq '.[] | {id, name, epoch, category}'
```

### Фильтр по эпохе

```bash
curl "http://localhost:8080/api/pois?epoch=soviet" | jq '.[] | {id, name}'
```

### Фильтр по категории

```bash
curl "http://localhost:8080/api/pois?category=architecture" | jq '.[] | {id, name}'
```

### Конкретное место

```bash
curl http://localhost:8080/api/pois/UUID | jq
```

---

## 📊 Сравнение способов

| Способ | Преимущества | Недостатки |
|--------|-------------|-----------|
| **Автопоиск** | Простой, быстрый | Нет контроля над местами |
| **POI ID** | Конкретные места из базы | Нужно знать ID |
| **Свои места** | Полный контроль, персонализация | Нужно указать координаты |

---

## ✨ Итог

Теперь вы можете:

✅ Создавать персональные аудиогиды  
✅ Добавлять свои любимые места  
✅ Рассказывать семейные истории  
✅ Делать уникальные туристические маршруты  

**Попробуйте прямо сейчас!**

```bash
curl -X POST http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {"lat": 55.7558, "lon": 37.6173},
    "duration_minutes": 60,
    "custom_pois": [
      {
        "name": "Ваше любимое место",
        "description": "Расскажите о нем",
        "latitude": 55.7558,
        "longitude": 37.6173
      }
    ]
  }' \
  -o my_guide.mp3
```
