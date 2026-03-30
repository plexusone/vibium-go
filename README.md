# W3Pilot

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/w3pilot/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/w3pilot/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/w3pilot/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/w3pilot/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/w3pilot/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/w3pilot/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/w3pilot
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/w3pilot
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/w3pilot
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/w3pilot
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fw3pilot
 [loc-svg]: https://tokei.rs/b1/github/plexusone/w3pilot
 [repo-url]: https://github.com/plexusone/w3pilot
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/w3pilot/blob/master/LICENSE

Go browser automation library using WebDriver BiDi for real-time bidirectional communication with browsers, ideal for AI-assisted automation.

## Overview

This project provides:

| Component | Description |
|-----------|-------------|
| **Go Client SDK** | Programmatic browser control |
| **MCP Server** | 161 tools across 21 namespaces for AI assistants |
| **CLI** | Command-line browser automation |
| **Script Runner** | Deterministic test execution |
| **Session Recording** | Capture actions as replayable scripts |

## Architecture

W3Pilot uses a **dual-protocol architecture** connecting to a single Chrome browser via both WebDriver BiDi and Chrome DevTools Protocol (CDP):

```
┌──────────────────────────────────────────────────────────────┐
│                          Applications                        │
├───────────────┬───────────────┬──────────────────────────────┤
│    w3pilot    │  w3pilot-mcp  │       Your Go App            │
│     (CLI)     │ (MCP Server)  │     import "w3pilot"         │
├───────────────┴───────────────┴──────────────────────────────┤
│                                                              │
│                        W3Pilot Go SDK                        │
│                 github.com/plexusone/w3pilot                 │
│                                                              │
│  ┌────────────────────────┐  ┌────────────────────────────┐  │
│  │      BiDi Client       │  │       CDP Client           │  │
│  │   (page automation)    │  │   (profiling/debugging)    │  │
│  │                        │  │                            │  │
│  │ • Navigation           │  │ • Heap snapshots           │  │
│  │ • Element interaction  │  │ • Network emulation        │  │
│  │ • Screenshots          │  │ • CPU throttling           │  │
│  │ • Tracing              │  │ • Code coverage            │  │
│  │ • Accessibility        │  │ • Console debugging        │  │
│  └───────────┬────────────┘  └─────────────┬──────────────┘  │
│              │                             │                 │
├──────────────┼─────────────────────────────┼─────────────────┤
│              ▼                             ▼                 │
│       WebDriver BiDi                Chrome DevTools          │
│       (stdio pipe)                  (CDP WebSocket)          │
├──────────────────────────────────────────────────────────────┤
│                       Chrome / Chromium                      │
└──────────────────────────────────────────────────────────────┘
```

### Why Dual-Protocol?

W3Pilot combines two complementary protocols for complete browser control:

| Protocol | Purpose | Strengths |
|----------|---------|-----------|
| **WebDriver BiDi** | Automation & Testing | Semantic selectors, real-time events, cross-browser potential, future-proof standard |
| **Chrome DevTools Protocol** | Inspection & Profiling | Heap profiling, network bodies, CPU/network emulation, coverage analysis |

**BiDi Client** excels at:

- Page automation (navigation, clicks, typing)
- Semantic element finding (by role, label, text, testid)
- Screenshots and accessibility trees
- Tracing and session recording
- Human-in-the-loop workflows (CAPTCHA, SSO)

**CDP Client** excels at:

- Memory profiling (heap snapshots)
- Network response body capture
- Performance emulation (Slow 3G, CPU throttling)
- Code coverage analysis
- Low-level debugging

Both protocols connect to the **same Chrome browser instance**, allowing you to automate with BiDi while profiling with CDP simultaneously.

### Protocol-Agnostic API

The SDK automatically handles protocol selection. Some methods try BiDi first and fall back to CDP when BiDi doesn't support the feature:

| Method | Tries First | Falls Back To |
|--------|-------------|---------------|
| `SetOffline()` | BiDi | CDP network emulation |
| `ConsoleMessages()` | BiDi | CDP console debugger |
| `ClearConsoleMessages()` | BiDi | CDP console debugger |

Users call the same method regardless of which protocol is used internally. When BiDi support is added upstream, the SDK will automatically use it without requiring code changes.

## Prerequisites

W3Pilot requires the **Clicker** binary, a WebDriver BiDi browser launcher from the [Vibium](https://github.com/VibiumDev/vibium) project.

### Install Clicker

**Option 1: Download from GitHub Releases**

Download the latest release for your platform from [vibium/releases](https://github.com/VibiumDev/vibium/releases) and add it to your PATH.

**Option 2: Build from Source**

```bash
git clone https://github.com/VibiumDev/vibium.git
cd vibium/clicker
go build -o clicker .
mv clicker /usr/local/bin/  # or add to PATH
```

**Option 3: Set Environment Variable**

If the binary is in a custom location:

```bash
export CLICKER_BIN_PATH=/path/to/clicker
```

### Verify Installation

```bash
clicker --version
```

### Browser Requirements

Clicker automatically manages Chrome/Chromium. If Chrome is not installed, download it from [google.com/chrome](https://www.google.com/chrome/).

## Installation

```bash
go get github.com/plexusone/w3pilot
```

## Quick Start

### Go Client SDK

```go
package main

import (
    "context"
    "log"

    "github.com/plexusone/w3pilot"
)

func main() {
    ctx := context.Background()

    // Launch browser
    pilot, err := w3pilot.Launch(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer pilot.Quit(ctx)

    // Navigate and interact
    pilot.Go(ctx, "https://example.com")

    link, _ := pilot.Find(ctx, "a", nil)
    link.Click(ctx, nil)
}
```

### Session Management

Manage persistent browser sessions that can be reused across CLI commands:

```go
import "github.com/plexusone/w3pilot/session"

// Create session manager
mgr := session.NewManager(session.Config{
    AutoReconnect: true,
})

// Get browser (launches if needed, reconnects if possible)
pilot, err := mgr.Pilot(ctx)

// Detach without closing browser
mgr.Detach()

// Later: reconnect to same browser
pilot, err = mgr.Pilot(ctx)

// When done: close browser
mgr.Close(ctx)
```

### MCP Server

Start the MCP server for AI assistant integration:

```bash
w3pilot mcp --headless
```

Configure in Claude Desktop (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "w3pilot": {
      "command": "w3pilot",
      "args": ["mcp", "--headless"]
    }
  }
}
```

### CLI Commands

```bash
# Browser lifecycle
w3pilot browser launch --headless
w3pilot browser quit

# Page navigation and capture
w3pilot page navigate https://example.com
w3pilot page back
w3pilot page screenshot result.png
w3pilot page title

# Element interactions
w3pilot element fill "#email" "user@example.com"
w3pilot element click "#submit"
w3pilot element text "#result"

# Wait for conditions
w3pilot wait selector "#modal"
w3pilot wait url "**/dashboard"

# JavaScript execution
w3pilot js eval "document.title"
```

### Script Runner

Execute deterministic test scripts:

```bash
w3pilot run test.json
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

### Client SDK

| Feature | Status |
|---------|:------:|
| Browser launch/quit | ✅ |
| Navigation (go, back, forward, reload) | ✅ |
| Element finding (CSS selectors) | ✅ |
| Click, type, fill | ✅ |
| Screenshots | ✅ |
| JavaScript evaluation | ✅ |
| Keyboard/mouse controllers | ✅ |
| Browser context management | ✅ |
| Network interception | ✅ |
| Tracing | ✅ |
| Clock control | ✅ |

### CDP Features (via Chrome DevTools Protocol)

| Feature | Status |
|---------|:------:|
| Heap snapshots | ✅ |
| Network emulation (Slow 3G, Fast 3G, 4G) | ✅ |
| CPU throttling | ✅ |
| Direct CDP command access | ✅ |

### Additional Features

| Feature | Description |
|---------|-------------|
| **MCP Server** | 161 tools across 21 namespaces for AI-assisted automation |
| **CLI** | `w3pilot` command with subcommands |
| **Script Runner** | Execute JSON/YAML test scripts |
| **Session Management** | Persistent browser sessions with reconnection support |
| **Session Recording** | Capture MCP actions as replayable scripts |
| **JSON Schema** | Validated script format |
| **Test Reporting** | Structured test results with diagnostics |

## MCP Server Tools

The MCP server provides **161 tools across 21 namespaces**. Export the full list as JSON with `w3pilot mcp --list-tools`.

**Namespaces:**

| Namespace | Tools | Examples |
|-----------|------:|----------|
| `accessibility_` | 1 | `accessibility_snapshot` |
| `batch_` | 1 | `batch_execute` |
| `browser_` | 2 | `browser_launch`, `browser_quit` |
| `cdp_` | 20 | `cdp_take_heap_snapshot`, `cdp_run_lighthouse`, `cdp_start_coverage` |
| `config_` | 1 | `config_get` |
| `console_` | 2 | `console_get_messages`, `console_clear` |
| `dialog_` | 2 | `dialog_handle`, `dialog_get` |
| `element_` | 33 | `element_click`, `element_fill`, `element_get_text`, `element_is_visible` |
| `frame_` | 2 | `frame_select`, `frame_select_main` |
| `http_` | 1 | `http_request` |
| `human_` | 1 | `human_pause` |
| `input_` | 12 | `input_keyboard_press`, `input_mouse_click`, `input_touch_tap` |
| `js_` | 4 | `js_evaluate`, `js_add_script`, `js_add_style`, `js_init_script` |
| `network_` | 6 | `network_get_requests`, `network_route`, `network_set_offline` |
| `page_` | 19 | `page_navigate`, `page_go_back`, `page_screenshot`, `page_emulate_media` |
| `record_` | 5 | `record_start`, `record_stop`, `record_export` |
| `storage_` | 17 | `storage_get_cookies`, `storage_local_get`, `storage_session_set` |
| `tab_` | 3 | `tab_list`, `tab_select`, `tab_close` |
| `test_` | 15 | `test_assert_text`, `test_verify_value`, `test_generate_locator` |
| `trace_` | 6 | `trace_start`, `trace_stop`, `trace_chunk_start` |
| `video_` | 2 | `video_start`, `video_stop` |
| `wait_` | 6 | `wait_for_state`, `wait_for_url`, `wait_for_load`, `wait_for_text` |

See [docs/reference/mcp-tools.md](docs/reference/mcp-tools.md) for the complete reference.

## Session Recording Workflow

Convert natural language test plans into deterministic scripts:

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│  Markdown Test   │     │   LLM + MCP      │     │   JSON Script    │
│  Plan (English)  │ ──▶ │   (exploration)  │ ──▶ │ (deterministic)  │
└──────────────────┘     └──────────────────┘     └──────────────────┘
```

1. Write test plan in Markdown
2. LLM executes via MCP with `record_start`
3. LLM explores, finds selectors, handles edge cases
4. Export with `record_export` to get JSON
5. Run deterministically with `w3pilot run`

## API Reference

See [pkg.go.dev](https://pkg.go.dev/github.com/plexusone/w3pilot) for full API documentation.

### Key Types

```go
// Launch browser
pilot, err := w3pilot.Launch(ctx)
pilot, err := w3pilot.LaunchHeadless(ctx)

// Navigation
pilot.Go(ctx, url)
pilot.Back(ctx)
pilot.Forward(ctx)
pilot.Reload(ctx)

// Finding elements by CSS selector
elem, err := pilot.Find(ctx, selector, nil)
elems, err := pilot.FindAll(ctx, selector, nil)

// Element interactions
elem.Click(ctx, nil)
elem.Fill(ctx, value, nil)
elem.Type(ctx, text, nil)

// Input controllers
pilot.Keyboard().Press(ctx, "Enter")
pilot.Mouse().Click(ctx, x, y)

// Capture
data, err := pilot.Screenshot(ctx)
```

## Semantic Selectors

Find elements by accessibility attributes instead of brittle CSS selectors. This is especially useful for AI-assisted automation where element structure may change but semantics remain stable.

### SDK Usage

```go
// Find by ARIA role and text content
elem, err := pilot.Find(ctx, "", &w3pilot.FindOptions{
    Role: "button",
    Text: "Submit",
})

// Find by label (for form inputs)
elem, err := pilot.Find(ctx, "", &w3pilot.FindOptions{
    Label: "Email address",
})

// Find by placeholder
elem, err := pilot.Find(ctx, "", &w3pilot.FindOptions{
    Placeholder: "Enter your email",
})

// Find by data-testid (recommended for testing)
elem, err := pilot.Find(ctx, "", &w3pilot.FindOptions{
    TestID: "login-button",
})

// Combine CSS selector with semantic filtering
elem, err := pilot.Find(ctx, "form", &w3pilot.FindOptions{
    Role: "textbox",
    Label: "Password",
})

// Find all buttons
buttons, err := pilot.FindAll(ctx, "", &w3pilot.FindOptions{Role: "button"})

// Find element near another element
elem, err := pilot.Find(ctx, "", &w3pilot.FindOptions{
    Role: "button",
    Near: "#username-input",
})
```

### MCP Tool Usage

Semantic selectors work with `element_click`, `element_type`, `element_fill`, and `element_press` tools:

```json
// Click a button by role and text
{"name": "element_click", "arguments": {"role": "button", "text": "Sign In"}}

// Fill input by label
{"name": "element_fill", "arguments": {"label": "Email", "value": "user@example.com"}}

// Type in input by placeholder
{"name": "element_type", "arguments": {"placeholder": "Search...", "text": "query"}}

// Click by data-testid
{"name": "element_click", "arguments": {"testid": "submit-btn"}}
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
err := pilot.AddInitScript(ctx, `window.testMode = true;`)

// Mock an API
err := pilot.AddInitScript(ctx, `
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
# Inject scripts when launching browser
w3pilot browser launch --init-script=./mock-api.js --init-script=./test-helpers.js

# Or with MCP server
w3pilot mcp --init-script=./mock-api.js --init-script=./test-helpers.js
```

### MCP Tool Usage

```json
{"name": "js_init_script", "arguments": {"script": "window.testMode = true;"}}
```

## Storage State

Save and restore complete browser state including cookies, localStorage, and sessionStorage. Essential for maintaining login sessions across browser restarts.

### SDK Usage

```go
// Get complete storage state
state, err := pilot.StorageState(ctx)

// Save to file
jsonBytes, _ := json.Marshal(state)
os.WriteFile("auth-state.json", jsonBytes, 0600)

// Restore from file
var savedState w3pilot.StorageState
json.Unmarshal(jsonBytes, &savedState)
err := pilot.SetStorageState(ctx, &savedState)

// Clear all storage
err := pilot.ClearStorage(ctx)
```

### MCP Tool Usage

```json
// Save session
{"name": "storage_get_state"}

// Restore session
{"name": "storage_set_state", "arguments": {"state": "<json from storage_get_state>"}}

// Clear all storage
{"name": "storage_clear_all"}
```

## Tracing

Record browser actions with screenshots and DOM snapshots for debugging and test creation.

### SDK Usage

```go
// Start tracing
tracing := pilot.Tracing()
err := tracing.Start(ctx, &w3pilot.TracingStartOptions{
    Screenshots: true,
    Snapshots:   true,
    Title:       "Login Flow Test",
})

// Perform actions...
pilot.Go(ctx, "https://example.com")
elem, _ := pilot.Find(ctx, "button", nil)
elem.Click(ctx, nil)

// Stop and save trace
data, err := tracing.Stop(ctx, nil)
os.WriteFile("trace.zip", data, 0600)
```

### MCP Tool Usage

```json
// Start trace
{"name": "trace_start", "arguments": {"screenshots": true, "title": "My Test"}}

// Stop and get trace data
{"name": "trace_stop", "arguments": {"path": "/tmp/trace.zip"}}
```

## CDP Features (Chrome DevTools Protocol)

W3Pilot provides direct CDP access for advanced profiling and emulation that isn't available through WebDriver BiDi.

### Heap Snapshots

Capture V8 heap snapshots for memory profiling:

```go
// Take heap snapshot
snapshot, err := pilot.TakeHeapSnapshot(ctx, "/tmp/snapshot.heapsnapshot")
fmt.Printf("Snapshot: %s (%d bytes)\n", snapshot.Path, snapshot.Size)

// Load in Chrome DevTools: Memory tab → Load
```

### Network Emulation

Simulate various network conditions:

```go
import "github.com/plexusone/w3pilot/cdp"

// Throttle to Slow 3G
err := pilot.EmulateNetwork(ctx, cdp.NetworkSlow3G)

// Or use presets
err := pilot.EmulateNetwork(ctx, cdp.NetworkFast3G)
err := pilot.EmulateNetwork(ctx, cdp.Network4G)

// Custom conditions
err := pilot.EmulateNetwork(ctx, cdp.NetworkConditions{
    Latency:            100,  // ms
    DownloadThroughput: 500 * 1024,  // 500 KB/s
    UploadThroughput:   250 * 1024,  // 250 KB/s
})

// Clear emulation
err := pilot.ClearNetworkEmulation(ctx)
```

### CPU Emulation

Simulate slower CPUs for performance testing:

```go
import "github.com/plexusone/w3pilot/cdp"

// 4x CPU slowdown (mid-tier mobile)
err := pilot.EmulateCPU(ctx, cdp.CPU4xSlowdown)

// Other presets
err := pilot.EmulateCPU(ctx, cdp.CPU2xSlowdown)
err := pilot.EmulateCPU(ctx, cdp.CPU6xSlowdown)

// Clear emulation
err := pilot.ClearCPUEmulation(ctx)
```

### Direct CDP Access

For advanced use cases, access the CDP client directly:

```go
if pilot.HasCDP() {
    cdpClient := pilot.CDP()

    // Send any CDP command
    result, err := cdpClient.Send(ctx, "Performance.getMetrics", nil)
}
```

## Testing

```bash
# Unit tests
go test -v ./...

# Integration tests
go test -tags=integration -v ./integration/...

# Headless mode
W3PILOT_HEADLESS=1 go test -tags=integration -v ./integration/...
```

## Debug Logging

```bash
W3PILOT_DEBUG=1 w3pilot mcp
```

## Related Projects

- [WebDriver BiDi](https://w3c.github.io/webdriver-bidi/) - Protocol specification

## License

MIT
