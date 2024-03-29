FROM golang:1.18-bullseye AS build
WORKDIR /app
COPY . .
RUN go build -o ./launcher .

FROM alpine AS healthcheck
RUN apk add wget
RUN wget -O minecraft-healthcheck.tar.gz "https://github.com/czM1K3/minecraft-healthcheck/releases/download/1.0.0/minecraft-healthcheck-1.0.0-linux-$([[ "$(uname -m)" = 'aarch64' ]] && echo 'arm64' || echo 'amd64').tar.gz"
RUN tar -xzf minecraft-healthcheck.tar.gz

FROM openjdk:17-slim-bullseye as mc
WORKDIR /app
COPY --from=build /app/launcher /app/launcher
COPY --from=healthcheck /minecraft-healthcheck /app/healthcheck

HEALTHCHECK --start-period=15s CMD /app/healthcheck
VOLUME ["/data"]
EXPOSE 25565/tcp
ENTRYPOINT ["./launcher"]
