cd /var/www/pterodactyl

php artisan p:location:make --short="pterodactyl-location" --long="pterodactyl-location automatically added by pterodactyl-installer" --no-interaction

php artisan p:node:make --name="pterodactyl-node" \
--description="pterodactyl-node automatically added by pterodactyl-installer" \
--locationId="1" \
--no-interaction \
--fqdn="{{url}}" \
--public="1" \
--scheme="https" \
--maxMemory="4000" \
--overallocateMemory="-1" \
--maxDisk="4000" \
--overallocateDisk="-1"

mkdir /etc/pterodactyl
touch /etc/pterodactyl/config.yml
php artisan p:node:configuration 1 > /etc/pterodactyl/config.yml
sed -i 's/\/etc\/letsencrypt\/live\/{{url}}\/fullchain.pem/\/etc\/ssl\/pterodactyl-cert.pem/g' /etc/pterodactyl/config.yml
sed -i 's/\/etc\/letsencrypt\/live\/{{url}}\/privkey.pem/\/etc\/ssl\/pterodactyl-key.pem/g' /etc/pterodactyl/config.yml
sed -i 's/http:\/\/localhost\/{{url}}/https:\/\/{{url}}/g' /etc/pterodactyl/config.yml