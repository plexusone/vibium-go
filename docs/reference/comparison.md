# Feature Comparison

This document compares w3pilot against VibiumDev clients (Java/JS/Python), the VibiumDev MCP server, Playwright MCP, and ChromeDevTools MCP.

## Overview

| Project | Language | Type | Repository |
|---------|----------|------|------------|
| **w3pilot** | Go | SDK + MCP Server | [plexusone/w3pilot](https://github.com/plexusone/w3pilot) |
| **VibiumDev** | Go + JS/Python/Java | Daemon + MCP + Clients | [VibiumDev/vibium](https://github.com/VibiumDev/vibium) |
| **Playwright MCP** | TypeScript | MCP Server | [microsoft/playwright-mcp](https://github.com/microsoft/playwright-mcp) |
| **ChromeDevTools MCP** | TypeScript | MCP Server | [ChromeDevTools/chrome-devtools-mcp](https://github.com/ChromeDevTools/chrome-devtools-mcp) |

---

## Part 1: MCP Server Comparison

Compare the three MCP servers for LLM agent integration.

### Architecture

```
┌───────────────────────────────────────────────────────────────────┐
│                             w3pilot MCP                           │
├───────────────────────────────────────────────────────────────────┤
│  w3pilot-mcp ──► w3pilot SDK ──► BiDi Client ──► Chrome           │
│       │                                                           │
│       └── Uses official Go MCP SDK                                │
│       └── 100+ tools, comprehensive automation                    │
└───────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────┐
│                           VibiumDev MCP                           │
├───────────────────────────────────────────────────────────────────┤
│  vibium mcp ──► clicker ──► BiDi Client ──► Chrome                │
│       │                                                           │
│       └── Hand-rolled JSON-RPC                                    │
│       └── ~25 tools, core automation                              │
└───────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────┐
│                          Playwright MCP                           │
├───────────────────────────────────────────────────────────────────┤
│  @playwright/mcp ──► Playwright ──► CDP/BiDi ──► Chromium         │
│       │                                                           │
│       └── Official TS MCP SDK                                     │
│       └── ~45 tools (with opt-in caps)                            │
│       └── Snapshot-based (accessibility tree)                     │
└───────────────────────────────────────────────────────────────────┘
```

### Implementation Details

| Aspect | w3pilot | VibiumDev | Playwright MCP |
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

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Launch browser | `browser_launch` | `browser_start` | (auto-launch) |
| Quit browser | `browser_quit` | `browser_stop` | `browser_close` |
| Resize viewport | `set_viewport` | :x: | `browser_resize` |
| Get viewport | `get_viewport` | :x: | :x: |

#### Navigation

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Navigate | `navigate` | `browser_navigate` | `browser_navigate` |
| Back | `back` | :x: | `browser_navigate_back` |
| Forward | `forward` | :x: | :x: |
| Reload | `reload` | :x: | :x: |
| Scroll | `scroll` | `browser_scroll` | :x: |

#### Element Interaction

| Tool | w3pilot | VibiumDev | Playwright MCP |
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

| Tool | w3pilot | VibiumDev | Playwright MCP |
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

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get title | `get_title` | `browser_get_title` | :x: |
| Get URL | `get_url` | `browser_get_url` | :x: |
| Get content | `get_content` | :x: | :x: |
| Set content | `set_content` | :x: | :x: |
| Accessibility snapshot | `accessibility_snapshot` | :x: | `browser_snapshot` |

#### Screenshots & PDF

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Screenshot | `screenshot` | `browser_screenshot` | `browser_take_screenshot` |
| Full page screenshot | `screenshot` (fullPage) | :x: | `browser_take_screenshot` (fullPage) |
| Element screenshot | `element_screenshot` | :x: | `browser_take_screenshot` (ref) |
| PDF | `pdf` | :x: | :x: |

#### JavaScript

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Evaluate | `evaluate` | `browser_evaluate` | `browser_evaluate` |
| Element eval | `element_eval` | :x: | `browser_evaluate` (ref) |
| Add script | `add_script` | :x: | :x: |
| Add style | `add_style` | :x: | :x: |
| Run Playwright code | :x: | :x: | `browser_run_code` |

#### Waiting

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Wait for element state | `wait_until` | `browser_wait` | `browser_wait_for` |
| Wait for selector | `wait_for_selector` | :x: | :x: |
| Wait for URL | `wait_for_url` | :x: | :x: |
| Wait for load | `wait_for_load` | :x: | :x: |
| Wait for function | `wait_for_function` | :x: | :x: |
| Wait for text | `wait_for_text` | :x: | `browser_wait_for` (text) |
| Wait for time | :x: | :x: | `browser_wait_for` (time) |

#### Input Controllers

| Tool | w3pilot | VibiumDev | Playwright MCP |
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

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| New page | `new_page` | `browser_new_page` | `browser_tabs` (new) |
| List pages | `get_pages`, `list_tabs` | `browser_list_pages` | `browser_tabs` (list) |
| Switch page | `select_tab` | `browser_switch_page` | `browser_tabs` (select) |
| Close page | `close_page`, `close_tab` | `browser_close_page` | `browser_tabs` (close) |
| Bring to front | `bring_to_front` | :x: | :x: |

#### Frame Management

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get frames | `get_frames` | :x: | :x: |
| Select frame | `select_frame` | :x: | :x: |
| Select main frame | `select_main_frame` | :x: | :x: |

#### Cookie Management

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get cookies | `get_cookies` | :x: | `browser_cookie_list` (opt) |
| Get cookie | :x: | :x: | `browser_cookie_get` (opt) |
| Set cookies | `set_cookies` | :x: | `browser_cookie_set` (opt) |
| Clear cookies | `clear_cookies` | :x: | `browser_cookie_clear` (opt) |
| Delete cookie | `delete_cookie` | :x: | `browser_cookie_delete` (opt) |

#### LocalStorage

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get item | `localstorage_get` | :x: | `browser_localstorage_get` (opt) |
| Set item | `localstorage_set` | :x: | `browser_localstorage_set` (opt) |
| List items | `localstorage_list` | :x: | `browser_localstorage_list` (opt) |
| Delete item | `localstorage_delete` | :x: | `browser_localstorage_delete` (opt) |
| Clear | `localstorage_clear` | :x: | `browser_localstorage_clear` (opt) |

#### SessionStorage

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get item | `sessionstorage_get` | :x: | `browser_sessionstorage_get` (opt) |
| Set item | `sessionstorage_set` | :x: | `browser_sessionstorage_set` (opt) |
| List items | `sessionstorage_list` | :x: | `browser_sessionstorage_list` (opt) |
| Delete item | `sessionstorage_delete` | :x: | `browser_sessionstorage_delete` (opt) |
| Clear | `sessionstorage_clear` | :x: | `browser_sessionstorage_clear` (opt) |

#### Storage State

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get storage state | `get_storage_state` | :x: | `browser_storage_state` (opt) |
| Set storage state | `set_storage_state` | :x: | `browser_set_storage_state` (opt) |
| Clear storage | `clear_storage` | :x: | :x: |

#### Network

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Route (mock) | `route` | :x: | `browser_route` (opt) |
| Route list | `route_list` | :x: | `browser_route_list` (opt) |
| Unroute | `unroute` | :x: | `browser_unroute` (opt) |
| Network offline | `network_state_set` | :x: | `browser_network_state_set` (opt) |
| Get requests | `get_network_requests` | :x: | `browser_network_requests` |
| Clear requests | `clear_network_requests` | :x: | :x: |

#### Console

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get console messages | `get_console_messages` | :x: | `browser_console_messages` |
| Clear console | `clear_console_messages` | :x: | :x: |

#### Dialogs

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Handle dialog | `handle_dialog` | :x: | `browser_handle_dialog` |
| Get dialog info | `get_dialog` | :x: | :x: |

#### Emulation

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Emulate media | `emulate_media` | :x: | :x: |
| Set geolocation | `set_geolocation` | :x: | :x: |
| Network throttling | `emulate_network` (CDP) | :x: | :x: |
| CPU throttling | `emulate_cpu` (CDP) | :x: | :x: |

#### Profiling (CDP)

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Heap snapshot | `take_heap_snapshot` (CDP) | :x: | :x: |
| Direct CDP access | `cdp_send` | :x: | :x: |

#### Recording & Tracing

| Tool | w3pilot | VibiumDev | Playwright MCP |
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

| Tool | w3pilot | VibiumDev | Playwright MCP |
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

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Pause for human | `pause_for_human` | :x: | :x: |

#### Configuration

| Tool | w3pilot | VibiumDev | Playwright MCP |
|------|:---------:|:---------:|:--------------:|
| Get config | `get_config` | :x: | `browser_get_config` (opt) |
| Add init script | `add_init_script` | :x: | :x: |

### Tool Count Summary

| MCP Server | Core Tools | Opt-in Tools | Total |
|------------|:----------:|:------------:|:-----:|
| **w3pilot** | 100+ | - | **100+** |
| Playwright MCP | ~20 | ~25 | ~45 |
| VibiumDev | ~25 | - | ~25 |

### Unique Features by Server

#### w3pilot Only

| Feature | Description |
|---------|-------------|
| Script Runner | YAML/JSON test execution via `w3pilot run` |
| Session Recording | Capture actions as replayable scripts |
| Human-in-the-Loop | `pause_for_human` for SSO, CAPTCHA, 2FA |
| Test Reports | Structured reports (box, diagnostic, JSON) |
| Verification Tools | `verify_*` tools with detailed output |
| Frame Selection | `select_frame`/`select_main_frame` |
| Trace Chunks/Groups | Fine-grained trace control |
| **CDP Integration** | Direct Chrome DevTools Protocol access |
| - Heap Snapshots | V8 memory profiling |
| - Network Emulation | Slow 3G, Fast 3G, 4G presets |
| - CPU Throttling | Simulate slower CPUs |

#### Playwright MCP Only

| Feature | Description |
|---------|-------------|
| `browser_run_code` | Execute arbitrary Playwright code |
| `browser_snapshot` | Accessibility tree snapshots |
| Multi-browser | Chromium, Firefox, WebKit |
| Opt-in capabilities | `--caps=network,storage,devtools` |

#### VibiumDev Only

| Feature | Description |
|---------|-------------|
| Daemon mode | HTTP API for multi-client scenarios |
| Multi-language clients | Official JS, Python, Java SDKs |

---

## ChromeDevTools MCP Comparison

[ChromeDevTools MCP](https://github.com/ChromeDevTools/chrome-devtools-mcp) is an official Chrome DevTools MCP server with 29 tools. W3Pilot provides equivalent functionality plus additional features.

### Tool Mapping

| ChromeDevTools MCP (29 tools) | W3Pilot Equivalent | Notes |
|-------------------------------|-------------------|-------|
| **Input Automation** | | |
| `click` | `element_click` | ✅ |
| `drag` | `element_drag_to` | ✅ |
| `fill` | `element_fill` | ✅ |
| `fill_form` | Multiple `element_fill` | Use multiple calls |
| `handle_dialog` | `dialog_handle` | ✅ |
| `hover` | `element_hover` | ✅ |
| `press_key` | `input_keyboard_press` | ✅ |
| `type_text` | `element_type` | ✅ |
| `upload_file` | `element_set_files` | ✅ |
| **Navigation** | | |
| `close_page` | `tab_close` | ✅ |
| `list_pages` | `tab_list` | ✅ |
| `navigate_page` | `page_navigate`, `page_go_back`, `page_go_forward`, `page_reload` | ✅ Split into specific tools |
| `new_page` | `page_new` | ✅ |
| `select_page` | `tab_select` | ✅ |
| `wait_for` | `wait_for_text`, `wait_for_selector` | ✅ Multiple wait tools |
| **Emulation** | | |
| `emulate` | `page_emulate_media`, `cdp_emulate_network`, `cdp_emulate_cpu` | ✅ Split into specific tools |
| `resize_page` | `page_set_viewport` | ✅ |
| **Performance** | | |
| `performance_start_trace` | `trace_start` | ✅ |
| `performance_stop_trace` | `trace_stop` | ✅ |
| `performance_analyze_insight` | `cdp_get_performance_metrics` | Core Web Vitals (LCP, CLS, INP) |
| `take_memory_snapshot` | `cdp_take_heap_snapshot` | ✅ |
| **Network** | | |
| `list_network_requests` | `network_get_requests` | ✅ |
| `get_network_request` | `cdp_get_response_body` | ✅ |
| **Debugging** | | |
| `evaluate_script` | `js_evaluate` | ✅ |
| `list_console_messages` | `console_get_messages` | ✅ |
| `get_console_message` | `console_get_messages` | ✅ Filter by ID |
| `lighthouse_audit` | `cdp_run_lighthouse` | ✅ |
| `take_screenshot` | `page_screenshot` | ✅ |
| `take_snapshot` | `accessibility_snapshot` | ✅ |

### Additional W3Pilot Features (not in ChromeDevTools MCP)

| Category | W3Pilot Tools |
|----------|---------------|
| **Session Recording** | `record_start`, `record_stop`, `record_export`, `record_status`, `record_clear` |
| **Video Capture** | `video_start`, `video_stop` |
| **Storage Management** | `storage_get_cookies`, `storage_set_cookies`, `storage_local_get`, `storage_local_set`, `storage_session_get`, `storage_session_set`, `storage_get_state`, `storage_set_state`, `storage_clear_all` |
| **Network Interception** | `network_route`, `network_unroute`, `network_route_list`, `network_set_offline` |
| **Code Coverage** | `cdp_start_coverage`, `cdp_stop_coverage` |
| **Test Assertions** | `test_assert_text`, `test_assert_element`, `test_verify_value`, `test_verify_visible`, `test_verify_enabled`, `test_verify_checked`, `test_generate_locator`, `test_get_report` |
| **Semantic Selectors** | Find by role, label, placeholder, testid, alt, title, xpath, near |
| **Frame Navigation** | `frame_select`, `frame_select_main`, `frame_list` |
| **Human-in-the-Loop** | `human_pause` |
| **Init Scripts** | `js_init_script`, `js_add_script`, `js_add_style` |
| **Trace Control** | `trace_chunk_start`, `trace_chunk_stop`, `trace_group_start`, `trace_group_stop` |
| **Clock Control** | Time manipulation for testing |
| **Geolocation** | Location spoofing |

### Tool Count Comparison

| MCP Server | Tools |
|------------|:-----:|
| **W3Pilot** | **169** |
| ChromeDevTools MCP | 29 |
| Playwright MCP | ~45 |
| VibiumDev MCP | ~25 |

### When to Use Which

| Use Case | Recommendation |
|----------|----------------|
| Comprehensive automation | W3Pilot (169 tools) |
| Simple debugging tasks | ChromeDevTools MCP |
| Performance tracing only | ChromeDevTools MCP |
| Test automation with assertions | W3Pilot |
| Session recording & export | W3Pilot |
| Storage/cookie management | W3Pilot |

---

## Part 2: Client Library Comparison

Compare w3pilot SDK with VibiumDev client libraries.

### Language & Integration

| Aspect | w3pilot | vibium-js | vibium-py | vibium-java |
|--------|-----------|-----------|-----------|-------------|
| Language | Go | JavaScript/TS | Python | Java |
| Async model | Context-based | Promises | async/await | CompletableFuture |
| Error handling | Error returns | try/catch | try/except | Exceptions |
| Package manager | go modules | npm | pip | Maven/Gradle |
| Type safety | Strong | TypeScript | Type hints | Strong |

### Core Features

| Feature | w3pilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Launch browser | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Headless mode | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Connect remote | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Multiple contexts | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Element Finding

| Feature | w3pilot | vibium-js | vibium-py | vibium-java |
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

| Feature | w3pilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Click/Type/Fill | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Check/Uncheck | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Select option | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Drag and drop | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| File upload | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| dispatchEvent | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Highlight element | :white_check_mark: | :x: | :x: | :white_check_mark: |

### Input Controllers

| Feature | w3pilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Keyboard | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Mouse | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Touch | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Event Listeners

| Feature | w3pilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| onConsole | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onError | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onRequest/onResponse | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onDialog | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onDownload | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onPage/onPopup | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onWebSocket | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Recording/Tracing

| Feature | w3pilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Trace recording | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Trace chunks | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Trace groups | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Video recording | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Storage

| Feature | w3pilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Cookies | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| LocalStorage | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| SessionStorage | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Full storage state | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Init scripts | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Delete single cookie | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Network

| Feature | w3pilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Route/mock | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Offline mode | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Extra headers | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Clock Control

| Feature | w3pilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| Install/fastForward | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| pauseAt/resume | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| setFixedTime | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| setTimezone | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Accessibility

| Feature | w3pilot | vibium-js | vibium-py | vibium-java |
|---------|:---------:|:---------:|:---------:|:-----------:|
| a11yTree | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| interestingOnly option | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| root element option | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

---

## When to Use Which

### MCP Server Selection

| Use Case | Recommendation |
|----------|----------------|
| **Comprehensive LLM automation** | w3pilot MCP (100+ tools) |
| **Simple browser control** | VibiumDev MCP or Playwright MCP |
| **Accessibility-focused** | Playwright MCP (`browser_snapshot`) |
| **Human-in-the-loop flows** | w3pilot MCP (`pause_for_human`) |
| **Test automation with reports** | w3pilot MCP (verification + reports) |
| **Multi-browser support** | Playwright MCP |

### Client Library Selection

| Use Case | Recommendation |
|----------|----------------|
| **Go application** | w3pilot SDK |
| **JavaScript/TypeScript app** | vibium-js |
| **Python application** | vibium-py |
| **Java/JVM application** | vibium-java |
| **Script-based automation** | w3pilot (`w3pilot run`) |

---

## Relationship Between Projects

```
┌────────────────────────────────────────────────────────────────────┐
│                      VibiumDev/vibium                              │
│                                                                    │
│   ┌─────────┐  ┌─────────┐  ┌─────────┐                            │
│   │vibium-js│  │vibium-py│  │vibium-  │                            │
│   │ Client  │  │ Client  │  │java Cli │                            │
│   └────┬────┘  └────┬────┘  └────┬────┘                            │
│        └────────────┼───────────┘                                  │
│                     │ HTTP API                                     │
│                     ▼                                              │
│               ┌──────────┐                                         │
│               │ clicker  │                                         │
│               │ (binary) │                                         │
│               └────┬─────┘                                         │
│                    │ BiDi                                          │
│                    ▼                                               │
│                 Chrome                                             │
└────────────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────────────┐
│                    plexusone/w3pilot                              │
│              (dual-protocol architecture)                          │
│                                                                    │
│   ┌────────────────────────────────────────────────────────────┐   │
│   │                    w3pilot SDK                            │   │
│   │         BiDi client (automation) + CDP client (profiling)  │   │
│   └────────────────────────┬───────────────────────────────────┘   │
│                            │                                       │
│         ┌──────────────────┼──────────────────┐                    │
│         │                  │                  │                    │
│         ▼                  ▼                  ▼                    │
│    w3pilot-mcp         w3pilot run        Direct SDK             │
│         │                  │                  │                    │
│         └──────────────────┼──────────────────┘                    │
│                    ┌───────┴───────┐                               │
│                    │               │                               │
│                    ▼               ▼                               │
│               BiDi (clicker)    CDP (direct)                       │
│                    │               │                               │
│                    └───────┬───────┘                               │
│                            ▼                                       │
│                         Chrome                                     │
│                  (single browser instance)                         │
└────────────────────────────────────────────────────────────────────┘
```

**Key point**: w3pilot uses clicker (from VibiumDev) for BiDi automation commands via pipe mode, and direct CDP for profiling features. This dual-protocol architecture provides comprehensive automation with advanced performance profiling capabilities.

---

## Legend

| Symbol | Meaning |
|--------|---------|
| :white_check_mark: | Supported |
| :x: | Not supported |
| (opt) | Requires opt-in flag |
