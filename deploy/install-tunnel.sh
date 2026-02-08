#!/bin/bash
# SSH reverse tunnel 서비스 설치 스크립트
# 사용법: ./deploy/install-tunnel.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/bastion.env"

# autossh 설치
if ! command -v autossh &> /dev/null; then
    echo "Installing autossh..."
    sudo apt install -y autossh
fi

# systemd user 디렉토리 생성
SYSTEMD_DIR="$HOME/.config/systemd/user"
mkdir -p "$SYSTEMD_DIR"

# bastion key 경로 확인
KEY_PATH="$HOME/git/bastion/bastion-key.pem"
if [ ! -f "$KEY_PATH" ]; then
    echo "WARNING: Bastion key not found at $KEY_PATH"
    echo "Update the service file with the correct key path"
fi

# 서비스 파일 복사 (IP와 key 경로를 bastion.env에서 반영)
sed -e "s|ubuntu@43.201.60.23|${BASTION_USER}@${BASTION_IP}|" \
    -e "s|%h/git/bastion/bastion-key.pem|${KEY_PATH}|" \
    "$SCRIPT_DIR/claribot-tunnel.service" > "$SYSTEMD_DIR/claribot-tunnel.service"

echo "Service file installed to $SYSTEMD_DIR/claribot-tunnel.service"

# systemd 리로드 및 활성화
systemctl --user daemon-reload
systemctl --user enable claribot-tunnel

echo ""
echo "=== Tunnel 서비스 설치 완료 ==="
echo "시작: systemctl --user start claribot-tunnel"
echo "상태: systemctl --user status claribot-tunnel"
echo "로그: journalctl --user -u claribot-tunnel -f"
