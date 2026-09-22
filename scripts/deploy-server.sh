#!/usr/bin/env bash
# Запускается на сервере: git pull → сборка → перезапуск сервиса.
#
# Первый запуск: bash deploy-server.sh
# Последующие:   bash deploy-server.sh  (из любой директории)

set -euo pipefail

APP_DIR="/opt/autohub"
REPO_DIR="/opt/autohub-src"
SERVICE_NAME="autohub"
GO_VERSION="1.23.4"

# ── Цвета ─────────────────────────────────────────────────────────────────────
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'
ok()   { echo -e "${GREEN}✓ $1${NC}"; }
fail() { echo -e "${RED}✗ $1${NC}"; exit 1; }

echo "▶ Деплой Auto Hub на сервере"
echo ""

# ── 1. Go ─────────────────────────────────────────────────────────────────────
if ! command -v go &> /dev/null; then
  echo "→ [1/5] Устанавливаем Go ${GO_VERSION}..."
  cd /tmp
  wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
  export PATH=$PATH:/usr/local/go/bin
  echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
  ok "Go $(go version | awk '{print $3}') установлен"
else
  export PATH=$PATH:/usr/local/go/bin
  echo "→ [1/5] Go уже установлен: $(go version | awk '{print $3}')"
fi

# ── 2. Git pull ───────────────────────────────────────────────────────────────
echo "→ [2/5] Получаем код из репозитория..."

if [ ! -d "$REPO_DIR/.git" ]; then
  # Первый запуск — клонируем
  if [ -z "${REPO_URL:-}" ]; then
    echo ""
    echo "  Репозиторий не найден. Укажи URL:"
    echo "  REPO_URL=https://github.com/user/repo bash deploy-server.sh"
    fail "REPO_URL не задан"
  fi
  sudo mkdir -p "$REPO_DIR"
  sudo chown "$(whoami)" "$REPO_DIR"
  git clone "$REPO_URL" "$REPO_DIR"
  ok "Репозиторий склонирован в $REPO_DIR"
else
  cd "$REPO_DIR"
  git pull
  COMMIT=$(git log -1 --format="%h %s")
  ok "Обновлено до: $COMMIT"
fi

# ── 3. Сборка ─────────────────────────────────────────────────────────────────
echo "→ [3/5] Сборка бинарника..."
cd "$REPO_DIR"

# Установить swag если нет
if ! command -v swag &> /dev/null && [ ! -f "$(go env GOPATH)/bin/swag" ]; then
  echo "   Устанавливаем swag..."
  go install github.com/swaggo/swag/cmd/swag@latest
fi
SWAG="$(go env GOPATH)/bin/swag"

# Генерация swagger docs
"$SWAG" init -g cmd/api/main.go --output docs -q

mkdir -p bin
go build -ldflags="-s -w" -o bin/api  ./cmd/api/main.go
go build -ldflags="-s -w" -o bin/seed ./cmd/seed/main.go
ok "Собрано: bin/api ($(du -sh bin/api | cut -f1))"

# ── 4. Обновление файлов ──────────────────────────────────────────────────────
echo "→ [4/5] Обновляем файлы в $APP_DIR..."
sudo mkdir -p "$APP_DIR"
sudo cp bin/api  "$APP_DIR/api"
sudo cp bin/seed "$APP_DIR/seed"
sudo cp -r migrations/ "$APP_DIR/migrations"
sudo chmod +x "$APP_DIR/api" "$APP_DIR/seed"
ok "Файлы обновлены"

# ── 5. Перезапуск ─────────────────────────────────────────────────────────────
echo "→ [5/5] Перезапускаем сервис..."
sudo systemctl restart "$SERVICE_NAME"
sleep 2

STATUS=$(systemctl is-active "$SERVICE_NAME")
if [ "$STATUS" = "active" ]; then
  ok "Сервис запущен"
else
  echo ""
  echo "Логи сервиса:"
  sudo journalctl -u "$SERVICE_NAME" -n 20 --no-pager
  fail "Сервис не запустился"
fi

# ── Итог ─────────────────────────────────────────────────────────────────────
echo ""
PORT=$(grep PORT "$APP_DIR/.env" 2>/dev/null | cut -d= -f2 || echo "8080")
IP=$(hostname -I | awk '{print $1}')
echo "  API:     http://${IP}:${PORT}/api/health"
echo "  Swagger: http://${IP}:${PORT}/swagger/index.html"
echo ""
sudo journalctl -u "$SERVICE_NAME" -n 5 --no-pager
