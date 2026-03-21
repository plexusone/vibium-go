# Feature Parity Tasks

Tasks for achieving feature parity with Vibium clients (Java/JS/Python) and Playwright MCP.

Reference: [Feature Comparison](docs/reference/comparison.md)

## Legend

- [ ] Not started
- [x] Completed
- [~] In progress

---

## High Priority - SDK Features

### Semantic Selectors

Find elements by accessibility attributes instead of just CSS selectors.

- [x] Add `FindOptions` struct with semantic fields (role, text, label, placeholder, alt, title, testid, xpath, near)
- [x] Update `Vibe.Find()` to accept semantic options via FindOptions
- [x] Update `Vibe.FindAll()` to accept semantic options (uses vibium:findAll BiDi method)
- [x] Add `Element.Find()` for scoped semantic search
- [x] Add `Element.FindAll()` for scoped semantic search
- [x] Add MCP tool parameters for semantic selectors (click, type, fill, press tools)
- [x] Add tests for semantic selectors
- [x] Update documentation (README, SDK guide)

### Recording/Tracing

Full trace recording for debugging and test creation.

- [x] Add `Tracing` type to SDK - already existed in tracing.go
- [x] Implement `Tracing.Start()` with options (name, screenshots, snapshots, sources, title) - already existed
- [x] Implement `Tracing.Stop()` returning zip data - already existed
- [x] Implement `Tracing.StartChunk()` / `Tracing.StopChunk()` - already existed
- [x] Implement `Tracing.StartGroup()` / `Tracing.StopGroup()` - already existed
- [x] Add `BrowserContext.Tracing()` accessor - already existed
- [x] Add `Vibe.Tracing()` accessor for default context
- [x] Add MCP tools: `start_trace`, `stop_trace`, `start_trace_chunk`, `stop_trace_chunk`, `start_trace_group`, `stop_trace_group`
- [x] Add tests
- [x] Update documentation

### Media Emulation

CSS media feature emulation for accessibility testing.

- [x] Extend `EmulateMedia()` to support colorScheme - already existed
- [x] Add reducedMotion support - already existed
- [x] Add forcedColors support - already existed
- [x] Add contrast support
- [x] Add MCP tool: `emulate_media` with all options
- [ ] Add tests
- [x] Update documentation

### Console/Error Collection

Capture and buffer console messages and page errors.

- [x] Add `ConsoleMessage` type (type, text, args, location) - exists in route.go
- [ ] Implement `Page.OnConsole()` listener
- [ ] Implement `Page.CollectConsole()` for buffered mode
- [x] Implement `Page.ConsoleMessages()` to retrieve buffer
- [x] Implement `Page.ClearConsoleMessages()` to clear buffer
- [ ] Implement `Page.OnError()` listener
- [ ] Implement `Page.CollectErrors()` for buffered mode
- [ ] Implement `Page.Errors()` to retrieve buffer
- [x] Add MCP tool: `get_console_messages` with filtering
- [x] Add MCP tool: `clear_console_messages`
- [ ] Add tests
- [x] Update documentation

### Request/Response Listeners

Network observation events.

- [x] Add `Request` type with full metadata - exists in route.go
- [x] Add `Response` type with body/json methods - exists in route.go
- [x] Add `NetworkRequest` type for captured requests
- [ ] Implement `Page.OnRequest()` listener
- [ ] Implement `Page.OnResponse()` listener
- [x] Implement `Page.NetworkRequests()` to retrieve buffer
- [x] Implement `Page.ClearNetworkRequests()` to clear buffer
- [x] Add MCP tool: `get_network_requests` with filtering
- [x] Add MCP tool: `clear_network_requests`
- [ ] Add tests
- [x] Update documentation

### Full Storage State

Complete browser storage management.

- [x] Add `StorageState` type (cookies, origins with localStorage/sessionStorage)
- [x] Add `StorageStateOrigin` type (origin, localStorage, sessionStorage)
- [x] Implement `Vibe.StorageState()` to get full state (cookies + localStorage + sessionStorage)
- [x] Implement `Vibe.SetStorageState()` to restore state
- [x] Implement `Vibe.ClearStorage()` to clear all
- [x] Update MCP `get_storage_state` to include sessionStorage
- [x] Update MCP `set_storage_state` to restore sessionStorage
- [x] Add MCP `clear_storage` tool
- [x] Add tests
- [x] Update documentation

### Init Scripts

Per-context initialization scripts.

- [x] Implement `BrowserContext.AddInitScript()` - already existed
- [x] Implement `Vibe.AddInitScript()` for default context
- [x] Scripts run before page scripts on every navigation
- [x] Add MCP tool: `add_init_script`
- [x] Add `--init-script` CLI flag to `vibium launch` and `vibium mcp`
- [x] Add tests
- [x] Update documentation

---

## High Priority - MCP Tools

### LocalStorage Tools

- [x] Add `localstorage_get` tool (key) -> value
- [x] Add `localstorage_set` tool (key, value)
- [x] Add `localstorage_list` tool () -> all items
- [x] Add `localstorage_delete` tool (key)
- [x] Add `localstorage_clear` tool ()
- [ ] Add tests
- [x] Update mcp-tools.md documentation

### SessionStorage Tools

- [x] Add `sessionstorage_get` tool (key) -> value
- [x] Add `sessionstorage_set` tool (key, value)
- [x] Add `sessionstorage_list` tool () -> all items
- [x] Add `sessionstorage_delete` tool (key)
- [x] Add `sessionstorage_clear` tool ()
- [ ] Add tests
- [x] Update mcp-tools.md documentation

### Network Mocking

- [x] Add `route` tool (pattern, response options)
- [x] Add `route_list` tool () -> active routes
- [x] Add `unroute` tool (pattern)
- [x] Add `network_state_set` tool (offline: bool)
- [ ] Add tests
- [x] Update mcp-tools.md documentation

### Tab Management

- [x] Add `list_tabs` tool () -> tab info array
- [x] Add `select_tab` tool (index or id)
- [x] Add `close_tab` tool (index or id)
- [ ] Add tests
- [x] Update mcp-tools.md documentation

### Dialog Handling

- [x] Add `handle_dialog` tool (action: accept/dismiss, promptText?)
- [x] Add `get_dialog` tool () -> dialog info
- [ ] Add tests
- [x] Update mcp-tools.md documentation

### Console Messages Tool

- [x] Add `get_console_messages` tool (level filter?)
- [x] Add `clear_console_messages` tool
- [ ] Add tests
- [x] Update mcp-tools.md documentation

### Network Requests Listing

- [x] Add `get_network_requests` tool (filter options?)
- [x] Add `clear_network_requests` tool
- [ ] Add tests
- [x] Update mcp-tools.md documentation

---

## Medium Priority - SDK Features

### Element Methods

- [x] `Element.InnerText()` - rendered text only (already implemented)
- [x] `Element.InnerHTML()` - inner HTML (already implemented)
- [x] `Element.DispatchEvent(eventType, eventInit)` - already implemented
- [x] Add `Element.HTML()` - outer HTML (outerHTML, not innerHTML)
- [x] Add MCP tool: `get_outer_html`

### Page Methods

- [x] Add `Page.Scroll(direction, amount, selector)` - general page scrolling
- [x] Add `Page.BringToFront()` - already implemented
- [ ] Add `Page.MainFrame()` - returns page itself
- [x] Add `Page.SetExtraHTTPHeaders()` - already implemented as SetExtraHTTPHeaders
- [ ] Add tests

### Page Events

- [ ] Implement `Browser.OnPage()` listener
- [ ] Implement `Browser.OnPopup()` listener
- [ ] Implement `Browser.RemoveAllListeners()`
- [ ] Add tests

### WebSocket Monitoring

- [ ] Add `WebSocketInfo` type (url, isClosed)
- [ ] Implement `Page.OnWebSocket()` listener
- [ ] Implement `WebSocketInfo.OnMessage()` listener
- [ ] Implement `WebSocketInfo.OnClose()` listener
- [ ] Add tests

---

## Medium Priority - MCP Tools

### Tracing Tools

- [ ] Add `start_tracing` tool (options)
- [ ] Add `stop_tracing` tool (path?) -> file path
- [ ] Add tests
- [ ] Update mcp-tools.md documentation

### Video Recording

- [ ] Add `start_video` tool (size options)
- [ ] Add `stop_video` tool () -> file path
- [ ] Add tests
- [ ] Update mcp-tools.md documentation

### Form Tools

- [x] Add `fill_form` tool (fields: array of {selector, value})
- [ ] Add tests
- [x] Update mcp-tools.md documentation

### Mouse Tools

- [x] Add `mouse_drag` tool (startX, startY, endX, endY)
- [ ] Add tests
- [x] Update mcp-tools.md documentation

### Testing Tools

- [x] Add `verify_value` tool (selector, expected)
- [x] Add `verify_list_visible` tool (items array)
- [x] Add `generate_locator` tool (selector) -> locator string
- [ ] Add tests
- [x] Update mcp-tools.md documentation

---

## Low Priority

### Miscellaneous SDK

- [ ] Add `Element.Highlight()` - visual debugging (Java only feature)
- [ ] Add accessibility tree options (interestingOnly, root)

### Miscellaneous MCP

- [x] Add `get_config` tool - return resolved configuration

---

## Completed

_Move completed tasks here with completion date._

---

## Notes

- Semantic selectors are the highest priority as they enable AI-native element finding
- Storage tools (localStorage/sessionStorage) unblock many authentication workflows
- Recording/tracing helps with debugging and test creation
- All new features should have corresponding MCP tools where applicable
