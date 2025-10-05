#!/bin/bash
# Быстрая сборка для продакшена

echo "🔨 Сборка Audioguid для Linux..."

cd backend

# Сборка для Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o audioguid-server cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✅ Сборка успешна!"
    echo ""
    echo "📦 Бинарник: backend/audioguid-server"
    ls -lh audioguid-server
    echo ""
    echo "📋 Для запуска на сервере:"
    echo "   1. Скопируйте файл на сервер: scp audioguid-server user@server:~/audioguid/"
    echo "   2. На сервере: chmod +x ~/audioguid/audioguid-server"
    echo "   3. Запустите: ./audioguid-server"
    echo ""
    echo "🎤 Голос: ermil (зафиксирован в коде)"
else
    echo "❌ Ошибка сборки!"
    exit 1
fi
