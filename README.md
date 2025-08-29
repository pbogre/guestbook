# Guestbook

A simple self-hostable guestbook for website interactivity (very early version)

## Table of Contents

- [Features](#features)
- [To do](#to-do)
- [Installation](#installation)
    - [Using docker compose](#using-docker-compose)
    - [Manually (dev)](#manually-dev)
- [Configuration](#configuration)
- [Acknowledgements](#acknowledgements)

## Features
- 📝 Promote interactivity by letting guests leave guestbook-like messages
- 📜 Public history of messages left by guests
- 1️⃣  Each guest is allowed to leave one message only (via IP)
- ⚙️  Fully configurable (WIP)
- 🙅 Global rate limiting feature
- 🪶 Extremely lightweight, responsive, and snappy interface

## To do
- [x] ~~User configuration (rate limits, unique messages, page size, ...)~~
- [x] ~~Docker image~~
- [ ] CI/CD github actions
- [ ] Skeuomorphic UI?

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

## Acknowledgements
- [water.css](https://watercss.kognise.dev/)
