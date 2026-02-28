# Testing

## Unit Tests

Run unit tests without browser:

```bash
go test -v ./...
```

## Integration Tests

Integration tests require the Vibium clicker binary.

!!! warning "Clicker Required"
    The clicker binary is not yet publicly distributed. See [Prerequisites](../getting-started/prerequisites.md) for details.

### Setup

```bash
# Set path to clicker binary
export VIBIUM_CLICKER_PATH=/path/to/clicker
```

### Running

```bash
# All integration tests
go test -tags=integration -v ./integration/...

# Headless mode (for CI)
VIBIUM_HEADLESS=1 go test -tags=integration -v ./integration/...

# Specific tests
go test -tags=integration -v ./integration/... -run TestExampleCom
```

### Test Sites

| Site | Description |
|------|-------------|
| `example.com` | Simple smoke tests |
| `the-internet.herokuapp.com` | Interactive UI patterns |

## MCP Server Tests

```bash
# Unit tests
go test -v ./mcp/...

# With verbose output
go test -v ./mcp/... -count=1
```

## Script Runner Tests

```bash
# Test script parsing
go test -v ./script/...

# Run example scripts
vibium run examples/basic.json --headless
```

## Linting

```bash
# Run linter
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix
```

## Coverage

```bash
# Generate coverage
go test -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out

# Coverage badge
gocoverbadge -dir . -exclude cmd -badge-only
```

## CI Configuration

See `.github/workflows/` for:

- `ci.yaml` - Build and unit tests
- `lint.yaml` - Linting
- `e2e.yaml` - E2E tests (manual trigger only, requires clicker URL)

## Writing Tests

### Unit Test Example

```go
func TestElementClick(t *testing.T) {
    // Setup mock BiDi client
    client := newMockBiDiClient()
    elem := &Element{client: client, selector: "#btn"}

    // Test
    err := elem.Click(context.Background(), nil)

    // Assert
    require.NoError(t, err)
    require.Equal(t, "vibium:click", client.lastCommand)
}
```

### Integration Test Example

```go
//go:build integration

func TestExampleComNavigation(t *testing.T) {
    ctx := context.Background()

    vibe, err := vibium.LaunchHeadless(ctx)
    require.NoError(t, err)
    defer vibe.Quit(ctx)

    err = vibe.Go(ctx, "https://example.com")
    require.NoError(t, err)

    title, err := vibe.Title(ctx)
    require.NoError(t, err)
    require.Equal(t, "Example Domain", title)
}
```
