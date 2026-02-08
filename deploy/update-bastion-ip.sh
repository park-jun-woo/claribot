#!/bin/bash
# Bastion IP 변경 시 실행하는 스크립트
# CloudFront origin은 bastion.parkjunwoo.com DNS를 사용하므로
# Route 53 A 레코드만 변경하면 됨
# 사용법: ./deploy/update-bastion-ip.sh <NEW_IP>
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/bastion.env"

NEW_IP=$1
if [ -z "$NEW_IP" ]; then
    echo "Usage: $0 <NEW_IP>"
    exit 1
fi

OLD_IP=$BASTION_IP
echo "Updating bastion IP: $OLD_IP -> $NEW_IP"

# 1. Route 53 bastion.parkjunwoo.com A 레코드 업데이트
echo "=== Updating Route 53 A record ==="
aws route53 change-resource-record-sets \
    --hosted-zone-id "$ROUTE53_ZONE_ID" \
    --change-batch "{
        \"Changes\": [{
            \"Action\": \"UPSERT\",
            \"ResourceRecordSet\": {
                \"Name\": \"bastion.parkjunwoo.com\",
                \"Type\": \"A\",
                \"TTL\": 60,
                \"ResourceRecords\": [{\"Value\": \"$NEW_IP\"}]
            }
        }]
    }"
echo "DNS updated: bastion.parkjunwoo.com -> $NEW_IP"

# 2. bastion.env 업데이트
echo "=== Updating bastion.env ==="
sed -i "s/^BASTION_IP=.*/BASTION_IP=$NEW_IP/" "$SCRIPT_DIR/bastion.env"
echo "bastion.env updated"

# 3. SSH tunnel 재연결
echo "=== Restarting SSH tunnel ==="
systemctl --user restart claribot-tunnel 2>/dev/null || echo "Tunnel service not found (manual restart needed)"

echo ""
echo "=== Done ==="
echo "New bastion IP: $NEW_IP"
echo "CloudFront config 변경 불필요 (bastion.parkjunwoo.com DNS 사용)"
