# Prerequisites

## System Requirements

- Go 1.21 or later
- Chrome, Chromium, or Chrome for Testing
- Node.js (for installing clicker via npm)

## Vibium Clicker

The clicker is a lightweight binary that bridges WebDriver BiDi with the browser.

### Install via npm

```bash
npm install -g vibium
```

### Manual Installation

Download from [Vibium releases](https://github.com/VibiumDev/vibium/releases) for your platform:

- `vibium-darwin-arm64` - macOS Apple Silicon
- `vibium-darwin-x64` - macOS Intel
- `vibium-linux-arm64` - Linux ARM64
- `vibium-linux-x64` - Linux x64
- `vibium-win32-x64.exe` - Windows x64

### Specify Path

If the clicker is not in your PATH:

```bash
export VIBIUM_CLICKER_PATH=/path/to/clicker
```

Or in Go code:

```go
vibe, err := vibium.Browser.Launch(ctx, &vibium.LaunchOptions{
    ExecutablePath: "/path/to/clicker",
})
```

## Browser

Chrome for Testing is recommended:

```bash
# Install via clicker
vibium install

# Or use existing Chrome/Chromium
export CHROME_PATH=/path/to/chrome
```
