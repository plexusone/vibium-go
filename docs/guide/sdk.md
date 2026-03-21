# Go Client SDK

The Go client SDK provides programmatic browser control with full feature parity to the JavaScript and Python clients.

## Installation

```bash
go get github.com/plexusone/vibium-go
```

## Basic Usage

```go
package main

import (
    "context"
    "log"

    vibium "github.com/plexusone/vibium-go"
)

func main() {
    ctx := context.Background()

    // Launch browser
    vibe, err := vibium.Launch(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer vibe.Quit(ctx)

    // Navigate and interact
    vibe.Go(ctx, "https://example.com")
    elem, _ := vibe.Find(ctx, "a", nil)
    elem.Click(ctx, nil)
}
```

## Browser Control

### Launching

```go
// Default launch (visible browser)
vibe, err := vibium.Launch(ctx)

// Headless launch
vibe, err := vibium.LaunchHeadless(ctx)

// Custom options
vibe, err := vibium.Browser.Launch(ctx, &vibium.LaunchOptions{
    Headless:       true,
    Port:           9515,
    ExecutablePath: "/path/to/clicker",
})
```

### Cleanup

```go
// Close browser
err := vibe.Quit(ctx)

// Check if closed
if vibe.IsClosed() {
    // ...
}
```

## Navigation

```go
// Navigate to URL
err := vibe.Go(ctx, "https://example.com")

// Get current URL
url, err := vibe.URL(ctx)

// Get page title
title, err := vibe.Title(ctx)

// History navigation
err := vibe.Back(ctx)
err := vibe.Forward(ctx)
err := vibe.Reload(ctx)

// Wait for navigation
err := vibe.WaitForNavigation(ctx, 30*time.Second)

// Wait for URL pattern
err := vibe.WaitForURL(ctx, "/dashboard", nil)

// Wait for load state
err := vibe.WaitForLoad(ctx, "networkidle", nil)
```

## Finding Elements

### By CSS Selector

```go
// Find single element
elem, err := vibe.Find(ctx, "button.submit", nil)

// Find with timeout
elem, err := vibe.Find(ctx, "button.submit", &vibium.FindOptions{
    Timeout: 10 * time.Second,
})

// Find all matching elements
elements, err := vibe.FindAll(ctx, "li.item", nil)
for _, elem := range elements {
    text, _ := elem.Text(ctx)
    fmt.Println(text)
}

// Must find (panics if not found)
elem := vibe.MustFind(ctx, "button.submit")
```

### By Semantic Selectors

Semantic selectors find elements by accessibility attributes instead of brittle CSS selectors. This is especially useful when:

- Page structure may change but semantics remain stable
- Working with AI assistants that describe elements by their purpose
- Writing more maintainable tests

```go
// Find by ARIA role and text
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    Role: "button",
    Text: "Submit",
})

// Find by associated label (great for form inputs)
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    Label: "Email address",
})

// Find by placeholder text
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    Placeholder: "Enter your email",
})

// Find by data-testid (recommended for testing)
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    TestID: "login-button",
})

// Find by image alt text
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    Alt: "Company logo",
})

// Find by title attribute
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    Title: "Close dialog",
})

// Find by XPath
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    XPath: "//button[@type='submit']",
})

// Find element near another element
elem, err := vibe.Find(ctx, "", &vibium.FindOptions{
    Role: "button",
    Near: "#username-input",
})
```

### Combining Selectors

You can combine CSS selectors with semantic filtering:

```go
// Find within a form, then filter by role and label
elem, err := vibe.Find(ctx, "form.login", &vibium.FindOptions{
    Role:  "textbox",
    Label: "Password",
})

// Find all buttons within a specific container
buttons, err := vibe.FindAll(ctx, ".dialog-footer", &vibium.FindOptions{
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
form, err := vibe.Find(ctx, "form.signup", nil)

// Then find within it
emailInput, err := form.Find(ctx, "", &vibium.FindOptions{
    Label: "Email",
})

// Find all checkboxes within the form
checkboxes, err := form.FindAll(ctx, "", &vibium.FindOptions{
    Role: "checkbox",
})
```

## Element Interactions

### Clicking

```go
// Click (waits for actionability)
err := elem.Click(ctx, nil)

// Click with timeout
err := elem.Click(ctx, &vibium.ActionOptions{
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
err := elem.SelectOption(ctx, vibium.SelectOptionValues{
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
keyboard := vibe.Keyboard()

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
mouse := vibe.Mouse()

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
touch := vibe.Touch()

// Tap at coordinates
err := touch.Tap(ctx, 100, 200)
```

## Screenshots and PDF

```go
// Screenshot
data, err := vibe.Screenshot(ctx)
os.WriteFile("page.png", data, 0644)

// Element screenshot
data, err := elem.Screenshot(ctx)

// PDF
data, err := vibe.PDF(ctx, nil)
```

## JavaScript

```go
// Evaluate script
result, err := vibe.Evaluate(ctx, "document.title")

// Evaluate with element
result, err := elem.Eval(ctx, "el => el.textContent")

// Add script tag
err := vibe.AddScript(ctx, "console.log('injected')", nil)

// Add stylesheet
err := vibe.AddStyle(ctx, "body { background: red }", nil)
```

## Page Management

```go
// Create new page
newVibe, err := vibe.NewPage(ctx)

// Get all pages
pages, err := vibe.Pages(ctx)

// Close current page
err := vibe.Close(ctx)

// Bring to front
err := vibe.BringToFront(ctx)

// Get frames
frames, err := vibe.Frames(ctx)

// Get frame by name/URL
frame, err := vibe.Frame(ctx, "iframe-name")
```

## Browser Context

```go
// Create new context (isolated session)
browserCtx, err := vibe.NewContext(ctx)

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
state, err := vibe.StorageState(ctx)

// Save to file for later restoration
jsonBytes, _ := json.Marshal(state)
os.WriteFile("storage.json", jsonBytes, 0600)

// Restore storage state from JSON
var savedState vibium.StorageState
json.Unmarshal(jsonBytes, &savedState)
err := vibe.SetStorageState(ctx, &savedState)

// Clear all storage (cookies, localStorage, sessionStorage)
err := vibe.ClearStorage(ctx)
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
err := vibe.AddInitScript(ctx, `window.testMode = true;`)

// Mock APIs before page scripts run
err := vibe.AddInitScript(ctx, `
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
err := vibe.AddInitScript(ctx, `
    window.gtag = () => {};
    window.analytics = { track: () => {} };
`)

// Via BrowserContext for isolated contexts
browserCtx, _ := vibe.NewContext(ctx)
err := browserCtx.AddInitScript(ctx, `window.contextId = 'isolated';`)
```

## Tracing

Record browser actions with screenshots and DOM snapshots:

```go
// Get tracing controller
tracing := vibe.Tracing()

// Start with options
err := tracing.Start(ctx, &vibium.TracingStartOptions{
    Name:        "login-test",
    Screenshots: true,
    Snapshots:   true,
    Sources:     true,
    Title:       "Login Flow Test",
})

// Perform actions to record
vibe.Go(ctx, "https://example.com")
elem, _ := vibe.Find(ctx, "button", nil)
elem.Click(ctx, nil)

// Stop and get trace data
data, err := tracing.Stop(ctx, nil)
os.WriteFile("trace.zip", data, 0600)

// Use chunks for segmented recording
err := tracing.StartChunk(ctx, &vibium.TracingChunkOptions{
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
err := vibe.SetViewport(ctx, vibium.Viewport{
    Width:  1920,
    Height: 1080,
})

// Media emulation
err := vibe.EmulateMedia(ctx, &vibium.EmulateMediaOptions{
    Media:       "print",
    ColorScheme: "dark",
})

// Geolocation
err := vibe.SetGeolocation(ctx, &vibium.Geolocation{
    Latitude:  37.7749,
    Longitude: -122.4194,
})
```

## Error Handling

```go
import "errors"

elem, err := vibe.Find(ctx, "#missing", nil)
if err != nil {
    if errors.Is(err, vibium.ErrElementNotFound) {
        // Element not found
    }
    if errors.Is(err, vibium.ErrTimeout) {
        // Timeout
    }
}
```

## Debug Logging

```bash
VIBIUM_DEBUG=1 go run main.go
```

```go
// Check debug mode
if vibium.Debug() {
    // ...
}

// Custom logger
logger := vibium.NewDebugLogger()
ctx = vibium.ContextWithLogger(ctx, logger)
```
