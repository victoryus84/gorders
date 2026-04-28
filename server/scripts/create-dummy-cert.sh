#!/bin/sh
set -e

: "${DOMAIN:?Need to set DOMAIN env var}"

CERT_DIR="/etc/nginx/ssl"

if [ -f "${CERT_DIR}/nginx-selfsigned.crt" ] && [ -f "${CERT_DIR}/nginx-selfsigned.key" ]; then
  echo "$(date): Real certs exist for ${DOMAIN}"
else
  echo "$(date): Creating dummy cert for ${DOMAIN}"
  mkdir -p "${CERT_DIR}"
  
  openssl req -x509 -nodes -newkey rsa:2048 -days 365 \
    -keyout "${CERT_DIR}/nginx-selfsigned.key" \
    -out "${CERT_DIR}/nginx-selfsigned.crt" \
    -subj "/C=MD/ST=Chisinau/L=Chisinau/O=Servidar SRL/OU=IT/CN=${DOMAIN}" \
    -addext "subjectAltName=DNS:${DOMAIN},IP:217.26.172.96"
  
  chmod 600 "${CERT_DIR}/nginx-selfsigned.key" "${CERT_DIR}/nginx-selfsigned.crt"
  echo "$(date): Dummy cert created at ${CERT_DIR}"
fi

# Запускаем nginx
exec nginx -g "daemon off;"
