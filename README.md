# Vibium Go

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

Go client and tooling for the [Vibium](https://github.com/VibiumDev/vibium) browser automation platform.

Vibium uses WebDriver BiDi for real-time bidirectional communication with browsers, making it ideal for AI-assisted automation.

## Overview

This project provides:

| Component | Description | Origin |
|-----------|-------------|--------|
| **Go Client SDK** | Programmatic browser control | Feature parity with JS/Python |
| **MCP Server** | 75+ tools for AI assistants | Go-specific |
| **CLI** | Command-line browser automation | Go-specific |
| **Script Runner** | Deterministic test execution | Go-specific |
| **Session Recording** | Capture actions as replayable scripts | Go-specific |

## Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                         vibium-go                              │
├─────────────┬─────────────┬─────────────┬──────────────────────┤
│  Go Client  │ MCP Server  │    CLI      │   Script Runner      │
│    SDK      │  (75 tools) │  (vibium)   │   (vibium run)       │
├─────────────┴─────────────┴─────────────┴──────────────────────┤
│                    WebDriver BiDi Protocol                     │
├────────────────────────────────────────────────────────────────┤
│                   Vibium Clicker (upstream)                    │
├────────────────────────────────────────────────────────────────┤
│                    Chrome / Chromium                           │
└────────────────────────────────────────────────────────────────┘
```

## Installation

```bash
go get github.com/plexusone/vibium-go
```

### Prerequisites

Install the Vibium clicker binary:

```bash
npm install -g vibium
```

Or set `VIBIUM_CLICKER_PATH` to point to the binary.

## Quick Start

### Go Client SDK

```go
package main

import (
    "context"
    "log"

    vibium "github.com/plexusone/vibium-go"
)

func main() {
    ctx := context.Background()

    // Launch browser
    vibe, err := vibium.Launch(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer vibe.Quit(ctx)

    // Navigate and interact
    vibe.Go(ctx, "https://example.com")

    link, _ := vibe.Find(ctx, "a", nil)
    link.Click(ctx, nil)
}
```

### MCP Server

Start the MCP server for AI assistant integration:

```bash
vibium mcp --headless
```

Configure in Claude Desktop (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "vibium": {
      "command": "vibium",
      "args": ["mcp", "--headless"]
    }
  }
}
```

### CLI Commands

```bash
# Launch browser and run commands
vibium launch --headless
vibium go https://example.com
vibium fill "#email" "user@example.com"
vibium click "#submit"
vibium screenshot result.png
vibium quit
```

### Script Runner

Execute deterministic test scripts:

```bash
vibium run test.json
```

Script format (JSON or YAML):

```json
{
  "name": "Login Test",
  "steps": [
    {"action": "navigate", "url": "https://example.com/login"},
    {"action": "fill", "selector": "#email", "value": "user@example.com"},
    {"action": "fill", "selector": "#password", "value": "secret"},
    {"action": "click", "selector": "#submit"},
    {"action": "assertUrl", "expected": "https://example.com/dashboard"}
  ]
}
```

## Feature Comparison

### Client SDK (Parity with JS/Python)

| Feature | JS | Python | Go |
|---------|:--:|:------:|:--:|
| Browser launch/quit | ✅ | ✅ | ✅ |
| Navigation (go, back, forward, reload) | ✅ | ✅ | ✅ |
| Element finding (CSS selectors) | ✅ | ✅ | ✅ |
| Click, type, fill | ✅ | ✅ | ✅ |
| Screenshots | ✅ | ✅ | ✅ |
| JavaScript evaluation | ✅ | ✅ | ✅ |
| Keyboard/mouse controllers | ✅ | ✅ | ✅ |
| Browser context management | ✅ | ✅ | ✅ |
| Network interception | ✅ | ✅ | ✅ |
| Tracing | ✅ | ✅ | ✅ |
| Clock control | ✅ | ✅ | ✅ |

### Go-Specific Features

| Feature | Description |
|---------|-------------|
| **MCP Server** | 75+ tools for AI-assisted automation |
| **CLI** | `vibium` command with subcommands |
| **Script Runner** | Execute JSON/YAML test scripts |
| **Session Recording** | Capture MCP actions as replayable scripts |
| **JSON Schema** | Validated script format |
| **Test Reporting** | Structured test results with diagnostics |

## MCP Server Tools

The MCP server provides 75+ tools organized by category:

| Category | Tools |
|----------|-------|
| Browser | `browser_launch`, `browser_quit` |
| Navigation | `navigate`, `back`, `forward`, `reload` |
| Interactions | `click`, `dblclick`, `type`, `fill`, `clear`, `press` |
| Forms | `check`, `uncheck`, `select_option`, `set_files` |
| Element State | `get_text`, `get_value`, `is_visible`, `is_enabled` |
| Page State | `get_title`, `get_url`, `get_content`, `screenshot` |
| Waiting | `wait_until`, `wait_for_url`, `wait_for_load` |
| Input | `keyboard_*`, `mouse_*`, `touch_*` |
| Recording | `start_recording`, `stop_recording`, `export_script` |
| Assertions | `assert_text`, `assert_element`, `assert_url` |

## Session Recording Workflow

Convert natural language test plans into deterministic scripts:

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│  Markdown Test   │     │   LLM + MCP      │     │   JSON Script    │
│  Plan (English)  │ ──▶ │   (exploration)  │ ──▶ │ (deterministic)  │
└──────────────────┘     └──────────────────┘     └──────────────────┘
```

1. Write test plan in Markdown
2. LLM executes via MCP with `start_recording`
3. LLM explores, finds selectors, handles edge cases
4. Export with `export_script` to get JSON
5. Run deterministically with `vibium run`

## API Reference

See [pkg.go.dev](https://pkg.go.dev/github.com/plexusone/vibium-go) for full API documentation.

### Key Types

```go
// Launch browser
vibe, err := vibium.Launch(ctx)
vibe, err := vibium.LaunchHeadless(ctx)

// Navigation
vibe.Go(ctx, url)
vibe.Back(ctx)
vibe.Forward(ctx)
vibe.Reload(ctx)

// Finding elements
elem, err := vibe.Find(ctx, selector, &vibium.FindOptions{})
elems, err := vibe.FindAll(ctx, selector)

// Element interactions
elem.Click(ctx, nil)
elem.Fill(ctx, value, nil)
elem.Type(ctx, text, nil)

// Input controllers
vibe.Keyboard().Press(ctx, "Enter")
vibe.Mouse().Click(ctx, x, y)

// Capture
data, err := vibe.Screenshot(ctx)
```

## Testing

```bash
# Unit tests
go test -v ./...

# Integration tests (requires clicker)
go test -tags=integration -v ./integration/...

# Headless mode
VIBIUM_HEADLESS=1 go test -tags=integration -v ./integration/...
```

## Debug Logging

```bash
VIBIUM_DEBUG=1 vibium mcp
```

## Related Projects

- [Vibium](https://github.com/VibiumDev/vibium) - Upstream platform
- [vibium-wcag](https://github.com/agentplexus/vibium-wcag) - WCAG 2.2 accessibility testing
- [omnillm](https://github.com/agentplexus/omnillm) - Unified LLM client
- [WebDriver BiDi](https://w3c.github.io/webdriver-bidi/) - Protocol specification

## License

MIT

 [go-ci-svg]: https://github.com/plexusone/vibium-go/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/vibium-go/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/vibium-go/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/vibium-go/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/vibium-go/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/vibium-go/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/vibium-go
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/vibium-go
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/vibium-go
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/vibium-go
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fcoreforge
 [loc-svg]: https://tokei.rs/b1/github/plexusone/vibium-go
 [repo-url]: https://github.com/plexusone/vibium-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/vibium-go/blob/master/LICENSE
