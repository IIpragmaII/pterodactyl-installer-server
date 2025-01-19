cd /var/www/pterodactyl

COMPOSER_ALLOW_SUPERUSER=1 composer install --no-dev --optimize-autoloader

# Only run the command below if you are installing this Panel for
# the first time and do not have any Pterodactyl Panel data in the database.
php artisan key:generate --force

php artisan p:environment:setup --author="{{email}}" \
    --url="{{url}}" \
    --timezone="{{timezone}}" \
    --cache="redis" \
    --session="redis" \
    --queue="redis" \
    --redis-host="localhost" \
    --redis-pass="null" \
    --redis-port="6379" \
    --settings-ui=true

php artisan p:environment:database  --host="127.0.0.1" \
    --port="3306" \
    --database="panel" \
    --username="pterodactyl" \
    --password="{{db_password}}"


php artisan migrate --seed --force

php artisan p:user:make --email="{{email}}" \
    --username="{{username}}" \
    --name-first="{{first_name}}" \
    --name-last="{{last_name}}" \
    --password="{{password}}" \
    --admin=1


# If using NGINX, Apache or Caddy (not on RHEL / Rocky Linux / AlmaLinux)
chown -R www-data:www-data /var/www/pterodactyl/*

cp /etc/ssl/pterodactyl-cert.pem /usr/local/share/ca-certificates/pterodactyl-cert.crt

update-ca-certificates
