//go:build integration

package integration

import (
	"testing"

	"github.com/plexusone/webpilot"
)

// TestEmulateMediaColorScheme tests color scheme emulation.
func TestEmulateMediaColorScheme(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate to a page that responds to color scheme
	bt.go_(`data:text/html,<!DOCTYPE html>
<html>
<head>
<style>
:root { --bg: white; --fg: black; }
@media (prefers-color-scheme: dark) {
  :root { --bg: black; --fg: white; }
}
body { background: var(--bg); color: var(--fg); }
</style>
</head>
<body><p id="test">Color scheme test</p></body>
</html>`)

	// Test dark mode
	err := bt.pilot.EmulateMedia(bt.ctx, webpilot.EmulateMediaOptions{
		ColorScheme: "dark",
	})
	if err != nil {
		t.Fatalf("EmulateMedia dark failed: %v", err)
	}

	// Verify dark mode is active via CSS custom property
	result, err := bt.pilot.Evaluate(bt.ctx, `getComputedStyle(document.documentElement).getPropertyValue('--bg').trim()`)
	if err != nil {
		t.Fatalf("Failed to get computed style: %v", err)
	}
	if result != "black" {
		t.Errorf("Expected dark mode bg=black, got %q", result)
	}

	// Test light mode
	err = bt.pilot.EmulateMedia(bt.ctx, webpilot.EmulateMediaOptions{
		ColorScheme: "light",
	})
	if err != nil {
		t.Fatalf("EmulateMedia light failed: %v", err)
	}

	result, err = bt.pilot.Evaluate(bt.ctx, `getComputedStyle(document.documentElement).getPropertyValue('--bg').trim()`)
	if err != nil {
		t.Fatalf("Failed to get computed style: %v", err)
	}
	if result != "white" {
		t.Errorf("Expected light mode bg=white, got %q", result)
	}
}

// TestEmulateMediaReducedMotion tests reduced motion preference emulation.
func TestEmulateMediaReducedMotion(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html>
<head>
<style>
:root { --motion: normal; }
@media (prefers-reduced-motion: reduce) {
  :root { --motion: reduced; }
}
</style>
</head>
<body><p>Reduced motion test</p></body>
</html>`)

	// Enable reduced motion
	err := bt.pilot.EmulateMedia(bt.ctx, webpilot.EmulateMediaOptions{
		ReducedMotion: "reduce",
	})
	if err != nil {
		t.Fatalf("EmulateMedia reduce motion failed: %v", err)
	}

	result, err := bt.pilot.Evaluate(bt.ctx, `getComputedStyle(document.documentElement).getPropertyValue('--motion').trim()`)
	if err != nil {
		t.Fatalf("Failed to get computed style: %v", err)
	}
	if result != "reduced" {
		t.Errorf("Expected motion=reduced, got %q", result)
	}

	// Disable reduced motion
	err = bt.pilot.EmulateMedia(bt.ctx, webpilot.EmulateMediaOptions{
		ReducedMotion: "no-preference",
	})
	if err != nil {
		t.Fatalf("EmulateMedia no-preference motion failed: %v", err)
	}

	result, err = bt.pilot.Evaluate(bt.ctx, `getComputedStyle(document.documentElement).getPropertyValue('--motion').trim()`)
	if err != nil {
		t.Fatalf("Failed to get computed style: %v", err)
	}
	if result != "normal" {
		t.Errorf("Expected motion=normal, got %q", result)
	}
}

// TestEmulateMediaPrint tests print media type emulation.
func TestEmulateMediaPrint(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html>
<head>
<style>
:root { --media: screen; }
@media print {
  :root { --media: print; }
}
</style>
</head>
<body><p>Print media test</p></body>
</html>`)

	// Enable print media
	err := bt.pilot.EmulateMedia(bt.ctx, webpilot.EmulateMediaOptions{
		Media: "print",
	})
	if err != nil {
		t.Fatalf("EmulateMedia print failed: %v", err)
	}

	result, err := bt.pilot.Evaluate(bt.ctx, `getComputedStyle(document.documentElement).getPropertyValue('--media').trim()`)
	if err != nil {
		t.Fatalf("Failed to get computed style: %v", err)
	}
	if result != "print" {
		t.Errorf("Expected media=print, got %q", result)
	}

	// Switch back to screen
	err = bt.pilot.EmulateMedia(bt.ctx, webpilot.EmulateMediaOptions{
		Media: "screen",
	})
	if err != nil {
		t.Fatalf("EmulateMedia screen failed: %v", err)
	}

	result, err = bt.pilot.Evaluate(bt.ctx, `getComputedStyle(document.documentElement).getPropertyValue('--media').trim()`)
	if err != nil {
		t.Fatalf("Failed to get computed style: %v", err)
	}
	if result != "screen" {
		t.Errorf("Expected media=screen, got %q", result)
	}
}

// TestEmulateMediaForcedColors tests forced colors emulation.
func TestEmulateMediaForcedColors(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html>
<head>
<style>
:root { --forced: none; }
@media (forced-colors: active) {
  :root { --forced: active; }
}
</style>
</head>
<body><p>Forced colors test</p></body>
</html>`)

	// Enable forced colors
	err := bt.pilot.EmulateMedia(bt.ctx, webpilot.EmulateMediaOptions{
		ForcedColors: "active",
	})
	if err != nil {
		t.Fatalf("EmulateMedia forced colors failed: %v", err)
	}

	result, err := bt.pilot.Evaluate(bt.ctx, `getComputedStyle(document.documentElement).getPropertyValue('--forced').trim()`)
	if err != nil {
		t.Fatalf("Failed to get computed style: %v", err)
	}
	if result != "active" {
		t.Errorf("Expected forced=active, got %q", result)
	}
}

// TestEmulateMediaContrast tests contrast preference emulation.
func TestEmulateMediaContrast(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html>
<head>
<style>
:root { --contrast: normal; }
@media (prefers-contrast: more) {
  :root { --contrast: more; }
}
@media (prefers-contrast: less) {
  :root { --contrast: less; }
}
</style>
</head>
<body><p>Contrast test</p></body>
</html>`)

	// Enable high contrast
	err := bt.pilot.EmulateMedia(bt.ctx, webpilot.EmulateMediaOptions{
		Contrast: "more",
	})
	if err != nil {
		t.Fatalf("EmulateMedia contrast more failed: %v", err)
	}

	result, err := bt.pilot.Evaluate(bt.ctx, `getComputedStyle(document.documentElement).getPropertyValue('--contrast').trim()`)
	if err != nil {
		t.Fatalf("Failed to get computed style: %v", err)
	}
	if result != "more" {
		t.Errorf("Expected contrast=more, got %q", result)
	}
}

// TestEmulateMediaCombined tests multiple media options at once.
func TestEmulateMediaCombined(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html>
<head>
<style>
:root { --scheme: light; --motion: normal; }
@media (prefers-color-scheme: dark) { :root { --scheme: dark; } }
@media (prefers-reduced-motion: reduce) { :root { --motion: reduced; } }
</style>
</head>
<body><p>Combined test</p></body>
</html>`)

	// Set multiple options
	err := bt.pilot.EmulateMedia(bt.ctx, webpilot.EmulateMediaOptions{
		ColorScheme:   "dark",
		ReducedMotion: "reduce",
	})
	if err != nil {
		t.Fatalf("EmulateMedia combined failed: %v", err)
	}

	// Check both applied
	scheme, err := bt.pilot.Evaluate(bt.ctx, `getComputedStyle(document.documentElement).getPropertyValue('--scheme').trim()`)
	if err != nil {
		t.Fatalf("Failed to get scheme: %v", err)
	}
	if scheme != "dark" {
		t.Errorf("Expected scheme=dark, got %q", scheme)
	}

	motion, err := bt.pilot.Evaluate(bt.ctx, `getComputedStyle(document.documentElement).getPropertyValue('--motion').trim()`)
	if err != nil {
		t.Fatalf("Failed to get motion: %v", err)
	}
	if motion != "reduced" {
		t.Errorf("Expected motion=reduced, got %q", motion)
	}
}
