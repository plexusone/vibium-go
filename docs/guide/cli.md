# CLI Reference

The `vibium` CLI provides command-line browser automation.

## Installation

```bash
go install github.com/grokify/vibium-go/cmd/vibium@latest
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
vibium launch [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--headless` | Run in headless mode |

**Example:**

```bash
vibium launch --headless
```

### go

Navigate to a URL.

```bash
vibium go <url> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--timeout` | Navigation timeout (default: 30s) |

**Example:**

```bash
vibium go https://example.com
```

### click

Click an element.

```bash
vibium click <selector> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--timeout` | Timeout (default: 10s) |

**Example:**

```bash
vibium click "#submit"
vibium click "button.login"
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
vibium fill <selector> <text> [flags]
```

**Example:**

```bash
vibium fill "#email" "user@example.com"
```

### screenshot

Capture a screenshot.

```bash
vibium screenshot <filename> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--selector` | Capture specific element |
| `--timeout` | Timeout (default: 30s) |

**Example:**

```bash
vibium screenshot page.png
vibium screenshot button.png --selector "#submit"
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
vibium quit
```

### mcp

Start MCP server.

```bash
vibium mcp [flags]
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
vibium run <script> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--headless` | Run headless |
| `--timeout` | Total script timeout |

**Example:**

```bash
vibium run test.yaml
vibium run login.json --headless
```

## Session Management

The CLI maintains session state in `~/.vibium/session.json`. This allows running commands across multiple invocations:

```bash
vibium launch
vibium go https://example.com
# ... later ...
vibium screenshot result.png
vibium quit
```

## Examples

### Login Flow

```bash
vibium launch --headless
vibium go https://example.com/login
vibium fill "#email" "user@example.com"
vibium fill "#password" "secret123"
vibium click "#submit"
vibium screenshot dashboard.png
vibium quit
```

### Form Automation

```bash
vibium launch
vibium go https://example.com/form
vibium fill "#name" "John Doe"
vibium fill "#email" "john@example.com"
vibium click "input[type='checkbox']"
vibium click "#submit"
vibium quit
```
