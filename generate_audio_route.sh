#!/bin/bash

# 🎧 Генерация аудиогида по маршруту (синхронно)
# Этот скрипт отправляет запрос и сразу получает готовый MP3 файл

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  🎧 Генерация аудиогида${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Параметры маршрута (можно изменить)
START_LAT=${1:-55.7558}
START_LON=${2:-37.6173}
DURATION=${3:-60}
EPOCHS=${4:-'["soviet","medieval"]'}
INTERESTS=${5:-'["architecture","history"]'}
MAX_POINTS=${6:-3}

echo -e "${YELLOW}📍 Параметры маршрута:${NC}"
echo "  Начальная точка: $START_LAT, $START_LON"
echo "  Длительность: $DURATION минут"
echo "  Эпохи: $EPOCHS"
echo "  Интересы: $INTERESTS"
echo "  Макс. точек: $MAX_POINTS"
echo ""

# Проверка сервера
echo -e "${YELLOW}📡 Проверка сервера...${NC}"
if ! curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
    echo -e "${RED}❌ Сервер не запущен!${NC}"
    echo ""
    echo "Запустите сервер:"
    echo "  cd backend"
    echo "  export DATABASE_URL=\"postgresql://nike:changeme@51.250.86.178:5432/audioguid?sslmode=disable\""
    echo "  export YANDEX_API_KEY=\"your-key\""
    echo "  export YANDEX_FOLDER_ID=\"your-folder\""
    echo "  go run cmd/server/main.go"
    exit 1
fi
echo -e "${GREEN}✅ Сервер работает${NC}"
echo ""

# Проверка Yandex API ключей
if [ -z "$YANDEX_API_KEY" ]; then
    echo -e "${YELLOW}⚠️  ВНИМАНИЕ: YANDEX_API_KEY не установлен${NC}"
    echo "   Для генерации аудио нужны Yandex API ключи."
    echo ""
    echo "   Получите ключи: https://console.cloud.yandex.ru/"
    echo "   Затем установите:"
    echo "     export YANDEX_API_KEY=\"your-key\""
    echo "     export YANDEX_FOLDER_ID=\"your-folder\""
    echo ""
    echo "   Продолжить без аудио? (y/n)"
    read -r CONTINUE
    if [ "$CONTINUE" != "y" ]; then
        exit 0
    fi
    echo ""
fi

# Генерация аудиогида
OUTPUT_FILE="route_audio_$(date +%Y%m%d_%H%M%S).mp3"

echo -e "${YELLOW}🚀 Генерация аудиогида...${NC}"
echo "   Это может занять 1-3 минуты в зависимости от количества точек"
echo "   (примерно 30-60 секунд на точку)"
echo ""

# Показываем прогресс
echo -e "${BLUE}⏳ Ожидание...${NC}"

# Отправляем запрос
HTTP_CODE=$(curl -w "%{http_code}" -o "$OUTPUT_FILE" -s -X POST \
  http://localhost:8080/api/routes/generate-audio \
  -H "Content-Type: application/json" \
  -d "{
    \"start_point\": {\"lat\": $START_LAT, \"lon\": $START_LON},
    \"duration_minutes\": $DURATION,
    \"epochs\": $EPOCHS,
    \"interests\": $INTERESTS,
    \"max_waypoints\": $MAX_POINTS
  }")

echo ""

# Проверяем результат
if [ "$HTTP_CODE" = "200" ]; then
    if [ -f "$OUTPUT_FILE" ] && [ -s "$OUTPUT_FILE" ]; then
        FILE_SIZE=$(ls -lh "$OUTPUT_FILE" | awk '{print $5}')
        echo -e "${GREEN}✅ Аудиогид успешно сгенерирован!${NC}"
        echo ""
        echo "📁 Файл: $OUTPUT_FILE"
        echo "📊 Размер: $FILE_SIZE"
        echo ""
        echo -e "${BLUE}🎵 Воспроизвести:${NC}"
        echo "  afplay $OUTPUT_FILE  # macOS"
        echo "  mpg123 $OUTPUT_FILE  # Linux"
        echo "  vlc $OUTPUT_FILE     # VLC player"
        echo ""
        
        # Автоматическое воспроизведение на macOS
        if command -v afplay &> /dev/null; then
            echo -e "${YELLOW}Воспроизвести сейчас? (y/n)${NC}"
            read -r PLAY
            if [ "$PLAY" = "y" ]; then
                afplay "$OUTPUT_FILE"
            fi
        fi
    else
        echo -e "${RED}❌ Файл пустой или не создан${NC}"
        rm -f "$OUTPUT_FILE"
        exit 1
    fi
else
    echo -e "${RED}❌ Ошибка: HTTP $HTTP_CODE${NC}"
    echo ""
    if [ -f "$OUTPUT_FILE" ]; then
        echo "Ответ сервера:"
        cat "$OUTPUT_FILE" | python3 -m json.tool 2>/dev/null || cat "$OUTPUT_FILE"
        rm -f "$OUTPUT_FILE"
    fi
    exit 1
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  ✅ Готово!${NC}"
echo -e "${BLUE}========================================${NC}"
