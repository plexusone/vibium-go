# Feature Comparison

This document compares webpilot against WebPilotDev clients (Java/JS/Python), the WebPilotDev MCP server, and Playwright MCP.

## Overview

| Project | Language | Type | Repository |
|---------|----------|------|------------|
| **webpilot** | Go | SDK + MCP Server | [plexusone/webpilot](https://github.com/plexusone/webpilot) |
| **WebPilotDev** | Go + JS/Python/Java | Daemon + MCP + Clients | [WebPilotDev/vibium](https://github.com/WebPilotDev/vibium) |
| **Playwright MCP** | TypeScript | MCP Server | [microsoft/playwright-mcp](https://github.com/microsoft/playwright-mcp) |

---

## Part 1: MCP Server Comparison

Compare the three MCP servers for LLM agent integration.

### Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         webpilot MCP                               │
├─────────────────────────────────────────────────────────────────────┤
│  webpilot-mcp ──► webpilot SDK ──► BiDi Client ──► Chrome            │
│       │                                                             │
│       └── Uses official Go MCP SDK                                  │
│       └── 100+ tools, comprehensive automation                      │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                        WebPilotDev MCP                                │
├─────────────────────────────────────────────────────────────────────┤
│  webpilot mcp ──► clicker ──► BiDi Client ──► Chrome                  │
│       │                                                             │
│       └── Hand-rolled JSON-RPC                                      │
│       └── ~25 tools, core automation                                │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                       Playwright MCP                                │
├─────────────────────────────────────────────────────────────────────┤
│  @playwright/mcp ──► Playwright ──► CDP/BiDi ──► Chromium           │
│       │                                                             │
│       └── Official TS MCP SDK                                       │
│       └── ~45 tools (with opt-in caps)                              │
│       └── Snapshot-based (accessibility tree)                       │
└─────────────────────────────────────────────────────────────────────┘
```

### Implementation Details

| Aspect | webpilot | WebPilotDev | Playwright MCP |
|--------|-----------|-----------|----------------|
| Language | Go | Go | TypeScript |
| MCP SDK | Official Go SDK | Hand-rolled | Official TS SDK |
| Protocol | WebDriver BiDi | WebDriver BiDi | CDP + BiDi |
| Tool prefix | None | `browser_` | `browser_` |
| Tool count | **100+** | ~25 | ~45 (with caps) |
| Browser | Chrome | Chrome | Chromium/Firefox/WebKit |

### Tool Comparison Matrix

Legend: :white_check_mark: = Supported, :x: = Not supported, (opt) = Requires opt-in flag

#### Browser Management

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Launch browser | `browser_launch` | `browser_start` | (auto-launch) |
| Quit browser | `browser_quit` | `browser_stop` | `browser_close` |
| Resize viewport | `set_viewport` | :x: | `browser_resize` |
| Get viewport | `get_viewport` | :x: | :x: |

#### Navigation

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Navigate | `navigate` | `browser_navigate` | `browser_navigate` |
| Back | `back` | :x: | `browser_navigate_back` |
| Forward | `forward` | :x: | :x: |
| Reload | `reload` | :x: | :x: |
| Scroll | `scroll` | `browser_scroll` | :x: |

#### Element Interaction

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Click | `click` | `browser_click` | `browser_click` |
| Double-click | `dblclick` | :x: | `browser_click` (doubleClick) |
| Type | `type` | `browser_type` | `browser_type` |
| Fill | `fill` | :x: | `browser_type` |
| Clear | `clear` | :x: | :x: |
| Press key | `press` | `browser_keys` | `browser_press_key` |
| Check/Uncheck | `check`, `uncheck` | :x: | :x: |
| Select option | `select_option` | `browser_select` | `browser_select_option` |
| Hover | `hover` | `browser_hover` | `browser_hover` |
| Focus | `focus` | :x: | :x: |
| Drag | `drag_to` | :x: | `browser_drag` |
| Tap (touch) | `tap` | :x: | :x: |
| Set files | `set_files` | :x: | `browser_file_upload` |
| Fill form | `fill_form` | :x: | `browser_fill_form` |
| Dispatch event | `dispatch_event` | :x: | :x: |

#### Element State

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get text | `get_text` | `browser_get_text` | `browser_snapshot` |
| Get value | `get_value` | :x: | :x: |
| Get innerHTML | `get_inner_html` | `browser_get_html` | :x: |
| Get outerHTML | `get_outer_html` | :x: | :x: |
| Get innerText | `get_inner_text` | :x: | :x: |
| Get attribute | `get_attribute` | :x: | :x: |
| Get bounding box | `get_bounding_box` | :x: | :x: |
| Is visible | `is_visible` | :x: | :x: |
| Is hidden | `is_hidden` | :x: | :x: |
| Is enabled | `is_enabled` | :x: | :x: |
| Is checked | `is_checked` | :x: | :x: |
| Is editable | `is_editable` | :x: | :x: |
| Get role | `get_role` | :x: | :x: |
| Get label | `get_label` | :x: | :x: |

#### Page State

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get title | `get_title` | `browser_get_title` | :x: |
| Get URL | `get_url` | `browser_get_url` | :x: |
| Get content | `get_content` | :x: | :x: |
| Set content | `set_content` | :x: | :x: |
| Accessibility snapshot | `accessibility_snapshot` | :x: | `browser_snapshot` |

#### Screenshots & PDF

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Screenshot | `screenshot` | `browser_screenshot` | `browser_take_screenshot` |
| Full page screenshot | `screenshot` (fullPage) | :x: | `browser_take_screenshot` (fullPage) |
| Element screenshot | `element_screenshot` | :x: | `browser_take_screenshot` (ref) |
| PDF | `pdf` | :x: | :x: |

#### JavaScript

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Evaluate | `evaluate` | `browser_evaluate` | `browser_evaluate` |
| Element eval | `element_eval` | :x: | `browser_evaluate` (ref) |
| Add script | `add_script` | :x: | :x: |
| Add style | `add_style` | :x: | :x: |
| Run Playwright code | :x: | :x: | `browser_run_code` |

#### Waiting

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Wait for element state | `wait_until` | `browser_wait` | `browser_wait_for` |
| Wait for selector | `wait_for_selector` | :x: | :x: |
| Wait for URL | `wait_for_url` | :x: | :x: |
| Wait for load | `wait_for_load` | :x: | :x: |
| Wait for function | `wait_for_function` | :x: | :x: |
| Wait for text | `wait_for_text` | :x: | `browser_wait_for` (text) |
| Wait for time | :x: | :x: | `browser_wait_for` (time) |

#### Input Controllers

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Keyboard press | `keyboard_press` | `browser_keys` | `browser_press_key` |
| Keyboard down/up | `keyboard_down`, `keyboard_up` | :x: | :x: |
| Keyboard type | `keyboard_type` | :x: | :x: |
| Mouse click | `mouse_click` | :x: | :x: |
| Mouse move | `mouse_move` | :x: | :x: |
| Mouse down/up | `mouse_down`, `mouse_up` | :x: | :x: |
| Mouse wheel | `mouse_wheel` | :x: | :x: |
| Mouse drag | `mouse_drag` | :x: | `browser_drag` |
| Touch tap | `touch_tap` | :x: | :x: |
| Touch swipe | `touch_swipe` | :x: | :x: |

#### Tab/Page Management

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| New page | `new_page` | `browser_new_page` | `browser_tabs` (new) |
| List pages | `get_pages`, `list_tabs` | `browser_list_pages` | `browser_tabs` (list) |
| Switch page | `select_tab` | `browser_switch_page` | `browser_tabs` (select) |
| Close page | `close_page`, `close_tab` | `browser_close_page` | `browser_tabs` (close) |
| Bring to front | `bring_to_front` | :x: | :x: |

#### Frame Management

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get frames | `get_frames` | :x: | :x: |
| Select frame | `select_frame` | :x: | :x: |
| Select main frame | `select_main_frame` | :x: | :x: |

#### Cookie Management

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get cookies | `get_cookies` | :x: | `browser_cookie_list` (opt) |
| Get cookie | :x: | :x: | `browser_cookie_get` (opt) |
| Set cookies | `set_cookies` | :x: | `browser_cookie_set` (opt) |
| Clear cookies | `clear_cookies` | :x: | `browser_cookie_clear` (opt) |
| Delete cookie | `delete_cookie` | :x: | `browser_cookie_delete` (opt) |

#### LocalStorage

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get item | `localstorage_get` | :x: | `browser_localstorage_get` (opt) |
| Set item | `localstorage_set` | :x: | `browser_localstorage_set` (opt) |
| List items | `localstorage_list` | :x: | `browser_localstorage_list` (opt) |
| Delete item | `localstorage_delete` | :x: | `browser_localstorage_delete` (opt) |
| Clear | `localstorage_clear` | :x: | `browser_localstorage_clear` (opt) |

#### SessionStorage

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get item | `sessionstorage_get` | :x: | `browser_sessionstorage_get` (opt) |
| Set item | `sessionstorage_set` | :x: | `browser_sessionstorage_set` (opt) |
| List items | `sessionstorage_list` | :x: | `browser_sessionstorage_list` (opt) |
| Delete item | `sessionstorage_delete` | :x: | `browser_sessionstorage_delete` (opt) |
| Clear | `sessionstorage_clear` | :x: | `browser_sessionstorage_clear` (opt) |

#### Storage State

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get storage state | `get_storage_state` | :x: | `browser_storage_state` (opt) |
| Set storage state | `set_storage_state` | :x: | `browser_set_storage_state` (opt) |
| Clear storage | `clear_storage` | :x: | :x: |

#### Network

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Route (mock) | `route` | :x: | `browser_route` (opt) |
| Route list | `route_list` | :x: | `browser_route_list` (opt) |
| Unroute | `unroute` | :x: | `browser_unroute` (opt) |
| Network offline | `network_state_set` | :x: | `browser_network_state_set` (opt) |
| Get requests | `get_network_requests` | :x: | `browser_network_requests` |
| Clear requests | `clear_network_requests` | :x: | :x: |

#### Console

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get console messages | `get_console_messages` | :x: | `browser_console_messages` |
| Clear console | `clear_console_messages` | :x: | :x: |

#### Dialogs

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Handle dialog | `handle_dialog` | :x: | `browser_handle_dialog` |
| Get dialog info | `get_dialog` | :x: | :x: |

#### Emulation

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Emulate media | `emulate_media` | :x: | :x: |
| Set geolocation | `set_geolocation` | :x: | :x: |

#### Recording & Tracing

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Start recording | `start_recording` | :x: | :x: |
| Stop recording | `stop_recording` | :x: | :x: |
| Export script | `export_script` | :x: | :x: |
| Start trace | `start_trace` | :x: | `browser_start_tracing` (opt) |
| Stop trace | `stop_trace` | :x: | `browser_stop_tracing` (opt) |
| Trace chunks | `start_trace_chunk`, `stop_trace_chunk` | :x: | :x: |
| Trace groups | `start_trace_group`, `stop_trace_group` | :x: | :x: |
| Start video | `start_video` | :x: | `browser_start_video` (opt) |
| Stop video | `stop_video` | :x: | `browser_stop_video` (opt) |

#### Testing & Assertions

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Assert text | `assert_text` | :x: | :x: |
| Assert element | `assert_element` | :x: | :x: |
| Verify value | `verify_value` | :x: | :x: |
| Verify text | `verify_text` | :x: | :x: |
| Verify visible | `verify_visible` | :x: | :x: |
| Verify hidden | `verify_hidden` | :x: | :x: |
| Verify enabled | `verify_enabled` | :x: | :x: |
| Verify disabled | `verify_disabled` | :x: | :x: |
| Verify checked | `verify_checked` | :x: | :x: |
| Verify list visible | `verify_list_visible` | :x: | :x: |
| Generate locator | `generate_locator` | :x: | :x: |
| Get test report | `get_test_report` | :x: | :x: |
| Reset session | `reset_session` | :x: | :x: |
| Set target | `set_target` | :x: | :x: |

#### Human-in-the-Loop

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Pause for human | `pause_for_human` | :x: | :x: |

#### Configuration

| Tool | webpilot | WebPilotDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get config | `get_config` | :x: | `browser_get_config` (opt) |
| Add init script | `add_init_script` | :x: | :x: |

### Tool Count Summary

| MCP Server | Core Tools | Opt-in Tools | Total |
|------------|:----------:|:------------:|:-----:|
| **webpilot** | 100+ | - | **100+** |
| Playwright MCP | ~20 | ~25 | ~45 |
| WebPilotDev | ~25 | - | ~25 |

### Unique Features by Server

#### webpilot Only

| Feature | Description |
|---------|-------------|
| Script Runner | YAML/JSON test execution via `webpilot run` |
| Session Recording | Capture actions as replayable scripts |
| Human-in-the-Loop | `pause_for_human` for SSO, CAPTCHA, 2FA |
| Test Reports | Structured reports (box, diagnostic, JSON) |
| Verification Tools | `verify_*` tools with detailed output |
| Frame Selection | `select_frame`/`select_main_frame` |
| Trace Chunks/Groups | Fine-grained trace control |

#### Playwright MCP Only

| Feature | Description |
|---------|-------------|
| `browser_run_code` | Execute arbitrary Playwright code |
| `browser_snapshot` | Accessibility tree snapshots |
| Multi-browser | Chromium, Firefox, WebKit |
| Opt-in capabilities | `--caps=network,storage,devtools` |

#### WebPilotDev Only

| Feature | Description |
|---------|-------------|
| Daemon mode | HTTP API for multi-client scenarios |
| Multi-language clients | Official JS, Python, Java SDKs |

---

## Part 2: Client Library Comparison

Compare webpilot SDK with WebPilotDev client libraries.

### Language & Integration

| Aspect | webpilot | vibium-js | vibium-py | vibium-java |
|--------|-----------|-----------|-----------|-------------|
| Language | Go | JavaScript/TS | Python | Java |
| Async model | Context-based | Promises | async/await | CompletableFuture |
| Error handling | Error returns | try/catch | try/except | Exceptions |
| Package manager | go modules | npm | pip | Maven/Gradle |
| Type safety | Strong | TypeScript | Type hints | Strong |

### Core Features

| Feature | webpilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Launch browser | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Headless mode | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Connect remote | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Multiple contexts | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Element Finding

| Feature | webpilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| CSS selector | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| **Semantic selectors** | | | | |
| - By role | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By text | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By label | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By placeholder | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By alt text | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By title | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By testid | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By xpath | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By proximity | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| FindAll | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Scoped find | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Interactions

| Feature | webpilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Click/Type/Fill | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Check/Uncheck | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Select option | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Drag and drop | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| File upload | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| dispatchEvent | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Highlight element | :white_check_mark: | :x: | :x: | :white_check_mark: |

### Input Controllers

| Feature | webpilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Keyboard | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Mouse | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Touch | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Event Listeners

| Feature | webpilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| onConsole | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onError | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onRequest/onResponse | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onDialog | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onDownload | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onPage/onPopup | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onWebSocket | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Recording/Tracing

| Feature | webpilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Trace recording | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Trace chunks | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Trace groups | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Video recording | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Storage

| Feature | webpilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Cookies | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| LocalStorage | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| SessionStorage | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Full storage state | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Init scripts | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Delete single cookie | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Network

| Feature | webpilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Route/mock | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Offline mode | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Extra headers | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Clock Control

| Feature | webpilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Install/fastForward | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| pauseAt/resume | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| setFixedTime | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| setTimezone | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Accessibility

| Feature | webpilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| a11yTree | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| interestingOnly option | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| root element option | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

---

## When to Use Which

### MCP Server Selection

| Use Case | Recommendation |
|----------|----------------|
| **Comprehensive LLM automation** | webpilot MCP (100+ tools) |
| **Simple browser control** | WebPilotDev MCP or Playwright MCP |
| **Accessibility-focused** | Playwright MCP (`browser_snapshot`) |
| **Human-in-the-loop flows** | webpilot MCP (`pause_for_human`) |
| **Test automation with reports** | webpilot MCP (verification + reports) |
| **Multi-browser support** | Playwright MCP |

### Client Library Selection

| Use Case | Recommendation |
|----------|----------------|
| **Go application** | webpilot SDK |
| **JavaScript/TypeScript app** | vibium-js |
| **Python application** | vibium-py |
| **Java/JVM application** | vibium-java |
| **Script-based automation** | webpilot (`webpilot run`) |

---

## Relationship Between Projects

```
┌─────────────────────────────────────────────────────────────────────┐
│                      WebPilotDev/vibium                               │
│                                                                     │
│   ┌─────────┐  ┌─────────┐  ┌─────────┐                            │
│   │vibium-js│  │vibium-py│  │vibium-  │                            │
│   │ Client  │  │ Client  │  │java Cli │                            │
│   └────┬────┘  └────┬────┘  └────┬────┘                            │
│        └────────────┼───────────┘                                  │
│                     │ HTTP API                                      │
│                     ▼                                               │
│               ┌──────────┐                                          │
│               │ clicker  │◄─── webpilot mcp                          │
│               │ (binary) │                                          │
│               └────┬─────┘                                          │
│                    │ BiDi                                           │
└────────────────────┼────────────────────────────────────────────────┘
                     │
                     │  webpilot uses clicker
                     │  ONLY for launching Chrome
                     ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    plexusone/webpilot                              │
│                                                                     │
│   ┌────────────────────────────────────────────────────────────┐   │
│   │                    webpilot SDK                            │   │
│   │              (own BiDi client implementation)               │   │
│   └────────────────────────┬───────────────────────────────────┘   │
│                            │                                        │
│         ┌──────────────────┼──────────────────┐                    │
│         │                  │                  │                    │
│         ▼                  ▼                  ▼                    │
│    webpilot-mcp         webpilot run        Direct SDK                 │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Key point**: webpilot is an independent implementation that uses WebPilotDev's `clicker` binary only as a Chrome launcher. The BiDi communication, MCP server, and all SDK functionality are implemented independently in Go.

---

## Legend

| Symbol | Meaning |
|--------|---------|
| :white_check_mark: | Supported |
| :x: | Not supported |
| (opt) | Requires opt-in flag |
