# MCP Tools Reference

Complete reference for all **169 MCP tools across 24 namespaces**.

## Naming Convention

All tool names follow the pattern: `{namespace}_{verb}_{target}`

**Naming Principles:**

1. **Explicit verbs**: `element_get_text` not `element_text`
2. **No abbreviations**: `wait_for_function` not `wait_fn`
3. **Full words**: `human_pause` not `hitl_pause`, `accessibility_snapshot` not `a11y_snapshot`
4. **Consistent verb patterns**: get/set, start/stop, is/has

| Namespace | Purpose | Count |
|-----------|---------|------:|
| `accessibility_` | Accessibility tree | 1 |
| `batch_` | Multi-step batch execution | 1 |
| `browser_` | Browser lifecycle | 2 |
| `cdp_` | Chrome DevTools Protocol | 20 |
| `config_` | Configuration | 1 |
| `console_` | Console messages | 2 |
| `dialog_` | Dialog handling | 2 |
| `element_` | Element interactions and state | 33 |
| `frame_` | Frame selection | 2 |
| `http_` | HTTP requests in browser context | 1 |
| `human_` | Human-in-the-loop | 1 |
| `input_` | Low-level keyboard/mouse/touch | 12 |
| `js_` | JavaScript execution | 4 |
| `network_` | Network requests and mocking | 6 |
| `page_` | Page navigation, state, screenshots, emulation | 20 |
| `record_` | Script recording | 5 |
| `state_` | Named state snapshots | 4 |
| `storage_` | Cookies, localStorage, sessionStorage | 17 |
| `tab_` | Tab management | 3 |
| `test_` | Assertions, verification, reporting | 16 |
| `trace_` | Tracing | 6 |
| `video_` | Video recording | 2 |
| `wait_` | Waiting operations | 6 |
| `workflow_` | High-level automation workflows | 2 |

## Machine-Readable Format

For programmatic access, use:

```bash
# Export as JSON
w3pilot mcp --list-tools > mcp-tools.json

# Or via the standalone binary
w3pilot-mcp --list-tools
```

The JSON format follows the MCP protocol structure:

```json
{
  "tools": [
    {
      "name": "tool_name",
      "description": "Tool description",
      "category": "Category Name"
    }
  ],
  "categories": {
    "Category Name": 5
  },
  "total": 169
}
```

A pre-generated version is available at [mcp-tools.json](mcp-tools.json).

## Browser Management

### browser_launch

Launch a browser instance.

**Input:**

| Field | Type | Description |
|-------|------|-------------|
| `headless` | boolean | Run without GUI |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

### browser_quit

Close the browser.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

## Navigation

### page_navigate

Navigate to a URL.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `url` | string | ✅ | Target URL |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Current URL |
| `title` | string | Page title |

### page_go_back

Navigate back in history.

### page_go_forward

Navigate forward in history.

### page_reload

Reload the current page.

### page_scroll

Scroll the page or a specific element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `direction` | string | ✅ | Scroll direction: up, down, left, right |
| `amount` | integer | | Pixels to scroll (0 for full page) |
| `selector` | string | | CSS selector to scroll within |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

## Element Interactions

### element_click

Click an element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `timeout_ms` | integer | | Timeout (default: 5000) |

### element_double_click

Double-click an element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `timeout_ms` | integer | | Timeout |

### element_type

Type text into an element (appends).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `text` | string | ✅ | Text to type |
| `timeout_ms` | integer | | Timeout |

### element_fill

Fill an input (replaces content).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `value` | string | ✅ | Value to fill |
| `timeout_ms` | integer | | Timeout |

### element_fill_form

Fill multiple form fields at once.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `fields` | array | ✅ | Array of {selector, value} objects |
| `timeout_ms` | integer | | Timeout per field |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `filled` | integer | Number of fields filled |
| `errors` | array | Any errors encountered |

### element_clear

Clear an input element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

### element_press

Press a key on an element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `key` | string | ✅ | Key (e.g., "Enter") |

## Form Controls

### element_check

Check a checkbox.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

### element_uncheck

Uncheck a checkbox.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

### element_select

Select dropdown option(s).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `values` | array | | Option values |
| `labels` | array | | Option labels |
| `indexes` | array | | Option indexes |

### element_set_files

Set files on a file input.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `files` | array | ✅ | File paths |

## Element State

### element_get_text

Get element text content.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `text` | string | Text content |

### element_get_value

Get input element value.

### element_get_inner_html

Get element innerHTML.

### element_get_outer_html

Get element outerHTML (including the element itself).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `html` | string | Element's outer HTML |

### element_get_inner_text

Get element innerText.

### element_get_attribute

Get element attribute.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `name` | string | ✅ | Attribute name |

### element_get_bounding_box

Get element bounding box.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `x` | number | X position |
| `y` | number | Y position |
| `width` | number | Width |
| `height` | number | Height |

### element_is_visible

Check if element is visible.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `visible` | boolean | Visibility state |

### element_is_hidden

Check if element is hidden.

### element_is_enabled

Check if element is enabled.

### element_is_checked

Check if checkbox/radio is checked.

### element_is_editable

Check if element is editable.

### element_get_role

Get ARIA role.

### element_get_label

Get accessible label.

## Page State

### page_get_title

Get page title.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Page title |

### page_get_url

Get current URL.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Current URL |

### page_get_content

Get page HTML content.

### page_set_content

Set page HTML content.

### page_get_viewport

Get viewport dimensions.

### page_set_viewport

Set viewport dimensions.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `width` | integer | ✅ | Width |
| `height` | integer | ✅ | Height |

## Screenshots & PDF

### page_screenshot

Capture page screenshot.

**Input:**

| Field | Type | Description |
|-------|------|-------------|
| `format` | string | "base64" or "file" |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `data` | string | Base64 image data |

### element_screenshot

Capture element screenshot.

### page_pdf

Generate PDF.

## JavaScript

### js_evaluate

Execute JavaScript.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `script` | string | ✅ | JavaScript code |
| `max_result_size` | integer | | Truncate result to this many characters (0 = no truncation) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `result` | any | Evaluation result |
| `truncated` | boolean | Whether result was truncated |

### element_evaluate

Evaluate JavaScript with element.

### js_add_script

Inject JavaScript.

### js_add_style

Inject CSS.

## HTTP Requests

### http_request

Make HTTP requests in the browser context with session cookies.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `url` | string | ✅ | Target URL |
| `method` | string | | HTTP method (default: GET) |
| `headers` | object | | Request headers |
| `body` | string | | Request body for POST/PUT/PATCH |
| `content_type` | string | | Content-Type header |
| `max_body_length` | integer | | Truncate response body (default: 8192) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `status` | integer | HTTP status code |
| `statusText` | string | HTTP status text |
| `headers` | object | Response headers |
| `body` | string | Response body |
| `truncated` | boolean | Whether body was truncated |
| `url` | string | Final URL (after redirects) |

## Batch Execution

### batch_execute

Execute multiple operations in a single MCP call to reduce round-trip latency.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `steps` | array | ✅ | Array of tool steps to execute |
| `stop_on_error` | boolean | | Stop on first error (default: true) |
| `continue_on_error` | boolean | | Continue even if steps fail |

Each step object:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `tool` | string | ✅ | Tool name to execute |
| `args` | object | | Arguments for the tool |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `results` | array | Results for each step |
| `total_steps` | integer | Total steps in batch |
| `success_count` | integer | Successful steps |
| `failure_count` | integer | Failed steps |
| `stopped_early` | boolean | Whether stopped before completion |
| `total_duration_ms` | integer | Total execution time |

**Supported Tools:**

| Category | Tools |
|----------|-------|
| Navigation | `page_navigate`, `page_go_back`, `page_go_forward`, `page_reload` |
| Page Info | `page_get_title`, `page_get_url`, `page_screenshot` |
| Elements | `element_click`, `element_fill`, `element_type`, `element_get_text` |
| JavaScript | `js_evaluate` |
| Waiting | `wait_for_selector`, `wait_for_load` |
| HTTP | `http_request` |

**Example:**

```json
{
  "tool": "batch_execute",
  "arguments": {
    "steps": [
      {"tool": "page_navigate", "args": {"url": "https://example.com/login"}},
      {"tool": "element_fill", "args": {"selector": "#username", "value": "user"}},
      {"tool": "element_fill", "args": {"selector": "#password", "value": "pass"}},
      {"tool": "element_click", "args": {"selector": "#login"}},
      {"tool": "wait_for_selector", "args": {"selector": "#dashboard"}}
    ]
  }
}
```

## Waiting

### wait_for_state

Wait for element state.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `state` | string | ✅ | visible/hidden/attached/detached |

### wait_for_url

Wait for URL pattern.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `pattern` | string | ✅ | URL pattern |

### wait_for_load

Wait for load state.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `state` | string | ✅ | load/domcontentloaded/networkidle |

### wait_for_function

Wait for JavaScript function.

## Input Controllers

### input_keyboard_press

Press a key.

### input_keyboard_down

Hold a key.

### input_keyboard_up

Release a key.

### input_keyboard_type

Type text via keyboard.

### input_mouse_click

Click at coordinates.

### input_mouse_move

Move mouse.

### input_mouse_down

Press mouse button.

### input_mouse_up

Release mouse button.

### input_mouse_wheel

Scroll mouse wheel.

### input_touch_tap

Tap at coordinates.

### input_touch_swipe

Swipe gesture.

### input_mouse_drag

Drag from one point to another using the mouse.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `start_x` | number | ✅ | Starting X coordinate |
| `start_y` | number | ✅ | Starting Y coordinate |
| `end_x` | number | ✅ | Ending X coordinate |
| `end_y` | number | ✅ | Ending Y coordinate |
| `steps` | integer | | Number of intermediate steps (default: 10) |

## Page Management

### page_new

Create new page/tab.

### page_get_count

Get page count.

### page_close

Close current page.

### page_bring_to_front

Activate page.

## Emulation

### page_emulate_media

Emulate CSS media features for accessibility testing.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `media` | string | | Media type: "screen" or "print" |
| `color_scheme` | string | | Color scheme: "light", "dark", "no-preference" |
| `reduced_motion` | string | | Reduced motion: "reduce", "no-preference" |
| `forced_colors` | string | | Forced colors: "active", "none" |
| `contrast` | string | | Contrast: "more", "less", "no-preference" |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `settings` | array | Applied settings |

**Accessibility Testing Use Cases:**

- **color_scheme**: Test dark/light mode support for users with light sensitivity
- **reduced_motion**: Test that animations are disabled for users with vestibular disorders
- **forced_colors**: Test Windows High Contrast Mode compatibility
- **contrast**: Test increased/decreased contrast for users with low vision

### page_set_geolocation

Set the browser's geolocation.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `latitude` | number | Yes | Latitude coordinate |
| `longitude` | number | Yes | Longitude coordinate |
| `accuracy` | number | | Accuracy in meters |

## Tab Management

### tab_list

List all open browser tabs.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `tabs` | array | Tab information (index, id, url, title) |
| `count` | integer | Number of open tabs |
| `current_tab` | integer | Index of the current tab |

### tab_select

Switch to a specific tab.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | integer | | Tab index (0-based) |
| `id` | string | | Tab ID from tab_list |

### tab_close

Close a specific tab.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | integer | | Tab index to close |
| `id` | string | | Tab ID to close |

## Dialog Handling

### dialog_handle

Handle a browser dialog (alert, confirm, prompt, beforeunload).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | Action: "accept" or "dismiss" |
| `prompt_text` | string | | Text for prompt dialogs |

### dialog_get

Get information about the current dialog.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `has_dialog` | boolean | Whether a dialog is open |
| `dialog_type` | string | alert, confirm, prompt, beforeunload |
| `message` | string | Dialog message |
| `default_value` | string | Default value for prompt dialogs |

## Console Messages

### console_get_messages

Get console messages from the page.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `level` | string | | Filter by level: log, info, warn, error, debug |
| `clear` | boolean | | Clear messages after retrieving |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `messages` | array | Console message objects |
| `count` | integer | Number of messages |

Each message contains:

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Message type (log, info, warn, error, debug) |
| `text` | string | Message text |
| `args` | array | Additional arguments |
| `url` | string | Source URL |
| `line` | integer | Source line number |

### console_clear

Clear all buffered console messages.

## Network Requests

### network_get_requests

Get captured network requests.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `url_pattern` | string | | Filter by URL pattern (glob or regex) |
| `method` | string | | Filter by HTTP method |
| `resource_type` | string | | Filter by resource type |
| `clear` | boolean | | Clear requests after retrieving |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `requests` | array | Network request objects |
| `count` | integer | Number of requests |

Each request contains:

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Request URL |
| `method` | string | HTTP method |
| `resource_type` | string | Resource type |
| `status` | integer | Response status code |
| `status_text` | string | Response status text |
| `response_size` | integer | Response size in bytes |

### network_clear

Clear all buffered network requests.

## Network Mocking

### network_route

Register a mock response for requests matching a URL pattern.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `pattern` | string | ✅ | URL pattern (glob or regex, e.g., `**/api/*`) |
| `status` | integer | | HTTP status code (default: 200) |
| `body` | string | | Response body content |
| `content_type` | string | | Content-Type header (default: application/json) |
| `headers` | object | | Additional response headers |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `pattern` | string | Registered pattern |

### network_list_routes

List all active route handlers.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `routes` | array | Route info objects |
| `count` | integer | Number of active routes |

Each route contains:

| Field | Type | Description |
|-------|------|-------------|
| `pattern` | string | URL pattern |
| `status` | integer | Response status code |
| `content_type` | string | Content-Type header |

### network_unroute

Remove a previously registered route handler.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `pattern` | string | ✅ | URL pattern to unregister |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

### network_set_offline

Set the browser's network state for offline mode testing.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `offline` | boolean | ✅ | Set to true for offline mode |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `offline` | boolean | Current offline state |

## Cookies & Storage

### storage_get_cookies

Get browser cookies.

### storage_set_cookies

Set browser cookies.

### storage_clear_cookies

Clear all cookies.

### storage_get_state

Get complete browser storage state including cookies, localStorage, and sessionStorage as JSON. This can be saved to a file and later restored using `storage_set_state` to resume a session.

**Output:**

Returns JSON containing:

- `cookies`: Array of all cookies
- `origins`: Array of origin storage, each with:
  - `origin`: The origin URL (e.g., "https://example.com")
  - `localStorage`: Key-value map of localStorage items
  - `sessionStorage`: Key-value map of sessionStorage items (for current page's origin)

### storage_set_state

Restore browser storage from JSON (output of `storage_get_state`). Restores cookies, localStorage, and sessionStorage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `state` | string | Yes | JSON from storage_get_state |

### storage_clear_all

Clear all browser storage including cookies, localStorage, and sessionStorage.

## LocalStorage

### storage_local_get

Get a value from localStorage by key.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to retrieve |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `key` | string | The key |
| `value` | string | The value (null if not found) |

### storage_local_set

Set a value in localStorage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to set |
| `value` | string | Yes | Value to store |

### storage_local_delete

Delete a key from localStorage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to delete |

### storage_local_clear

Clear all localStorage data for the current origin.

### storage_local_list

List all keys and values in localStorage.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `items` | object | Key-value pairs |
| `count` | integer | Number of items |

## SessionStorage

### storage_session_get

Get a value from sessionStorage by key.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to retrieve |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `key` | string | The key |
| `value` | string | The value (null if not found) |

### storage_session_set

Set a value in sessionStorage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to set |
| `value` | string | Yes | Value to store |

### storage_session_delete

Delete a key from sessionStorage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to delete |

### storage_session_clear

Clear all sessionStorage data for the current origin.

### storage_session_list

List all keys and values in sessionStorage.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `items` | object | Key-value pairs |
| `count` | integer | Number of items |

## Script Recording

### record_start

Begin recording actions.

**Input:**

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Script name |
| `description` | string | Description |
| `baseUrl` | string | Base URL |

### record_stop

Stop recording.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `stepCount` | integer | Steps recorded |

### record_export

Export recorded script.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `script` | string | JSON script |
| `stepCount` | integer | Steps |
| `format` | string | Output format |

### record_get_status

Check recording state.

### record_clear

Clear recorded steps.

## Tracing

Tools for recording detailed execution traces for debugging and analysis. Traces include screenshots, DOM snapshots, and network activity. View traces with `npx playwright show-trace <trace.zip>`.

### trace_start

Start trace recording with screenshots and DOM snapshots.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | | Trace name (used for file naming) |
| `title` | string | | Title shown in trace viewer |
| `screenshots` | boolean | | Include screenshots (default: true) |
| `snapshots` | boolean | | Include DOM snapshots (default: true) |
| `sources` | boolean | | Include source files |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `name` | string | Trace name |

### trace_stop

Stop trace recording and save or return the trace data.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `path` | string | | File path to save trace ZIP (returns base64 if omitted) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `path` | string | File path if saved |
| `data` | string | Base64-encoded ZIP if no path specified |
| `size_kb` | integer | Trace size in KB |
| `view_hint` | string | Command to view trace |

### trace_chunk_start

Start a new trace chunk within an active trace for segmenting recordings.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | | Chunk name |
| `title` | string | | Chunk title shown in trace viewer |

### trace_chunk_stop

Stop the current trace chunk.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `path` | string | | File path to save chunk ZIP |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `path` | string | File path if saved |
| `data` | string | Base64-encoded ZIP if no path |
| `size_kb` | integer | Chunk size in KB |

### trace_group_start

Start a trace group for logical grouping of actions in the trace viewer.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Group name |
| `location` | string | | Source location to associate |

### trace_group_stop

Stop the current trace group.

## Init Scripts

### js_init_script

Add JavaScript that runs before page scripts on every navigation. Useful for mocking APIs, injecting test helpers, or setting up authentication.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `script` | string | ✅ | JavaScript code to inject before page scripts |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

**Example Use Cases:**

```javascript
// Mock fetch API
window.fetch = async (url) => {
  if (url.includes('/api/user')) {
    return { json: () => ({ id: 1, name: 'Test User' }) };
  }
  return originalFetch(url);
};

// Disable analytics
window.gtag = () => {};
window.analytics = { track: () => {}, identify: () => {} };

// Inject test helpers
window.testHelpers = {
  fillForm: (data) => { /* ... */ },
  waitForElement: (sel) => { /* ... */ }
};
```

## Assertions

### test_assert_text

Assert text exists.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `text` | string | ✅ | Expected text |
| `selector` | string | | Limit to element |

### test_assert_element

Assert element exists.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

## Testing Tools

### test_verify_value

Verify that an input element has the expected value.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector for the input element |
| `expected` | string | ✅ | Expected value to verify |
| `timeout_ms` | integer | | Timeout in milliseconds (default: 5000) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `passed` | boolean | Whether the verification passed |
| `actual` | string | The actual value found |
| `message` | string | Status message |

### test_verify_list

Verify that a list of text items are all visible on the page.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `items` | array | ✅ | List of text items that should be visible |
| `selector` | string | | Optional CSS selector to scope the search |
| `timeout_ms` | integer | | Timeout in milliseconds (default: 5000) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `passed` | boolean | Whether all items were found |
| `found` | array | Items that were found |
| `missing` | array | Items that were not found |
| `message` | string | Status message |

### test_generate_locator

Generate a locator string for a given element using a specific strategy.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector for the element |
| `strategy` | string | | Locator strategy: css, xpath, testid, role, text (default: css) |
| `timeout_ms` | integer | | Timeout in milliseconds (default: 5000) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `locator` | string | Generated locator string |
| `strategy` | string | Strategy used |
| `metadata` | object | Additional metadata (role, label, text, etc.) |

## Test Reporting

### test_get_report

Get test execution report.

### test_reset

Clear test results.

### test_set_target

Set test target description.

## State Management

Named state snapshots for saving and restoring browser state across sessions.

### state_save

Save current browser state (cookies, localStorage, sessionStorage) to a named snapshot.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Name for the state snapshot |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `path` | string | Path to saved state file |

### state_load

Load a previously saved state snapshot.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Name of the state snapshot to load |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

### state_list

List all saved state snapshots.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `states` | array | List of state names |
| `count` | integer | Number of saved states |

### state_delete

Delete a saved state snapshot.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Name of the state snapshot to delete |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

## Workflow Automation

High-level workflow tools for common automation patterns.

### workflow_login

Automated login workflow with credential input and success verification.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `url` | string | ✅ | Login page URL |
| `username_selector` | string | ✅ | CSS selector for username field |
| `password_selector` | string | ✅ | CSS selector for password field |
| `submit_selector` | string | ✅ | CSS selector for submit button |
| `username` | string | ✅ | Username to enter |
| `password` | string | ✅ | Password to enter |
| `success_url` | string | | Expected URL after successful login |
| `success_selector` | string | | Element that indicates successful login |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Whether login succeeded |
| `url` | string | Final URL after login |
| `message` | string | Status message |

### workflow_extract_table

Extract data from an HTML table to structured JSON.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector for the table |
| `headers` | boolean | | Whether first row contains headers (default: true) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `rows` | array | Array of row objects |
| `columns` | array | Column names |
| `count` | integer | Number of rows |

## Configuration

### config_get

Get the resolved MCP server configuration.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `headless` | boolean | Whether browser runs headless |
| `project` | string | Project name for reports |
| `default_timeout_ms` | integer | Default timeout in milliseconds |
| `browser_launched` | boolean | Whether browser has been launched |

---

## CDP Tools (Chrome DevTools Protocol)

These tools use direct CDP connection for advanced browser profiling and debugging that isn't available through WebDriver BiDi.

### Performance & Profiling

#### cdp_get_performance_metrics

Get Core Web Vitals (LCP, CLS, INP) and navigation timing metrics.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `lcp` | object | Largest Contentful Paint data |
| `cls` | number | Cumulative Layout Shift score |
| `inp` | object | Interaction to Next Paint data |
| `navigation` | object | Navigation timing metrics |

#### cdp_get_memory_stats

Get JavaScript heap memory statistics.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `usedJSHeapSize` | integer | Used heap size in bytes |
| `totalJSHeapSize` | integer | Total heap size in bytes |
| `jsHeapSizeLimit` | integer | Maximum heap size in bytes |

#### cdp_take_heap_snapshot

Capture a V8 heap snapshot for memory profiling.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `path` | string | | File path to save snapshot (default: auto-generated) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `path` | string | Path to saved .heapsnapshot file |
| `size` | integer | Snapshot size in bytes |

### Network Emulation (CDP)

#### cdp_emulate_network

Emulate network conditions.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `preset` | string | | Preset: "slow3g", "fast3g", "4g", "wifi", "offline" |
| `latency` | number | | Custom latency in ms |
| `download` | number | | Custom download throughput in bytes/sec |
| `upload` | number | | Custom upload throughput in bytes/sec |

#### cdp_emulate_cpu

Emulate CPU throttling.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `rate` | number | ✅ | Throttling rate (2, 4, or 6 for 2x, 4x, 6x slowdown) |

#### cdp_clear_network_emulation

Clear network emulation and restore normal conditions.

#### cdp_clear_cpu_emulation

Clear CPU emulation and restore normal speed.

### Quality Auditing

#### cdp_run_lighthouse

Run Lighthouse audit for performance, accessibility, SEO, and best practices.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `categories` | array | | Categories: "performance", "accessibility", "seo", "best-practices" |
| `device` | string | | Device: "desktop" or "mobile" (default: desktop) |
| `output_dir` | string | | Directory for HTML/JSON reports |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Audited URL |
| `scores` | object | Scores by category (0-100) |
| `passed_audits` | integer | Number of passing audits |
| `failed_audits` | integer | Number of failing audits |
| `report_paths` | object | Paths to generated reports |

### Code Coverage (CDP)

#### cdp_start_coverage

Start collecting JavaScript and CSS code coverage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `js` | boolean | | Include JavaScript coverage (default: true) |
| `css` | boolean | | Include CSS coverage (default: true) |
| `call_count` | boolean | | Include call counts for JS functions |
| `detailed` | boolean | | Include detailed block coverage |

#### cdp_stop_coverage

Stop coverage collection and return the coverage report.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `js` | object | JavaScript coverage summary |
| `css` | object | CSS coverage summary |
| `scripts` | array | Detailed script coverage |

### Console Debugging (CDP)

#### cdp_enable_console_debugger

Enable enhanced console debugging with stack traces.

#### cdp_get_console_entries

Get console messages with full stack traces.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `entries` | array | Console entries with type, args, stack trace |
| `count` | integer | Number of entries |

#### cdp_get_browser_logs

Get browser logs including deprecations, interventions, and violations.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `logs` | array | Log entries with source, level, text |
| `count` | integer | Number of logs |

#### cdp_disable_console_debugger

Disable enhanced console debugging.

### Screencast (CDP)

#### cdp_start_screencast

Start streaming screen frames.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `format` | string | | Image format: "jpeg" or "png" (default: jpeg) |
| `quality` | integer | | Image quality 0-100 (default: 80) |
| `max_width` | integer | | Maximum width in pixels |
| `max_height` | integer | | Maximum height in pixels |

#### cdp_stop_screencast

Stop screen streaming.

### Extensions Management (CDP)

#### cdp_install_extension

Install a Chrome extension from a local path.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `path` | string | ✅ | Path to unpacked extension directory |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Installed extension ID |

#### cdp_uninstall_extension

Uninstall a Chrome extension by ID.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | ✅ | Extension ID to uninstall |

#### cdp_list_extensions

List all installed Chrome extensions.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `extensions` | array | Extension info (id, name, enabled) |
| `count` | integer | Number of extensions |

### Network Request Bodies (CDP)

#### cdp_get_response_body

Get the response body for a specific network request by ID.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `request_id` | string | ✅ | Network request ID from network_get_requests |
| `save_to_file` | string | | Optional file path to save binary content |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `body` | string | Response body content |
| `base64_encoded` | boolean | Whether body is base64 encoded |
| `path` | string | File path if saved to file |
