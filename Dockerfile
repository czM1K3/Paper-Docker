FROM golang:1.17-stretch AS build
WORKDIR /app
COPY launcher/ .
RUN go build -o ./launcher ./paper-launcher.go

FROM openjdk:17-slim-bullseye as mc
COPY src/ /default/
WORKDIR /app
COPY --from=build /app/launcher /app/launcher

ENTRYPOINT ["./launcher"]
