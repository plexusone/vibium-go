# Vibium Go

Go client and tooling for the [Vibium](https://github.com/VibiumDev/vibium) browser automation platform.

## What is Vibium?

Vibium is a browser automation platform built for AI agents. It uses the WebDriver BiDi protocol for real-time bidirectional communication with browsers, enabling:

- **Instant feedback** - No polling, real-time events
- **AI-native** - Designed for LLM tool use
- **Cross-browser** - Works with Chrome, Firefox, Edge

## What is Vibium Go?

This project provides Go tooling for Vibium:

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
│                         vibium-go                                │
├─────────────┬─────────────┬─────────────┬──────────────────────┤
│  Go Client  │ MCP Server  │    CLI      │   Script Runner      │
│    SDK      │  (75 tools) │  (vibium)   │   (vibium run)       │
├─────────────┴─────────────┴─────────────┴──────────────────────┤
│                    WebDriver BiDi Protocol                       │
├─────────────────────────────────────────────────────────────────┤
│                   Vibium Clicker (upstream)                      │
├─────────────────────────────────────────────────────────────────┤
│                    Chrome / Chromium                             │
└─────────────────────────────────────────────────────────────────┘
```

## Feature Origin

| Feature | Source |
|---------|--------|
| Core browser automation | Upstream Vibium (JS/Python parity) |
| MCP server | Go-specific |
| CLI | Go-specific |
| Script runner | Go-specific |
| Session recording | Go-specific |

## Quick Links

- [Installation](getting-started/installation.md)
- [Quick Start](getting-started/quickstart.md)
- [MCP Server Guide](guide/mcp-server.md)
- [CLI Reference](guide/cli.md)
- [API Reference](reference/api.md)

## Related Projects

| Project | Description |
|---------|-------------|
| [vibium-wcag](https://github.com/agentplexus/vibium-wcag) | WCAG 2.2 accessibility testing using vibium-go |
| [omnillm](https://github.com/agentplexus/omnillm) | Unified LLM client for multiple providers |
