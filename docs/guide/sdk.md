# Go Client SDK

The Go client SDK provides programmatic browser control with full feature parity to the JavaScript and Python clients.

## Installation

```bash
go get github.com/plexusone/webpilot
```

## Basic Usage

```go
package main

import (
    "context"
    "log"

    webpilot "github.com/plexusone/webpilot"
)

func main() {
    ctx := context.Background()

    // Launch browser
    pilot, err := webpilot.Launch(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer pilot.Quit(ctx)

    // Navigate and interact
    pilot.Go(ctx, "https://example.com")
    elem, _ := pilot.Find(ctx, "a", nil)
    elem.Click(ctx, nil)
}
```

## Browser Control

### Launching

```go
// Default launch (visible browser)
pilot, err := webpilot.Launch(ctx)

// Headless launch
pilot, err := webpilot.LaunchHeadless(ctx)

// Custom options
pilot, err := webpilot.Browser.Launch(ctx, &webpilot.LaunchOptions{
    Headless:       true,
    Port:           9515,
    ExecutablePath: "/path/to/clicker",
})
```

### Cleanup

```go
// Close browser
err := pilot.Quit(ctx)

// Check if closed
if pilot.IsClosed() {
    // ...
}
```

## Navigation

```go
// Navigate to URL
err := pilot.Go(ctx, "https://example.com")

// Get current URL
url, err := pilot.URL(ctx)

// Get page title
title, err := pilot.Title(ctx)

// History navigation
err := pilot.Back(ctx)
err := pilot.Forward(ctx)
err := pilot.Reload(ctx)

// Wait for navigation
err := pilot.WaitForNavigation(ctx, 30*time.Second)

// Wait for URL pattern
err := pilot.WaitForURL(ctx, "/dashboard", nil)

// Wait for load state
err := pilot.WaitForLoad(ctx, "networkidle", nil)
```

## Finding Elements

### By CSS Selector

```go
// Find single element
elem, err := pilot.Find(ctx, "button.submit", nil)

// Find with timeout
elem, err := pilot.Find(ctx, "button.submit", &webpilot.FindOptions{
    Timeout: 10 * time.Second,
})

// Find all matching elements
elements, err := pilot.FindAll(ctx, "li.item", nil)
for _, elem := range elements {
    text, _ := elem.Text(ctx)
    fmt.Println(text)
}

// Must find (panics if not found)
elem := pilot.MustFind(ctx, "button.submit")
```

### By Semantic Selectors

Semantic selectors find elements by accessibility attributes instead of brittle CSS selectors. This is especially useful when:

- Page structure may change but semantics remain stable
- Working with AI assistants that describe elements by their purpose
- Writing more maintainable tests

```go
// Find by ARIA role and text
elem, err := pilot.Find(ctx, "", &webpilot.FindOptions{
    Role: "button",
    Text: "Submit",
})

// Find by associated label (great for form inputs)
elem, err := pilot.Find(ctx, "", &webpilot.FindOptions{
    Label: "Email address",
})

// Find by placeholder text
elem, err := pilot.Find(ctx, "", &webpilot.FindOptions{
    Placeholder: "Enter your email",
})

// Find by data-testid (recommended for testing)
elem, err := pilot.Find(ctx, "", &webpilot.FindOptions{
    TestID: "login-button",
})

// Find by image alt text
elem, err := pilot.Find(ctx, "", &webpilot.FindOptions{
    Alt: "Company logo",
})

// Find by title attribute
elem, err := pilot.Find(ctx, "", &webpilot.FindOptions{
    Title: "Close dialog",
})

// Find by XPath
elem, err := pilot.Find(ctx, "", &webpilot.FindOptions{
    XPath: "//button[@type='submit']",
})

// Find element near another element
elem, err := pilot.Find(ctx, "", &webpilot.FindOptions{
    Role: "button",
    Near: "#username-input",
})
```

### Combining Selectors

You can combine CSS selectors with semantic filtering:

```go
// Find within a form, then filter by role and label
elem, err := pilot.Find(ctx, "form.login", &webpilot.FindOptions{
    Role:  "textbox",
    Label: "Password",
})

// Find all buttons within a specific container
buttons, err := pilot.FindAll(ctx, ".dialog-footer", &webpilot.FindOptions{
    Role: "button",
})
```

### Semantic Selector Reference

| Selector | Description | Example Values |
|----------|-------------|----------------|
| `Role` | ARIA role | `"button"`, `"textbox"`, `"link"`, `"checkbox"`, `"menuitem"` |
| `Text` | Visible text content | `"Submit"`, `"Learn more"`, `"Cancel"` |
| `Label` | Associated label text | `"Email address"`, `"Password"`, `"Remember me"` |
| `Placeholder` | Input placeholder | `"Enter email"`, `"Search..."` |
| `TestID` | `data-testid` attribute | `"login-btn"`, `"user-avatar"` |
| `Alt` | Image alt text | `"Company logo"`, `"Profile picture"` |
| `Title` | Element title attribute | `"Close"`, `"More options"` |
| `XPath` | XPath expression | `"//button[@type='submit']"` |
| `Near` | CSS selector of nearby element | `"#username"`, `".form-group"` |

### Scoped Element Search

Find elements within a parent element:

```go
// Find a container first
form, err := pilot.Find(ctx, "form.signup", nil)

// Then find within it
emailInput, err := form.Find(ctx, "", &webpilot.FindOptions{
    Label: "Email",
})

// Find all checkboxes within the form
checkboxes, err := form.FindAll(ctx, "", &webpilot.FindOptions{
    Role: "checkbox",
})
```

## Element Interactions

### Clicking

```go
// Click (waits for actionability)
err := elem.Click(ctx, nil)

// Click with timeout
err := elem.Click(ctx, &webpilot.ActionOptions{
    Timeout: 5 * time.Second,
})

// Double-click
err := elem.DblClick(ctx, nil)
```

### Text Input

```go
// Type text (appends)
err := elem.Type(ctx, "hello", nil)

// Fill text (clears first)
err := elem.Fill(ctx, "hello", nil)

// Clear input
err := elem.Clear(ctx, nil)

// Press key
err := elem.Press(ctx, "Enter", nil)
```

### Form Controls

```go
// Checkbox
err := elem.Check(ctx, nil)
err := elem.Uncheck(ctx, nil)

// Select dropdown
err := elem.SelectOption(ctx, webpilot.SelectOptionValues{
    Values: []string{"option1"},
}, nil)

// File input
err := elem.SetFiles(ctx, []string{"/path/to/file.pdf"}, nil)
```

### Other Interactions

```go
// Hover
err := elem.Hover(ctx, nil)

// Focus
err := elem.Focus(ctx, nil)

// Scroll into view
err := elem.ScrollIntoView(ctx, nil)

// Drag and drop
err := source.DragTo(ctx, target, nil)

// Tap (touch)
err := elem.Tap(ctx, nil)
```

## Element State

```go
// Get text content
text, err := elem.Text(ctx)

// Get input value
value, err := elem.Value(ctx)

// Get innerHTML
html, err := elem.InnerHTML(ctx)

// Get attribute
href, err := elem.GetAttribute(ctx, "href")

// Get bounding box
box, err := elem.BoundingBox(ctx)
// box.X, box.Y, box.Width, box.Height

// State checks
visible, err := elem.IsVisible(ctx)
hidden, err := elem.IsHidden(ctx)
enabled, err := elem.IsEnabled(ctx)
checked, err := elem.IsChecked(ctx)
editable, err := elem.IsEditable(ctx)

// Accessibility
role, err := elem.Role(ctx)
label, err := elem.Label(ctx)

// Wait for state
err := elem.WaitUntil(ctx, "visible", nil)
```

## Input Controllers

### Keyboard

```go
keyboard := pilot.Keyboard()

// Press key
err := keyboard.Press(ctx, "Enter")

// Key down/up
err := keyboard.Down(ctx, "Shift")
err := keyboard.Up(ctx, "Shift")

// Type text
err := keyboard.Type(ctx, "hello world")
```

### Mouse

```go
mouse := pilot.Mouse()

// Click at coordinates
err := mouse.Click(ctx, 100, 200)

// Move mouse
err := mouse.Move(ctx, 100, 200)

// Mouse button
err := mouse.Down(ctx)
err := mouse.Up(ctx)

// Scroll
err := mouse.Wheel(ctx, 0, 100)
```

### Touch

```go
touch := pilot.Touch()

// Tap at coordinates
err := touch.Tap(ctx, 100, 200)
```

## Screenshots and PDF

```go
// Screenshot
data, err := pilot.Screenshot(ctx)
os.WriteFile("page.png", data, 0644)

// Element screenshot
data, err := elem.Screenshot(ctx)

// PDF
data, err := pilot.PDF(ctx, nil)
```

## JavaScript

```go
// Evaluate script
result, err := pilot.Evaluate(ctx, "document.title")

// Evaluate with element
result, err := elem.Eval(ctx, "el => el.textContent")

// Add script tag
err := pilot.AddScript(ctx, "console.log('injected')", nil)

// Add stylesheet
err := pilot.AddStyle(ctx, "body { background: red }", nil)
```

## Page Management

```go
// Create new page
newVibe, err := pilot.NewPage(ctx)

// Get all pages
pages, err := pilot.Pages(ctx)

// Close current page
err := pilot.Close(ctx)

// Bring to front
err := pilot.BringToFront(ctx)

// Get frames
frames, err := pilot.Frames(ctx)

// Get frame by name/URL
frame, err := pilot.Frame(ctx, "iframe-name")
```

## Browser Context

```go
// Create new context (isolated session)
browserCtx, err := pilot.NewContext(ctx)

// Cookies
cookies, err := browserCtx.Cookies(ctx)
err := browserCtx.SetCookies(ctx, cookies)
err := browserCtx.ClearCookies(ctx)

// Storage state (cookies + localStorage only)
state, err := browserCtx.StorageState(ctx)
```

## Storage State

Full storage state management including cookies, localStorage, and sessionStorage:

```go
// Get complete storage state (cookies + localStorage + sessionStorage)
state, err := pilot.StorageState(ctx)

// Save to file for later restoration
jsonBytes, _ := json.Marshal(state)
os.WriteFile("storage.json", jsonBytes, 0600)

// Restore storage state from JSON
var savedState webpilot.StorageState
json.Unmarshal(jsonBytes, &savedState)
err := pilot.SetStorageState(ctx, &savedState)

// Clear all storage (cookies, localStorage, sessionStorage)
err := pilot.ClearStorage(ctx)
```

The `StorageState` type contains:

```go
type StorageState struct {
    Cookies []Cookie             `json:"cookies"`
    Origins []StorageStateOrigin `json:"origins"`
}

type StorageStateOrigin struct {
    Origin         string            `json:"origin"`
    LocalStorage   map[string]string `json:"localStorage"`
    SessionStorage map[string]string `json:"sessionStorage,omitempty"`
}
```

## Init Scripts

Inject JavaScript that runs before any page scripts on every navigation:

```go
// Add script that runs before page loads
err := pilot.AddInitScript(ctx, `window.testMode = true;`)

// Mock APIs before page scripts run
err := pilot.AddInitScript(ctx, `
    const originalFetch = window.fetch;
    window.fetch = async (url, opts) => {
        if (url.includes('/api/user')) {
            return {
                ok: true,
                json: () => Promise.resolve({ id: 1, name: 'Test User' })
            };
        }
        return originalFetch(url, opts);
    };
`)

// Disable analytics
err := pilot.AddInitScript(ctx, `
    window.gtag = () => {};
    window.analytics = { track: () => {} };
`)

// Via BrowserContext for isolated contexts
browserCtx, _ := pilot.NewContext(ctx)
err := browserCtx.AddInitScript(ctx, `window.contextId = 'isolated';`)
```

## Tracing

Record browser actions with screenshots and DOM snapshots:

```go
// Get tracing controller
tracing := pilot.Tracing()

// Start with options
err := tracing.Start(ctx, &webpilot.TracingStartOptions{
    Name:        "login-test",
    Screenshots: true,
    Snapshots:   true,
    Sources:     true,
    Title:       "Login Flow Test",
})

// Perform actions to record
pilot.Go(ctx, "https://example.com")
elem, _ := pilot.Find(ctx, "button", nil)
elem.Click(ctx, nil)

// Stop and get trace data
data, err := tracing.Stop(ctx, nil)
os.WriteFile("trace.zip", data, 0600)

// Use chunks for segmented recording
err := tracing.StartChunk(ctx, &webpilot.TracingChunkOptions{
    Name: "step-1",
})
// ... actions ...
chunkData, err := tracing.StopChunk(ctx, nil)

// Use groups for logical organization
err := tracing.StartGroup(ctx, "Login Flow", nil)
// ... login actions ...
err := tracing.StopGroup(ctx)
```

## Emulation

```go
// Viewport
err := pilot.SetViewport(ctx, webpilot.Viewport{
    Width:  1920,
    Height: 1080,
})

// Media emulation
err := pilot.EmulateMedia(ctx, &webpilot.EmulateMediaOptions{
    Media:       "print",
    ColorScheme: "dark",
})

// Geolocation
err := pilot.SetGeolocation(ctx, &webpilot.Geolocation{
    Latitude:  37.7749,
    Longitude: -122.4194,
})
```

## Error Handling

```go
import "errors"

elem, err := pilot.Find(ctx, "#missing", nil)
if err != nil {
    if errors.Is(err, webpilot.ErrElementNotFound) {
        // Element not found
    }
    if errors.Is(err, webpilot.ErrTimeout) {
        // Timeout
    }
}
```

## Debug Logging

```bash
WEBPILOT_DEBUG=1 go run main.go
```

```go
// Check debug mode
if webpilot.Debug() {
    // ...
}

// Custom logger
logger := webpilot.NewDebugLogger()
ctx = webpilot.ContextWithLogger(ctx, logger)
```
