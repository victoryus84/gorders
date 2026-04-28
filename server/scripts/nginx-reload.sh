#!/bin/sh
set -e

echo "Nginx reload script started"

# Ждем появления SSL конфигурации от certbot
echo "Waiting for SSL configuration..."
while [ ! -f /etc/letsencrypt/nginx-ssl.conf ]; do
    sleep 5
    echo "Still waiting for SSL config..."
done

echo "SSL configuration found, applying to nginx config..."

# Вставляем SSL конфигурацию в nginx.conf
sed -i '/# SSL_CONFIG_PLACEHOLDER/r /etc/letsencrypt/nginx-ssl.conf' /etc/nginx/nginx.conf

echo "Nginx configuration updated with SSL settings"