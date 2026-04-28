#!/bin/bash

# Загружает переменные окружения из файла .env, если он существует
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Проверяет наличие docker и docker compose, выбирает подходящую команду
if command -v docker >/dev/null 2>&1; then
  if docker compose version >/dev/null 2>&1; then
    DCMD="docker compose"
  elif command -v docker-compose >/dev/null 2>&1; then
    DCMD="docker-compose"
  else
    echo 'Error: docker compose is not installed.' >&2
    exit 1
  fi
else
  echo 'Error: docker is not installed.' >&2
  exit 1
fi

# Устанавливает переменные: домены, размер ключа, путь к данным, email, режим тестирования
domains=(${DOMAIN:-example.com})
rsa_key_size=4096
data_path="./nginx/certbot"
email="${EMAIL:-hello@example.com}" # Рекомендуется указать реальный email
staging=0 # Если 1 — используется тестовый режим Let's Encrypt

# Проверяет, есть ли уже сертификаты, и спрашивает, нужно ли их заменить
if [ -d "$data_path" ]; then
    read -p "Existing data found for $domains. Continue and replace existing certificate? (y/N) " decision
    if [ "$decision" != "Y" ] && [ "$decision" != "y" ]; then
        exit
    fi
fi

# Скачивает рекомендуемые параметры TLS, если их нет
if [ ! -e "$data_path/conf/options-ssl-nginx.conf" ] || [ ! -e "$data_path/conf/ssl-dhparams.pem" ]; then
    echo "### Downloading recommended TLS parameters ..."
    mkdir -p "$data_path/conf"
    curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot-nginx/certbot_nginx/_internal/tls_configs/options-ssl-nginx.conf >"$data_path/conf/options-ssl-nginx.conf"
    curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot/certbot/ssl-dhparams.pem >"$data_path/conf/ssl-dhparams.pem"
    echo
fi

# Создаёт временный (dummy) сертификат для запуска nginx
echo "### Creating dummy certificate for $domains ..."
path="/etc/letsencrypt/live/$domains"
mkdir -p "$data_path/conf/live/$domains"
docker compose -f "docker-compose.yml" run --rm --entrypoint "\
  openssl req -x509 -nodes -newkey rsa:$rsa_key_size -days 1\
    -keyout '$path/privkey.pem' \
    -out '$path/fullchain.pem' \
    -subj '/CN=localhost'" certbot
echo

# Запускает nginx с временным сертификатом
echo "### Starting nginx ..."
docker compose  -f "docker-compose.yml" up --force-recreate -d nginx
echo

# Удаляет временный сертификат
echo "### Deleting dummy certificate for $domains ..."
docker compose  -f "docker-compose.yml" run --rm --entrypoint "\
  rm -Rf /etc/letsencrypt/live/$domains && \
  rm -Rf /etc/letsencrypt/archive/$domains && \
  rm -Rf /etc/letsencrypt/renewal/$domains.conf" certbot
echo

# Запрашивает реальный сертификат Let's Encrypt
echo "### Requesting Let's Encrypt certificate for $domains ..."
# Формирует аргументы для доменов
domain_args=""
for domain in "${domains[@]}"; do
    domain_args="$domain_args -d $domain"
done

# Формирует аргумент для email
case "$email" in
"") email_arg="--register-unsafely-without-email" ;;
*) email_arg="--email $email" ;;
esac

# Включает тестовый режим, если нужно
if [ $staging != "0" ]; then staging_arg="--staging"; fi

# Запускает certbot для получения сертификата
docker compose -f "docker-compose.yml" run --rm --entrypoint "\
  certbot certonly --webroot -w /var/www/certbot \
    $staging_arg \
    $email_arg \
    $domain_args \
    --rsa-key-size $rsa_key_size \
    --agree-tos \
    --force-renewal" certbot
echo

# Перезапускает nginx, чтобы он использовал новый сертификат
docker compose -f "docker-compose.yml" exec nginx nginx -s reload