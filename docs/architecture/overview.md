# Architecture Overview

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              User Layer                                  │
├─────────────────┬─────────────────┬─────────────────┬──────────────────┤
│    Go Client    │   MCP Server    │      CLI        │  Script Runner   │
│      SDK        │   (75+ tools)   │    (vibium)     │  (vibium run)    │
├─────────────────┴─────────────────┴─────────────────┴──────────────────┤
│                           vibium-go Core                                │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐     │
│  │   Vibe   │ │ Element  │ │ Keyboard │ │  Mouse   │ │  Touch   │     │
│  │ (page)   │ │ (DOM)    │ │ (input)  │ │ (input)  │ │ (input)  │     │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐                   │
│  │ Context  │ │  Clock   │ │ Tracing  │ │  Route   │                   │
│  │(session) │ │ (time)   │ │(capture) │ │(network) │                   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘                   │
├─────────────────────────────────────────────────────────────────────────┤
│                        BiDi Client Layer                                │
│           WebSocket connection to Vibium Clicker                        │
├─────────────────────────────────────────────────────────────────────────┤
│                        Vibium Clicker                                   │
│          Custom commands (vibium:*) + WebDriver BiDi                    │
├─────────────────────────────────────────────────────────────────────────┤
│                   WebDriver BiDi Protocol                               │
│               Bidirectional browser communication                       │
├─────────────────────────────────────────────────────────────────────────┤
│                    Chrome / Chromium                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

## Component Descriptions

### Go Client SDK

The core programmatic API for browser automation:

- **Vibe**: Page-level operations (navigation, screenshots, JS evaluation)
- **Element**: DOM element interactions (click, type, fill, state queries)
- **Input Controllers**: Low-level keyboard, mouse, touch control
- **Context**: Isolated browser sessions with separate cookies/storage
- **Network**: Request interception and modification

### MCP Server

Model Context Protocol server for AI assistant integration:

- 75+ browser automation tools
- Session management with test reporting
- Script recording capability
- Structured error messages with suggestions

### CLI

Command-line interface for scripted automation:

- Subcommand structure (`vibium launch`, `vibium click`, etc.)
- Session persistence between commands
- YAML/JSON script execution

### Script Runner

Deterministic test execution:

- JSON/YAML script format with JSON Schema
- Variable interpolation
- Assertions and data extraction
- Error handling with `continueOnError`

## Data Flow

### MCP Tool Call Flow

```
Claude                    MCP Server              Vibe              Browser
  │                           │                    │                   │
  │──── navigate ────────────▶│                    │                   │
  │                           │──── Go(url) ──────▶│                   │
  │                           │                    │── BiDi request ──▶│
  │                           │                    │◀── BiDi event ────│
  │                           │◀─── url, title ────│                   │
  │◀─── NavigateOutput ───────│                    │                   │
```

### Session Recording Flow

```
Claude                    MCP Server              Recorder
  │                           │                      │
  │── start_recording ───────▶│                      │
  │                           │──── Start() ────────▶│
  │                           │                      │
  │──── navigate ────────────▶│                      │
  │                           │── RecordNavigate() ─▶│
  │                           │                      │
  │──── click ───────────────▶│                      │
  │                           │── RecordClick() ────▶│
  │                           │                      │
  │──── export_script ───────▶│                      │
  │                           │◀── ExportJSON() ────│
  │◀─── JSON script ──────────│                      │
```

## Feature Origin

| Component | Origin | Notes |
|-----------|--------|-------|
| BiDi client | Upstream | WebDriver BiDi protocol |
| Vibe API | Upstream | Parity with JS/Python |
| Element API | Upstream | Parity with JS/Python |
| Input controllers | Upstream | Parity with JS/Python |
| MCP server | Go-specific | AI assistant integration |
| CLI | Go-specific | Command-line automation |
| Script runner | Go-specific | Deterministic replay |
| Session recording | Go-specific | LLM action capture |
| JSON Schema | Go-specific | Script validation |
| Test reporting | Go-specific | Structured diagnostics |

## Key Design Decisions

### WebDriver BiDi

Vibium uses WebDriver BiDi instead of Chrome DevTools Protocol (CDP) for:

- Standardization across browsers
- Bidirectional events (no polling)
- Future-proof design

### Custom Commands

Vibium extends BiDi with `vibium:*` commands for:

- High-level actions (fill, check, selectOption)
- Actionability checks (wait for visible, enabled, stable)
- Page-level operations (screenshot, PDF, evaluate)

### MCP Architecture

The MCP server uses the Model Context Protocol for:

- Standardized tool definitions
- Structured input/output
- Easy AI assistant integration

### Session Recording

Recording captures tool calls (not raw BiDi) for:

- Human-readable scripts
- Portability (same format as CLI scripts)
- Easy editing and customization
