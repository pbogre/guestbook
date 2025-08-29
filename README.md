# Guestbook

A simple self-hostable guestbook for website interactivity (very early version)

## Table of Contents

- [Features](#features)
- [Installation](#installation)
    - [Using docker compose](#using-docker-compose)
    - [Manually (dev)](#manually-dev)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
    - [All requests coming from same IP](#all-requests-coming-from-same-ip)
- [Acknowledgements](#acknowledgements)

## Features
- 📝 Promote interactivity by letting guests leave guestbook-like messages
- 📜 Public history of messages left by guests
- 1️⃣  Each guest is allowed to leave one message only (via IP)
- ⚙️  Fully configurable (WIP)
- 🙅 Global rate limiting feature
- 🪶 Extremely lightweight, responsive, and snappy interface

## Installation

### Using docker compose

Sample `docker-compose.yml`:
```yml
services:
  guestbook:
    image: pbogre/guestbook:latest
    volumes:
      - /path/to/data:/data
      - /path/to/config/guestbook.yml:/config/guestbook.yml:ro
    restart: unless-stopped
    environment:
      PORT: 8080 # 8080 by default
    ports:
      - 8080:8080
```

### Manually (dev)
1. `git clone https://github.com/pbogre/guestbook`
2. `go mod tidy`
3. `go run .`

## Configuration

Default `guestbook.yml` config:
```yml
title: "Guestbook"
globalRateLimit: 5 # 1 request every x seconds (int only)
globalBurstLimit: 2
entriesPerPage: 10
```

You can mount your own config to the container, or if you
are happy with the default values, you can simply not mount
anything in your `docker-compose.yml` and Guestbook will
use the default config.

## Troubleshooting

### All requests coming from same IP

This means that the unique remote address computation is not
working properly. This is most likely due to how your setup
works. For example, if you are using nginx as reverse proxy,
you should make sure to add these lines to your configuration:

```nginx
proxy_set_header Host $host;
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
```

Also, if you are running nginx in docker, you should make sure
that it has the real IP information in the first place. You 
can do so by adding this line to your nginx's `docker-compose.yml`:

```yml
network_mode: "host"
```

You would have to do similar changes for other reverse proxies such 
as Caddy or Traefik.

## Acknowledgements
- [water.css](https://watercss.kognise.dev/)
