#!/bin/bash
# CloudFront + ACM + Route 53 설정 스크립트
# 사용법: ./deploy/setup-cloudfront.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/bastion.env"

if [ -z "$BASTION_IP" ]; then
    echo "ERROR: BASTION_IP not set in bastion.env"
    exit 1
fi

# ── Phase 1: ACM 인증서 발급 ──────────────────────────────────────────

echo "=== Phase 1: ACM 인증서 발급 ==="

if [ -z "$ACM_CERTIFICATE_ARN" ]; then
    echo "Requesting ACM certificate for clari.parkjunwoo.com..."
    ACM_CERTIFICATE_ARN=$(aws acm request-certificate \
        --domain-name clari.parkjunwoo.com \
        --validation-method DNS \
        --region us-east-1 \
        --query 'CertificateArn' \
        --output text)
    echo "Certificate ARN: $ACM_CERTIFICATE_ARN"
    sed -i "s|^ACM_CERTIFICATE_ARN=.*|ACM_CERTIFICATE_ARN=$ACM_CERTIFICATE_ARN|" "$SCRIPT_DIR/bastion.env"

    # DNS 검증 레코드 추가
    echo "Waiting for certificate details..."
    sleep 5

    CERT_DETAILS=$(aws acm describe-certificate \
        --certificate-arn "$ACM_CERTIFICATE_ARN" \
        --region us-east-1 \
        --query 'Certificate.DomainValidationOptions[0].ResourceRecord')

    CNAME_NAME=$(echo "$CERT_DETAILS" | jq -r '.Name')
    CNAME_VALUE=$(echo "$CERT_DETAILS" | jq -r '.Value')

    if [ -z "$ROUTE53_ZONE_ID" ]; then
        ROUTE53_ZONE_ID=$(aws route53 list-hosted-zones-by-name \
            --dns-name parkjunwoo.com \
            --query 'HostedZones[0].Id' \
            --output text | sed 's|/hostedzone/||')
        sed -i "s|^ROUTE53_ZONE_ID=.*|ROUTE53_ZONE_ID=$ROUTE53_ZONE_ID|" "$SCRIPT_DIR/bastion.env"
    fi

    echo "Adding DNS validation record: $CNAME_NAME -> $CNAME_VALUE"
    aws route53 change-resource-record-sets \
        --hosted-zone-id "$ROUTE53_ZONE_ID" \
        --change-batch "{
            \"Changes\": [{
                \"Action\": \"UPSERT\",
                \"ResourceRecordSet\": {
                    \"Name\": \"$CNAME_NAME\",
                    \"Type\": \"CNAME\",
                    \"TTL\": 300,
                    \"ResourceRecords\": [{\"Value\": \"$CNAME_VALUE\"}]
                }
            }]
        }"

    echo "Waiting for certificate validation (this may take a few minutes)..."
    aws acm wait certificate-validated \
        --certificate-arn "$ACM_CERTIFICATE_ARN" \
        --region us-east-1
    echo "Certificate validated!"
else
    echo "Using existing certificate: $ACM_CERTIFICATE_ARN"
fi

# ── Phase 2: CloudFront 배포 생성 ─────────────────────────────────────

echo ""
echo "=== Phase 2: CloudFront 배포 생성 ==="

# CloudFront 시크릿 헤더 생성
if [ -z "$CF_SECRET" ]; then
    CF_SECRET=$(openssl rand -hex 32)
    sed -i "s|^CF_SECRET=.*|CF_SECRET=$CF_SECRET|" "$SCRIPT_DIR/bastion.env"
    echo "Generated CF_SECRET"
fi

if [ -z "$CF_DISTRIBUTION_ID" ]; then
    echo "Creating CloudFront distribution..."

    CF_CONFIG=$(cat <<CFEOF
{
    "CallerReference": "claribot-$(date +%s)",
    "Aliases": {
        "Quantity": 1,
        "Items": ["clari.parkjunwoo.com"]
    },
    "DefaultRootObject": "",
    "Origins": {
        "Quantity": 1,
        "Items": [{
            "Id": "bastion-origin",
            "DomainName": "$BASTION_IP",
            "CustomOriginConfig": {
                "HTTPPort": 9847,
                "HTTPSPort": 443,
                "OriginProtocolPolicy": "http-only",
                "OriginSslProtocols": {"Quantity": 1, "Items": ["TLSv1.2"]}
            },
            "CustomHeaders": {
                "Quantity": 1,
                "Items": [{
                    "HeaderName": "X-CF-Secret",
                    "HeaderValue": "$CF_SECRET"
                }]
            }
        }]
    },
    "DefaultCacheBehavior": {
        "TargetOriginId": "bastion-origin",
        "ViewerProtocolPolicy": "redirect-to-https",
        "AllowedMethods": {
            "Quantity": 7,
            "Items": ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"],
            "CachedMethods": {"Quantity": 2, "Items": ["GET", "HEAD"]}
        },
        "CachePolicyId": "4135ea2d-6df8-44a3-9df3-4b5a84be39ad",
        "OriginRequestPolicyId": "216adef6-5c7f-47e4-b989-5492eafa07d3",
        "Compress": true
    },
    "Comment": "Claribot Bastion Relay",
    "Enabled": true,
    "ViewerCertificate": {
        "ACMCertificateArn": "$ACM_CERTIFICATE_ARN",
        "SSLSupportMethod": "sni-only",
        "MinimumProtocolVersion": "TLSv1.2_2021"
    },
    "HttpVersion": "http2"
}
CFEOF
    )

    RESULT=$(aws cloudfront create-distribution \
        --distribution-config "$CF_CONFIG" \
        --query 'Distribution.{Id:Id,DomainName:DomainName}' \
        --output json)

    CF_DISTRIBUTION_ID=$(echo "$RESULT" | jq -r '.Id')
    CF_DOMAIN=$(echo "$RESULT" | jq -r '.DomainName')

    sed -i "s|^CF_DISTRIBUTION_ID=.*|CF_DISTRIBUTION_ID=$CF_DISTRIBUTION_ID|" "$SCRIPT_DIR/bastion.env"

    echo "Distribution ID: $CF_DISTRIBUTION_ID"
    echo "Domain: $CF_DOMAIN"
else
    echo "Using existing distribution: $CF_DISTRIBUTION_ID"
    CF_DOMAIN=$(aws cloudfront get-distribution \
        --id "$CF_DISTRIBUTION_ID" \
        --query 'Distribution.DomainName' \
        --output text)
    echo "Domain: $CF_DOMAIN"
fi

# ── Phase 3: Route 53 DNS ─────────────────────────────────────────────

echo ""
echo "=== Phase 3: Route 53 DNS ==="

if [ -z "$ROUTE53_ZONE_ID" ]; then
    ROUTE53_ZONE_ID=$(aws route53 list-hosted-zones-by-name \
        --dns-name parkjunwoo.com \
        --query 'HostedZones[0].Id' \
        --output text | sed 's|/hostedzone/||')
    sed -i "s|^ROUTE53_ZONE_ID=.*|ROUTE53_ZONE_ID=$ROUTE53_ZONE_ID|" "$SCRIPT_DIR/bastion.env"
fi

echo "Setting up ALIAS record: clari.parkjunwoo.com -> $CF_DOMAIN"
aws route53 change-resource-record-sets \
    --hosted-zone-id "$ROUTE53_ZONE_ID" \
    --change-batch "{
        \"Changes\": [{
            \"Action\": \"UPSERT\",
            \"ResourceRecordSet\": {
                \"Name\": \"clari.parkjunwoo.com\",
                \"Type\": \"A\",
                \"AliasTarget\": {
                    \"HostedZoneId\": \"Z2FDTNDATAQYW2\",
                    \"DNSName\": \"$CF_DOMAIN\",
                    \"EvaluateTargetHealth\": false
                }
            }
        }]
    }"

echo "DNS record created!"

echo ""
echo "=== Setup Complete ==="
echo "Distribution ID: $CF_DISTRIBUTION_ID"
echo "CF Domain: $CF_DOMAIN"
echo "Custom Domain: clari.parkjunwoo.com"
echo "CF Secret: $CF_SECRET"
echo ""
echo "Next steps:"
echo "  1. Run deploy/setup-bastion.sh to configure nginx on bastion"
echo "  2. Start the SSH tunnel: systemctl --user start claribot-tunnel"
echo "  3. Wait for CloudFront deployment (may take 5-15 minutes)"
echo "  4. Run deploy/verify-relay.sh to verify"
