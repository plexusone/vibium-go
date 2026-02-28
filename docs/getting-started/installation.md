# Installation

## Go Client SDK

Install the Go module:

```bash
go get github.com/grokify/vibium-go
```

## CLI Tool

Build and install the CLI:

```bash
go install github.com/grokify/vibium-go/cmd/vibium@latest
```

Or build from source:

```bash
git clone https://github.com/grokify/vibium-go
cd vibium-go
go build -o vibium ./cmd/vibium
```

## Prerequisites

### Vibium Clicker Binary

The Go client requires the Vibium clicker binary. Install via npm:

```bash
npm install -g vibium
```

Or download from [Vibium releases](https://github.com/VibiumDev/vibium/releases).

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VIBIUM_CLICKER_PATH` | Path to clicker binary | Auto-detected |
| `VIBIUM_DEBUG` | Enable debug logging | `false` |
| `VIBIUM_HEADLESS` | Run headless by default | `false` |

## Verify Installation

```bash
# Check CLI
vibium --help

# Check clicker
vibium mcp --help
```
