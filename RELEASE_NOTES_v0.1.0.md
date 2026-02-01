# Release Notes - v0.1.0

Initial release of the Vibium Go Client SDK.

## Highlights

**Go client for Vibium browser automation platform with full feature parity with JavaScript and Python clients.**

This SDK provides a native Go interface to the [Vibium](https://github.com/VibiumDev/vibium) browser automation platform, which uses the WebDriver BiDi protocol for bidirectional communication with the browser.

## Features

### Browser Control

- `Browser.Launch()` - Launch browser with configurable options (headless, port, executable path)
- `Launch()` / `LaunchHeadless()` - Convenience functions for common launch patterns
- `Quit()` - Graceful browser shutdown

### Navigation

- `Go(url)` - Navigate to URL with automatic wait for page load
- `Reload()` - Reload current page
- `Back()` / `Forward()` - Browser history navigation
- `URL()` / `Title()` - Get current page URL and title
- `WaitForNavigation()` - Wait for navigation to complete

### Element Interaction

- `Find(selector)` - Find element by CSS selector with actionability waits
- `FindAll(selector)` - Find all matching elements
- `MustFind(selector)` - Find or panic (for concise test code)
- `Click()` - Click with automatic actionability checks (visible, stable, enabled)
- `Type(text)` - Type text with editability checks
- `Text()` - Get element text content
- `GetAttribute(name)` - Get element attribute value
- `BoundingBox()` / `Center()` - Get element geometry

### Screenshots & JavaScript

- `Screenshot()` - Capture page as PNG
- `Evaluate(script)` - Execute JavaScript in page context

### Debug Logging

- Set `VIBIUM_DEBUG=1` to enable structured JSON logging via Go's `slog` package

## Installation

```bash
go get github.com/grokify/vibium-go
```

## Prerequisites

Install the Vibium clicker binary:

```bash
npm install -g vibium
```

Or set `VIBIUM_CLICKER_PATH` to point to the binary.

## Quick Example

```go
package main

import (
    "context"
    "log"

    "github.com/grokify/vibium-go"
)

func main() {
    ctx := context.Background()

    vibe, err := vibium.Launch(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer vibe.Quit(ctx)

    if err := vibe.Go(ctx, "https://example.com"); err != nil {
        log.Fatal(err)
    }

    link, err := vibe.Find(ctx, "a", nil)
    if err != nil {
        log.Fatal(err)
    }

    if err := link.Click(ctx, nil); err != nil {
        log.Fatal(err)
    }
}
```

## Feature Parity

| Feature | JS | Python | Go |
|---------|:--:|:------:|:--:|
| `browser.Launch()` | ✅ | ✅ | ✅ |
| `vibe.Go(url)` | ✅ | ✅ | ✅ |
| `vibe.Screenshot()` | ✅ | ✅ | ✅ |
| `vibe.Find(selector)` | ✅ | ✅ | ✅ |
| `vibe.FindAll(selector)` | ❌ | ❌ | ✅ |
| `vibe.Evaluate(script)` | ✅ | ✅ | ✅ |
| `vibe.Reload()` | ❌ | ❌ | ✅ |
| `vibe.Back()` / `vibe.Forward()` | ❌ | ❌ | ✅ |
| `vibe.Quit()` | ✅ | ✅ | ✅ |
| `element.Click()` | ✅ | ✅ | ✅ |
| `element.Type(text)` | ✅ | ✅ | ✅ |
| `element.Text()` | ✅ | ✅ | ✅ |
| `element.GetAttribute()` | ✅ | ✅ | ✅ |
| `element.BoundingBox()` | ✅ | ✅ | ✅ |
| Actionability waits | ✅ | ✅ | ✅ |
| Debug logging | ✅ | ❌ | ✅ |

## Documentation

- [README](https://github.com/grokify/vibium-go/blob/main/README.md) - Full documentation
- [GoDoc](https://pkg.go.dev/github.com/grokify/vibium-go) - API reference
- [CHANGELOG](https://github.com/grokify/vibium-go/blob/main/CHANGELOG.md) - Version history
