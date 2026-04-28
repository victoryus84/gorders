#!/bin/sh
# Health check для nginx с поддержкой reload
if curl -f http://api:8080/health > /dev/null 2>&1; then
    exit 0
else
    exit 1
fi