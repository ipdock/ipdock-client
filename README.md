# ipdock-client

> The official dynamic DNS (DDNS) client for [ipdock.io](https://ipdock.io) — Dynamic DNS for self-hosters.

[![Docker Pulls](https://img.shields.io/docker/pulls/ipdockrepo/ipdock-client)](https://hub.docker.com/r/ipdockrepo/ipdock-client)
[![GitHub Container Registry](https://img.shields.io/badge/ghcr.io-ipdock%2Fipdock--client-blue)](https://ghcr.io/ipdock/ipdock-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## What is ipdock?

**ipdock.io** is a modern Dynamic DNS (DDNS) service built for self-hosters. It keeps your domain name pointing at your home or server IP — even when your ISP changes it.

- **Free to use** — get a hostname on `porthole.sh` or `ipdock.me` instantly
- **Bring your own domain** — point your custom domain at ipdock for DDNS updates
- **Self-hosting friendly** — lightweight Docker client, simple API, no bloat
- **Horizontally scalable** — built for reliability with globally distributed DNS nodes

**What this client does:**  
Runs silently in the background, detects your current public IP, and updates your ipdock hostname automatically whenever it changes.

---

## Quick Start

### 1. Sign up and create a hostname

Go to [ipdock.io](https://ipdock.io), create an account, verify your email, and add a hostname (e.g. `myhome.porthole.sh`). You'll get an API token starting with `pth_`.

### 2. Run the client

**GitHub Container Registry (recommended):**
```bash
docker run -d \
  --name ipdock \
  --restart unless-stopped \
  -e IPDOCK_TOKEN=pth_your_token_here \
  ghcr.io/ipdock/ipdock-client:latest
```

**Docker Hub:**
```bash
docker run -d \
  --name ipdock \
  --restart unless-stopped \
  -e IPDOCK_TOKEN=pth_your_token_here \
  ipdockrepo/ipdock-client:latest
```

That's it. Your hostname will update automatically whenever your IP changes.

---

## Docker Compose

```yaml
services:
  ipdock:
    image: ghcr.io/ipdock/ipdock-client:latest
    container_name: ipdock
    restart: unless-stopped
    environment:
      IPDOCK_TOKEN: pth_your_token_here
      IPDOCK_INTERVAL: "60"   # optional, default 60 seconds
```

---

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `IPDOCK_TOKEN` | ✅ | — | Your API token from the ipdock.io dashboard |
| `IPDOCK_API_URL` | ❌ | `https://api.ipdock.io` | API base URL (change for self-hosted deployments) |
| `IPDOCK_INTERVAL` | ❌ | `60` | IP check interval in seconds (minimum: 10) |

---

## How It Works

1. On startup, detects your public IP by querying one of several IP echo services (`ipify`, `icanhazip`, AWS checkip)
2. Calls `GET /api/v1/update?ip=<your-ip>` with your token
3. ipdock updates your DNS record if the IP has changed (`good`) or skips if unchanged (`nochg`)
4. Repeats every `IPDOCK_INTERVAL` seconds — only sending updates when the IP actually changes

The client is stateless and minimal — a single ~5MB Go binary in a scratch container.

---

## Supported Platforms

| Platform | Supported |
|---|---|
| `linux/amd64` | ✅ |
| `linux/arm64` | ✅ (Raspberry Pi 4+, Apple Silicon VMs) |

---

## API Compatibility

ipdock uses a DynDNS-compatible update API. If you're already running a DynDNS client, you can often point it at `api.ipdock.io` with your token instead of using this client.

**Update endpoint:**
```
GET https://api.ipdock.io/api/v1/update?ip=<your-ip>
Authorization: Bearer pth_your_token_here
```

**Responses:**
- `{"status":"good","ip":"..."}` — IP updated successfully
- `{"status":"nochg","ip":"..."}` — IP unchanged, no update needed

---

## Building From Source

```bash
git clone https://github.com/ipdock/ipdock-client.git
cd ipdock-client
go build -o ipdock-client .

# Or with Docker:
docker build -t ipdock-client .
```

---

## Self-Hosted / Development

If you're running your own ipdock instance:

```bash
docker run -d \
  --name ipdock \
  --restart unless-stopped \
  -e IPDOCK_TOKEN=pth_your_token_here \
  -e IPDOCK_API_URL=https://your-ipdock-instance.example.com \
  ghcr.io/ipdock/ipdock-client:latest
```

---

## Links

- 🌐 [ipdock.io](https://ipdock.io) — sign up and manage your hostnames
- 📦 [GitHub Container Registry](https://ghcr.io/ipdock/ipdock-client)
- 🐳 [Docker Hub](https://hub.docker.com/r/ipdockrepo/ipdock-client)
- 🐛 [Issues](https://github.com/ipdock/ipdock-client/issues)

---

## License

MIT — see [LICENSE](LICENSE) for details.
