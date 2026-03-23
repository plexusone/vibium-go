# vibium-go vs VibiumDev

This document compares vibium-go with the official VibiumDev project, covering both MCP servers and client libraries.

## Overview

vibium-go is an independent Go implementation that uses VibiumDev's `clicker` binary for browser launching but provides its own BiDi client and MCP server.

| Project | Language | MCP Server | Client Library | Protocol |
|---------|----------|------------|----------------|----------|
| [VibiumDev/vibium](https://github.com/VibiumDev/vibium) | Go (clicker) + JS/Python/Java (clients) | `vibium mcp` | vibium-js, vibium-py, vibium-java | WebDriver BiDi |
| [plexusone/vibium-go](https://github.com/plexusone/vibium-go) | Go | `vibium-mcp` | vibium-go SDK | WebDriver BiDi |
| [microsoft/playwright-mcp](https://github.com/microsoft/playwright-mcp) | TypeScript | `playwright-mcp` | Playwright | CDP + BiDi |

---

## Architecture

### VibiumDev Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        VibiumDev/vibium                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │
│  │ vibium-js   │    │ vibium-py   │    │ vibium-java │             │
│  │ (Client)    │    │ (Client)    │    │ (Client)    │             │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘             │
│         │                  │                  │                     │
│         └──────────────────┼──────────────────┘                     │
│                            │ HTTP API                               │
│                            ▼                                        │
│                   ┌────────────────┐                                │
│                   │    clicker     │◄──── vibium mcp (MCP Server)   │
│                   │   (Go binary)  │                                │
│                   └────────┬───────┘                                │
│                            │ WebDriver BiDi                         │
│                            ▼                                        │
│                   ┌────────────────┐                                │
│                   │     Chrome     │                                │
│                   └────────────────┘                                │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### vibium-go Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        plexusone/vibium-go                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                      vibium-go SDK                           │   │
│  │  (Go client library with full browser automation API)        │   │
│  └──────────────────────────┬──────────────────────────────────┘   │
│                             │                                       │
│         ┌───────────────────┼───────────────────┐                   │
│         │                   │                   │                   │
│         ▼                   ▼                   ▼                   │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │
│  │ vibium-mcp  │    │ vibium run  │    │ Direct SDK  │             │
│  │ (MCP Server)│    │ (Script)    │    │   Usage     │             │
│  └─────────────┘    └─────────────┘    └─────────────┘             │
│                             │                                       │
│                             │ WebDriver BiDi (direct)               │
│                             ▼                                       │
│  ┌─────────────┐    ┌────────────────┐                             │
│  │   clicker   │───►│     Chrome     │                             │
│  │ (launcher)  │    └────────────────┘                             │
│  └─────────────┘                                                    │
│        ▲                                                            │
│        │ Used only for launching Chrome                             │
│        │ (from VibiumDev/vibium)                                    │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Key Architectural Difference

```
VibiumDev:    Client ──HTTP──► clicker ──BiDi──► Chrome
vibium-go:    SDK ──────────────BiDi (direct)──► Chrome
                                    ▲
                              clicker only launches
```

vibium-go communicates directly with Chrome via WebDriver BiDi, bypassing the HTTP API layer. The `clicker` binary is only used to launch Chrome with the correct flags and obtain the WebSocket URL.

---

## MCP Server Comparison

### Tool Inventory

| Category | VibiumDev MCP | Playwright MCP | vibium-go MCP |
|----------|:-------------:|:--------------:|:-------------:|
| **Browser Management** | | | |
| Launch browser | `browser_start` | `browser_launch` | `browser_launch` |
| Quit browser | `browser_stop` | `browser_close` | `browser_quit` |
| **Navigation** | | | |
| Navigate | `browser_navigate` | `browser_navigate` | `navigate` |
| Back | - | `browser_back` | `back` |
| Forward | - | `browser_forward` | `forward` |
| Reload | - | `browser_reload` | `reload` |
| Scroll | `browser_scroll` | `browser_scroll` | `scroll` |
| **Element Interaction** | | | |
| Click | `browser_click` | `browser_click` | `click` |
| Double-click | - | - | `dblclick` |
| Type | `browser_type` | `browser_type` | `type` |
| Fill | - | `browser_fill` | `fill` |
| Clear | - | - | `clear` |
| Press key | `browser_keys` | `browser_press_key` | `press` |
| Check/Uncheck | - | - | `check`, `uncheck` |
| Select option | `browser_select` | `browser_select_option` | `select_option` |
| Hover | `browser_hover` | `browser_hover` | `hover` |
| Drag | - | `browser_drag` | `drag_to` |
| Fill form | - | - | `fill_form` |
| **Element State** | | | |
| Get text | `browser_get_text` | `browser_snapshot` | `get_text` |
| Get value | - | - | `get_value` |
| Get HTML | `browser_get_html` | `browser_snapshot` | `get_inner_html`, `get_outer_html` |
| Get attribute | - | - | `get_attribute` |
| Get bounding box | - | - | `get_bounding_box` |
| Is visible | - | - | `is_visible` |
| Is enabled | - | - | `is_enabled` |
| Is checked | - | - | `is_checked` |
| **Page State** | | | |
| Get title | `browser_get_title` | - | `get_title` |
| Get URL | `browser_get_url` | - | `get_url` |
| Get content | - | `browser_snapshot` | `get_content` |
| Set content | - | - | `set_content` |
| Accessibility snapshot | - | `browser_snapshot` | `accessibility_snapshot` |
| **Screenshots & PDF** | | | |
| Screenshot | `browser_screenshot` | `browser_screenshot` | `screenshot` |
| Element screenshot | - | `browser_screenshot` | `element_screenshot` |
| PDF | - | `browser_pdf_save` | `pdf` |
| **JavaScript** | | | |
| Evaluate | `browser_evaluate` | `browser_console_exec` | `evaluate` |
| Element eval | - | - | `element_eval` |
| Add script | - | - | `add_script` |
| Add style | - | - | `add_style` |
| **Waiting** | | | |
| Wait for element | `browser_wait` | `browser_wait_for` | `wait_until` |
| Wait for URL | - | - | `wait_for_url` |
| Wait for load | - | - | `wait_for_load` |
| Wait for function | - | - | `wait_for_function` |
| Wait for text | - | `browser_wait_for` | `wait_for_text` |
| Wait for selector | - | `browser_wait_for` | `wait_for_selector` |
| **Input Controllers** | | | |
| Keyboard press | `browser_keys` | `browser_press_key` | `keyboard_press` |
| Keyboard down/up | - | - | `keyboard_down`, `keyboard_up` |
| Keyboard type | - | - | `keyboard_type` |
| Mouse click | - | `browser_click` | `mouse_click` |
| Mouse move | - | `browser_move_mouse` | `mouse_move` |
| Mouse down/up | - | - | `mouse_down`, `mouse_up` |
| Mouse wheel | - | - | `mouse_wheel` |
| Mouse drag | - | `browser_drag` | `mouse_drag` |
| Touch tap | - | - | `touch_tap` |
| Touch swipe | - | - | `touch_swipe` |
| **Tab/Page Management** | | | |
| New page | `browser_new_page` | `browser_tab_new` | `new_page` |
| List pages | `browser_list_pages` | `browser_tab_list` | `get_pages`, `list_tabs` |
| Switch page | `browser_switch_page` | `browser_tab_select` | `select_tab` |
| Close page | `browser_close_page` | `browser_tab_close` | `close_page`, `close_tab` |
| Bring to front | - | - | `bring_to_front` |
| **Frame Management** | | | |
| Get frames | - | - | `get_frames` |
| Select frame | - | - | `select_frame` |
| Select main frame | - | - | `select_main_frame` |
| **Cookie Management** | | | |
| Get cookies | - | - | `get_cookies` |
| Set cookies | - | - | `set_cookies` |
| Clear cookies | - | - | `clear_cookies` |
| Delete cookie | - | - | `delete_cookie` |
| **LocalStorage** | | | |
| Get item | - | - | `localstorage_get` |
| Set item | - | - | `localstorage_set` |
| List items | - | - | `localstorage_list` |
| Delete item | - | - | `localstorage_delete` |
| Clear | - | - | `localstorage_clear` |
| **SessionStorage** | | | |
| Get item | - | - | `sessionstorage_get` |
| Set item | - | - | `sessionstorage_set` |
| List items | - | - | `sessionstorage_list` |
| Delete item | - | - | `sessionstorage_delete` |
| Clear | - | - | `sessionstorage_clear` |
| **Storage State** | | | |
| Get storage state | - | - | `get_storage_state` |
| Set storage state | - | - | `set_storage_state` |
| Clear storage | - | - | `clear_storage` |
| **Network** | | | |
| Route (mock) | - | `browser_network_mock` | `route` |
| Route list | - | - | `route_list` |
| Unroute | - | `browser_network_unmock` | `unroute` |
| Network offline | - | - | `network_state_set` |
| Get network requests | - | - | `get_network_requests` |
| Clear network requests | - | - | `clear_network_requests` |
| **Console** | | | |
| Get console messages | - | `browser_console_messages` | `get_console_messages` |
| Clear console messages | - | - | `clear_console_messages` |
| **Dialogs** | | | |
| Handle dialog | - | `browser_dialog_handle` | `handle_dialog` |
| Get dialog info | - | - | `get_dialog` |
| **Emulation** | | | |
| Emulate media | - | - | `emulate_media` |
| Set geolocation | - | - | `set_geolocation` |
| **Recording & Tracing** | | | |
| Start recording | - | - | `start_recording` |
| Stop recording | - | - | `stop_recording` |
| Export script | - | - | `export_script` |
| Start trace | - | - | `start_trace` |
| Stop trace | - | - | `stop_trace` |
| Start video | - | - | `start_video` |
| Stop video | - | - | `stop_video` |
| **Testing & Assertions** | | | |
| Assert text | - | - | `assert_text` |
| Assert element | - | - | `assert_element` |
| Verify value | - | - | `verify_value` |
| Verify text | - | - | `verify_text` |
| Verify visible | - | - | `verify_visible` |
| Verify hidden | - | - | `verify_hidden` |
| Verify enabled | - | - | `verify_enabled` |
| Verify disabled | - | - | `verify_disabled` |
| Verify checked | - | - | `verify_checked` |
| Verify list visible | - | - | `verify_list_visible` |
| Generate locator | - | - | `generate_locator` |
| Get test report | - | - | `get_test_report` |
| **Human-in-the-Loop** | | | |
| Pause for human | - | - | `pause_for_human` |
| **Configuration** | | | |
| Get config | - | - | `get_config` |
| Add init script | - | - | `add_init_script` |

### Tool Count Summary

| MCP Server | Tool Count | Focus |
|------------|:----------:|-------|
| VibiumDev | ~25 | Core automation |
| Playwright MCP | ~30 | Snapshot-based, AI-friendly |
| vibium-go | **100+** | Comprehensive automation + testing |

### MCP Implementation

| Aspect | VibiumDev MCP | Playwright MCP | vibium-go MCP |
|--------|---------------|----------------|---------------|
| Language | Go | TypeScript | Go |
| MCP SDK | Hand-rolled JSON-RPC | Official TS SDK | Official Go SDK |
| Tool naming | `browser_*` prefix | `browser_*` prefix | No prefix |
| Schema generation | Manual | Manual | Auto from Go structs |

---

## Client Library Comparison

### VibiumDev Clients vs vibium-go SDK

| Feature | vibium-js | vibium-py | vibium-java | vibium-go |
|---------|:---------:|:---------:|:-----------:|:---------:|
| **Core** | | | | |
| Launch browser | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Headless mode | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Connect remote | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Multiple contexts | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| **Element Finding** | | | | |
| CSS selector | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Semantic selectors | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By role | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By text | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By label | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By testid | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| - By proximity | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| **Interactions** | | | | |
| Click/Type/Fill | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Check/Uncheck | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Select option | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Drag and drop | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| File upload | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| **Input Controllers** | | | | |
| Keyboard | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Mouse | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Touch | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| **Event Listeners** | | | | |
| onConsole | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onError | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onRequest/onResponse | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onDialog | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onDownload | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onPage/onPopup | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| onWebSocket | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| **Recording/Tracing** | | | | |
| Trace recording | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Video recording | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| **Storage** | | | | |
| Cookies | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| LocalStorage | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| SessionStorage | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Full storage state | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Init scripts | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| **Network** | | | | |
| Route/mock | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Offline mode | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Extra headers | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| **Clock Control** | | | | |
| Install/fastForward | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| pauseAt/resume | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Language-Specific Considerations

| Aspect | vibium-js | vibium-py | vibium-java | vibium-go |
|--------|-----------|-----------|-------------|-----------|
| Async model | Promises/async-await | async/await | CompletableFuture | Context-based |
| Error handling | try/catch | try/except | try/catch | error returns |
| Package manager | npm | pip | Maven/Gradle | go modules |
| Type safety | TypeScript optional | Type hints | Strong typing | Strong typing |

---

## Unique Features

### vibium-go Only

| Feature | Description |
|---------|-------------|
| Script Runner | YAML/JSON deterministic test execution via `vibium run` |
| Session Recording | Capture MCP actions as replayable scripts |
| Human-in-the-Loop | `pause_for_human` for SSO, CAPTCHA, 2FA handling |
| Test Reports | Structured test execution reports (box, diagnostic, JSON) |
| Verification Tools | `verify_*` tools for assertions with detailed output |
| Frame Selection | `select_frame`/`select_main_frame` for iframe navigation |

### VibiumDev Only

| Feature | Description |
|---------|-------------|
| Multi-language clients | Official JS, Python, Java clients |
| Daemon mode | HTTP API server for multi-client scenarios |
| Element highlight | Visual debugging overlay (Java client) |

---

## When to Use Which

| Use Case | Recommendation |
|----------|----------------|
| **Go application** | vibium-go SDK |
| **JavaScript/Python/Java app** | VibiumDev clients |
| **LLM agent (Claude, etc.)** | vibium-go MCP (most comprehensive) |
| **Simple MCP automation** | VibiumDev MCP or Playwright MCP |
| **E2E testing with reports** | vibium-go (test reports, verification tools) |
| **Human-in-the-loop flows** | vibium-go (`pause_for_human`) |
| **Script-based automation** | vibium-go (`vibium run`) |

---

## Summary

vibium-go provides the most comprehensive MCP server with 100+ tools, while VibiumDev offers a simpler ~25 tool set with multi-language client support. Both use WebDriver BiDi protocol and can work together (vibium-go uses VibiumDev's clicker for browser launching).

For Go developers or LLM agents requiring extensive browser automation, vibium-go is the recommended choice. For polyglot teams needing JS/Python/Java clients, VibiumDev provides official support.
