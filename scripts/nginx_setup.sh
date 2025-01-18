rm /etc/nginx/sites-enabled/default

# You do not need to symlink this file if you are using RHEL, Rocky Linux, or AlmaLinux.
sudo ln -s /etc/nginx/sites-available/pterodactyl.conf /etc/nginx/sites-enabled/pterodactyl.conf

# You need to restart nginx regardless of OS.
sudo systemctl restart nginx