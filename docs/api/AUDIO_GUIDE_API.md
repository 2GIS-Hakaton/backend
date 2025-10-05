# 🎧 API для получения аудиогида

## Новая функциональность

Добавлена ручка для получения полного аудиогида маршрута в формате MP3.

## Эндпоинт

```
GET /api/routes/:route_id/audio
```

### Параметры

- `route_id` (UUID) - ID маршрута, полученный при создании через `/api/routes/generate`

### Ответы

#### ✅ 200 OK - Аудиогид готов

Возвращает MP3 файл со всеми точками маршрута, объединенными в один аудиогид.

**Headers:**
```
Content-Type: audio/mpeg
Content-Disposition: attachment; filename=route_<short_id>.mp3
```

#### ⏳ 206 Partial Content - Частично готов

Некоторые аудиофайлы еще генерируются.

**Response:**
```json
{
  "error": "Some audio files are not ready yet",
  "ready": 2,
  "pending": ["Красная площадь", "ГУМ"],
  "message": "Only 2 of 4 audio files are ready"
}
```

#### ⚠️ 404 Not Found - Не готов

Аудио еще не сгенерировано.

**Response:**
```json
{
  "error": "No audio files generated yet",
  "pending": ["Красная площадь", "ГУМ", "Собор Василия Блаженного"],
  "message": "Audio is being generated. Please try again in 1-2 minutes"
}
```

## Примеры использования

### 1. Базовое использование

```bash
# Получить аудиогид
curl http://localhost:8080/api/routes/550e8400-e29b-41d4-a716-446655440000/audio \
  -o my_audio_guide.mp3

# Воспроизвести
afplay my_audio_guide.mp3  # macOS
mpg123 my_audio_guide.mp3  # Linux
```

### 2. С проверкой статуса

```bash
# Скачать с выводом HTTP кода
curl -w "\nHTTP Code: %{http_code}\n" \
  http://localhost:8080/api/routes/550e8400-e29b-41d4-a716-446655440000/audio \
  -o audio_guide.mp3
```

### 3. Использование тестового скрипта

```bash
# Сначала создайте маршрут
./test_route.sh

# Затем скачайте аудиогид (используйте route_id из вывода)
./test_audio_guide.sh <route_id>
```

## Workflow

```
1. Создать маршрут
   POST /api/routes/generate
   → Получить route_id

2. Подождать генерации контента (30-60 сек на точку)
   
3. Скачать аудиогид
   GET /api/routes/:route_id/audio
   → Получить MP3 файл

4. Воспроизвести
```

## Технические детали

### Объединение аудиофайлов

Если в маршруте несколько точек, их аудиофайлы объединяются в один MP3:

```go
// Последовательное объединение MP3 файлов
for each waypoint in route:
    append waypoint.audio to merged_file
```

### Кэширование

Объединенный файл сохраняется как:
```
./audio/route_<route_id>_full.mp3
```

При повторном запросе возвращается готовый файл.

### Асинхронная генерация

Аудио генерируется асинхронно после создания маршрута:

```
POST /routes/generate
  → Создает маршрут
  → Запускает фоновую генерацию аудио
  → Возвращает route_id сразу

GET /routes/:id/audio
  → Проверяет готовность
  → Возвращает MP3 или статус
```

## Сравнение с существующими ручками

| Ручка | Что возвращает | Когда использовать |
|-------|----------------|-------------------|
| `GET /api/audio/:waypoint_id` | MP3 одной точки | Для прослушивания отдельной точки |
| `GET /api/routes/:route_id/audio` | MP3 всего маршрута | Для полного аудиогида |
| `GET /api/routes/:route_id` | JSON с метаданными | Для получения информации о маршруте |

## Ограничения

- Максимальный размер файла зависит от количества точек (обычно 5-50 MB)
- Генерация занимает 30-60 секунд на точку
- Объединение происходит простой конкатенацией (без нормализации громкости)

## Возможные улучшения

1. **Нормализация громкости** - выровнять уровень звука между точками
2. **Добавление пауз** - вставить тишину между точками
3. **Фоновая музыка** - добавить ambient музыку
4. **Прогресс генерации** - WebSocket для real-time статуса
5. **Кэширование** - проверять существование файла перед объединением
6. **Форматы** - поддержка других форматов (OGG, AAC)

## Troubleshooting

### Аудио не готово

```bash
# Проверьте статус маршрута
curl http://localhost:8080/api/routes/<route_id> | jq '.waypoints[].content'

# Если content == null, подождите 1-2 минуты
```

### Ошибка при объединении

```bash
# Проверьте логи сервера
# Убедитесь, что директория ./audio существует и доступна для записи
mkdir -p ./audio
chmod 755 ./audio
```

### Файл поврежден

```bash
# Проверьте отдельные аудиофайлы
ls -lh ./audio/

# Попробуйте скачать отдельные точки
curl http://localhost:8080/api/audio/<waypoint_id> -o test.mp3
```

## Примеры интеграции

### JavaScript/TypeScript

```typescript
async function downloadAudioGuide(routeId: string): Promise<Blob> {
  const response = await fetch(
    `http://localhost:8080/api/routes/${routeId}/audio`
  );
  
  if (response.status === 200) {
    return await response.blob();
  } else if (response.status === 206 || response.status === 404) {
    const json = await response.json();
    throw new Error(json.message);
  } else {
    throw new Error('Failed to download audio guide');
  }
}

// Использование
try {
  const audioBlob = await downloadAudioGuide(routeId);
  const audioUrl = URL.createObjectURL(audioBlob);
  const audio = new Audio(audioUrl);
  audio.play();
} catch (error) {
  console.error('Audio not ready:', error.message);
  // Повторить через 30 секунд
}
```

### Python

```python
import requests
import time

def download_audio_guide(route_id: str, max_retries: int = 5):
    url = f"http://localhost:8080/api/routes/{route_id}/audio"
    
    for attempt in range(max_retries):
        response = requests.get(url)
        
        if response.status_code == 200:
            # Сохранить файл
            with open(f"route_{route_id[:8]}.mp3", "wb") as f:
                f.write(response.content)
            return True
        elif response.status_code in [206, 404]:
            print(f"Audio not ready, waiting... (attempt {attempt + 1}/{max_retries})")
            time.sleep(30)
        else:
            raise Exception(f"Error: {response.status_code}")
    
    return False

# Использование
if download_audio_guide(route_id):
    print("Audio guide downloaded!")
else:
    print("Audio guide not ready after retries")
```

### Swift (iOS)

```swift
func downloadAudioGuide(routeId: String, completion: @escaping (URL?, Error?) -> Void) {
    let url = URL(string: "http://localhost:8080/api/routes/\(routeId)/audio")!
    
    let task = URLSession.shared.downloadTask(with: url) { localURL, response, error in
        guard let httpResponse = response as? HTTPURLResponse else {
            completion(nil, error)
            return
        }
        
        switch httpResponse.statusCode {
        case 200:
            completion(localURL, nil)
        case 206, 404:
            // Аудио еще не готово
            completion(nil, NSError(domain: "AudioGuide", code: httpResponse.statusCode))
        default:
            completion(nil, NSError(domain: "AudioGuide", code: httpResponse.statusCode))
        }
    }
    
    task.resume()
}

// Использование
downloadAudioGuide(routeId: routeId) { url, error in
    if let url = url {
        // Воспроизвести аудио
        let player = AVPlayer(url: url)
        player.play()
    } else {
        print("Audio not ready, retry in 30 seconds")
    }
}
```
