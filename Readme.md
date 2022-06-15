# Paper Docker

## About
This container easily run Paper Minecraft Server. On startup it automatically downloads the latest Paper version for specific Minecraft version. Imprortant files are located at /data. I recommend you to make it persistent. By default, it runs on port 25565, so you should also port forward this port. If you are going to be using any kind of Map or web interface, you have to also forward this port. If you do so, I highly recommend to use reverse proxy (like Nginx) and Certbot to use secured connection. By default, it uses 2GB of RAM, you can (probably should) tweak it to your needs with environment variable MEMORY. Value has to be integer in MB for example 3072 for 3GB of ram.

When using, you agree with Mojang's [EULA](https://account.mojang.com/documents/minecraft_eula)

## Instalation
Running in foreground
```bash
docker run -v /your/data/path:/data -p 25565:25565 -e MEMORY=2048 --name my-paper-server --rm -ti czm1k3/paper-docker
```

Running in background
```bash
docker run -v /your/data/path:/data -p 25565:25565 -e MEMORY=2048 --restart unless-stopped --name my-paper-server -tid czm1k3/paper-docker
```

## Getting into prompt
If you want to access server's prompt (to give op, ban, etc), you firstly need to find id of running container (in first column):
```bash
docker ps
```
Then you can attack to it with command:
```bash
docker attack <container-id>
```
When you want to leave, just press: CTRL+P than CTRL+Q

# Build
```bash
docker buildx build --platform=linux/amd64,linux/arm64 --push --tag czm1k3/paper-docker .
```
