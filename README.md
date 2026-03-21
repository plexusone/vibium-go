# Vibium Go

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

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
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fvibium-go
 [loc-svg]: https://tokei.rs/b1/github/plexusone/vibium-go
 [repo-url]: https://github.com/plexusone/vibium-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/vibium-go/blob/master/LICENSE

Go client and tooling for the [Vibium](https://github.com/VibiumDev/vibium) browser automation platform.

Vibium uses WebDriver BiDi for real-time bidirectional communication with browsers, making it ideal for AI-assisted automation.

## Overview

This project provides:

| Component | Description | Origin |
|-----------|-------------|--------|
| **Go Client SDK** | Programmatic browser control | Feature parity with JS/Python |
| **MCP Server** | 85+ tools for AI assistants | Go-specific |
| **CLI** | Command-line browser automation | Go-specific |
| **Script Runner** | Deterministic test execution | Go-specific |
| **Session Recording** | Capture actions as replayable scripts | Go-specific |

## Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                         vibium-go                              │
├─────────────┬─────────────┬─────────────┬──────────────────────┤
│  Go Client  │ MCP Server  │    CLI      │   Script Runner      │
│    SDK      │  (85 tools) │  (vibium)   │   (vibium run)       │
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

The MCP server provides 85+ tools organized by category:

| Category | Tools |
|----------|-------|
| Browser | `browser_launch`, `browser_quit` |
| Navigation | `navigate`, `back`, `forward`, `reload` |
| Interactions | `click`, `dblclick`, `type`, `fill`, `clear`, `press` |
| Forms | `check`, `uncheck`, `select_option`, `set_files`, `fill_form` |
| Element State | `get_text`, `get_value`, `is_visible`, `is_enabled` |
| Page State | `get_title`, `get_url`, `get_content`, `screenshot` |
| Waiting | `wait_until`, `wait_for_url`, `wait_for_load` |
| Storage | `get_storage_state`, `set_storage_state`, `clear_storage` |
| LocalStorage | `localstorage_get`, `localstorage_set`, `localstorage_list`, `localstorage_delete`, `localstorage_clear` |
| SessionStorage | `sessionstorage_get`, `sessionstorage_set`, `sessionstorage_list`, `sessionstorage_delete`, `sessionstorage_clear` |
| Network | `get_network_requests`, `clear_network_requests`, `route`, `unroute`, `route_list` |
| Tabs | `list_tabs`, `select_tab`, `close_tab` |
| Dialogs | `handle_dialog`, `get_dialog` |
| HITL | `pause_for_human` |
| Input | `keyboard_*`, `mouse_*`, `touch_*` |
| Tracing | `start_trace`, `stop_trace`, `start_trace_chunk`, `stop_trace_chunk` |
| Init Scripts | `add_init_script` |
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

// Finding elements by CSS selector
elem, err := vibe.Find(ctx, selector, nil)
elems, err := vibe.FindAll(ctx, selector, nil)

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

## Semantic Selectors

Find elements by accessibility attributes instead of brittle CSS selectors. This is especially useful for AI-assisted automation where element structure may change but semantics remain stable.

### SDK Usage

```go
// Find by ARIA role and text content
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    Role: "button",
    Text: "Submit",
})

// Find by label (for form inputs)
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    Label: "Email address",
})

// Find by placeholder
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    Placeholder: "Enter your email",
})

// Find by data-testid (recommended for testing)
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    TestID: "login-button",
})

// Combine CSS selector with semantic filtering
elem, err := vibe.Find(ctx, "form", &vibium.FindOptions{
    Role: "textbox",
    Label: "Password",
})

// Find all buttons
buttons, err := vibe.FindAll(ctx, "", &vibium.FindOptions{Role: "button"})

// Find element near another element
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    Role: "button",
    Near: "#username-input",
})
```

### MCP Tool Usage

Semantic selectors work with `click`, `type`, `fill`, and `press` tools:

```json
// Click a button by role and text
{"name": "click", "arguments": {"role": "button", "text": "Sign In"}}

// Fill input by label
{"name": "fill", "arguments": {"label": "Email", "value": "user@example.com"}}

// Type in input by placeholder
{"name": "type", "arguments": {"placeholder": "Search...", "text": "query"}}

// Click by data-testid
{"name": "click", "arguments": {"testid": "submit-btn"}}
```

### Available Selectors

| Selector | Description | Example |
|----------|-------------|---------|
| `role` | ARIA role | `button`, `textbox`, `link`, `checkbox` |
| `text` | Visible text content | `"Submit"`, `"Learn more"` |
| `label` | Associated label text | `"Email address"`, `"Password"` |
| `placeholder` | Input placeholder | `"Enter email"` |
| `testid` | `data-testid` attribute | `"login-btn"` |
| `alt` | Image alt text | `"Company logo"` |
| `title` | Element title attribute | `"Close dialog"` |
| `xpath` | XPath expression | `"//button[@type='submit']"` |
| `near` | CSS selector of nearby element | `"#username"` |

## Init Scripts

Inject JavaScript that runs before any page scripts on every navigation. Useful for mocking APIs, injecting test helpers, or setting up authentication.

### SDK Usage

```go
// Add init script to inject before page scripts
err := vibe.AddInitScript(ctx, `window.testMode = true;`)

// Mock an API
err := vibe.AddInitScript(ctx, `
    window.fetch = async (url, opts) => {
        if (url.includes('/api/user')) {
            return { json: () => ({ id: 1, name: 'Test User' }) };
        }
        return originalFetch(url, opts);
    };
`)
```

### CLI Usage

```bash
# Inject scripts when launching
vibium mcp --init-script=./mock-api.js --init-script=./test-helpers.js

# Or with the standalone binary
vibium-mcp -init-script=./mock-api.js
```

### MCP Tool Usage

```json
{"name": "add_init_script", "arguments": {"script": "window.testMode = true;"}}
```

## Storage State

Save and restore complete browser state including cookies, localStorage, and sessionStorage. Essential for maintaining login sessions across browser restarts.

### SDK Usage

```go
// Get complete storage state
state, err := vibe.StorageState(ctx)

// Save to file
jsonBytes, _ := json.Marshal(state)
os.WriteFile("auth-state.json", jsonBytes, 0600)

// Restore from file
var savedState vibium.StorageState
json.Unmarshal(jsonBytes, &savedState)
err := vibe.SetStorageState(ctx, &savedState)

// Clear all storage
err := vibe.ClearStorage(ctx)
```

### MCP Tool Usage

```json
// Save session
{"name": "get_storage_state"}

// Restore session
{"name": "set_storage_state", "arguments": {"state": "<json from get_storage_state>"}}

// Clear all storage
{"name": "clear_storage"}
```

## Tracing

Record browser actions with screenshots and DOM snapshots for debugging and test creation.

### SDK Usage

```go
// Start tracing
tracing := vibe.Tracing()
err := tracing.Start(ctx, &vibium.TracingStartOptions{
    Screenshots: true,
    Snapshots:   true,
    Title:       "Login Flow Test",
})

// Perform actions...
vibe.Go(ctx, "https://example.com")
elem, _ := vibe.Find(ctx, "button", nil)
elem.Click(ctx, nil)

// Stop and save trace
data, err := tracing.Stop(ctx, nil)
os.WriteFile("trace.zip", data, 0600)
```

### MCP Tool Usage

```json
// Start trace
{"name": "start_trace", "arguments": {"screenshots": true, "title": "My Test"}}

// Stop and get trace data
{"name": "stop_trace", "arguments": {"path": "/tmp/trace.zip"}}
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
