FROM scratch

ENV GIN_MODE=release

EXPOSE 8080

USER 100

COPY configs .
COPY scripts .
COPY pterodactyl-installer-server .

ENTRYPOINT ["./pterodactyl-installer-server"]