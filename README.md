# go-http-server

[![Go Version](https://img.shields.io/badge/go-1.26.2-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/Irayago/go-http-server/pulls)

A simple Go HTTP server for testing client/server request behavior, response delays, and TCP connection handling. Includes a PowerShell client for generating test requests.

## Demo

![Demo of go-http-server handling a delayed request](demo.gif)

*Run `go run main.go`, then hit the server with `curl` or the included PowerShell client and watch the configurable delay and request logging in action.*

## Features

- Lightweight HTTP server built on Go's standard `net/http` package
- Configurable response delay (`httpTimer`), adjustable at runtime via API — useful for testing timeouts, slow clients, or keep-alive behavior
- TCP keep-alive disabled by default for predictable connection testing
- Logs request metadata (method, user agent, IP, headers) to the console
- PowerShell client script for quickly firing test requests at the server

## Requirements

- Go 1.26.2 or later
- PowerShell (optional, only needed for the included client script)

## Getting Started

Clone the repo and run the server:

```bash
git clone https://github.com/Irayago/go-http-server.git
cd go-http-server
go run main.go
```

By default, the server listens on port `9999`.

## API Endpoints

### `GET /`

Waits for `httpTimer` seconds, then responds with a plain-text summary of the request (timestamp, user agent, remote IP, and elapsed time). Also logs request metadata to the server console.

**Example:**

```bash
curl http://localhost:9999/
```

### `GET /httptimer`

Returns the current `httpTimer` value (in seconds) as plain text.

**Example:**

```bash
curl http://localhost:9999/httptimer
```

### `PUT /httptimer`

Sets a new `httpTimer` value. Expects a JSON body with an `httpTimer` key holding a non-negative number.

**Example:**

```bash
curl -X PUT http://localhost:9999/httptimer \
  -H "Content-Type: application/json" \
  -d '{"httpTimer": 5}'
```

Returns `400 Bad Request` if the body is empty, isn't valid JSON, is missing the `httpTimer` key, the value isn't a number, or the value is negative.

## PowerShell Client

`http_client.ps1` is a small script for sending test requests to the server. Replace `REPLACE_WITH_IP` with your server's host/IP before running:

```powershell
Invoke-WebRequest -Uri "http://REPLACE_WITH_IP" | Select-Object -ExpandProperty Content
```

A commented-out loop is included in the script for sending multiple requests over the same TCP connection, if you want to test connection reuse.

## Coming Soon
- API endpoints for adjusting HTTP server timeouts

## Project Structure

```
go-http-server/
├── main.go           # HTTP server implementation
├── http_client.ps1    # PowerShell test client
├── go.mod
└── LICENSE
```

## License

MIT License — see [LICENSE](LICENSE) for details.
