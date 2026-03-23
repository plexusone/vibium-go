# CLI Reference

The `vibium` CLI provides command-line browser automation.

## Installation

```bash
go install github.com/grokify/webpilot/cmd/vibium@latest
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--session` | Session file path (default: `~/.vibium/session.json`) |
| `-v, --verbose` | Verbose output |

## Commands

### launch

Launch a browser instance.

```bash
webpilot launch [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--headless` | Run in headless mode |

**Example:**

```bash
webpilot launch --headless
```

### go

Navigate to a URL.

```bash
webpilot go <url> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--timeout` | Navigation timeout (default: 30s) |

**Example:**

```bash
webpilot go https://example.com
```

### click

Click an element.

```bash
webpilot click <selector> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--timeout` | Timeout (default: 10s) |

**Example:**

```bash
webpilot click "#submit"
webpilot click "button.login"
```

### type

Type text into an element (appends).

```bash
vibium type <selector> <text> [flags]
```

**Example:**

```bash
vibium type "#search" "hello world"
```

### fill

Fill an input (replaces existing content).

```bash
webpilot fill <selector> <text> [flags]
```

**Example:**

```bash
webpilot fill "#email" "user@example.com"
```

### screenshot

Capture a screenshot.

```bash
webpilot screenshot <filename> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--selector` | Capture specific element |
| `--timeout` | Timeout (default: 30s) |

**Example:**

```bash
webpilot screenshot page.png
webpilot screenshot button.png --selector "#submit"
```

### eval

Execute JavaScript.

```bash
vibium eval <javascript> [flags]
```

**Example:**

```bash
vibium eval "document.title"
vibium eval "document.querySelectorAll('a').length"
```

### quit

Close the browser.

```bash
webpilot quit
```

### mcp

Start MCP server.

```bash
webpilot mcp [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--headless` | Run headless |
| `--timeout` | Default timeout |
| `--project` | Project name for reports |

### run

Run a YAML/JSON script.

```bash
webpilot run <script> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--headless` | Run headless |
| `--timeout` | Total script timeout |

**Example:**

```bash
webpilot run test.yaml
webpilot run login.json --headless
```

## Session Management

The CLI maintains session state in `~/.vibium/session.json`. This allows running commands across multiple invocations:

```bash
webpilot launch
webpilot go https://example.com
# ... later ...
webpilot screenshot result.png
webpilot quit
```

## Examples

### Login Flow

```bash
webpilot launch --headless
webpilot go https://example.com/login
webpilot fill "#email" "user@example.com"
webpilot fill "#password" "secret123"
webpilot click "#submit"
webpilot screenshot dashboard.png
webpilot quit
```

### Form Automation

```bash
webpilot launch
webpilot go https://example.com/form
webpilot fill "#name" "John Doe"
webpilot fill "#email" "john@example.com"
webpilot click "input[type='checkbox']"
webpilot click "#submit"
webpilot quit
```
