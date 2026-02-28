# Prerequisites

## System Requirements

- Go 1.21 or later
- Chrome, Chromium, or Chrome for Testing
- Vibium Clicker binary (see below)

## Vibium Clicker

The clicker is a lightweight binary that bridges WebDriver BiDi with the browser.

!!! warning "Clicker Availability"
    The clicker binary is not yet publicly distributed. Contact the maintainers for access, or check the [releases page](https://github.com/agentplexus/vibium-go/releases) for updates.

### Specifying the Clicker Path

Once you have the clicker binary, specify its location:

```bash
export VIBIUM_CLICKER_PATH=/path/to/clicker
```

Or in Go code:

```go
vibe, err := vibium.Browser.Launch(ctx, &vibium.LaunchOptions{
    ExecutablePath: "/path/to/clicker",
})
```

### Platform Binaries

When available, binaries will be provided for:

| Platform | Binary Name |
|----------|-------------|
| macOS Apple Silicon | `clicker-darwin-arm64` |
| macOS Intel | `clicker-darwin-x64` |
| Linux ARM64 | `clicker-linux-arm64` |
| Linux x64 | `clicker-linux-x64` |
| Windows x64 | `clicker-win32-x64.exe` |

## Browser

Chrome for Testing is recommended:

```bash
# Use existing Chrome/Chromium
export CHROME_PATH=/path/to/chrome
```

The clicker will automatically detect Chrome installations in standard locations.
