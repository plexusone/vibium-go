# Quick Start

## Go Client SDK

```go
package main

import (
    "context"
    "fmt"
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

    // Navigate
    if err := vibe.Go(ctx, "https://example.com"); err != nil {
        log.Fatal(err)
    }

    // Get page title
    title, _ := vibe.Title(ctx)
    fmt.Println("Title:", title)

    // Find and click a link
    link, err := vibe.Find(ctx, "a", nil)
    if err != nil {
        log.Fatal(err)
    }

    if err := link.Click(ctx, nil); err != nil {
        log.Fatal(err)
    }

    // Take screenshot
    data, _ := vibe.Screenshot(ctx)
    os.WriteFile("screenshot.png", data, 0644)
}
```

## MCP Server

Start the server:

```bash
vibium mcp --headless
```

Configure in Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "vibium": {
      "command": "vibium",
      "args": ["mcp", "--headless"]
    }
  }
}
```

Then ask Claude: "Navigate to example.com and take a screenshot"

## CLI

Interactive browser control:

```bash
# Launch browser
vibium launch

# Navigate
vibium go https://example.com

# Interact
vibium fill "#search" "hello world"
vibium click "#submit"

# Capture
vibium screenshot result.png

# Cleanup
vibium quit
```

## Script Runner

Create `test.json`:

```json
{
  "name": "Example Test",
  "steps": [
    {"action": "navigate", "url": "https://example.com"},
    {"action": "assertTitle", "expected": "Example Domain"},
    {"action": "screenshot", "file": "result.png"}
  ]
}
```

Run:

```bash
vibium run test.json
```
