# Guestbook

A simple self-hostable guestbook for website interactivity (very early version)

## Table of Contents

- [Features](#features)
- [Installation](#installation)
    - [Using docker compose](#using-docker-compose)
    - [Manually (dev)](#manually-dev)
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
    restart: unless-stopped
    environment:
      GB_TITLE: My Title        # default: Guestbook
      GB_FOOTER: <b>hey</b>     # default: ""
      GB_ENTRIES_PER_PAGE: 20   # default: 10
      GB_PAPER_CSS: true        # default: false
      GB_USE_RATELIMIT: true    # default: true
      GB_RATELIMIT: 0.16        # default: 0.2 (one per 5 seconds)
      GB_BURSTLIMIT: 2          # default: 2 (max burst of requests)
      PORT: 8080                # default: 8080
    ports:
      - 8080:8080
```

### Manually (dev)
1. `git clone https://github.com/pbogre/guestbook`
2. `go mod tidy`
3. Manually update the database file path in `main.go`
4. `go run .`

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
