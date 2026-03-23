# MCP Server

The MCP (Model Context Protocol) server provides 80+ browser automation tools for AI assistants like Claude.

## Installation

The MCP server can be run two ways:

1. **Standalone binary** (recommended for MCP clients):

   ```bash
   go install github.com/plexusone/webpilot/cmd/webpilot-mcp@latest
   ```

2. **Via the vibium CLI**:

   ```bash
   go install github.com/plexusone/webpilot/cmd/vibium@latest
   ```

## Starting the Server

### Standalone Binary

```bash
# Default (headless browser)
webpilot-mcp

# Visible browser (for debugging)
webpilot-mcp -headless=false

# Custom timeout
webpilot-mcp -timeout=60s
```

### Via CLI

```bash
# Default (visible browser)
webpilot mcp

# Headless mode
webpilot mcp --headless

# Custom timeout
webpilot mcp --timeout 60s
```

## Client Configuration

### Claude Desktop

Edit the config file:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "webpilot": {
      "command": "webpilot-mcp",
      "args": ["-headless=false"]
    }
  }
}
```

### Claude Code (CLI)

Add to your Claude Code MCP settings:

```json
{
  "mcpServers": {
    "webpilot": {
      "command": "webpilot-mcp",
      "args": ["-headless=false"]
    }
  }
}
```

Or use the CLI command:

```bash
claude mcp add vibium webpilot-mcp -- -headless=false
```

### Kiro CLI

```bash
kiro-cli mcp add --name vibium --command webpilot-mcp --args "-headless=false"
```

### Cursor

Edit `.cursor/mcp.json` in your project or home directory:

```json
{
  "mcpServers": {
    "webpilot": {
      "command": "webpilot-mcp",
      "args": ["-headless=false"]
    }
  }
}
```

### Windsurf

Edit the MCP configuration in Windsurf settings:

```json
{
  "mcpServers": {
    "webpilot": {
      "command": "webpilot-mcp",
      "args": ["-headless=false"]
    }
  }
}
```

### Generic MCP Client

For any MCP-compatible client, use:

- **Command**: `webpilot-mcp`
- **Args**: `["-headless=false"]` (visible browser) or `[]` (headless)

## Command-Line Options

| Option | Default | Description |
|--------|---------|-------------|
| `-headless` | `true` | Run browser without GUI |
| `-project` | `"vibium-tests"` | Project name for reports |
| `-timeout` | `30s` | Default timeout for operations |
| `-init-script` | | JavaScript file to inject before page scripts (repeatable) |

### Init Scripts

Inject JavaScript that runs before any page scripts on every navigation:

```bash
webpilot-mcp -init-script=./mock-api.js -init-script=./test-helpers.js
```

Use cases:

- Mock APIs before page loads
- Disable analytics/tracking
- Inject test utilities
- Set up authentication tokens

## Environment Variables

| Variable | Description |
|----------|-------------|
| `WEBPILOT_DEBUG` | Enable debug logging |
| `WEBPILOT_CLICKER_PATH` | Path to clicker binary |

## Tool Categories

### Browser Management

| Tool | Description |
|------|-------------|
| `browser_launch` | Launch browser instance |
| `browser_quit` | Close browser |

### Navigation

| Tool | Description |
|------|-------------|
| `navigate` | Go to URL |
| `back` | Navigate back |
| `forward` | Navigate forward |
| `reload` | Reload page |

### Element Interactions

| Tool | Description |
|------|-------------|
| `click` | Click element |
| `dblclick` | Double-click element |
| `type` | Type text (append) |
| `fill` | Fill input (replace) |
| `clear` | Clear input |
| `press` | Press key |
| `hover` | Hover over element |
| `focus` | Focus element |

### Form Controls

| Tool | Description |
|------|-------------|
| `check` | Check checkbox |
| `uncheck` | Uncheck checkbox |
| `select_option` | Select dropdown option |
| `set_files` | Set file input |

### Element State

| Tool | Description |
|------|-------------|
| `get_text` | Get element text |
| `get_value` | Get input value |
| `get_attribute` | Get attribute |
| `is_visible` | Check visibility |
| `is_enabled` | Check enabled state |
| `is_checked` | Check checkbox state |

### Page State

| Tool | Description |
|------|-------------|
| `get_title` | Get page title |
| `get_url` | Get current URL |
| `get_content` | Get page HTML |
| `screenshot` | Capture screenshot |
| `pdf` | Generate PDF |

### Waiting

| Tool | Description |
|------|-------------|
| `wait_until` | Wait for element state |
| `wait_for_url` | Wait for URL pattern |
| `wait_for_load` | Wait for load state |

### Human-in-the-Loop

| Tool | Description |
|------|-------------|
| `pause_for_human` | Pause for human action (SSO, CAPTCHA) |
| `get_storage_state` | Export session (cookies + localStorage) |
| `set_storage_state` | Restore saved session |

### Input Controllers

| Tool | Description |
|------|-------------|
| `keyboard_press` | Press key |
| `keyboard_type` | Type text |
| `mouse_click` | Click at coordinates |
| `mouse_move` | Move mouse |

### Script Recording

| Tool | Description |
|------|-------------|
| `start_recording` | Begin recording |
| `stop_recording` | End recording |
| `export_script` | Export as JSON |
| `recording_status` | Check status |

### Tracing

| Tool | Description |
|------|-------------|
| `start_trace` | Start trace with screenshots/snapshots |
| `stop_trace` | Stop and save/return trace ZIP |
| `start_trace_chunk` | Start a trace segment |
| `stop_trace_chunk` | Stop trace segment |
| `start_trace_group` | Group actions logically |
| `stop_trace_group` | End action group |

### Init Scripts

| Tool | Description |
|------|-------------|
| `add_init_script` | Inject JS before page scripts |

### Assertions

| Tool | Description |
|------|-------------|
| `assert_text` | Assert text exists |
| `assert_element` | Assert element exists |
| `assert_url` | Assert URL matches |

## Example Conversation

**User:** Navigate to example.com and click the "More information" link

**Claude:** I'll help you navigate and click the link.

```
[Calls browser_launch]
[Calls navigate with url="https://example.com"]
[Calls click with selector="a"]
```

Done! I've navigated to example.com and clicked the link.

## Session Recording

Record your actions for deterministic replay:

```
[Calls start_recording with name="Example Test"]
[Calls navigate with url="https://example.com"]
[Calls click with selector="a"]
[Calls export_script]
```

The exported JSON can be run with `webpilot run`.
