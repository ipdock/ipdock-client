# ipdock-client

> Docker-native DDNS update client for [ipdock.io](https://ipdock.io)

Automatically keeps your ipdock.io hostname updated with your current public IP address.

## Quick Start

```bash
docker run -d \
  --name ipdock \
  --restart=always \
  -e IPDOCK_TOKEN=your_token_here \
  ipdockrepo/ipdock-client:latest
```

Get your token from the [ipdock.io dashboard](https://ipdock.io/dashboard).

## Docker Compose

```yaml
services:
  ipdock:
    image: ipdockrepo/ipdock-client:latest
    restart: always
    environment:
      - IPDOCK_TOKEN=your_token_here
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `IPDOCK_TOKEN` | ✅ | Your hostname token from the ipdock.io dashboard |
| `IPDOCK_API_URL` | ❌ | Override API URL (default: `https://api.ipdock.io`) |
| `IPDOCK_INTERVAL` | ❌ | Update interval in seconds (default: `300`) |

## How it Works

The client:
1. Detects your current public IP address
2. Sends an authenticated update to the ipdock.io API
3. Waits for the configured interval, then repeats
4. Only sends updates when your IP actually changes (no unnecessary API calls)

## Links

- 🌐 Dashboard: [ipdock.io/dashboard](https://ipdock.io/dashboard)
- 📖 Docs: [ipdock.io/docs](https://ipdock.io/docs)
- 🐛 Issues: [github.com/ipdock/ipdock-client/issues](https://github.com/ipdock/ipdock-client/issues)
- 📧 Support: [support@ipdock.io](mailto:support@ipdock.io)

## License

MIT
