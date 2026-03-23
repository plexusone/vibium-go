# Installation

## Go Client SDK

Install the Go module:

```bash
go get github.com/grokify/webpilot
```

## CLI Tool

Build and install the CLI:

```bash
go install github.com/grokify/webpilot/cmd/vibium@latest
```

Or build from source:

```bash
git clone https://github.com/grokify/webpilot
cd webpilot
go build -o vibium ./cmd/vibium
```

## Prerequisites

### WebPilot Clicker Binary

The Go client requires the WebPilot clicker binary. Install via npm:

```bash
npm install -g vibium
```

Or download from [WebPilot releases](https://github.com/WebPilotDev/vibium/releases).

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `WEBPILOT_CLICKER_PATH` | Path to clicker binary | Auto-detected |
| `WEBPILOT_DEBUG` | Enable debug logging | `false` |
| `WEBPILOT_HEADLESS` | Run headless by default | `false` |

## Verify Installation

```bash
# Check CLI
vibium --help

# Check clicker
webpilot mcp --help
```
