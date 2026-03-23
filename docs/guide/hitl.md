# Human-in-the-Loop Workflows

Handle authentication flows, CAPTCHAs, and other scenarios requiring human intervention.

## Overview

Some automation tasks require human interaction:

- **SSO/OAuth logins** (Google, Microsoft, Okta)
- **CAPTCHA challenges**
- **2FA/MFA verification**
- **Manual data entry or review**

WebPilot supports these workflows with two key tools:

| Tool | Purpose |
|------|---------|
| `pause_for_human` | Pause automation, wait for human to complete action |
| `set_storage_state` | Restore saved session (cookies + localStorage) |

## Workflow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     First Run (with human)                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  AI Agent                              Human                    │
│  ─────────                             ─────                    │
│  browser_launch {headless: false}                               │
│  navigate → login page                                          │
│  pause_for_human ─────────────────────→ Sees overlay            │
│       ↓                                 Completes SSO           │
│    (waiting)                            Clicks "Continue"       │
│       ↓ ←──────────────────────────────────────┘                │
│  get_storage_state → save to file                               │
│  ... continue automation ...                                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                  Subsequent Runs (automated)                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  AI Agent                                                       │
│  ─────────                                                      │
│  browser_launch {headless: true}                                │
│  set_storage_state ← load from file                             │
│  navigate → dashboard (already logged in!)                      │
│  ... continue automation ...                                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## MCP Tools

### pause_for_human

Pauses automation and displays a visual overlay. The human completes their action and clicks "Continue" to resume.

**Parameters:**

| Name | Type | Default | Description |
|------|------|---------|-------------|
| `message` | string | "Please complete the required action..." | Message displayed to the human |
| `timeout_ms` | int | 300000 (5 min) | Maximum wait time in milliseconds |

**Example:**

```json
{
  "name": "pause_for_human",
  "arguments": {
    "message": "Please complete the Google SSO login",
    "timeout_ms": 300000
  }
}
```

**Visual overlay:**

```
┌──────────────────────────────────────────────────────────────┐
│ Please complete the Google SSO login          [ Continue ]   │
└──────────────────────────────────────────────────────────────┘
```

The overlay appears at the top of the browser window with a gradient background. When the human clicks "Continue", automation resumes.

### get_storage_state

Exports cookies and localStorage as JSON for session persistence.

**Returns:**

```json
{
  "state": "{\"cookies\":[...],\"origins\":[{\"origin\":\"https://example.com\",\"localStorage\":{...}}]}"
}
```

### set_storage_state

Restores a previously saved session state (cookies and localStorage).

**Parameters:**

| Name | Type | Description |
|------|------|-------------|
| `state` | string | JSON from `get_storage_state` |

**Example:**

```json
{
  "name": "set_storage_state",
  "arguments": {
    "state": "{\"cookies\":[...],\"origins\":[...]}"
  }
}
```

## Complete Example: Google SSO

### First Run (Human Login)

**User prompt:**

> Navigate to my-app.example.com and log me in. I'll complete the Google SSO manually. Then scrape the dashboard data.

**Claude's execution:**

```
[browser_launch: headless=false]
Browser launched in visible mode.

[navigate: url="https://my-app.example.com"]
Navigated to login page.

[pause_for_human: message="Please complete Google SSO login, then click Continue"]
Waiting for human...

// Human completes SSO, clicks Continue

[get_storage_state]
Session saved with 12 cookies and localStorage for 1 origin.

[get_text: selector=".dashboard-content"]
Dashboard data retrieved.

[browser_quit]
```

### Subsequent Run (Automated)

**User prompt:**

> Use my saved session to scrape the dashboard again.

**Claude's execution:**

```
[browser_launch: headless=true]
Browser launched in headless mode.

[set_storage_state: state="<saved JSON>"]
Restored 12 cookies and localStorage for 1 origin.

[navigate: url="https://my-app.example.com/dashboard"]
Already logged in! Navigated to dashboard.

[get_text: selector=".dashboard-content"]
Dashboard data retrieved.
```

## Session Persistence

### Saving Sessions

Save the session state to a file for reuse:

```
User: Save my session to google-session.json

Claude:
[get_storage_state]
Session state retrieved. I'll save this to google-session.json.
[Write file: google-session.json]
```

### Loading Sessions

Restore a saved session:

```
User: Load my session from google-session.json and go to the dashboard

Claude:
[Read file: google-session.json]
[browser_launch: headless=true]
[set_storage_state: state="<file contents>"]
[navigate: url="https://my-app.example.com/dashboard"]
Successfully logged in using saved session.
```

### Session Expiration

Sessions expire based on the application's policies. Common patterns:

| Token Type | Typical Lifetime |
|------------|------------------|
| Session cookie | Browser close or hours |
| Remember me | Days to weeks |
| OAuth access token | Hours |
| OAuth refresh token | Days to months |

When a session expires, you'll need to do another human login:

```
User: Go to the dashboard

Claude:
[set_storage_state: state="<saved JSON>"]
[navigate: url="https://my-app.example.com/dashboard"]
Redirected to login page - session has expired.

The saved session is no longer valid. Would you like me to:
1. Launch a visible browser so you can log in again
2. Try a different saved session
```

## Configuration

### MCP Server Setup

For HITL workflows, start the MCP server with visible browser:

```bash
webpilot mcp --headless=false
```

Or configure in Claude Desktop:

```json
{
  "mcpServers": {
    "webpilot": {
      "command": "webpilot",
      "args": ["mcp"]
    }
  }
}
```

!!! note
    The `browser_launch` tool accepts a `headless` parameter that overrides the server default. For HITL workflows, always use `headless: false`.

### Timeout Configuration

For long authentication flows (e.g., hardware security keys), increase the timeout:

```json
{
  "name": "pause_for_human",
  "arguments": {
    "message": "Please complete 2FA with your security key",
    "timeout_ms": 600000
  }
}
```

## Best Practices

### 1. Use Descriptive Messages

Tell the human exactly what to do:

```
❌ "Please continue"
✅ "Please complete the Google SSO login, then click Continue"
✅ "Please solve the CAPTCHA, then click Continue"
✅ "Please approve the 2FA prompt on your phone, then click Continue"
```

### 2. Wait for the Right State

After human action, verify you're in the expected state:

```
[pause_for_human: message="Please log in"]
[wait_for_url: pattern="**/dashboard*", timeout_ms=10000]
[assert_element: selector=".user-menu"]
Login confirmed.
```

### 3. Handle Session Expiration Gracefully

Check if the session is still valid before proceeding:

```
[set_storage_state: state="<saved>"]
[navigate: url="https://app.example.com/dashboard"]
[get_url]

If URL contains "/login":
  "Session expired. Let me launch a visible browser for re-authentication."
  [browser_quit]
  [browser_launch: headless=false]
  [navigate: url="https://app.example.com"]
  [pause_for_human: message="Please log in again"]
```

### 4. Secure Session Storage

Session files contain sensitive authentication data:

- Store in a secure location (not in git)
- Use appropriate file permissions
- Consider encrypting at rest
- Delete when no longer needed

```bash
# Secure permissions
chmod 600 session.json

# Add to .gitignore
echo "*.session.json" >> .gitignore
```

### 5. Separate Sessions by Environment

Maintain different sessions for different environments:

```
sessions/
├── prod-session.json
├── staging-session.json
└── dev-session.json
```

## Supported Authentication Flows

| Flow | Support |
|------|---------|
| Google SSO | Full |
| Microsoft/Azure AD | Full |
| Okta/Auth0 | Full |
| GitHub OAuth | Full |
| SAML | Full |
| Username/Password | Full (but use `fill` directly) |
| CAPTCHA | Full |
| SMS 2FA | Full |
| TOTP (Authenticator app) | Full |
| Hardware keys (WebAuthn) | Full |

## Troubleshooting

### Overlay Not Appearing

If the overlay doesn't appear:

1. Ensure browser is in visible mode (`headless: false`)
2. Check if page has strict CSP blocking inline styles
3. Try refreshing the page before `pause_for_human`

### Session Not Restoring

If `set_storage_state` doesn't restore the session:

1. Verify the JSON is valid
2. Check if cookies have expired
3. Some sites use additional fingerprinting beyond cookies
4. Try navigating to the origin before setting storage state

### Timeout During Human Action

If the timeout is too short:

```json
{
  "name": "pause_for_human",
  "arguments": {
    "message": "Please complete login",
    "timeout_ms": 600000
  }
}
```

Default is 5 minutes (300000ms). Maximum recommended is 10 minutes (600000ms).
