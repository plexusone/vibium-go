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

	// === Emulation ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "emulate_media",
		Description: "Emulate media features (print, color scheme, etc).",
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
		Name:        "get_storage_state",
		Description: "Get cookies and localStorage as JSON.",
	}, s.handleGetStorageState)

	// === Assertions ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "assert_text",
		Description: "Assert that text exists on the page.",
	}, s.handleAssertText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "assert_element",
		Description: "Assert that an element exists on the page.",
	}, s.handleAssertElement)

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
}

// Run starts the MCP server.
func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

// Close closes the server and browser session.
func (s *Server) Close(ctx context.Context) error {
	return s.session.Close(ctx)
}
