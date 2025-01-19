mariadb -u root -e "CREATE USER 'pterodactyl'@'127.0.0.1' IDENTIFIED BY '{{db_password}}';
CREATE DATABASE panel;
GRANT ALL PRIVILEGES ON panel.* TO 'pterodactyl'@'127.0.0.1' WITH GRANT OPTION;"