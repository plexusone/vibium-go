# Error Types

## Sentinel Errors

```go
import vibium "github.com/grokify/vibium-go"

// Check specific error types
if errors.Is(err, vibium.ErrElementNotFound) {
    // Element not found
}
```

| Error | Description |
|-------|-------------|
| `ErrElementNotFound` | Element matching selector not found |
| `ErrTimeout` | Operation timed out |
| `ErrConnectionFailed` | Failed to connect to browser |
| `ErrConnectionClosed` | Connection to browser closed |
| `ErrBrowserCrashed` | Browser process crashed |
| `ErrClickerNotFound` | Clicker binary not found |

## Error Types

### ElementNotFoundError

```go
type ElementNotFoundError struct {
    Selector string
}

func (e ElementNotFoundError) Error() string
func (e ElementNotFoundError) Is(target error) bool
```

**Example:**

```go
elem, err := vibe.Find(ctx, "#missing", nil)
if err != nil {
    var notFound vibium.ElementNotFoundError
    if errors.As(err, &notFound) {
        fmt.Printf("Selector not found: %s\n", notFound.Selector)
    }
}
```

### TimeoutError

```go
type TimeoutError struct {
    Selector string
    Timeout  time.Duration
    Reason   string
}

func (e TimeoutError) Error() string
func (e TimeoutError) Is(target error) bool
```

**Example:**

```go
err := elem.Click(ctx, &vibium.ActionOptions{Timeout: time.Second})
if err != nil {
    var timeout vibium.TimeoutError
    if errors.As(err, &timeout) {
        fmt.Printf("Timed out after %v: %s\n", timeout.Timeout, timeout.Reason)
    }
}
```

### ConnectionError

```go
type ConnectionError struct {
    URL   string
    Cause error
}

func (e ConnectionError) Error() string
func (e ConnectionError) Is(target error) bool
func (e ConnectionError) Unwrap() error
```

### BrowserCrashedError

```go
type BrowserCrashedError struct {
    ExitCode int
    Output   string
}

func (e BrowserCrashedError) Error() string
func (e BrowserCrashedError) Is(target error) bool
```

### BiDiError

WebDriver BiDi protocol errors.

```go
type BiDiError struct {
    ErrorType string
    Message   string
}

func (e BiDiError) Error() string
```

## Error Handling Patterns

### Check Specific Error

```go
elem, err := vibe.Find(ctx, selector, nil)
if err != nil {
    if errors.Is(err, vibium.ErrElementNotFound) {
        // Try alternative selector
        elem, err = vibe.Find(ctx, altSelector, nil)
    }
    if err != nil {
        return err
    }
}
```

### Extract Error Details

```go
err := elem.Click(ctx, nil)
if err != nil {
    var timeout vibium.TimeoutError
    if errors.As(err, &timeout) {
        log.Printf("Click timed out on %s after %v",
            timeout.Selector, timeout.Timeout)
    }
    return err
}
```

### Retry on Timeout

```go
func clickWithRetry(ctx context.Context, elem *vibium.Element, attempts int) error {
    var err error
    for i := 0; i < attempts; i++ {
        err = elem.Click(ctx, nil)
        if err == nil {
            return nil
        }
        if !errors.Is(err, vibium.ErrTimeout) {
            return err // Don't retry non-timeout errors
        }
        time.Sleep(time.Second)
    }
    return err
}
```

### Handle Connection Issues

```go
vibe, err := vibium.Launch(ctx)
if err != nil {
    if errors.Is(err, vibium.ErrClickerNotFound) {
        return fmt.Errorf("install vibium: npm install -g vibium")
    }
    if errors.Is(err, vibium.ErrConnectionFailed) {
        return fmt.Errorf("browser failed to start")
    }
    return err
}
```
