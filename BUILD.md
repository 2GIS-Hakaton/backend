# 🔨 Сборка проекта для продакшена

Краткое руководство по сборке и запуску проекта на сервере.

---

## 🚀 Быстрая сборка

### На локальной машине (для деплоя на Linux сервер)

```bash
cd /Users/dimmy-kor/Projects/2gis_hackaton/backend

# Сборка для Linux (если вы на macOS/Windows)
GOOS=linux GOARCH=amd64 go build -o audioguid-server cmd/server/main.go

# Проверка
file audioguid-server
# Должно быть: audioguid-server: ELF 64-bit LSB executable, x86-64

# Размер
ls -lh audioguid-server
```

### На сервере Linux

```bash
cd /path/to/backend

# Простая сборка
go build -o audioguid-server cmd/server/main.go

# Или с оптимизацией (уменьшает размер бинарника)
go build -ldflags="-s -w" -o audioguid-server cmd/server/main.go
```

---

## 📦 Что нужно на сервере

### Минимальный набор файлов:

```
audioguid/
├── audioguid-server          # Скомпилированный бинарник
├── config.yaml                # Конфигурация (из backend/config/)
├── .env                       # Переменные окружения
├── keys/                      # API ключи
│   └── 2gis.key
├── audio/                     # Создастся автоматически
└── logs/                      # Для логов (опционально)
```

### Переменные окружения (.env):

```env
# Database
DATABASE_URL=postgresql://user:password@host:port/dbname?sslmode=disable

# Yandex Cloud
YANDEX_API_KEY=AQVN...
YANDEX_FOLDER_ID=b1g...

# 2GIS API
GAPIS_API_KEY=your_key
GAPIS_APP_ID=your_app_id
GAPIS_KEY_FILE=./keys/2gis.key

# Server
PORT=8080
GIN_MODE=release
```

---

## ▶️ Запуск на сервере

### Вариант 1: Простой запуск

```bash
cd ~/audioguid
./audioguid-server
```

### Вариант 2: В фоне с nohup

```bash
cd ~/audioguid
nohup ./audioguid-server > logs/app.log 2>&1 &
echo $! > audioguid.pid

# Проверить
ps aux | grep audioguid-server

# Остановить
kill $(cat audioguid.pid)
```

### Вариант 3: Systemd service (рекомендуется)

Создайте файл `/etc/systemd/system/audioguid.service`:

```ini
[Unit]
Description=Audioguid API Server
After=network.target

[Service]
Type=simple
User=your_user
WorkingDirectory=/home/your_user/audioguid
EnvironmentFile=/home/your_user/audioguid/.env
ExecStart=/home/your_user/audioguid/audioguid-server
Restart=on-failure
RestartSec=10
StandardOutput=append:/home/your_user/audioguid/logs/app.log
StandardError=append:/home/your_user/audioguid/logs/error.log

[Install]
WantedBy=multi-user.target
```

Затем:

```bash
sudo systemctl daemon-reload
sudo systemctl enable audioguid
sudo systemctl start audioguid
sudo systemctl status audioguid
```

---

## ✅ Проверка работы

```bash
# Health check
curl http://localhost:8080/api/health

# Список POI
curl http://localhost:8080/api/pois

# Swagger UI
curl http://localhost:8080/swagger/index.html
```

---

## 🎤 Голос аудиогида

**Голос зафиксирован в коде: `ermil` (мужской, эмоциональный)**

Переменная окружения `YANDEX_VOICE` больше не используется. Голос всегда будет `ermil`.

Если нужно изменить голос, отредактируйте файл:
```
backend/internal/services/tts_service.go
```

Строка 27:
```go
voice:    "ermil", // Измените здесь на нужный голос
```

Доступные голоса:
- **Мужские**: `ermil` (эмоциональный), `filipp` (нейтральный), `omazh` (глубокий), `zahar` (спокойный)
- **Женские**: `alena` (мягкий), `oksana` (нейтральный), `jane` (эмоциональный)

---

## 🔄 Обновление на сервере

```bash
# 1. Собрать новую версию локально
cd /Users/dimmy-kor/Projects/2gis_hackaton/backend
GOOS=linux GOARCH=amd64 go build -o audioguid-server cmd/server/main.go

# 2. Загрузить на сервер
scp audioguid-server user@server:~/audioguid/audioguid-server.new

# 3. На сервере - заменить
ssh user@server
cd ~/audioguid
sudo systemctl stop audioguid  # или: kill $(cat audioguid.pid)
mv audioguid-server audioguid-server.old
mv audioguid-server.new audioguid-server
chmod +x audioguid-server
sudo systemctl start audioguid  # или: nohup ./audioguid-server > logs/app.log 2>&1 &

# 4. Проверить
curl http://localhost:8080/api/health
```

---

## 📚 Дополнительная информация

- **Полная инструкция по деплою**: [DEPLOY.md](DEPLOY.md) (если создан)
- **Документация API**: http://localhost:8080/swagger/index.html
- **Руководство по запуску**: [GETTING_STARTED.md](GETTING_STARTED.md)

---

## 🐛 Частые проблемы

### "Failed to connect to database"
```bash
# Проверить доступность БД
psql -h host -p port -U user -d dbname -c "SELECT 1;"

# Проверить DATABASE_URL в .env
cat ~/audioguid/.env | grep DATABASE_URL
```

### "Permission denied"
```bash
chmod +x ~/audioguid/audioguid-server
```

### "Port already in use"
```bash
# Найти процесс на порту 8080
sudo lsof -i :8080

# Убить процесс
sudo kill -9 <PID>
```

### "YANDEX_API_KEY not set"
```bash
# Проверить .env
cat ~/audioguid/.env | grep YANDEX

# Убедиться, что файл загружается
# Если используете systemd, проверьте EnvironmentFile в service файле
```

---

<div align="center">
  <b>Успешной сборки! 🚀</b>
</div>
