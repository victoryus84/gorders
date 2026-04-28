#!/bin/sh
set -e

echo "$(date): Starting certbot renew check..."

# Проверяем что сертификаты существуют
if [ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
    echo "ERROR: Certificate not found at /etc/letsencrypt/live/${DOMAIN}/fullchain.pem"
    echo "Waiting for certbot-init to complete..."
    exit 1
fi

# Проверяем срок действия сертификата
DAYS_REMAINING=$(openssl x509 -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" -noout -text | grep -A1 "Not After" | tail -1 | awk '{print $1,$2,$4}')
echo "Certificate valid until: $DAYS_REMAINING"

# Выполняем renew
echo "Running certbot renew..."
if certbot renew --quiet; then
    echo "$(date): Certificate renewal successful"
else
    echo "$(date): Certificate renewal failed"
fi
