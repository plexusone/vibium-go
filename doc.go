// Package vibium provides a Go client for the Vibium browser automation platform.
//
// Vibium is a browser automation platform built for AI agents that uses the
// WebDriver BiDi protocol for bidirectional communication with the browser.
//
// # Quick Start
//
// Launch a browser and navigate to a page:
//
//	ctx := context.Background()
//	vibe, err := vibium.Launch(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer vibe.Quit(ctx)
//
//	if err := vibe.Go(ctx, "https://example.com"); err != nil {
//	    log.Fatal(err)
//	}
//
// # Finding and Interacting with Elements
//
// Find elements using CSS selectors and interact with them:
//
//	link, err := vibe.Find(ctx, "a.my-link", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if err := link.Click(ctx, nil); err != nil {
//	    log.Fatal(err)
//	}
//
// Type text into input fields:
//
//	input, err := vibe.Find(ctx, "input[name='search']", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if err := input.Type(ctx, "search query", nil); err != nil {
//	    log.Fatal(err)
//	}
//
// # Screenshots
//
// Capture screenshots as PNG data:
//
//	data, err := vibe.Screenshot(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("screenshot.png", data, 0644)
//
// # Headless Mode
//
// Launch in headless mode for CI/server environments:
//
//	vibe, err := vibium.LaunchHeadless(ctx)
//
// Or with explicit options:
//
//	vibe, err := vibium.Browser.Launch(ctx, &vibium.LaunchOptions{
//	    Headless: true,
//	    Port:     9515,
//	})
//
// # Requirements
//
// This client requires the Vibium clicker binary to be available. Install it via:
//
//	npm install -g vibium
//
// Or set the VIBIUM_CLICKER_PATH environment variable to point to the binary.
package vibium
