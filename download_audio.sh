#!/bin/bash
# 🎵 Скрипт автоматического скачивания аудио для маршрута
# Использование: ./download_audio.sh

set -e

# Цвета для вывода
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8080/api"

echo -e "${BLUE}════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  🎵 Генератор Аудиомаршрутов${NC}"
echo -e "${BLUE}════════════════════════════════════════════════${NC}"
echo ""

# Параметры маршрута (можно изменить)
START_LAT=${START_LAT:-55.7558}
START_LON=${START_LON:-37.6173}
DURATION=${DURATION:-60}
EPOCHS=${EPOCHS:-'["soviet","medieval"]'}
INTERESTS=${INTERESTS:-'["architecture","history"]'}
MAX_POINTS=${MAX_POINTS:-3}

echo -e "${YELLOW}Параметры маршрута:${NC}"
echo "  Начальная точка: $START_LAT, $START_LON"
echo "  Длительность: $DURATION минут"
echo "  Эпохи: $EPOCHS"
echo "  Интересы: $INTERESTS"
echo "  Макс. точек: $MAX_POINTS"
echo ""

# Проверка доступности API
echo -e "${YELLOW}Проверка сервера...${NC}"
if ! curl -s -f "${API_BASE}/health" > /dev/null; then
    echo -e "${RED}❌ Сервер недоступен на ${API_BASE}${NC}"
    echo "Запустите сервер: cd backend && go run cmd/server/main.go"
    exit 1
fi
echo -e "${GREEN}✅ Сервер работает${NC}"
echo ""

# 1. Создание маршрута
echo -e "${YELLOW}🚀 Создание маршрута...${NC}"

ROUTE_JSON=$(curl -s -X POST "${API_BASE}/routes/generate" \
  -H "Content-Type: application/json" \
  -d "{
    \"start_point\": {\"lat\": $START_LAT, \"lon\": $START_LON},
    \"duration_minutes\": $DURATION,
    \"epochs\": $EPOCHS,
    \"interests\": $INTERESTS,
    \"max_waypoints\": $MAX_POINTS
  }")

# Проверка на ошибки
if echo "$ROUTE_JSON" | grep -q '"error"'; then
    echo -e "${RED}❌ Ошибка создания маршрута:${NC}"
    echo "$ROUTE_JSON" | python3 -m json.tool 2>/dev/null || echo "$ROUTE_JSON"
    exit 1
fi

ROUTE_ID=$(echo "$ROUTE_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin).get('route_id',''))" 2>/dev/null)

if [ -z "$ROUTE_ID" ]; then
    echo -e "${RED}❌ Не удалось получить route_id${NC}"
    echo "$ROUTE_JSON" | python3 -m json.tool 2>/dev/null || echo "$ROUTE_JSON"
    exit 1
fi

echo -e "${GREEN}✅ Маршрут создан: ${ROUTE_ID}${NC}"
echo ""

# 2. Показать точки маршрута
echo -e "${YELLOW}📍 Точки маршрута:${NC}"
echo "$ROUTE_JSON" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    for i, wp in enumerate(data['waypoints'], 1):
        print(f'  {i}. {wp[\"name\"]} ({wp[\"epoch\"]}, {wp[\"category\"]})')
except Exception as e:
    print(f'Ошибка: {e}')
" 2>/dev/null
echo ""

# 3. Проверка наличия Yandex API ключей
echo -e "${YELLOW}🔑 Проверка API ключей...${NC}"
if [ -z "$YANDEX_API_KEY" ]; then
    echo -e "${YELLOW}⚠️  YANDEX_API_KEY не установлен${NC}"
    echo "   Аудио не будет сгенерировано автоматически."
    echo "   Установите: export YANDEX_API_KEY=\"your-key\""
    echo ""
    echo -e "${YELLOW}Продолжить без генерации аудио? (y/n)${NC}"
    read -r CONTINUE
    if [ "$CONTINUE" != "y" ]; then
        echo "Отменено."
        exit 0
    fi
    NO_AUDIO=true
else
    echo -e "${GREEN}✅ YANDEX_API_KEY найден${NC}"
    NO_AUDIO=false
fi
echo ""

# 4. Ожидание генерации контента
if [ "$NO_AUDIO" = false ]; then
    WAIT_TIME=$((MAX_POINTS * 30))
    echo -e "${YELLOW}⏳ Ожидание генерации аудио (${WAIT_TIME}с)...${NC}"
    echo "   Генерация занимает ~30 секунд на точку"
    
    for i in $(seq $WAIT_TIME -1 1); do
        printf "\r   Осталось: %3d сек" $i
        sleep 1
    done
    echo ""
    echo ""
fi

# 5. Получение данных маршрута
echo -e "${YELLOW}📥 Получение данных маршрута...${NC}"
ROUTE_DATA=$(curl -s "${API_BASE}/routes/${ROUTE_ID}")

# 6. Проверка статуса контента
echo -e "${YELLOW}🔍 Статус генерации:${NC}"
echo ""

STATUS_OUTPUT=$(echo "$ROUTE_DATA" | python3 << 'EOF'
import sys, json

try:
    data = json.load(sys.stdin)
    ready = 0
    total = len(data['waypoints'])
    
    for i, wp in enumerate(data['waypoints'], 1):
        has_content = wp.get('content') is not None
        status = "✅ готово" if has_content else "⏳ генерируется"
        print(f"  {i}. {wp['name']}: {status}")
        if has_content:
            ready += 1
    
    print(f"\nГотово: {ready}/{total}")
    
    # Exit code: 0 если все готово, 1 если нет
    if ready < total:
        exit(1)
except Exception as e:
    print(f"Ошибка: {e}")
    exit(1)
EOF
)

echo "$STATUS_OUTPUT"
echo ""

# 7. Скачивание аудиофайлов
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Все аудио готовы!${NC}"
    echo ""
    echo -e "${YELLOW}🔊 Скачивание аудиофайлов...${NC}"
    
    # Создать директорию для аудио
    AUDIO_DIR="route_${ROUTE_ID:0:8}"
    mkdir -p "$AUDIO_DIR"
    
    echo "$ROUTE_DATA" | python3 << EOF
import sys, json, subprocess, os

data = json.load(sys.stdin)
downloaded = 0
failed = 0
audio_dir = "$AUDIO_DIR"

for i, wp in enumerate(data['waypoints'], 1):
    if wp.get('content') and wp['content']:
        waypoint_id = wp['id']
        # Безопасное имя файла
        safe_name = "".join(c if c.isalnum() or c in (' ','-','_') else '_' for c in wp['name'])
        filename = os.path.join(audio_dir, f"{i:02d}_{safe_name}.mp3")
        
        print(f"Скачивание {i}/{len(data['waypoints'])}: {wp['name']}")
        
        # Скачать файл
        result = subprocess.run([
            'curl', '-s', '-f',
            f'http://localhost:8080/api/audio/{waypoint_id}',
            '-o', filename
        ], capture_output=True)
        
        if result.returncode == 0 and os.path.exists(filename) and os.path.getsize(filename) > 0:
            size = os.path.getsize(filename) / 1024  # KB
            duration = wp['content'].get('duration_seconds', 0)
            print(f"  ✅ Сохранено: {filename} ({size:.1f} KB, {duration}s)\n")
            downloaded += 1
        else:
            print(f"  ❌ Ошибка скачивания\n")
            if os.path.exists(filename):
                os.remove(filename)
            failed += 1

print(f"{'='*60}")
print(f"Скачано успешно: {downloaded}")
if failed > 0:
    print(f"Ошибки: {failed}")
EOF
    
    echo ""
    echo -e "${GREEN}✨ Готово!${NC}"
    echo ""
    echo -e "${YELLOW}📂 Аудиофайлы сохранены в: ${AUDIO_DIR}/${NC}"
    echo ""
    ls -lh "$AUDIO_DIR"/*.mp3 2>/dev/null
    
    # Сохранить информацию о маршруте
    echo "$ROUTE_DATA" | python3 -m json.tool > "${AUDIO_DIR}/route_info.json"
    echo ""
    echo "Информация о маршруте: ${AUDIO_DIR}/route_info.json"
    
else
    echo -e "${YELLOW}⏳ Аудио еще генерируется${NC}"
    echo ""
    echo "Подождите 1-2 минуты и повторите попытку:"
    echo "  curl ${API_BASE}/routes/${ROUTE_ID} | python3 -m json.tool"
    echo ""
    echo "Или запустите этот скрипт снова для скачивания:"
    echo "  ROUTE_ID=${ROUTE_ID} ./download_audio.sh"
fi

echo ""
echo -e "${BLUE}════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  Route ID: ${ROUTE_ID}${NC}"
echo -e "${BLUE}════════════════════════════════════════════════${NC}"
