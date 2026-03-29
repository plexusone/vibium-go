# Feature Parity Tasks

Tasks for achieving feature parity with Vibium clients (Java/JS/Python), Playwright MCP, and Chrome DevTools MCP.

Reference: [Feature Comparison](docs/reference/comparison.md)

## Legend

- [ ] Not started
- [~] In progress
- [x] Completed

---

## Open Tasks

### P0 - Critical

#### Live Session Mode for CLI

**Problem**: `w3pilot run` launches a browser, runs the script, and quits. We need a live session mode where the browser stays open for interactive use by AI agents, similar to how the MCP server works.

**Use Case**: AI assistant + human working together - the AI uses MCP tools while the human can also interact with the browser and use CLI commands interactively.

**Implementation Status**:

SDK Session Package (completed):
- [x] `session/types.go` - Session info, config, status types
- [x] `session/persistence.go` - Session file I/O (Save, Load, Clear, Exists)
- [x] `session/manager.go` - SessionManager with LaunchIfNeeded, Reconnect, Detach, Close
- [x] Unit tests for session package

Core SDK Support (completed):
- [x] `w3pilot.Connect(ctx, wsURL)` - Connect to existing clicker instance
- [x] `Pilot.Clicker()` - Access clicker process info
- [x] `ClickerProcess.Process()` - Access underlying os.Process

CLI Commands (TODO):
- [ ] `w3pilot session start` - Launch browser and keep session alive
- [ ] `w3pilot session attach` - Attach to existing session
- [ ] `w3pilot session detach` - Detach without closing browser
- [ ] `w3pilot session stop` - Close browser and end session
- [ ] All CLI commands check for active session before launching new browser
- [ ] `w3pilot run --keep-alive` flag to keep browser open after script

**Architecture**:
```
┌─────────────────────────────────────────────────────────────┐
│                    Live Session Mode                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  w3pilot session start ──► Browser launched                 │
│          │                     │                             │
│          ▼                     ▼                             │
│  session.json ◄────── Connection info saved                 │
│          │                                                   │
│          ▼                                                   │
│  w3pilot page navigate ...  ──► Reuses existing session     │
│  w3pilot element click ...  ──► Reuses existing session     │
│  w3pilot run script.yaml    ──► Reuses existing session     │
│          │                                                   │
│          ▼                                                   │
│  w3pilot session stop ──► Browser closed                    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

#### SDK Assertion & Verification Methods

**Problem**: Test assertion and verification logic exists only in MCP handlers (`mcp/tools_testing.go`). This prevents CLI from using these features and creates code duplication.

**Goal**: Move assertion/verification logic to the SDK so both CLI and MCP can use it.

**SDK Methods to Add**:

Pilot Assertion Methods:
- [x] `Pilot.AssertText(ctx, text, opts)` - Assert text exists on page
- [x] `Pilot.AssertElement(ctx, selector, opts)` - Assert element exists
- [x] `Pilot.AssertURL(ctx, pattern, opts)` - Assert URL matches pattern
- [x] `Pilot.GenerateLocator(ctx, selector, opts)` - Generate robust locator for element

Element Verification Methods:
- [x] `Element.VerifyValue(ctx, expected)` - Verify input value matches
- [x] `Element.VerifyText(ctx, expected, opts)` - Verify element text matches
- [x] `Element.VerifyVisible(ctx)` - Verify element is visible
- [x] `Element.VerifyHidden(ctx)` - Verify element is hidden
- [x] `Element.VerifyEnabled(ctx)` - Verify element is enabled
- [x] `Element.VerifyDisabled(ctx)` - Verify element is disabled
- [x] `Element.VerifyChecked(ctx)` - Verify checkbox/radio is checked
- [x] `Element.VerifyUnchecked(ctx)` - Verify checkbox/radio is unchecked (bonus)

**Files Created/Modified**:
- [x] `assert.go` - New file with Pilot assertion methods and types
- [x] `element.go` - Added Element verification methods
- [x] `assert_test.go` - Unit tests for assertion methods
- [x] `element_verify_test.go` - Unit tests for verification types

**Remaining Work**:
- [ ] `mcp/tools_testing.go` - Refactor to use SDK methods
- [ ] `cmd/w3pilot/cmd/test_*.go` - Add CLI commands using SDK

---

### P1 - CLI/MCP Parity

#### Current Status

| Category | MCP Tools | CLI Commands | Gap |
|----------|-----------|--------------|-----|
| browser | 2 | 2 | ✅ Parity |
| page | 20 | 14 | ⚠️ 6 missing |
| element | 33 | 19 | ⚠️ 14 missing |
| js | 4 | 4 | ✅ Parity |
| wait | 6 | 5 | ⚠️ 1 missing |
| state | 4 | 4 | ✅ Parity |
| test | 16 | 1 | ⚠️ 15 missing |
| cdp | 20 | 0 | ❌ Stub only |
| storage | 17 | 0 | ❌ Stub only |
| input | 12 | 0 | ❌ Stub only |
| network | 6 | 0 | ❌ Stub only |
| trace | 6 | 0 | ❌ Stub (disabled) |
| record | 5 | 0 | ❌ Stub only |
| tab | 3 | 0 | ❌ Stub only |
| dialog | 2 | 0 | ❌ Stub only |
| console | 2 | 0 | ❌ Stub only |
| video | 2 | 0 | ❌ Stub only |
| frame | 2 | 0 | ❌ Stub only |
| workflow | 2 | 0 | ❌ Not in CLI |
| a11y | 1 | 0 | ❌ Stub only |
| config | 1 | 0 | ❌ Not in CLI |
| human | 1 | 0 | ❌ MCP-only |
| **Total** | **167** | **~45** | **122 missing** |

#### Priority Order for CLI Implementation

1. **High Priority** (most useful for scripting):
   - [ ] `storage` commands (localStorage, sessionStorage, cookies)
   - [ ] `network` commands (offline, requests, routes)
   - [ ] `console` commands (messages, clear)
   - [ ] `test` commands (assertions, verification)

2. **Medium Priority** (useful for debugging):
   - [ ] `cdp` commands (performance, memory, coverage)
   - [ ] `input` commands (keyboard, mouse, touch)
   - [ ] `tab` commands (list, select, close)
   - [ ] `dialog` commands (handle, get)

3. **Lower Priority** (specialized use cases):
   - [ ] `frame` commands (select, main)
   - [ ] `record` commands (start, stop, export)
   - [ ] `video` commands (start, stop)
   - [ ] `a11y` commands (snapshot)

---

### Blocked - Pending Clicker Support

Features that are implemented but disabled because clicker doesn't support the required commands.

#### Tracing (vibium:tracing.*)

Action recording for debugging and test creation. Commented out in v0.7.0.

- [ ] `vibium:tracing.start` - Start trace recording
- [ ] `vibium:tracing.stop` - Stop and get trace data
- [ ] `vibium:tracing.startChunk` / `stopChunk` - Chunk recording
- [ ] `vibium:tracing.startGroup` / `stopGroup` - Logical grouping

**Files affected**: `tracing.go`, `pilot.go`, `context.go`, `mcp/tools_tracing.go`, `mcp/server.go`

**Note**: CDP has `Tracing.start/end` but it's for performance tracing, not action recording.

#### Storage State (vibium:context.storageState)

Full browser state capture including localStorage and sessionStorage.

- [ ] `vibium:context.storageState` - Get complete storage state

**Workaround**: Use JavaScript evaluation to access localStorage/sessionStorage directly.

---

### P0 - Critical (Chrome DevTools MCP Parity)

Features that provide significant differentiation from other browser automation tools.

#### Performance Insights ✅ (JS APIs available) - COMPLETED v0.5.2

Core Web Vitals via JavaScript PerformanceObserver APIs - **no clicker changes needed**.

- [x] `get_performance_metrics` MCP tool
- [x] LCP via `PerformanceObserver('largest-contentful-paint')`
- [x] CLS via `PerformanceObserver('layout-shift')`
- [x] INP via `PerformanceObserver('event')` with interactionId
- [x] Navigation timing via `performance.timing`
- [x] SDK `Pilot.GetPerformanceMetrics()` method
- [x] SDK `Pilot.ObserveWebVitals()` method (real-time)

#### Memory Stats ✅ (JS APIs available) - COMPLETED v0.5.2

Basic memory info via `performance.memory` - **no clicker changes needed**.

- [x] `get_memory_stats` MCP tool (usedJSHeapSize, totalJSHeapSize, jsHeapSizeLimit)
- [x] SDK `Pilot.GetMemoryStats()` method
- [ ] Memory threshold alerts (optional)

#### Lighthouse Integration ✅ (External CLI) - COMPLETED v0.5.2

Run quality audits - requires **lighthouse CLI** (Node.js).

- [x] `lighthouse_audit` MCP tool (categories: accessibility, seo, best-practices, performance)
- [x] Shell to `npx lighthouse` or `lighthouse` CLI
- [x] SDK `Pilot.LighthouseAudit()` method
- [x] Return structured audit results with scores
- [ ] Integration with a11y-lab test fixtures (stretch goal)

#### Heap Snapshots ✅ (Direct CDP) - COMPLETED v0.5.2

Full heap profiling via **direct CDP connection** to same browser session.

- [x] Add `CDPClient` to Pilot struct (parallel to BiDi)
- [x] Discover CDP port from Chrome's `DevToolsActivePort` file
- [x] `take_heap_snapshot` MCP tool (returns .heapsnapshot file path)
- [x] SDK `Pilot.TakeHeapSnapshot()` method
- [x] CDP `HeapProfiler.takeHeapSnapshot` implementation
- [ ] Memory diff between snapshots (stretch goal)

### P1 - High Priority (AI Agent Ergonomics) - COMPLETED v0.6.0

Features that make w3pilot easier for AI agents to use effectively.

#### Page Inspection ✅

Single-call overview of interactive elements for AI decision-making.

- [x] `page_inspect` MCP tool - returns structured page overview
  - Buttons: text, selector, enabled state
  - Links: text, href, selector
  - Form fields: type, label, name, current value, selector
  - Headings: level, text, selector (for page structure)
- [x] SDK `Pilot.Inspect()` method
- [x] `w3pilot page inspect` CLI command with `--format json` output

#### JSON Output Mode for CLI ✅

Machine-readable output for scripting and AI tool use.

- [x] Global `--format` flag on root command (text, json)
- [x] `--format json` on all read commands (title, url, text, attr, etc.)
- [x] Structured JSON output with consistent schema
- [x] Error responses in JSON format when `--format json` enabled

#### Better Error Context ✅

Provide actionable context when operations fail.

- [x] Include current page URL in element errors
- [x] Include page title in navigation errors
- [x] Show visible text near expected selector location
- [x] Suggest alternative selectors when element not found
- [x] `ElementNotFoundError` includes `PageContext`

#### Selector Validation Tool ✅

Pre-validate selectors before running actions.

- [x] `test_validate_selectors` MCP tool (accepts array of selectors)
- [x] Returns: found/not found, count, visible state for each
- [x] SDK `Pilot.ValidateSelectors()` method
- [x] `w3pilot test validate-selectors` CLI command

#### Workflow Recipes ✅

Pre-built patterns for common automation tasks.

- [x] `workflow_login` MCP tool (username_selector, password_selector, submit_selector, credentials)
- [x] `workflow_extract_table` MCP tool (table_selector) → JSON array of rows
- [ ] `workflow_fill_form` MCP tool (fields array with label-based matching) - future
- [x] SDK `Pilot.Login()`, `Pilot.ExtractTable()` methods
- [x] Documentation with common workflow examples

#### Named State Snapshots ✅

Save/restore named browser states for session management.

- [x] `w3pilot state save <name>` CLI command
- [x] `w3pilot state load <name>` CLI command
- [x] `w3pilot state list` CLI command
- [x] `w3pilot state delete <name>` CLI command
- [x] `state_save` MCP tool (name) - saves to ~/.w3pilot/states/
- [x] `state_load` MCP tool (name) - restores from saved state
- [x] `state_list` MCP tool - list available named states
- [x] `state_delete` MCP tool - delete named state
- [x] State management via `state.Manager`

---

### P1 - High Priority (Debugging)

Features that improve debugging and testing capabilities.

#### Network Request Bodies ✅ (Direct CDP) - COMPLETED v0.5.2

Retrieve full request/response content via **direct CDP connection**.

- [x] CDP `Network.enable` with response body capture
- [x] CDP `Network.getResponseBody` implementation
- [x] `get_network_request_body` MCP tool (request_id, save_to_file)
- [x] SDK `Pilot.GetNetworkResponseBody()` method
- [x] Support for binary content (images, etc.)

#### Emulation Presets ✅ (Direct CDP) - COMPLETED v0.5.2

Network/CPU throttling via **direct CDP connection**.

- [x] CDP `Network.emulateNetworkConditions` implementation
- [x] CDP `Emulation.setCPUThrottlingRate` implementation
- [x] Network throttling presets (Slow 3G, Fast 3G, 4G, Offline)
- [x] CPU throttling (1x, 2x, 4x, 6x slowdown)
- [x] `emulate_network` MCP tool (preset or custom latency/bandwidth)
- [x] `emulate_cpu` MCP tool (throttle factor)
- [x] `clear_network_emulation` MCP tool
- [x] `clear_cpu_emulation` MCP tool
- [x] SDK `Pilot.EmulateNetwork()` / `Pilot.EmulateCPU()` methods

#### Enhanced Console Debugging ✅ (Direct CDP) - COMPLETED v0.5.2

Source-mapped stack traces via **direct CDP connection**.

- [x] CDP `Runtime.enable` for console events
- [x] CDP `Runtime.consoleAPICalled` events
- [x] `enable_console_debugger` MCP tool
- [x] `get_console_entries_with_stack` MCP tool (full stack traces)
- [x] `get_browser_logs` MCP tool (deprecations, interventions, violations)
- [x] `disable_console_debugger` MCP tool
- [x] SDK methods: `EnableConsoleDebugger`, `ConsoleEntries`, `BrowserLogs`

### P2 - Medium Priority

Nice-to-have features for comprehensive tooling.

#### File Upload ✅ (Clicker has CLI support) - COMPLETED

Upload files - already implemented via `vibium:el.setFiles`.

- [x] `set_files` MCP tool (selector, file_paths)
- [x] SDK `Element.SetFiles()` method
- [x] Multiple file support

#### Drag and Drop ✅ (Clicker has CLI support) - COMPLETED

Drag operations - already implemented via `vibium:dragTo`.

- [x] `drag_to` MCP tool (source_selector, target_selector)
- [x] SDK `Element.DragTo(target)` method

#### Wait for Text ✅ (Already implemented)

Wait for text - already exists in MCP tools (v0.5.0).

- [x] `wait_for_text` MCP tool - **already implemented**
- [ ] Regex pattern support (enhancement)
- [ ] Case-insensitive option (enhancement)

### P3 - Low Priority (Future)

Features for specialized use cases.

#### CrUX Integration

Real User Experience data from Chrome UX Report.

- [ ] `get_crux_data` MCP tool (url)
- [ ] Origin-level metrics
- [ ] Field data vs lab data comparison

#### Coverage Analysis ✅ - COMPLETED v0.5.2

Code coverage for JavaScript and CSS.

- [x] `start_coverage` / `stop_coverage` MCP tools
- [x] JS and CSS coverage reports
- [x] Unused code detection (CSS usage percent)
- [x] SDK methods: `StartCoverage`, `StopCoverage`, `StartJSCoverage`, `StartCSSCoverage`
- [x] CDP `Profiler.startPreciseCoverage` / `CSS.startRuleUsageTracking`

---

## In Progress

### SDK/Clicker Compatibility - v0.7.0

SDK methods must match actual clicker command implementations.

#### Command Name Mismatches

SDK sends wrong command names that clicker doesn't recognize:

| SDK Method | Sends | Clicker Has | Fix |
|------------|-------|-------------|-----|
| `HandleDialog()` | `vibium:dialog.handle` | `dialog.accept`, `dialog.dismiss` | [x] Split into two calls |
| `SetExtraHTTPHeaders()` | `vibium:network.setHeaders` | `vibium:page.setHeaders` | [x] Fix command name |

#### Protocol-Agnostic SDK (BiDi first, CDP fallback)

SDK methods try BiDi first, fall back to CDP when BiDi doesn't support the command:

| SDK Method | BiDi Command | CDP Fallback | Status |
|------------|--------------|--------------|--------|
| `SetOffline()` | `vibium:network.setOffline` | `EmulateNetwork(NetworkOffline)` | [x] Implemented |
| `ConsoleMessages()` | `vibium:console.messages` | `ConsoleEntries()` | [x] Implemented |
| `ClearConsoleMessages()` | `vibium:console.clear` | `consoleDebugger.Clear()` | [x] Implemented |
| `NetworkRequests()` | `vibium:network.requests` | *(needs CDP impl)* | [ ] TODO |
| `ClearNetworkRequests()` | `vibium:network.clearRequests` | *(needs CDP impl)* | [ ] TODO |

#### CDP-Only Features (no BiDi equivalent)

These features only work via CDP and are documented as such:

| Feature | CDP Method | Notes |
|---------|------------|-------|
| `EmulateNetwork()` | `Network.emulateNetworkConditions` | Fine-grained latency/bandwidth control |
| `EmulateCPU()` | `Emulation.setCPUThrottlingRate` | CPU throttling |
| `TakeHeapSnapshot()` | `HeapProfiler.takeHeapSnapshot` | Memory profiling |
| `ConsoleEntries()` | `Runtime.consoleAPICalled` | Full stack traces (richer than ConsoleMessages) |
| Coverage, Screencast, Extensions | Various | See CDP guide |

#### Working but Flaky

Features that exist in clicker but have issues:

| Feature | Issue | Status |
|---------|-------|--------|
| `EmulateMedia()` | CSS changes not persisting | [ ] Investigate |
| `Fill()` | Timeout issues | [ ] Investigate |
| `FindAll` with semantic selectors | Not working | [ ] Investigate |
| `role="alert"`, `role="dialog"` | Not found | [ ] Investigate |

#### Tasks

- [x] Fix `HandleDialog()` to use `dialog.accept`/`dialog.dismiss`
- [x] Fix `SetExtraHTTPHeaders()` command name
- [x] Make SDK protocol-agnostic (BiDi first, CDP fallback)
- [x] Implement `SetOffline()` with CDP fallback
- [x] Implement `ConsoleMessages()` with CDP fallback
- [x] Implement `ClearConsoleMessages()` with CDP fallback
- [x] Re-enable integration tests for CDP-backed features
- [ ] Implement `NetworkRequests()` with CDP fallback
- [ ] Document SDK/clicker compatibility matrix

---

### v0.5.1 - Bug Fixes

- [x] Fix clicker WebSocket transport (pipe.go → clicker.go + transport_ws.go)
- [ ] Verify all `vibium:*` commands work with current clicker version

### v0.6.0 - Chrome DevTools MCP Parity

New features from Chrome DevTools MCP analysis.

#### Screencast ✅ (P2) - COMPLETED v0.5.2

Live screen streaming (from Chrome DevTools MCP).

- [x] `start_screencast` MCP tool
- [x] `stop_screencast` MCP tool
- [x] SDK `Pilot.StartScreencast()` / `Pilot.StopScreencast()` methods
- [x] SDK `Pilot.IsScreencasting()` status method
- [x] CDP `Page.startScreencast` / `Page.stopScreencast` implementation

#### Extensions Management ✅ (P3) - COMPLETED v0.5.2

Browser extension control (from Chrome DevTools MCP).

- [x] `install_extension` MCP tool (path)
- [x] `uninstall_extension` MCP tool (id)
- [x] `list_extensions` MCP tool
- [ ] `trigger_extension_action` MCP tool (id) - stretch goal

---

## Completed - v0.5.2 (2026-03-24)

### Lighthouse Integration

Run quality audits via external lighthouse CLI.

- [x] `lighthouse.go` - LighthouseAudit SDK method
- [x] `LighthouseOptions` (Categories, Device, OutputDir, Port)
- [x] `LighthouseResult` (URL, Scores, PassedAudits, FailedAudits, ReportPaths)
- [x] `findLighthouseBinary()` - locate lighthouse/npx CLI
- [x] `lighthouse_audit` MCP tool
- [x] Support for desktop/mobile device emulation
- [x] JSON and HTML report generation

### Network Request Bodies MCP Tool

- [x] `get_network_request_body` MCP tool
- [x] Retrieve response body by request ID
- [x] Optional save to file for binary content
- [x] Base64 encoding indicator

### Screencast

- [x] `cdp/screencast.go` - Screencast CDP implementation
- [x] `Pilot.StartScreencast()` / `Pilot.StopScreencast()` SDK methods
- [x] `Pilot.IsScreencasting()` status method
- [x] `start_screencast` / `stop_screencast` MCP tools
- [x] Configurable format (jpeg/png), quality, dimensions

### Extensions Management

- [x] `cdp/extensions.go` - Extensions CDP implementation
- [x] `Pilot.InstallExtension()` / `Pilot.UninstallExtension()` SDK methods
- [x] `Pilot.ListExtensions()` SDK method
- [x] `install_extension` / `uninstall_extension` / `list_extensions` MCP tools

### Code Coverage

- [x] `cdp/coverage.go` - Coverage CDP implementation (JS + CSS)
- [x] `Pilot.StartCoverage()` / `Pilot.StopCoverage()` SDK methods
- [x] `Pilot.StartJSCoverage()` / `Pilot.StartCSSCoverage()` SDK methods
- [x] `start_coverage` / `stop_coverage` MCP tools
- [x] Coverage summary with JS scripts, functions, and CSS usage percentage

### Enhanced Console Debugging

- [x] `cdp/debugger.go` - Console debugger CDP implementation
- [x] `Pilot.EnableConsoleDebugger()` / `Pilot.DisableConsoleDebugger()` SDK methods
- [x] `Pilot.ConsoleEntries()` / `Pilot.BrowserLogs()` SDK methods
- [x] `enable_console_debugger` / `disable_console_debugger` MCP tools
- [x] `get_console_entries_with_stack` / `get_browser_logs` MCP tools
- [x] Full stack traces, deprecations, interventions, violations

### CDP Client Infrastructure (Dual Protocol Architecture)

Direct Chrome DevTools Protocol access alongside existing BiDi.

#### CDP Client Package

- [x] `cdp/client.go` - CDP WebSocket client with Send/OnEvent
- [x] `cdp/protocol.go` - CDP message types, network presets, constants
- [x] `cdp/discovery.go` - CDP port discovery from `DevToolsActivePort` file
- [x] `cdp/process.go` - Chrome user-data-dir detection from running processes
- [x] `cdp/heap.go` - HeapProfiler.takeHeapSnapshot implementation
- [x] `cdp/network.go` - Network emulation and response body capture
- [x] `cdp/emulation.go` - CPU throttling implementation

#### Pilot Integration

- [x] Add `cdpClient *cdp.Client` field to `Pilot` struct
- [x] Auto-discover and connect CDP on `Launch()`
- [x] `Pilot.CDP()` for direct CDP access
- [x] `Pilot.HasCDP()` / `Pilot.CDPPort()` status methods
- [x] Graceful degradation if CDP unavailable

#### SDK Methods (CDP-based)

- [x] `Pilot.TakeHeapSnapshot(ctx, path)` - Capture V8 heap snapshot
- [x] `Pilot.EmulateNetwork(ctx, conditions)` - Network throttling (Slow3G, Fast3G, 4G)
- [x] `Pilot.ClearNetworkEmulation(ctx)` - Remove network throttling
- [x] `Pilot.EmulateCPU(ctx, rate)` - CPU throttling (2x, 4x, 6x slowdown)
- [x] `Pilot.ClearCPUEmulation(ctx)` - Remove CPU throttling

#### Network Presets

- [x] `cdp.NetworkSlow3G` - 400ms latency, 400 Kbps
- [x] `cdp.NetworkFast3G` - 150ms latency, 1.5 Mbps
- [x] `cdp.Network4G` - 50ms latency, 4 Mbps
- [x] `cdp.NetworkWifi` - 10ms latency, 30 Mbps
- [x] `cdp.NetworkOffline` - Offline mode

#### CPU Presets

- [x] `cdp.CPUNoThrottle` (1x)
- [x] `cdp.CPU2xSlowdown`
- [x] `cdp.CPU4xSlowdown` (mid-tier mobile)
- [x] `cdp.CPU6xSlowdown` (low-end mobile)

#### Documentation

- [x] README.md - Architecture diagram, CDP features section
- [x] docs/architecture/overview.md - Dual-protocol architecture
- [x] docs/reference/api.md - CDP methods and types
- [x] docs/reference/comparison.md - CDP features comparison
- [x] docs/guide/cdp.md - Comprehensive CDP guide (new)
- [x] mkdocs.yml - Navigation updated

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

- All Vibium/Playwright parity tasks complete as of v0.5.0
- Chrome DevTools MCP analysis added 2026-03-24 (see P0-P3 tasks above)
- Reference: `/Users/johnwang/go/src/github.com/ChromeDevTools/chrome-devtools-mcp`
- Goal: W3Pilot as superset of VibiumCLI-MCP, Playwright MCP, and Chrome DevTools MCP

## Implementation Approach (2026-03-24)

### Dual-Protocol Architecture

**Key Discovery**: BiDi (via clicker) and CDP can connect to the **same browser session**.

Chrome exposes both protocols simultaneously:
- **BiDi**: Via chromedriver WebSocket (clicker uses this)
- **CDP**: Via Chrome's DevTools port (found in `DevToolsActivePort` file)

```
┌─────────────────────────────────────────────────────────────┐
│                      Chrome Browser                          │
│                   (launched by clicker)                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  BiDi WebSocket ◄──── clicker ◄──── W3Pilot (existing)     │
│  (chromedriver)       (serve)       - Semantic selectors     │
│                                     - Navigation, clicks     │
│                                     - Screenshots, a11y      │
│                                                              │
│  CDP WebSocket  ◄──── W3Pilot (new CDPClient)              │
│  (DevTools port)      - Heap snapshots                       │
│                       - Network bodies                       │
│                       - CPU/Network throttling               │
│                       - Coverage analysis                    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Feature Implementation Map

| Feature | Protocol | Implementation |
|---------|----------|----------------|
| Performance Metrics | JS | `Pilot.Evaluate()` with PerformanceObserver |
| Memory Stats | JS | `Pilot.Evaluate()` with `performance.memory` |
| Lighthouse | External | Shell to `npx lighthouse` CLI |
| Heap Snapshots | **CDP** | `HeapProfiler.takeHeapSnapshot` |
| Network Bodies | **CDP** | `Network.getResponseBody` |
| Emulation Presets | **CDP** | `Network.emulateNetworkConditions`, `Emulation.setCPUThrottlingRate` |
| Console Debugging | **CDP** | `Debugger.enable`, `Runtime.consoleAPICalled` |
| Coverage | **CDP** | `Profiler.startPreciseCoverage` |
| File Upload | BiDi | Existing `clicker upload` |
| Drag and Drop | BiDi | Existing `clicker drag` |

### CDP Connection Discovery

When clicker launches Chrome, CDP port is available via:

```
/path/to/user-data-dir/DevToolsActivePort
```

Contents:
```
59726                                              ← CDP port
/devtools/browser/59b41735-8eaa-4112-bd09-...     ← Browser endpoint
```

W3Pilot can read this file to establish parallel CDP connection.

### Interface Design

```go
// Pilot provides unified access to both protocols
type Pilot struct {
    bidi    *BiDiClient    // Via clicker (semantic selectors, navigation)
    cdp     *CDPClient     // Direct to Chrome (profiling, debugging)
    cdpPort int            // Discovered from DevToolsActivePort
}

// BiDi operations (existing)
pilot.Find(ctx, "button", &FindOptions{Role: "button"})
pilot.Click(ctx, selector)
pilot.Screenshot(ctx)

// CDP operations (new)
pilot.TakeHeapSnapshot(ctx)           // CDP HeapProfiler
pilot.GetNetworkResponseBody(ctx, id) // CDP Network
pilot.EmulateNetwork(ctx, Slow3G)     // CDP Network
pilot.EmulateCPU(ctx, 4)              // CDP Emulation
```

### No Clicker Enhancements Required

Previous approach required adding CDP passthrough to clicker. With dual-protocol:
- ✅ W3Pilot connects directly to Chrome's CDP endpoint
- ✅ Same browser session, no coordination issues
- ✅ No upstream dependency for new features
- ✅ Enhancement requests in `docs/enhancement-requests/` are now optional

## Architecture Comparison

| Tool | Protocol | Transport | Strengths |
|------|----------|-----------|-----------|
| W3Pilot | **BiDi + CDP** | Dual WebSocket | Best of both: semantic selectors + full debugging |
| Chrome DevTools MCP | CDP only | Puppeteer | Performance insights, memory profiling |
| Playwright MCP | CDP only | Playwright | Cross-browser, codegen, fixtures |
| VibiumCLI | BiDi only | clicker | Semantic selectors, RPA |

**W3Pilot advantage**: Unified interface to both protocols on single browser.

## Priority Rationale

- **P0**: Unique Chrome DevTools features that provide significant value (Lighthouse, memory, perf insights)
- **P1**: Debugging improvements that enhance developer experience
- **P2**: Quality-of-life features already available in other tools
- **P3**: Specialized features for advanced use cases
