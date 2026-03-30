# W3Pilot

Go browser automation library using WebDriver BiDi and Chrome DevTools Protocol for real-time bidirectional communication with browsers.

## What is W3Pilot?

W3Pilot is a browser automation library built for AI agents. It uses a **dual-protocol architecture** connecting to Chrome via both WebDriver BiDi and CDP:

- **Instant feedback** - No polling, real-time events
- **AI-native** - Designed for LLM tool use
- **Advanced profiling** - Heap snapshots, network/CPU emulation via CDP

## Features

| Component | Description |
|-----------|-------------|
| **Go Client SDK** | Programmatic browser control with full feature parity |
| **MCP Server** | 169 tools across 24 namespaces for AI assistants |
| **CLI** | Command-line browser automation |
| **Script Runner** | Deterministic JSON/YAML test execution |
| **Session Recording** | Capture LLM actions as replayable scripts |
| **CDP Integration** | Heap profiling, network/CPU emulation |

## Architecture

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

## Why Dual-Protocol?

W3Pilot combines two complementary protocols for complete browser control:

| Protocol | Purpose | Use Cases |
|----------|---------|-----------|
| **WebDriver BiDi** | Automation & Testing | Page navigation, element interactions, semantic selectors, screenshots, tracing |
| **Chrome DevTools Protocol** | Inspection & Profiling | Heap snapshots, network bodies, CPU/network emulation, coverage |

### BiDi for Automation

WebDriver BiDi (via VibiumDev clicker) provides:

- **Semantic selectors** - Find elements by role, label, text, testid
- **Real-time events** - No polling, instant feedback
- **Cross-browser potential** - W3C standard
- **Human-in-the-loop** - Handle CAPTCHA, SSO, 2FA

### CDP for Profiling

Chrome DevTools Protocol provides:

- **Memory profiling** - Heap snapshots for leak detection
- **Network emulation** - Simulate Slow 3G, Fast 3G, 4G
- **CPU throttling** - Test on low-powered devices
- **Response bodies** - Capture full network content

Both protocols connect to the **same browser instance** - automate with BiDi while profiling with CDP.

## Quick Links

- [Installation](getting-started/installation.md)
- [Quick Start](getting-started/quickstart.md)
- [MCP Server Guide](guide/mcp-server.md)
- [CDP Features](guide/cdp.md)
- [CLI Reference](guide/cli.md)
- [API Reference](reference/api.md)

## Related Projects

| Project | Description |
|---------|-------------|
| [WebDriver BiDi](https://w3c.github.io/webdriver-bidi/) | Protocol specification |
