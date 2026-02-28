# MCP Tools Reference

Complete reference for all 75+ MCP tools.

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

### navigate

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

### back

Navigate back in history.

### forward

Navigate forward in history.

### reload

Reload the current page.

## Element Interactions

### click

Click an element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `timeout_ms` | integer | | Timeout (default: 5000) |

### dblclick

Double-click an element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `timeout_ms` | integer | | Timeout |

### type

Type text into an element (appends).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `text` | string | ✅ | Text to type |
| `timeout_ms` | integer | | Timeout |

### fill

Fill an input (replaces content).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `value` | string | ✅ | Value to fill |
| `timeout_ms` | integer | | Timeout |

### clear

Clear an input element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

### press

Press a key on an element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `key` | string | ✅ | Key (e.g., "Enter") |

## Form Controls

### check

Check a checkbox.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

### uncheck

Uncheck a checkbox.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

### select_option

Select dropdown option(s).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `values` | array | | Option values |
| `labels` | array | | Option labels |
| `indexes` | array | | Option indexes |

### set_files

Set files on a file input.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `files` | array | ✅ | File paths |

## Element State

### get_text

Get element text content.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `text` | string | Text content |

### get_value

Get input element value.

### get_inner_html

Get element innerHTML.

### get_inner_text

Get element innerText.

### get_attribute

Get element attribute.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `name` | string | ✅ | Attribute name |

### get_bounding_box

Get element bounding box.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `x` | number | X position |
| `y` | number | Y position |
| `width` | number | Width |
| `height` | number | Height |

### is_visible

Check if element is visible.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `visible` | boolean | Visibility state |

### is_hidden

Check if element is hidden.

### is_enabled

Check if element is enabled.

### is_checked

Check if checkbox/radio is checked.

### is_editable

Check if element is editable.

### get_role

Get ARIA role.

### get_label

Get accessible label.

## Page State

### get_title

Get page title.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Page title |

### get_url

Get current URL.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Current URL |

### get_content

Get page HTML content.

### set_content

Set page HTML content.

### get_viewport

Get viewport dimensions.

### set_viewport

Set viewport dimensions.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `width` | integer | ✅ | Width |
| `height` | integer | ✅ | Height |

## Screenshots & PDF

### screenshot

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

### pdf

Generate PDF.

## JavaScript

### evaluate

Execute JavaScript.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `script` | string | ✅ | JavaScript code |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `result` | any | Evaluation result |

### element_eval

Evaluate JavaScript with element.

### add_script

Inject JavaScript.

### add_style

Inject CSS.

## Waiting

### wait_until

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

### keyboard_press

Press a key.

### keyboard_down

Hold a key.

### keyboard_up

Release a key.

### keyboard_type

Type text via keyboard.

### mouse_click

Click at coordinates.

### mouse_move

Move mouse.

### mouse_down

Press mouse button.

### mouse_up

Release mouse button.

### mouse_wheel

Scroll mouse wheel.

### touch_tap

Tap at coordinates.

### touch_swipe

Swipe gesture.

## Page Management

### new_page

Create new page/tab.

### get_pages

Get page count.

### close_page

Close current page.

### bring_to_front

Activate page.

## Cookies & Storage

### get_cookies

Get browser cookies.

### set_cookies

Set browser cookies.

### clear_cookies

Clear all cookies.

### get_storage_state

Get cookies and localStorage.

## Script Recording

### start_recording

Begin recording actions.

**Input:**

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Script name |
| `description` | string | Description |
| `baseUrl` | string | Base URL |

### stop_recording

Stop recording.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `stepCount` | integer | Steps recorded |

### export_script

Export recorded script.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `script` | string | JSON script |
| `stepCount` | integer | Steps |
| `format` | string | Output format |

### recording_status

Check recording state.

### clear_recording

Clear recorded steps.

## Assertions

### assert_text

Assert text exists.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `text` | string | ✅ | Expected text |
| `selector` | string | | Limit to element |

### assert_element

Assert element exists.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

## Test Reporting

### get_test_report

Get test execution report.

### reset_session

Clear test results.

### set_target

Set test target description.
