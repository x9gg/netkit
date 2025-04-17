# NetKit

**ALPHA VERSION - NOT FOR PRODUCTION USE**

A simple tool for diagnosing network infrastructure.

## What it does

- Shows your IP addresses (IPv4/IPv6)
- Displays HTTP headers
- Works in browser or command line

## Usage

```bash
# Run locally
go run main.go

# with make
make dev

# Access in browser
http://localhost:8080

# with curl
curl http://localhost:8080
```

## Endpoints

- `/` - Web view
- `/plain` - Text output
- `/ipv4` - IPv4 only
- `/ipv6` - IPv6 only
- `/debug` - All headers

## Planed

- Ping
- Traceroute
- DNS tools
- FTP/SSH tests

## License

MIT