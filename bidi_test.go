package webpilot

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
)

// mockTransport records all calls for verification.
type mockTransport struct {
	mu       sync.Mutex
	calls    []mockCall
	handlers map[string][]EventHandler

	// Response to return for Send calls
	response json.RawMessage
	err      error
}

type mockCall struct {
	Method string
	Params interface{}
}

func newMockTransport() *mockTransport {
	return &mockTransport{
		handlers: make(map[string][]EventHandler),
		// Default response for most calls
		response: json.RawMessage(`{}`),
	}
}

func (m *mockTransport) Send(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, mockCall{Method: method, Params: params})
	return m.response, m.err
}

func (m *mockTransport) OnEvent(method string, handler EventHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[method] = append(m.handlers[method], handler)
}

func (m *mockTransport) RemoveEventHandlers(method string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.handlers, method)
}

func (m *mockTransport) Close() error {
	return nil
}

func (m *mockTransport) getCalls() []mockCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]mockCall, len(m.calls))
	copy(result, m.calls)
	return result
}

func (m *mockTransport) setResponse(resp json.RawMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.response = resp
}

// TestPilotFind_SendsVibiumPageFind verifies that Find sends the vibium:page.find method.
func TestPilotFind_SendsVibiumPageFind(t *testing.T) {
	mock := newMockTransport()
	// Response with element info
	mock.setResponse(json.RawMessage(`{"tag":"button","text":"Click me","box":{"x":10,"y":20,"width":100,"height":40}}`))

	client := NewBiDiClient(mock)
	pilot := &Pilot{
		client:          client,
		browsingContext: "ctx-123",
	}

	ctx := context.Background()
	_, err := pilot.Find(ctx, "button", nil)
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	calls := mock.getCalls()
	if len(calls) == 0 {
		t.Fatal("Expected at least one call")
	}

	// Find the vibium:page.find call
	found := false
	for _, call := range calls {
		if call.Method == "vibium:page.find" {
			found = true
			// Verify params contain selector
			params, ok := call.Params.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected map params, got %T", call.Params)
			}
			if params["selector"] != "button" {
				t.Errorf("Expected selector 'button', got %v", params["selector"])
			}
			if params["context"] != "ctx-123" {
				t.Errorf("Expected context 'ctx-123', got %v", params["context"])
			}
			break
		}
	}

	if !found {
		t.Errorf("Expected vibium:page.find call, got: %v", calls)
	}
}

// TestPilotFindAll_SendsVibiumPageFindAll verifies that FindAll sends vibium:page.findAll.
func TestPilotFindAll_SendsVibiumPageFindAll(t *testing.T) {
	mock := newMockTransport()
	mock.setResponse(json.RawMessage(`[]`))

	client := NewBiDiClient(mock)
	pilot := &Pilot{
		client:          client,
		browsingContext: "ctx-123",
	}

	ctx := context.Background()
	_, err := pilot.FindAll(ctx, "button", nil)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	calls := mock.getCalls()
	found := false
	for _, call := range calls {
		if call.Method == "vibium:page.findAll" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected vibium:page.findAll call, got: %v", calls)
	}
}

// TestPilotContent_SendsVibiumPageContent verifies Content sends vibium:page.content.
func TestPilotContent_SendsVibiumPageContent(t *testing.T) {
	mock := newMockTransport()
	mock.setResponse(json.RawMessage(`{"content":"<html></html>"}`))

	client := NewBiDiClient(mock)
	pilot := &Pilot{
		client:          client,
		browsingContext: "ctx-123",
	}

	ctx := context.Background()
	_, err := pilot.Content(ctx)
	if err != nil {
		t.Fatalf("Content failed: %v", err)
	}

	calls := mock.getCalls()
	found := false
	for _, call := range calls {
		if call.Method == "vibium:page.content" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected vibium:page.content call, got: %v", calls)
	}
}

// TestPilotSetContent_SendsVibiumPageSetContent verifies SetContent sends correct method.
func TestPilotSetContent_SendsVibiumPageSetContent(t *testing.T) {
	mock := newMockTransport()

	client := NewBiDiClient(mock)
	pilot := &Pilot{
		client:          client,
		browsingContext: "ctx-123",
	}

	ctx := context.Background()
	err := pilot.SetContent(ctx, "<h1>Hello</h1>")
	if err != nil {
		t.Fatalf("SetContent failed: %v", err)
	}

	calls := mock.getCalls()
	found := false
	for _, call := range calls {
		if call.Method == "vibium:page.setContent" {
			found = true
			params, ok := call.Params.(map[string]interface{})
			if ok && params["html"] != "<h1>Hello</h1>" {
				t.Errorf("Expected html param, got %v", params)
			}
			break
		}
	}

	if !found {
		t.Errorf("Expected vibium:page.setContent call, got: %v", calls)
	}
}

// TestPilotScroll_SendsVibiumPageScroll verifies Scroll sends correct method.
func TestPilotScroll_SendsVibiumPageScroll(t *testing.T) {
	mock := newMockTransport()

	client := NewBiDiClient(mock)
	pilot := &Pilot{
		client:          client,
		browsingContext: "ctx-123",
	}

	ctx := context.Background()
	err := pilot.Scroll(ctx, "down", 500, nil)
	if err != nil {
		t.Fatalf("Scroll failed: %v", err)
	}

	calls := mock.getCalls()
	found := false
	for _, call := range calls {
		if call.Method == "vibium:page.scroll" {
			found = true
			params, ok := call.Params.(map[string]interface{})
			if ok {
				if params["direction"] != "down" {
					t.Errorf("Expected direction 'down', got %v", params["direction"])
				}
				if params["amount"] != 500 {
					t.Errorf("Expected amount 500, got %v", params["amount"])
				}
			}
			break
		}
	}

	if !found {
		t.Errorf("Expected vibium:page.scroll call, got: %v", calls)
	}
}

// TestPilotEmulateMedia_SendsVibiumPageEmulateMedia verifies EmulateMedia sends correct method.
func TestPilotEmulateMedia_SendsVibiumPageEmulateMedia(t *testing.T) {
	mock := newMockTransport()

	client := NewBiDiClient(mock)
	pilot := &Pilot{
		client:          client,
		browsingContext: "ctx-123",
	}

	ctx := context.Background()
	err := pilot.EmulateMedia(ctx, EmulateMediaOptions{ColorScheme: "dark"})
	if err != nil {
		t.Fatalf("EmulateMedia failed: %v", err)
	}

	calls := mock.getCalls()
	found := false
	for _, call := range calls {
		if call.Method == "vibium:page.emulateMedia" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected vibium:page.emulateMedia call, got: %v", calls)
	}
}

// TestElement_Click_SendsVibiumElementClick verifies Element.Click sends vibium:element.click.
func TestElement_Click_SendsVibiumElementClick(t *testing.T) {
	mock := newMockTransport()

	client := NewBiDiClient(mock)
	elem := NewElement(client, "ctx-123", "button.submit", ElementInfo{
		Tag:  "button",
		Text: "Submit",
		Box:  BoundingBox{X: 100, Y: 200, Width: 80, Height: 30},
	})

	ctx := context.Background()
	err := elem.Click(ctx, nil)
	if err != nil {
		t.Fatalf("Click failed: %v", err)
	}

	calls := mock.getCalls()
	found := false
	for _, call := range calls {
		if call.Method == "vibium:element.click" {
			found = true
			params, ok := call.Params.(map[string]interface{})
			if ok && params["selector"] != "button.submit" {
				t.Errorf("Expected selector 'button.submit', got %v", params["selector"])
			}
			break
		}
	}

	if !found {
		t.Errorf("Expected vibium:element.click call, got: %v", calls)
	}
}

// TestElement_Type_SendsVibiumElementType verifies Element.Type sends vibium:element.type.
func TestElement_Type_SendsVibiumElementType(t *testing.T) {
	mock := newMockTransport()

	client := NewBiDiClient(mock)
	elem := NewElement(client, "ctx-123", "input[name='email']", ElementInfo{
		Tag: "input",
	})

	ctx := context.Background()
	err := elem.Type(ctx, "test@example.com", nil)
	if err != nil {
		t.Fatalf("Type failed: %v", err)
	}

	calls := mock.getCalls()
	found := false
	for _, call := range calls {
		if call.Method == "vibium:element.type" {
			found = true
			params, ok := call.Params.(map[string]interface{})
			if ok {
				if params["text"] != "test@example.com" {
					t.Errorf("Expected text 'test@example.com', got %v", params["text"])
				}
			}
			break
		}
	}

	if !found {
		t.Errorf("Expected vibium:element.type call, got: %v", calls)
	}
}

// TestElement_Fill_SendsVibiumElementFill verifies Element.Fill sends vibium:element.fill.
func TestElement_Fill_SendsVibiumElementFill(t *testing.T) {
	mock := newMockTransport()

	client := NewBiDiClient(mock)
	elem := NewElement(client, "ctx-123", "input[name='name']", ElementInfo{
		Tag: "input",
	})

	ctx := context.Background()
	err := elem.Fill(ctx, "John Doe", nil)
	if err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	calls := mock.getCalls()
	found := false
	for _, call := range calls {
		if call.Method == "vibium:element.fill" {
			found = true
			params, ok := call.Params.(map[string]interface{})
			if ok && params["value"] != "John Doe" {
				t.Errorf("Expected value 'John Doe', got %v", params["value"])
			}
			break
		}
	}

	if !found {
		t.Errorf("Expected vibium:element.fill call, got: %v", calls)
	}
}

// TestVibiumMethodPrefix_AllVibiumCommandsHavePrefix is a meta-test that documents
// which methods should use vibium: prefix vs standard BiDi methods.
func TestVibiumMethodPrefix_DocumentedMethods(t *testing.T) {
	// This test documents which Pilot/Element methods use vibium: commands.
	// If the protocol changes, update this list.
	vibiumMethods := []string{
		// Page finding
		"vibium:page.find",
		"vibium:page.findAll",

		// Element finding (child elements)
		"vibium:element.find",
		"vibium:element.findAll",

		// Element actions
		"vibium:element.click",
		"vibium:element.type",
		"vibium:element.fill",
		"vibium:element.press",
		"vibium:element.clear",
		"vibium:element.check",
		"vibium:element.uncheck",
		"vibium:element.selectOption",
		"vibium:element.focus",
		"vibium:element.hover",
		"vibium:element.scrollIntoView",
		"vibium:element.dblclick",
		"vibium:element.dragTo",
		"vibium:element.tap",
		"vibium:element.dispatchEvent",

		// Element state
		"vibium:element.text",
		"vibium:element.value",
		"vibium:element.attr",
		"vibium:element.bounds",
		"vibium:element.html",
		"vibium:element.outerHTML",
		"vibium:element.innerText",
		"vibium:element.isVisible",
		"vibium:element.isHidden",
		"vibium:element.isEnabled",
		"vibium:element.isChecked",
		"vibium:element.isEditable",
		"vibium:element.role",
		"vibium:element.label",
		"vibium:element.waitFor",
		"vibium:element.setFiles",
		"vibium:element.screenshot",
		"vibium:element.eval",
		"vibium:element.highlight",

		// Page methods
		"vibium:page.content",
		"vibium:page.setContent",
		"vibium:page.viewport",
		"vibium:page.setViewport",
		"vibium:page.window",
		"vibium:page.setWindow",
		"vibium:page.pdf",
		"vibium:page.frames",
		"vibium:page.frame",
		"vibium:page.a11yTree",
		"vibium:page.emulateMedia",
		"vibium:page.setGeolocation",
		"vibium:page.addScript",
		"vibium:page.addStyle",
		"vibium:page.expose",
		"vibium:page.waitForURL",
		"vibium:page.waitForLoad",
		"vibium:page.waitForFunction",
		"vibium:page.scroll",

		// Network
		"vibium:network.route",
		"vibium:network.unroute",
		"vibium:network.mockRoute",
		"vibium:network.listRoutes",
		"vibium:network.setOffline",
		"vibium:network.setHeaders",
		"vibium:network.onRequest",
		"vibium:network.onResponse",
		"vibium:network.requests",
		"vibium:network.clearRequests",
		"vibium:network.fulfill",
		"vibium:network.continue",
		"vibium:network.abort",

		// Console
		"vibium:console.on",
		"vibium:console.collect",
		"vibium:console.messages",
		"vibium:console.clear",

		// Dialog
		"vibium:dialog.on",
		"vibium:dialog.handle",
		"vibium:dialog.get",

		// Download
		"vibium:download.on",
		"vibium:download.path",
		"vibium:download.saveAs",
		"vibium:download.cancel",
		"vibium:download.failure",

		// Page errors
		"vibium:page.onError",
		"vibium:page.collectErrors",
		"vibium:page.errors",
		"vibium:page.clearErrors",

		// Context
		"vibium:context.addInitScript",
		"vibium:context.storageState",
		"vibium:context.grantPermissions",
		"vibium:context.clearPermissions",

		// Clock
		"vibium:clock.install",
		"vibium:clock.fastForward",
		"vibium:clock.runFor",
		"vibium:clock.pauseAt",
		"vibium:clock.resume",
		"vibium:clock.setFixedTime",
		"vibium:clock.setSystemTime",
		"vibium:clock.setTimezone",

		// Video
		"vibium:video.start",
		"vibium:video.stop",
		"vibium:video.delete",

		// Keyboard
		"vibium:keyboard.press",
		"vibium:keyboard.down",
		"vibium:keyboard.up",
		"vibium:keyboard.type",
		"vibium:keyboard.insertText",

		// Mouse
		"vibium:mouse.click",
		"vibium:mouse.move",
		"vibium:mouse.down",
		"vibium:mouse.up",
		"vibium:mouse.wheel",

		// Touch
		"vibium:touch.tap",
		"vibium:touch.swipe",
		"vibium:touch.pinch",

		// Tracing
		"vibium:tracing.start",
		"vibium:tracing.stop",
		"vibium:tracing.startChunk",
		"vibium:tracing.stopChunk",
		"vibium:tracing.startGroup",
		"vibium:tracing.stopGroup",

		// Lifecycle
		"vibium:lifecycle.ready",
	}

	// Standard BiDi methods (no vibium: prefix)
	standardBiDiMethods := []string{
		"browsingContext.navigate",
		"browsingContext.reload",
		"browsingContext.traverseHistory",
		"browsingContext.captureScreenshot",
		"browsingContext.getTree",
		"browsingContext.activate",
		"browsingContext.close",
		"browsingContext.create",
		"browser.createUserContext",
		"browser.removeUserContext",
		"browser.getUserContexts",
		"script.callFunction",
		"session.subscribe",
		"storage.getCookies",
		"storage.setCookie",
		"storage.deleteCookies",
	}

	t.Logf("Documented %d vibium: methods", len(vibiumMethods))
	t.Logf("Documented %d standard BiDi methods", len(standardBiDiMethods))

	// This is a documentation test - it always passes but serves as a reference
}
