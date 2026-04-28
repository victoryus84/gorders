#!/bin/bash
set -e

cd /opt/orders

echo "=== SSL Manager ==="
echo "Domain: $DOMAIN"
echo "Email: $EMAIL"

case "$1" in
    "init")
        echo "Starting initial SSL setup..."
        docker compose up -d nginx
        docker compose run --rm certbot
        docker compose restart nginx
        docker compose up -d certbot-renew
        echo "✅ SSL setup completed"
        ;;
    "renew")
        echo "Manually renewing certificates..."
        docker compose run --rm certbot-renew /renew-certificates.sh
        ;;
    "logs")
        echo "Showing renewal logs:"
        docker compose logs certbot-renew
        ;;
    "test")
        echo "Testing certificate:"
        curl -I "https://$DOMAIN/health"
        ;;
    *)
        echo "Usage: $0 {init|renew|logs|test}"
        exit 1
        ;;
esac