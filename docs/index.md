# WebPilot

Go browser automation library using WebDriver BiDi for real-time bidirectional communication with browsers.

## What is WebPilot?

WebPilot is a browser automation library built for AI agents. It uses the WebDriver BiDi protocol for real-time bidirectional communication with browsers, enabling:

- **Instant feedback** - No polling, real-time events
- **AI-native** - Designed for LLM tool use
- **Cross-browser** - Works with Chrome, Firefox, Edge

## Features

| Component | Description |
|-----------|-------------|
| **Go Client SDK** | Programmatic browser control with full feature parity |
| **MCP Server** | 75+ tools for AI assistant integration |
| **CLI** | Command-line browser automation |
| **Script Runner** | Deterministic JSON/YAML test execution |
| **Session Recording** | Capture LLM actions as replayable scripts |

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         webpilot                                │
├─────────────┬─────────────┬─────────────┬──────────────────────┤
│  Go Client  │ MCP Server  │    CLI      │   Script Runner      │
│    SDK      │  (75 tools) │  (webpilot) │   (webpilot run)     │
├─────────────┴─────────────┴─────────────┴──────────────────────┤
│                    WebDriver BiDi Protocol                     │
├─────────────────────────────────────────────────────────────────┤
│                    Chrome / Chromium                           │
└─────────────────────────────────────────────────────────────────┘
```

## Quick Links

- [Installation](getting-started/installation.md)
- [Quick Start](getting-started/quickstart.md)
- [MCP Server Guide](guide/mcp-server.md)
- [CLI Reference](guide/cli.md)
- [API Reference](reference/api.md)

## Related Projects

| Project | Description |
|---------|-------------|
| [WebDriver BiDi](https://w3c.github.io/webdriver-bidi/) | Protocol specification |
