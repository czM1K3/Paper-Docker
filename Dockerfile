FROM golang:1.17-bullseye AS build
WORKDIR /app
COPY launcher/ .
RUN go build -o ./launcher ./paper-launcher.go

FROM alpine AS healthcheck
RUN wget -O minecraft-healthcheck.tar.gz "https://github.com/czM1K3/minecraft-healthcheck/releases/download/1.0.0/minecraft-healthcheck-1.0.0-linux-$([[ "$(uname -m)" = 'aarch64' ]] && echo 'arm64' || echo 'amd64').tar.gz"
RUN tar -xzf minecraft-healthcheck.tar.gz

FROM openjdk:17-slim-bullseye as mc
COPY src/ /default/
WORKDIR /app
COPY --from=build /app/launcher /app/launcher
COPY --from=healthcheck /minecraft-healthcheck /app/healthcheck

HEALTHCHECK --start-period=15s CMD /app/healthcheck
ENTRYPOINT ["./launcher"]
