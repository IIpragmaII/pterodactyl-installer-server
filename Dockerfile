FROM scratch

ENV GIN_MODE=release

EXPOSE 8080

USER 100

COPY configs configs
COPY scripts scripts
COPY pterodactyl-installer-server .

ENTRYPOINT ["./pterodactyl-installer-server"]