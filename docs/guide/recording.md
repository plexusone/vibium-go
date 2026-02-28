# Session Recording

Convert AI-assisted exploration into deterministic test scripts.

## Overview

Session recording captures MCP tool calls and exports them as JSON scripts that can be replayed without an LLM.

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│  Markdown Test   │     │   LLM + MCP      │     │   JSON Script    │
│  Plan (English)  │ ──▶ │   (exploration)  │ ──▶ │ (deterministic)  │
└──────────────────┘     └──────────────────┘     └──────────────────┘
```

## Workflow

### 1. Write Test Plan (Markdown)

```markdown
# Login Test

1. Navigate to the login page
2. Enter email: test@example.com
3. Enter password: secret123
4. Click the login button
5. Verify we reach the dashboard
```

### 2. Execute with LLM

Ask Claude to execute your test plan while recording:

> "Start recording, then follow my test plan: [paste plan]. When done, export the script."

Claude will:

1. Call `start_recording`
2. Execute each step using MCP tools
3. Handle selectors, waits, and edge cases
4. Call `export_script` when done

### 3. Get JSON Script

The exported script:

```json
{
  "name": "Login Test",
  "version": 1,
  "steps": [
    {"action": "navigate", "url": "https://example.com/login"},
    {"action": "fill", "selector": "#email", "value": "test@example.com"},
    {"action": "fill", "selector": "#password", "value": "secret123"},
    {"action": "click", "selector": "button[type='submit']"},
    {"action": "waitForUrl", "pattern": "/dashboard"}
  ]
}
```

### 4. Run Deterministically

```bash
vibium run login-test.json
```

## MCP Recording Tools

### start_recording

Begin recording actions.

**Parameters:**

| Name | Type | Description |
|------|------|-------------|
| `name` | string | Script name |
| `description` | string | What the script tests |
| `baseUrl` | string | Base URL for relative paths |

**Example:**

```json
{
  "name": "start_recording",
  "arguments": {
    "name": "Login Test",
    "description": "Verify user can log in"
  }
}
```

### stop_recording

Stop recording actions.

**Returns:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `stepCount` | int | Number of steps recorded |

### export_script

Export recorded actions as JSON.

**Parameters:**

| Name | Type | Description |
|------|------|-------------|
| `format` | string | `json` or `yaml` (default: json) |

**Returns:**

| Field | Type | Description |
|-------|------|-------------|
| `script` | string | The script content |
| `stepCount` | int | Number of steps |
| `format` | string | Output format |

### recording_status

Check if recording is active.

**Returns:**

| Field | Type | Description |
|-------|------|-------------|
| `recording` | bool | Is recording active |
| `stepCount` | int | Steps recorded so far |

### clear_recording

Clear recorded steps without stopping.

## Recorded Actions

The following MCP tools are recorded:

| Category | Actions |
|----------|---------|
| Navigation | `navigate` |
| Clicks | `click`, `dblclick` |
| Input | `type`, `fill`, `clear`, `press` |
| Forms | `check`, `uncheck`, `select_option` |
| Element | `hover`, `focus`, `scroll_into_view` |
| Capture | `screenshot`, `evaluate` |

## Best Practices

### Let the LLM Find Selectors

Don't specify exact selectors in your test plan:

```markdown
❌ Click the element with selector "#submit-btn-v2"
✅ Click the submit button
```

The LLM will find the appropriate selector and it will be captured.

### Include Assertions

Ask the LLM to verify outcomes:

```markdown
1. Fill in the login form
2. Click submit
3. Verify the dashboard appears
```

### Handle Dynamic Content

The LLM can add appropriate waits:

```markdown
1. Click "Load More"
2. Wait for new items to appear
3. Verify at least 10 items are shown
```

### Review Before Committing

Always review exported scripts:

- Check selectors are stable (prefer IDs, data-testid)
- Verify waits are sufficient
- Add `continueOnError` where appropriate

## Example Session

**User:** Start recording a test called "Search Test", then: navigate to Google, search for "vibium", and verify results appear.

**Claude:** I'll start recording and execute your test.

```
[start_recording: name="Search Test"]
[browser_launch]
[navigate: url="https://google.com"]
[fill: selector="input[name='q']", value="vibium"]
[press: selector="input[name='q']", key="Enter"]
[wait_until: selector="#search", state="visible"]
[export_script]
```

Here's your recorded script:

```json
{
  "name": "Search Test",
  "steps": [
    {"action": "navigate", "url": "https://google.com"},
    {"action": "fill", "selector": "input[name='q']", "value": "vibium"},
    {"action": "press", "selector": "input[name='q']", "key": "Enter"},
    {"action": "waitForSelector", "selector": "#search", "state": "visible"}
  ]
}
```

You can run this with `vibium run search-test.json`.
