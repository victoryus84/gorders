#!/bin/sh

echo "=== Certificate Status ==="
certbot certificates

echo "=== Renewal Logs ==="
tail -20 /var/log/certbot-renewal.log 2>/dev/null || echo "No renewal logs yet"