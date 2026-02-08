#!/bin/bash
# Bastion nginx 설정 스크립트
# 사용법: ./deploy/setup-bastion.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/bastion.env"

if [ -z "$BASTION_IP" ] || [ -z "$CF_SECRET" ]; then
    echo "ERROR: BASTION_IP and CF_SECRET must be set in bastion.env"
    echo "Run setup-cloudfront.sh first to generate CF_SECRET"
    exit 1
fi

SSH_CMD="ssh -i $BASTION_KEY ${BASTION_USER}@${BASTION_IP}"

echo "=== Bastion nginx 설정 ==="

# nginx 설치
echo "Installing nginx..."
$SSH_CMD "sudo apt update && sudo apt install -y nginx"

# nginx 설정 파일 전송
echo "Configuring nginx..."
NGINX_CONF=$(cat <<NGINXEOF
server {
    listen 9847;

    # CloudFront 시크릿 헤더 검증
    if (\$http_x_cf_secret != "$CF_SECRET") {
        return 403;
    }

    location / {
        proxy_pass http://127.0.0.1:19847;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";

        # 타임아웃 설정 (긴 작업 대비)
        proxy_read_timeout 300s;
        proxy_send_timeout 300s;
    }
}
NGINXEOF
)

$SSH_CMD "echo '$NGINX_CONF' | sudo tee /etc/nginx/sites-available/claribot > /dev/null"

# 심볼릭 링크 생성 (기존 있으면 덮어쓰기)
$SSH_CMD "sudo ln -sf /etc/nginx/sites-available/claribot /etc/nginx/sites-enabled/"

# default 사이트 비활성화 (충돌 방지)
$SSH_CMD "sudo rm -f /etc/nginx/sites-enabled/default"

# nginx 설정 테스트 및 리로드
echo "Testing and reloading nginx..."
$SSH_CMD "sudo nginx -t && sudo systemctl reload nginx"

# UFW 포트 오픈
echo "Opening port 9847 in UFW..."
$SSH_CMD "sudo ufw allow 9847/tcp" || echo "UFW not active or already allowed"

echo ""
echo "=== Bastion 설정 완료 ==="
echo "nginx listening on :9847 -> 127.0.0.1:19847"
echo "CloudFront secret header validation enabled"
