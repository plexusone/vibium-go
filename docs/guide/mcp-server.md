# MCP Server

The MCP (Model Context Protocol) server provides 75+ browser automation tools for AI assistants like Claude.

## Starting the Server

```bash
# Default (visible browser)
vibium mcp

# Headless mode
vibium mcp --headless

# Custom timeout
vibium mcp --timeout 60s
```

## Configuration

### Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

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

### Environment Variables

| Variable | Description |
|----------|-------------|
| `VIBIUM_DEBUG` | Enable debug logging |
| `VIBIUM_CLICKER_PATH` | Path to clicker binary |

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

The exported JSON can be run with `vibium run`.
