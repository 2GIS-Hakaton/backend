#!/bin/bash
# Скрипт получения аудио для маршрута

set -e

# Параметры маршрута (измени по необходимости)
START_LAT=55.7558
START_LON=37.6173
DURATION=60
EPOCHS='["soviet","medieval"]'
INTERESTS='["architecture","history"]'
MAX_POINTS=3

# Проверка Yandex API ключей
if [ -z "$YANDEX_API_KEY" ]; then
    echo "⚠️  ВНИМАНИЕ: YANDEX_API_KEY не установлен"
    echo "   Аудио не будет сгенерировано автоматически."
    echo ""
    echo "   Для генерации аудио установите:"
    echo "   export YANDEX_API_KEY=\"your-key\""
    echo "   export YANDEX_FOLDER_ID=\"your-folder\""
    echo ""
    echo "   Продолжить без генерации аудио? (y/n)"
    read -r CONTINUE
    if [ "$CONTINUE" != "y" ]; then
        echo "Отменено."
        exit 0
    fi
    echo ""
fi

echo "🚀 Создание маршрута..."

# 1. Создать маршрут
ROUTE_JSON=$(curl -s -X POST http://localhost:8080/api/routes/generate \
  -H "Content-Type: application/json" \
  -d "{
    \"start_point\": {\"lat\": $START_LAT, \"lon\": $START_LON},
    \"duration_minutes\": $DURATION,
    \"epochs\": $EPOCHS,
    \"interests\": $INTERESTS,
    \"max_waypoints\": $MAX_POINTS
  }")

ROUTE_ID=$(echo $ROUTE_JSON | python3 -c "import sys,json; print(json.load(sys.stdin).get('route_id',''))")

if [ -z "$ROUTE_ID" ]; then
    echo "❌ Ошибка создания маршрута"
    echo $ROUTE_JSON | python3 -m json.tool
    exit 1
fi

echo "✅ Маршрут создан: $ROUTE_ID"
echo "📍 Точки маршрута:"
echo $ROUTE_JSON | python3 -c "
import sys, json
data = json.load(sys.stdin)
for i, wp in enumerate(data['waypoints'], 1):
    print(f'  {i}. {wp[\"name\"]} ({wp[\"epoch\"]}, {wp[\"category\"]})')
"

# 2. Подождать генерации
WAIT_TIME=$((MAX_POINTS * 30))
echo ""
echo "⏳ Ожидание генерации аудио ($WAIT_TIME секунд)..."
for i in $(seq $WAIT_TIME -1 1); do
    printf "\r   Осталось: %3d сек" $i
    sleep 1
done
echo ""

# 3. Скачать аудио
echo ""
echo "🔊 Скачивание аудиофайлов..."

ROUTE_DATA=$(curl -s http://localhost:8080/api/routes/$ROUTE_ID)

# Проверка что получили данные
if [ -z "$ROUTE_DATA" ]; then
    echo "❌ Ошибка: не удалось получить данные маршрута"
    echo "Проверьте что сервер запущен: curl http://localhost:8080/api/health"
    exit 1
fi

# Проверка на ошибку в ответе
if echo "$ROUTE_DATA" | grep -q '"error"'; then
    echo "❌ Ошибка от сервера:"
    echo "$ROUTE_DATA" | python3 -m json.tool 2>/dev/null || echo "$ROUTE_DATA"
    exit 1
fi

echo "$ROUTE_DATA" | python3 << 'EOF'
import sys, json, subprocess, os

data = json.load(sys.stdin)
downloaded = 0
pending = 0

for i, wp in enumerate(data['waypoints'], 1):
    if wp.get('content') and wp['content']:
        waypoint_id = wp['id']
        # Безопасное имя файла
        safe_name = "".join(c if c.isalnum() or c in (' ','-','_') else '_' for c in wp['name'])
        filename = f"audio_{i:02d}_{safe_name}.mp3"
        
        # Скачать файл
        result = subprocess.run([
            'curl', '-s', '-f',
            f'http://localhost:8080/api/audio/{waypoint_id}',
            '-o', filename
        ], capture_output=True)
        
        if result.returncode == 0 and os.path.getsize(filename) > 0:
            size = os.path.getsize(filename) / 1024  # KB
            duration = wp['content'].get('duration_seconds', 0)
            print(f"✅ {i}. {wp['name']}")
            print(f"   Файл: {filename} ({size:.1f} KB, {duration}s)")
            downloaded += 1
        else:
            print(f"❌ {i}. {wp['name']}: не удалось скачать")
            if os.path.exists(filename):
                os.remove(filename)
    else:
        print(f"⏳ {i}. {wp['name']}: аудио еще не сгенерировано")
        pending += 1

print(f"\n{'='*60}")
print(f"Скачано: {downloaded} файл(ов)")
if pending > 0:
    print(f"Ожидают генерации: {pending}")
    if downloaded == 0:
        print("\n⚠️  Аудио не сгенерировано.")
        print("Причина: не установлены Yandex API ключи")
        print("\nДля генерации аудио:")
        print("1. Получите ключи: https://console.cloud.yandex.ru/")
        print("2. Установите переменные:")
        print("   export YANDEX_API_KEY=\"your-key\"")
        print("   export YANDEX_FOLDER_ID=\"your-folder\"")
        print("3. Перезапустите сервер")
        print("4. Запустите скрипт снова")
    else:
        print("Повторите попытку через 1-2 минуты")
EOF

echo ""
if [ -f audio_*.mp3 ]; then
    echo "✨ Готово!"
    echo ""
    echo "📂 Список файлов:"
    ls -lh audio_*.mp3 2>/dev/null
else
    echo "⚠️  Аудиофайлы не найдены"
    echo ""
    echo "Маршрут создан успешно: $ROUTE_ID"
    echo "Для получения аудио нужны Yandex API ключи."
    echo ""
    echo "Без аудио вы можете:"
    echo "  - Посмотреть маршрут: curl http://localhost:8080/api/routes/$ROUTE_ID"
    echo "  - Получить список точек с описанием"
    echo "  - Использовать координаты для навигации"
fi