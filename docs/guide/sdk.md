# Go Client SDK

The Go client SDK provides programmatic browser control with full feature parity to the JavaScript and Python clients.

## Installation

```bash
go get github.com/grokify/vibium-go
```

## Basic Usage

```go
package main

import (
    "context"
    "log"

    vibium "github.com/grokify/vibium-go"
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

```go
// Find single element
elem, err := vibe.Find(ctx, "button.submit", nil)

// Find with timeout
elem, err := vibe.Find(ctx, "button.submit", &vibium.FindOptions{
    Timeout: 10 * time.Second,
})

// Find with semantic selectors
elem, err := vibe.Find(ctx, "button", &vibium.FindOptions{
    Role:   "button",
    Text:   "Submit",
    TestID: "submit-btn",
})

// Find all matching elements
elements, err := vibe.FindAll(ctx, "li.item")
for _, elem := range elements {
    text, _ := elem.Text(ctx)
    fmt.Println(text)
}

// Must find (panics if not found)
elem := vibe.MustFind(ctx, "button.submit")
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

// Storage state
state, err := browserCtx.StorageState(ctx)
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
