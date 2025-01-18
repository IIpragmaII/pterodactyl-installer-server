FROM ubuntu:latest

USER 100

COPY pterodactyl-installer-server /bin/pterodactyl-installer-server

ENTRYPOINT ["/bin/pterodactyl-installer-server"]