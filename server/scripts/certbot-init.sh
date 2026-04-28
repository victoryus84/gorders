#!/bin/sh
set -e

echo "Starting initial certificate request for $DOMAIN"

certbot certonly \
  --webroot \
  -w /var/www/certbot \
  --email "$EMAIL" \
  -d "$DOMAIN" \
  --agree-tos \
  --no-eff-email \
  --force-renewal

echo "Certificate successfully obtained for $DOMAIN"