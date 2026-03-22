package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server is the Vibium MCP server.
type Server struct {
	session   *Session
	mcpServer *mcp.Server
	config    Config
}

// NewServer creates a new MCP server.
func NewServer(config Config) *Server {
	s := &Server{
		config: config,
		session: NewSession(SessionConfig{
			Headless:       config.Headless,
			DefaultTimeout: config.DefaultTimeout,
			Project:        config.Project,
			InitScripts:    config.InitScripts,
		}),
	}

	s.mcpServer = mcp.NewServer(
		&mcp.Implementation{
			Name:    "vibium-mcp",
			Version: "0.2.0",
		},
		nil,
	)

	s.registerTools()
	return s
}

// registerTools registers all MCP tools.
func (s *Server) registerTools() {
	// === Browser Management ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "browser_launch",
		Description: "Launch a browser instance. Call this before any other browser operations.",
	}, s.handleBrowserLaunch)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "browser_quit",
		Description: "Close the browser and cleanup resources.",
	}, s.handleBrowserQuit)

	// === Navigation ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "navigate",
		Description: "Navigate to a URL.",
	}, s.handleNavigate)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "back",
		Description: "Navigate back in browser history.",
	}, s.handleBack)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "forward",
		Description: "Navigate forward in browser history.",
	}, s.handleForward)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "reload",
		Description: "Reload the current page.",
	}, s.handleReload)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "scroll",
		Description: "Scroll the page or a specific element in a direction (up, down, left, right).",
	}, s.handleScroll)

	// === Basic Interactions ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "click",
		Description: "Click an element by CSS selector.",
	}, s.handleClick)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "dblclick",
		Description: "Double-click an element by CSS selector.",
	}, s.handleDblClick)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "type",
		Description: "Type text into an input element (appends to existing content).",
	}, s.handleType)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "fill",
		Description: "Clear an input and fill it with text (replaces existing content).",
	}, s.handleFill)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "fill_form",
		Description: "Fill multiple form fields at once. Provide an array of {selector, value} pairs.",
	}, s.handleFillForm)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "clear",
		Description: "Clear the content of an input element.",
	}, s.handleClear)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "press",
		Description: "Press a key on an element (e.g., Enter, Tab, ArrowDown).",
	}, s.handlePress)

	// === Form Controls ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "check",
		Description: "Check a checkbox element.",
	}, s.handleCheck)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "uncheck",
		Description: "Uncheck a checkbox element.",
	}, s.handleUncheck)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "select_option",
		Description: "Select option(s) in a <select> element by value, label, or index.",
	}, s.handleSelectOption)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "set_files",
		Description: "Set files on a file input element.",
	}, s.handleSetFiles)

	// === Element Interaction ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "hover",
		Description: "Hover over an element.",
	}, s.handleHover)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "focus",
		Description: "Focus an element.",
	}, s.handleFocus)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "scroll_into_view",
		Description: "Scroll an element into view.",
	}, s.handleScrollIntoView)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "drag_to",
		Description: "Drag an element to another element.",
	}, s.handleDragTo)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "tap",
		Description: "Tap an element (touch gesture).",
	}, s.handleTap)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "dispatch_event",
		Description: "Dispatch a DOM event on an element.",
	}, s.handleDispatchEvent)

	// === Element State ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_text",
		Description: "Get the text content of an element.",
	}, s.handleGetText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_value",
		Description: "Get the value of an input element.",
	}, s.handleGetValue)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_inner_html",
		Description: "Get the innerHTML of an element.",
	}, s.handleGetInnerHTML)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_outer_html",
		Description: "Get the outerHTML of an element (including the element itself).",
	}, s.handleGetOuterHTML)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_inner_text",
		Description: "Get the innerText of an element.",
	}, s.handleGetInnerText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_attribute",
		Description: "Get an attribute value of an element.",
	}, s.handleGetAttribute)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_bounding_box",
		Description: "Get the bounding box of an element.",
	}, s.handleGetBoundingBox)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "is_visible",
		Description: "Check if an element is visible.",
	}, s.handleIsVisible)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "is_hidden",
		Description: "Check if an element is hidden.",
	}, s.handleIsHidden)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "is_enabled",
		Description: "Check if an element is enabled.",
	}, s.handleIsEnabled)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "is_checked",
		Description: "Check if a checkbox/radio is checked.",
	}, s.handleIsChecked)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "is_editable",
		Description: "Check if an element is editable.",
	}, s.handleIsEditable)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_role",
		Description: "Get the ARIA role of an element.",
	}, s.handleGetRole)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_label",
		Description: "Get the accessible label of an element.",
	}, s.handleGetLabel)

	// === Page State ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_title",
		Description: "Get the current page title.",
	}, s.handleGetTitle)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_url",
		Description: "Get the current page URL.",
	}, s.handleGetURL)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_content",
		Description: "Get the full HTML content of the page.",
	}, s.handleGetContent)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "set_content",
		Description: "Set the HTML content of the page.",
	}, s.handleSetContent)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_viewport",
		Description: "Get the viewport dimensions.",
	}, s.handleGetViewport)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "set_viewport",
		Description: "Set the viewport dimensions.",
	}, s.handleSetViewport)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_frames",
		Description: "Get all frames on the page.",
	}, s.handleGetFrames)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "select_frame",
		Description: "Switch to a frame by name or URL pattern. Subsequent commands will target this frame.",
	}, s.handleSelectFrame)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "select_main_frame",
		Description: "Switch back to the main frame (top-level page).",
	}, s.handleSelectMainFrame)

	// === Screenshots & PDF ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "screenshot",
		Description: "Capture a screenshot of the current page.",
	}, s.handleScreenshot)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_screenshot",
		Description: "Capture a screenshot of a specific element.",
	}, s.handleElementScreenshot)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "pdf",
		Description: "Generate a PDF of the page.",
	}, s.handlePDF)

	// === JavaScript ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "evaluate",
		Description: "Execute JavaScript on the page and return the result.",
	}, s.handleEvaluate)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_eval",
		Description: "Evaluate JavaScript with an element as the first argument.",
	}, s.handleElementEval)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "add_script",
		Description: "Inject JavaScript into the page.",
	}, s.handleAddScript)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "add_style",
		Description: "Inject CSS into the page.",
	}, s.handleAddStyle)

	// === Waiting ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_until",
		Description: "Wait for an element to reach a state (attached, detached, visible, hidden).",
	}, s.handleWaitUntil)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_url",
		Description: "Wait for the URL to match a pattern.",
	}, s.handleWaitForURL)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_load",
		Description: "Wait for page load state (load, domcontentloaded, networkidle).",
	}, s.handleWaitForLoad)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_function",
		Description: "Wait for a JavaScript function to return truthy.",
	}, s.handleWaitForFunction)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_text",
		Description: "Wait for text to appear on the page. Optionally scope to a specific element.",
	}, s.handleWaitForText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "accessibility_snapshot",
		Description: "Get an accessibility tree snapshot of the page. Useful for understanding page structure and testing accessibility.",
	}, s.handleAccessibilitySnapshot)

	// === Input Controllers ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "keyboard_press",
		Description: "Press a key on the keyboard.",
	}, s.handleKeyboardPress)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "keyboard_down",
		Description: "Hold down a key.",
	}, s.handleKeyboardDown)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "keyboard_up",
		Description: "Release a held key.",
	}, s.handleKeyboardUp)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "keyboard_type",
		Description: "Type text using the keyboard.",
	}, s.handleKeyboardType)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "mouse_click",
		Description: "Click at coordinates.",
	}, s.handleMouseClick)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "mouse_move",
		Description: "Move the mouse to coordinates.",
	}, s.handleMouseMove)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "mouse_down",
		Description: "Press the mouse button.",
	}, s.handleMouseDown)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "mouse_up",
		Description: "Release the mouse button.",
	}, s.handleMouseUp)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "mouse_wheel",
		Description: "Scroll the mouse wheel.",
	}, s.handleMouseWheel)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "touch_tap",
		Description: "Tap at coordinates (touch).",
	}, s.handleTouchTap)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "touch_swipe",
		Description: "Swipe from one point to another (touch).",
	}, s.handleTouchSwipe)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "mouse_drag",
		Description: "Drag from one point to another using the mouse.",
	}, s.handleMouseDrag)

	// === Page Management ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "new_page",
		Description: "Create a new page/tab.",
	}, s.handleNewPage)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_pages",
		Description: "Get the number of open pages.",
	}, s.handleGetPages)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "close_page",
		Description: "Close the current page.",
	}, s.handleClosePage)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "bring_to_front",
		Description: "Bring the page to the front.",
	}, s.handleBringToFront)

	// === Tab Management ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "list_tabs",
		Description: "List all open browser tabs with their index, ID, URL, and title.",
	}, s.handleListTabs)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "select_tab",
		Description: "Switch to a specific tab by index (0-based) or tab ID.",
	}, s.handleSelectTab)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "close_tab",
		Description: "Close a specific tab by index or ID. Defaults to current tab if not specified.",
	}, s.handleCloseTab)

	// === Emulation ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "emulate_media",
		Description: "Emulate CSS media features for accessibility testing: color scheme (dark/light mode), reduced motion (disable animations), forced colors (high contrast mode), and contrast preferences.",
	}, s.handleEmulateMedia)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "set_geolocation",
		Description: "Set the browser's geolocation.",
	}, s.handleSetGeolocation)

	// === Cookies & Storage ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_cookies",
		Description: "Get browser cookies.",
	}, s.handleGetCookies)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "set_cookies",
		Description: "Set browser cookies.",
	}, s.handleSetCookies)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "clear_cookies",
		Description: "Clear all cookies.",
	}, s.handleClearCookies)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "delete_cookie",
		Description: "Delete a specific cookie by name. Optionally filter by domain and path.",
	}, s.handleDeleteCookie)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_storage_state",
		Description: "Get complete browser storage state (cookies, localStorage, and sessionStorage) as JSON.",
	}, s.handleGetStorageState)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "set_storage_state",
		Description: "Restore browser storage from JSON (output of get_storage_state). Restores cookies, localStorage, and sessionStorage.",
	}, s.handleSetStorageState)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "clear_storage",
		Description: "Clear all browser storage (cookies, localStorage, and sessionStorage).",
	}, s.handleClearStorage)

	// === LocalStorage ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "localstorage_get",
		Description: "Get a value from localStorage by key.",
	}, s.handleLocalStorageGet)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "localstorage_set",
		Description: "Set a value in localStorage.",
	}, s.handleLocalStorageSet)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "localstorage_delete",
		Description: "Delete a key from localStorage.",
	}, s.handleLocalStorageDelete)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "localstorage_clear",
		Description: "Clear all localStorage data for the current origin.",
	}, s.handleLocalStorageClear)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "localstorage_list",
		Description: "List all keys and values in localStorage.",
	}, s.handleLocalStorageList)

	// === SessionStorage ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "sessionstorage_get",
		Description: "Get a value from sessionStorage by key.",
	}, s.handleSessionStorageGet)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "sessionstorage_set",
		Description: "Set a value in sessionStorage.",
	}, s.handleSessionStorageSet)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "sessionstorage_delete",
		Description: "Delete a key from sessionStorage.",
	}, s.handleSessionStorageDelete)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "sessionstorage_clear",
		Description: "Clear all sessionStorage data for the current origin.",
	}, s.handleSessionStorageClear)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "sessionstorage_list",
		Description: "List all keys and values in sessionStorage.",
	}, s.handleSessionStorageList)

	// === Dialog Handling ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "handle_dialog",
		Description: "Handle a browser dialog (alert, confirm, prompt, beforeunload) by accepting or dismissing it.",
	}, s.handleHandleDialog)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_dialog",
		Description: "Get information about the current dialog, if any is open.",
	}, s.handleGetDialog)

	// === Console Messages ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_console_messages",
		Description: "Get console messages from the page. Optionally filter by level (log, info, warn, error, debug).",
	}, s.handleGetConsoleMessages)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "clear_console_messages",
		Description: "Clear the buffered console messages.",
	}, s.handleClearConsoleMessages)

	// === Network Requests ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_network_requests",
		Description: "Get captured network requests. Optionally filter by URL pattern, HTTP method, or resource type.",
	}, s.handleGetNetworkRequests)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "clear_network_requests",
		Description: "Clear the buffered network requests.",
	}, s.handleClearNetworkRequests)

	// === Network Mocking ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "route",
		Description: "Register a mock response for requests matching a URL pattern. Use glob patterns (e.g., **/api/*) or regex (e.g., /api/.*).",
	}, s.handleRoute)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "route_list",
		Description: "List all active route handlers.",
	}, s.handleRouteList)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "unroute",
		Description: "Remove a previously registered route handler.",
	}, s.handleUnroute)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "network_state_set",
		Description: "Set the browser's network state. Use offline=true to simulate offline mode for testing.",
	}, s.handleNetworkStateSet)

	// === Human-in-the-Loop ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "pause_for_human",
		Description: "Pause automation and wait for human to complete an action (e.g., SSO login, CAPTCHA). Shows a visual overlay that the human dismisses when done.",
	}, s.handlePauseForHuman)

	// === Assertions ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "assert_text",
		Description: "Assert that text exists on the page.",
	}, s.handleAssertText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "assert_element",
		Description: "Assert that an element exists on the page.",
	}, s.handleAssertElement)

	// === Testing Tools ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "verify_value",
		Description: "Verify that an input element has the expected value.",
	}, s.handleVerifyValue)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "verify_list_visible",
		Description: "Verify that a list of text items are all visible on the page.",
	}, s.handleVerifyListVisible)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "generate_locator",
		Description: "Generate a locator string for a given element using a specific strategy (css, xpath, testid, role, text).",
	}, s.handleGenerateLocator)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_selector",
		Description: "Wait for an element to reach a specific state (attached, detached, visible, hidden).",
	}, s.handleWaitForSelector)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "verify_text",
		Description: "Verify that an element's text content matches the expected value.",
	}, s.handleVerifyText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "verify_visible",
		Description: "Verify that an element is visible on the page.",
	}, s.handleVerifyVisible)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "verify_enabled",
		Description: "Verify that an element is enabled (not disabled).",
	}, s.handleVerifyEnabled)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "verify_checked",
		Description: "Verify that a checkbox or radio button is checked or unchecked.",
	}, s.handleVerifyChecked)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "verify_hidden",
		Description: "Verify that an element is hidden (not visible) on the page.",
	}, s.handleVerifyHidden)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "verify_disabled",
		Description: "Verify that an element is disabled.",
	}, s.handleVerifyDisabled)

	// === Test Reporting ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_test_report",
		Description: "Get the test execution report in the specified format (box, diagnostic, or json).",
	}, s.handleGetTestReport)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "reset_session",
		Description: "Clear test results and start a new test session.",
	}, s.handleResetSession)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "set_target",
		Description: "Set the test target description for reports.",
	}, s.handleSetTarget)

	// === Script Recording ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "start_recording",
		Description: "Start recording browser actions to create a replayable test script.",
	}, s.handleStartRecording)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "stop_recording",
		Description: "Stop recording browser actions.",
	}, s.handleStopRecording)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "export_script",
		Description: "Export recorded actions as a JSON test script that can be run with 'vibium run'.",
	}, s.handleExportScript)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "recording_status",
		Description: "Check if recording is active and how many steps have been recorded.",
	}, s.handleRecordingStatus)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "clear_recording",
		Description: "Clear all recorded steps without stopping recording.",
	}, s.handleClearRecording)

	// === Tracing ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "start_trace",
		Description: "Start trace recording with screenshots and DOM snapshots for debugging. The trace can be viewed with 'npx playwright show-trace <trace.zip>'.",
	}, s.handleStartTrace)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "stop_trace",
		Description: "Stop trace recording and save or return the trace data as a ZIP file.",
	}, s.handleStopTrace)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "start_trace_chunk",
		Description: "Start a new trace chunk within an active trace. Useful for segmenting traces into logical sections.",
	}, s.handleStartTraceChunk)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "stop_trace_chunk",
		Description: "Stop the current trace chunk and optionally save it to a file.",
	}, s.handleStopTraceChunk)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "start_trace_group",
		Description: "Start a trace group for logical grouping of actions in the trace viewer.",
	}, s.handleStartTraceGroup)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "stop_trace_group",
		Description: "Stop the current trace group.",
	}, s.handleStopTraceGroup)

	// === Video Recording ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "start_video",
		Description: "Start recording video of the browser page. Video is saved when stop_video is called.",
	}, s.handleStartVideo)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "stop_video",
		Description: "Stop video recording and return the path to the video file.",
	}, s.handleStopVideo)

	// === Init Scripts ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "add_init_script",
		Description: "Add JavaScript that runs before page scripts on every navigation. Useful for mocking APIs, injecting test helpers, or setting up authentication.",
	}, s.handleAddInitScript)

	// === Configuration ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_config",
		Description: "Get the resolved MCP server configuration including headless mode, project name, and timeouts.",
	}, s.handleGetConfig)
}

// Run starts the MCP server.
func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

// Close closes the server and browser session.
func (s *Server) Close(ctx context.Context) error {
	return s.session.Close(ctx)
}
