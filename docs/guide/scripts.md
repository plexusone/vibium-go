# Script Runner

Execute deterministic test scripts in JSON or YAML format.

## Running Scripts

```bash
vibium run test.json
vibium run test.yaml --headless
```

## Script Format

### Basic Structure

```json
{
  "name": "Test Name",
  "description": "What this test does",
  "version": 1,
  "headless": true,
  "baseUrl": "https://example.com",
  "timeout": "30s",
  "steps": [
    {"action": "navigate", "url": "/login"},
    {"action": "fill", "selector": "#email", "value": "user@example.com"}
  ]
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Test name (required) |
| `description` | string | Test description |
| `version` | int | Schema version (1) |
| `headless` | bool | Run headless |
| `baseUrl` | string | Prepended to relative URLs |
| `timeout` | string | Default step timeout |
| `variables` | object | Reusable values |
| `steps` | array | Test steps (required) |

## Actions

### Navigation

```json
{"action": "navigate", "url": "https://example.com"}
{"action": "go", "url": "/page"}
{"action": "back"}
{"action": "forward"}
{"action": "reload"}
```

### Interactions

```json
{"action": "click", "selector": "#button"}
{"action": "dblclick", "selector": "#item"}
{"action": "type", "selector": "#input", "text": "hello"}
{"action": "fill", "selector": "#input", "value": "hello"}
{"action": "clear", "selector": "#input"}
{"action": "press", "selector": "#input", "key": "Enter"}
```

### Form Controls

```json
{"action": "check", "selector": "#checkbox"}
{"action": "uncheck", "selector": "#checkbox"}
{"action": "select", "selector": "#dropdown", "value": "option1"}
```

### Element Interactions

```json
{"action": "hover", "selector": "#menu"}
{"action": "focus", "selector": "#input"}
{"action": "scrollIntoView", "selector": "#footer"}
{"action": "dragTo", "selector": "#source", "target": "#dest"}
{"action": "tap", "selector": "#button"}
```

### Capture

```json
{"action": "screenshot", "file": "page.png"}
{"action": "screenshot", "file": "full.png", "fullPage": true}
{"action": "pdf", "file": "page.pdf"}
```

### JavaScript

```json
{"action": "eval", "script": "document.title"}
```

### Waiting

```json
{"action": "wait", "duration": "1s"}
{"action": "waitForSelector", "selector": "#loaded", "state": "visible"}
{"action": "waitForUrl", "pattern": "/dashboard"}
{"action": "waitForLoad", "loadState": "networkidle"}
```

### Page Actions

```json
{"action": "setViewport", "width": 1920, "height": 1080}
{"action": "newPage"}
{"action": "closePage"}
```

### Assertions

```json
{"action": "assertText", "selector": "#message", "expected": "Success"}
{"action": "assertElement", "selector": "#result"}
{"action": "assertVisible", "selector": "#modal"}
{"action": "assertHidden", "selector": "#loading"}
{"action": "assertUrl", "expected": "https://example.com/success"}
{"action": "assertTitle", "expected": "Dashboard"}
```

### Data Extraction

```json
{"action": "getText", "selector": "#value", "store": "myValue"}
{"action": "getValue", "selector": "#input", "store": "inputValue"}
{"action": "getAttribute", "selector": "#link", "attribute": "href", "store": "url"}
```

## Step Options

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique step ID |
| `name` | string | Step description |
| `timeout` | string | Step timeout override |
| `continueOnError` | bool | Continue if step fails |
| `store` | string | Store result in variable |

## Variables

Define reusable values:

```json
{
  "variables": {
    "email": "test@example.com",
    "password": "secret123"
  },
  "steps": [
    {"action": "fill", "selector": "#email", "value": "${email}"},
    {"action": "fill", "selector": "#password", "value": "${password}"}
  ]
}
```

## YAML Format

```yaml
name: Login Test
headless: true
steps:
  - action: navigate
    url: https://example.com/login
  - action: fill
    selector: "#email"
    value: user@example.com
  - action: fill
    selector: "#password"
    value: secret123
  - action: click
    selector: "#submit"
  - action: assertUrl
    expected: https://example.com/dashboard
```

## JSON Schema

The script format is defined by a JSON Schema:

```bash
# View schema
cat script/vibium-script.schema.json

# Validate script
# (use your preferred JSON Schema validator)
```

Schema location: `script/vibium-script.schema.json`
