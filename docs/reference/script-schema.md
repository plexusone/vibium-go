# Script Schema Reference

The script format is defined by a JSON Schema generated from Go types.

## Schema Location

```
script/vibium-script.schema.json
```

## Schema URL

```
https://github.com/grokify/vibium-go/script/vibium-script.schema.json
```

## Top-Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Human-readable name |
| `description` | string | | Additional context |
| `version` | integer | | Schema version (default: 1) |
| `headless` | boolean | | Run in headless mode |
| `baseUrl` | string | | Prepended to relative URLs |
| `timeout` | string | | Default timeout (e.g., "30s") |
| `variables` | object | | Reusable values |
| `steps` | array | ✅ | Automation steps |

## Step Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | ✅ | Action type (see below) |
| `id` | string | | Unique step identifier |
| `name` | string | | Human-readable description |
| `selector` | string | | CSS selector |
| `url` | string | | Target URL |
| `value` | string | | Input value |
| `text` | string | | Alias for value |
| `key` | string | | Key to press |
| `script` | string | | JavaScript code |
| `file` | string | | Output file path |
| `files` | array | | File paths for input |
| `timeout` | string | | Timeout override |
| `duration` | string | | Wait duration |
| `fullPage` | boolean | | Full page screenshot |
| `target` | string | | Drag target selector |
| `x` | number | | X coordinate |
| `y` | number | | Y coordinate |
| `width` | integer | | Viewport width |
| `height` | integer | | Viewport height |
| `state` | string | | Wait state |
| `pattern` | string | | URL pattern |
| `loadState` | string | | Load state |
| `expected` | string | | Expected value |
| `attribute` | string | | Attribute name |
| `store` | string | | Variable to store result |
| `continueOnError` | boolean | | Continue on failure |

## Action Types

### Navigation

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `navigate` | `url` | Navigate to URL |
| `go` | `url` | Alias for navigate |
| `back` | | Navigate back |
| `forward` | | Navigate forward |
| `reload` | | Reload page |

### Interactions

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `click` | `selector` | Click element |
| `dblclick` | `selector` | Double-click |
| `type` | `selector`, `text` | Type text (append) |
| `fill` | `selector`, `value` | Fill input (replace) |
| `clear` | `selector` | Clear input |
| `press` | `selector`, `key` | Press key |

### Form Controls

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `check` | `selector` | Check checkbox |
| `uncheck` | `selector` | Uncheck checkbox |
| `select` | `selector`, `value` | Select option |
| `setFiles` | `selector`, `files` | Set file input |

### Element Interactions

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `hover` | `selector` | Hover over element |
| `focus` | `selector` | Focus element |
| `scrollIntoView` | `selector` | Scroll into view |
| `dragTo` | `selector`, `target` | Drag to target |
| `tap` | `selector` | Touch tap |

### Capture

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `screenshot` | `file` | Page screenshot |
| `pdf` | `file` | Generate PDF |

### JavaScript

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `eval` | `script` | Execute JavaScript |

### Waiting

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `wait` | `duration` | Wait for duration |
| `waitForSelector` | `selector` | Wait for element |
| `waitForUrl` | `pattern` | Wait for URL |
| `waitForLoad` | `loadState` | Wait for load state |

### Page Actions

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `setViewport` | `width`, `height` | Set viewport |
| `newPage` | | Create new page |
| `closePage` | | Close page |

### Input Controllers

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `keyboardPress` | `key` | Press key |
| `keyboardType` | `text` | Type text |
| `mouseClick` | `x`, `y` | Click at coords |
| `mouseMove` | `x`, `y` | Move mouse |

### Assertions

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `assertText` | `selector`, `expected` | Assert text |
| `assertElement` | `selector` | Assert exists |
| `assertValue` | `selector`, `expected` | Assert value |
| `assertVisible` | `selector` | Assert visible |
| `assertHidden` | `selector` | Assert hidden |
| `assertUrl` | `expected` | Assert URL |
| `assertTitle` | `expected` | Assert title |
| `assertAttribute` | `selector`, `attribute`, `expected` | Assert attr |

### Data Extraction

| Action | Required Fields | Description |
|--------|-----------------|-------------|
| `getText` | `selector`, `store` | Get text |
| `getValue` | `selector`, `store` | Get value |
| `getAttribute` | `selector`, `attribute`, `store` | Get attr |
| `getUrl` | `store` | Get URL |
| `getTitle` | `store` | Get title |

## State Values

For `waitForSelector`:

- `visible` - Element is visible
- `hidden` - Element is hidden
- `attached` - Element is in DOM
- `detached` - Element removed from DOM

## Load State Values

For `waitForLoad`:

- `load` - Window load event
- `domcontentloaded` - DOMContentLoaded event
- `networkidle` - No network activity

## Examples

### Minimal

```json
{
  "name": "Simple Test",
  "steps": [
    {"action": "navigate", "url": "https://example.com"}
  ]
}
```

### With Variables

```json
{
  "name": "Login Test",
  "variables": {
    "baseUrl": "https://example.com",
    "email": "user@example.com"
  },
  "steps": [
    {"action": "navigate", "url": "${baseUrl}/login"},
    {"action": "fill", "selector": "#email", "value": "${email}"}
  ]
}
```

### With Assertions

```json
{
  "name": "Form Test",
  "steps": [
    {"action": "navigate", "url": "https://example.com"},
    {"action": "assertTitle", "expected": "Example Domain"},
    {"action": "assertElement", "selector": "h1"},
    {"action": "getText", "selector": "h1", "store": "heading"},
    {"action": "assertText", "selector": "h1", "expected": "Example Domain"}
  ]
}
```

## Regenerating Schema

The schema is generated from Go types:

```bash
go run ./cmd/genscriptschema > script/vibium-script.schema.json
schemago lint script/vibium-script.schema.json
```

Source types: `script/types.go`
