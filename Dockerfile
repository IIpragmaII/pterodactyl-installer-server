FROM scratch

RUN mkdir /installer
WORKDIR /installer

USER 100

COPY configs .
COPY scripts .
COPY pterodactyl-installer-server .

ENTRYPOINT ["/bin/pterodactyl-installer-server"]