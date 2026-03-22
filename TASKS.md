# Feature Parity Tasks

Tasks for achieving feature parity with Vibium clients (Java/JS/Python) and Playwright MCP.

Reference: [Feature Comparison](docs/reference/comparison.md)

## Legend

- [ ] Not started
- [x] Completed

---

## Open Tasks

No open tasks remaining.

---

## Completed - v0.5.0 (2026-03-21)

### Tests

Integration tests for implemented features:

- [x] Media emulation tests - `integration/media_emulation_test.go`
- [x] LocalStorage MCP tools tests - `integration/storage_tools_test.go`
- [x] SessionStorage MCP tools tests - `integration/storage_tools_test.go`
- [x] Network mocking MCP tools tests (`route`, `unroute`, `network_state_set`) - `integration/network_tools_test.go`
- [x] Tab management MCP tools tests (`NewPage`, `Pages`, `Close`, `BringToFront`) - `integration/dialog_tab_test.go`
- [x] Dialog handling MCP tools tests (`HandleDialog`) - `integration/dialog_tab_test.go`
- [x] Console messages MCP tools tests (`ConsoleMessages`, `ClearConsoleMessages`) - `integration/page_methods_test.go`
- [x] Network requests MCP tools tests (`NetworkRequests`, `ClearNetworkRequests`) - `integration/network_tools_test.go`
- [x] Form tools tests (`Fill`) - `integration/page_methods_test.go`
- [x] Mouse tools tests (`Mouse.Move`, `Mouse.Down`, `Mouse.Up`) - `integration/page_methods_test.go`
- [x] Testing tools tests (`verify_value`, `verify_list_visible`, `generate_locator`) - `integration/verify_test.go`
- [x] Page methods tests (`Scroll`, `SetExtraHTTPHeaders`) - `integration/page_methods_test.go`

### Event Listeners

Real-time event callbacks for SDK users:

- [x] `Vibe.OnConsole()` - console message listener
- [x] `Vibe.CollectConsole()` - enable buffered console collection
- [x] `Vibe.OnError()` - page error listener
- [x] `Vibe.CollectErrors()` - enable buffered error collection
- [x] `Vibe.Errors()` - retrieve buffered errors
- [x] `Vibe.ClearErrors()` - clear buffered errors
- [x] `Vibe.OnRequest()` - network request listener
- [x] `Vibe.OnResponse()` - network response listener
- [x] `Vibe.OnDialog()` - dialog event listener
- [x] `Vibe.OnDownload()` - download event listener
- [x] BiDi client event dispatch infrastructure (`bidi.go`)

### Page Events

Browser-level event listeners:

- [x] `Vibe.OnPage()` - new page created listener
- [x] `Vibe.OnPopup()` - popup window listener
- [x] `Vibe.RemoveAllListeners()` - cleanup all listeners

### WebSocket Monitoring

WebSocket connection observation:

- [x] `WebSocketInfo` type (URL, IsClosed, socketID)
- [x] `WebSocketMessage` type (Data, IsBinary, Direction)
- [x] `Vibe.OnWebSocket()` - WebSocket connection listener
- [x] `WebSocketInfo.OnMessage()` - message listener
- [x] `WebSocketInfo.OnClose()` - close listener

### Video Recording

Screen recording for debugging:

- [x] `start_video` MCP tool (size options)
- [x] `stop_video` MCP tool () -> file path
- [x] SDK `Vibe.StartVideo()` / `Vibe.StopVideo()` methods
- [x] `Video.Path()` / `Video.Delete()` methods

### SDK Methods

API compatibility and debugging features:

- [x] `Vibe.MainFrame()` - returns self for API compatibility
- [x] `Element.Highlight()` - visual debugging overlay
- [x] `A11yTreeOptions` - interestingOnly and root options for accessibility tree
- [x] `BrowserContext.DeleteCookie()` - delete single cookie by name

### MCP Tools - Additional

- [x] `delete_cookie` - delete a specific cookie by name
- [x] `wait_for_text` - wait for text to appear on the page
- [x] `accessibility_snapshot` - get accessibility tree snapshot
- [x] `wait_for_selector` - wait for element to appear/disappear with state option
- [x] `verify_text` - verify element text matches expected value
- [x] `verify_visible` - verify element is visible
- [x] `verify_enabled` - verify element is enabled
- [x] `verify_checked` - verify checkbox/radio is checked
- [x] `verify_hidden` - verify element is hidden
- [x] `verify_disabled` - verify element is disabled

### MCP Tools - Frame Selection

- [x] `select_frame` - switch to a frame by name or URL pattern
- [x] `select_main_frame` - switch back to the main frame

---

## Completed - v0.4.0 (2026-03-21)

### Semantic Selectors

Find elements by accessibility attributes instead of CSS selectors.

- [x] `FindOptions` struct with semantic fields (role, text, label, placeholder, alt, title, testid, xpath, near)
- [x] `Vibe.Find()` accepts semantic options via FindOptions
- [x] `Vibe.FindAll()` accepts semantic options
- [x] `Element.Find()` for scoped semantic search
- [x] `Element.FindAll()` for scoped semantic search
- [x] MCP tool parameters for semantic selectors (click, type, fill, press)
- [x] Integration tests
- [x] Documentation (README, SDK guide)

### Recording/Tracing

Full trace recording for debugging and test creation.

- [x] `Tracing` type with `Start()`, `Stop()`, `StartChunk()`, `StopChunk()`, `StartGroup()`, `StopGroup()`
- [x] `BrowserContext.Tracing()` accessor
- [x] `Vibe.Tracing()` accessor for default context
- [x] MCP tools: `start_trace`, `stop_trace`, `start_trace_chunk`, `stop_trace_chunk`, `start_trace_group`, `stop_trace_group`
- [x] Integration tests
- [x] Documentation

### Full Storage State

Complete browser storage management including sessionStorage.

- [x] `StorageState` type (cookies, origins with localStorage/sessionStorage)
- [x] `StorageStateOrigin` type (origin, localStorage, sessionStorage)
- [x] `Vibe.StorageState()` - get full state
- [x] `Vibe.SetStorageState()` - restore state
- [x] `Vibe.ClearStorage()` - clear all
- [x] MCP tools: `get_storage_state`, `set_storage_state`, `clear_storage`
- [x] Integration tests
- [x] Documentation

### Init Scripts

Per-context initialization scripts that run before page scripts.

- [x] `BrowserContext.AddInitScript()`
- [x] `Vibe.AddInitScript()` for default context
- [x] MCP tool: `add_init_script`
- [x] `--init-script` CLI flag for `vibium launch` and `vibium mcp`
- [x] Integration tests
- [x] Documentation

### MCP Tools - Dialog Handling

- [x] `handle_dialog` tool (action: accept/dismiss, promptText)
- [x] `get_dialog` tool

### MCP Tools - Network

- [x] `get_network_requests` tool (filter options)
- [x] `clear_network_requests` tool
- [x] `route` tool (pattern, response options)
- [x] `route_list` tool
- [x] `unroute` tool (pattern)
- [x] `network_state_set` tool (offline: bool)

### MCP Tools - Storage

- [x] `localstorage_get`, `localstorage_set`, `localstorage_list`, `localstorage_delete`, `localstorage_clear`
- [x] `sessionstorage_get`, `sessionstorage_set`, `sessionstorage_list`, `sessionstorage_delete`, `sessionstorage_clear`

### MCP Tools - Tabs

- [x] `list_tabs`, `select_tab`, `close_tab`

### MCP Tools - Console

- [x] `get_console_messages` (level filter)
- [x] `clear_console_messages`

### MCP Tools - Testing

- [x] `fill_form` (fields array)
- [x] `verify_value` (selector, expected)
- [x] `verify_list_visible` (items array)
- [x] `generate_locator` (selector)
- [x] `mouse_drag` (startX, startY, endX, endY)

---

## Completed - v0.3.0 (2026-03-16)

### Human-in-the-Loop

- [x] `pause_for_human` MCP tool
- [x] `set_storage_state` MCP tool (initial version)

---

## Completed - Pre-v0.3.0

### Media Emulation

- [x] `EmulateMedia()` with colorScheme, reducedMotion, forcedColors, contrast
- [x] `emulate_media` MCP tool

### Element Methods

- [x] `Element.InnerText()`, `Element.InnerHTML()`, `Element.HTML()` (outerHTML)
- [x] `Element.DispatchEvent()`
- [x] `get_outer_html` MCP tool

### Page Methods

- [x] `Page.Scroll()` (direction, amount, selector)
- [x] `Page.BringToFront()`
- [x] `Page.SetExtraHTTPHeaders()`

### Console/Network Buffering

- [x] `ConsoleMessage` type
- [x] `Page.ConsoleMessages()`, `Page.ClearConsoleMessages()`
- [x] `Request`, `Response`, `NetworkRequest` types
- [x] `Page.NetworkRequests()`, `Page.ClearNetworkRequests()`

### Miscellaneous MCP

- [x] `get_config` tool

---

## Notes

- All feature parity tasks are complete as of v0.5.0
- Future work can focus on performance optimization, additional platform support, or new features
