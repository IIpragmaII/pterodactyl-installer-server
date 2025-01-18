FROM scratch

USER 100

COPY configs .
COPY scripts .
COPY pterodactyl-installer-server .

ENTRYPOINT ["./pterodactyl-installer-server"]