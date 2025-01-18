FROM scratch

USER 100

RUN mkdir /installer
WORKDIR /installer

COPY configs .
COPY scripts .
COPY pterodactyl-installer-server .

ENTRYPOINT ["/bin/pterodactyl-installer-server"]