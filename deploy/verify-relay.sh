#!/bin/bash
# Bastion Relay 전체 검증 스크립트
# 사용법: ./deploy/verify-relay.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/bastion.env"

SSH_CMD="ssh -i $BASTION_KEY ${BASTION_USER}@${BASTION_IP}"
PASS=0
FAIL=0

check() {
    local name=$1
    local result=$2
    if [ $result -eq 0 ]; then
        echo "  PASS: $name"
        PASS=$((PASS + 1))
    else
        echo "  FAIL: $name"
        FAIL=$((FAIL + 1))
    fi
}

echo "=== Bastion Relay 검증 ==="
echo ""

# 1. SSH tunnel 동작 확인
echo "[1] SSH tunnel (bastion:19847 -> local:9847)"
TUNNEL_RESULT=$($SSH_CMD "curl -sf http://127.0.0.1:19847/api/health" 2>/dev/null) && RC=0 || RC=1
check "tunnel responds" $RC
if [ $RC -eq 0 ]; then
    echo "       Response: $TUNNEL_RESULT"
fi

# 2. nginx 동작 확인 (시크릿 헤더 포함)
echo "[2] nginx reverse proxy (bastion:9847 with secret)"
NGINX_RESULT=$($SSH_CMD "curl -sf -H 'X-CF-Secret: $CF_SECRET' http://127.0.0.1:9847/api/health" 2>/dev/null) && RC=0 || RC=1
check "nginx with secret header" $RC
if [ $RC -eq 0 ]; then
    echo "       Response: $NGINX_RESULT"
fi

# 3. nginx 시크릿 없이 접근 차단 확인
echo "[3] nginx without secret (should be 403)"
HTTP_CODE=$($SSH_CMD "curl -sf -o /dev/null -w '%{http_code}' http://127.0.0.1:9847/api/health" 2>/dev/null) || HTTP_CODE="000"
if [ "$HTTP_CODE" = "403" ]; then
    check "secret header enforcement" 0
else
    check "secret header enforcement (got $HTTP_CODE, expected 403)" 1
fi

# 4. CloudFront 동작 확인
echo "[4] CloudFront -> HTTPS endpoint"
CF_RESULT=$(curl -sf https://clari.parkjunwoo.com/api/health 2>/dev/null) && RC=0 || RC=1
check "CloudFront endpoint" $RC
if [ $RC -eq 0 ]; then
    echo "       Response: $CF_RESULT"
fi

# 5. IP 직접 접근 차단 확인
echo "[5] Direct IP access (should be blocked)"
DIRECT_CODE=$(curl -sf -o /dev/null -w '%{http_code}' --connect-timeout 5 "http://${BASTION_IP}:9847/api/health" 2>/dev/null) || DIRECT_CODE="000"
if [ "$DIRECT_CODE" = "403" ] || [ "$DIRECT_CODE" = "000" ]; then
    check "direct IP blocked (code: $DIRECT_CODE)" 0
else
    check "direct IP blocked (got $DIRECT_CODE, expected 403 or timeout)" 1
fi

# 6. systemd tunnel 서비스 상태
echo "[6] Tunnel service status"
systemctl --user is-active claribot-tunnel &>/dev/null && RC=0 || RC=1
check "tunnel service active" $RC

echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="
