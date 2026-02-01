# vibium-go

A Go client for the [Vibium](https://github.com/VibiumDev/vibium) browser automation platform.

Vibium is a browser automation platform built for AI agents that uses the WebDriver BiDi protocol for bidirectional communication with the browser.

## Feature Parity

This Go client has full feature parity with the official JavaScript and Python clients.

| Feature | JS | Python | Go |
|---------|:--:|:------:|:--:|
| `browser.Launch()` | ✅ | ✅ | ✅ |
| `vibe.Go(url)` | ✅ | ✅ | ✅ |
| `vibe.Screenshot()` | ✅ | ✅ | ✅ |
| `vibe.Find(selector)` | ✅ | ✅ | ✅ |
| `vibe.Evaluate(script)` | ✅ | ✅ | ✅ |
| `vibe.Quit()` | ✅ | ✅ | ✅ |
| `element.Click()` | ✅ | ✅ | ✅ |
| `element.Type(text)` | ✅ | ✅ | ✅ |
| `element.Text()` | ✅ | ✅ | ✅ |
| `element.GetAttribute()` | ✅ | ✅ | ✅ |
| `element.BoundingBox()` | ✅ | ✅ | ✅ |
| Browsing context management | ✅ | ✅ | ✅ |
| Actionability waits | ✅ | ✅ | ✅ |

## Installation

```bash
go get github.com/grokify/vibium-go
```

## Prerequisites

This client requires the Vibium clicker binary. Install it via npm:

```bash
npm install -g vibium
```

Or set the `VIBIUM_CLICKER_PATH` environment variable to point to the binary.

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/grokify/vibium-go"
)

func main() {
    ctx := context.Background()

    // Launch browser
    vibe, err := vibium.Launch(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer func() { _ = vibe.Quit(ctx) }()

    // Navigate to a page
    if err := vibe.Go(ctx, "https://example.com"); err != nil {
        log.Fatal(err)
    }

    // Find and click a link
    link, err := vibe.Find(ctx, "a", nil)
    if err != nil {
        log.Fatal(err)
    }

    if err := link.Click(ctx, nil); err != nil {
        log.Fatal(err)
    }
}
```

## API Reference

### Browser Control

```go
// Launch with default options
vibe, err := vibium.Launch(ctx)

// Launch headless
vibe, err := vibium.LaunchHeadless(ctx)

// Launch with custom options
vibe, err := vibium.Browser.Launch(ctx, &vibium.LaunchOptions{
    Headless:       true,
    Port:           9515,
    ExecutablePath: "/path/to/clicker",
})
```

### Navigation

```go
// Navigate to URL
err := vibe.Go(ctx, "https://example.com")

// Get current URL
url, err := vibe.URL(ctx)

// Get page title
title, err := vibe.Title(ctx)

// Wait for navigation to complete
err := vibe.WaitForNavigation(ctx, 30*time.Second)
```

### Finding Elements

```go
// Find element by CSS selector
elem, err := vibe.Find(ctx, "button.submit", nil)

// Find with custom timeout
elem, err := vibe.Find(ctx, "button.submit", &vibium.FindOptions{
    Timeout: 10 * time.Second,
})

// Must find (panics if not found)
elem := vibe.MustFind(ctx, "button.submit")
```

### Element Interaction

```go
// Click element (waits for actionability)
err := elem.Click(ctx, nil)

// Type text into element (waits for editability)
err := elem.Type(ctx, "Hello, World!", nil)

// Get element text content
text, err := elem.Text(ctx)

// Get element attribute
href, err := elem.GetAttribute(ctx, "href")

// Get bounding box
box, err := elem.BoundingBox(ctx)
// box.X, box.Y, box.Width, box.Height

// Get element info
info := elem.Info()
// info.Tag, info.Text, info.Box

// Get element center point
x, y := elem.Center()
```

### Screenshots

```go
// Capture screenshot as PNG data
data, err := vibe.Screenshot(ctx)

// Save to file
os.WriteFile("screenshot.png", data, 0644)
```

### JavaScript Evaluation

```go
// Execute JavaScript
result, err := vibe.Evaluate(ctx, "return document.title")
```

### Cleanup

```go
// Close browser and cleanup
err := vibe.Quit(ctx)

// Check if closed
closed := vibe.IsClosed()
```

## Actionability

The client automatically waits for elements to be actionable before performing actions:

- **Click**: Waits for element to be visible, stable, receive events, and enabled
- **Type**: Same as click, plus waits for element to be editable

Default timeout is 30 seconds, configurable via `ActionOptions`:

```go
err := elem.Click(ctx, &vibium.ActionOptions{
    Timeout: 10 * time.Second,
})
```

## Error Types

```go
// Connection failure
vibium.ErrConnectionFailed
vibium.ConnectionError{URL, Cause}

// Element not found
vibium.ErrElementNotFound
vibium.ElementNotFoundError{Selector}

// Timeout
vibium.ErrTimeout
vibium.TimeoutError{Selector, Timeout, Reason}

// Browser crashed
vibium.ErrBrowserCrashed
vibium.BrowserCrashedError{ExitCode, Output}

// Clicker binary not found
vibium.ErrClickerNotFound

// Connection closed
vibium.ErrConnectionClosed

// BiDi protocol error
vibium.BiDiError{ErrorType, Message}
```

## License

Apache-2.0
