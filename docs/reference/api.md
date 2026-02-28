# API Reference

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/grokify/vibium-go).

## Package Structure

```
github.com/grokify/vibium-go
├── vibium.go       # Main Vibe type, browser control
├── element.go      # Element interactions
├── types.go        # Options and configuration
├── errors.go       # Error types
├── keyboard.go     # Keyboard controller
├── mouse.go        # Mouse controller
├── touch.go        # Touch controller
├── context.go      # Browser context
├── clock.go        # Clock control
├── tracing.go      # Trace recording
├── mcp/            # MCP server
│   ├── server.go
│   ├── session.go
│   ├── recorder.go
│   └── tools*.go
├── script/         # Script format
│   ├── types.go
│   └── schema.go
└── cmd/vibium/     # CLI
    └── cmd/
```

## Core Types

### Vibe

The main browser controller.

```go
type Vibe struct {
    // ...
}

// Launch
func Launch(ctx context.Context) (*Vibe, error)
func LaunchHeadless(ctx context.Context) (*Vibe, error)

// Navigation
func (v *Vibe) Go(ctx context.Context, url string) error
func (v *Vibe) URL(ctx context.Context) (string, error)
func (v *Vibe) Title(ctx context.Context) (string, error)
func (v *Vibe) Back(ctx context.Context) error
func (v *Vibe) Forward(ctx context.Context) error
func (v *Vibe) Reload(ctx context.Context) error

// Finding elements
func (v *Vibe) Find(ctx context.Context, selector string, opts *FindOptions) (*Element, error)
func (v *Vibe) FindAll(ctx context.Context, selector string) ([]*Element, error)
func (v *Vibe) MustFind(ctx context.Context, selector string) *Element

// Screenshots
func (v *Vibe) Screenshot(ctx context.Context) ([]byte, error)
func (v *Vibe) PDF(ctx context.Context, opts *PDFOptions) ([]byte, error)

// JavaScript
func (v *Vibe) Evaluate(ctx context.Context, script string) (any, error)

// Input controllers
func (v *Vibe) Keyboard() *Keyboard
func (v *Vibe) Mouse() *Mouse
func (v *Vibe) Touch() *Touch

// Cleanup
func (v *Vibe) Quit(ctx context.Context) error
func (v *Vibe) IsClosed() bool
```

### Element

Represents a DOM element.

```go
type Element struct {
    // ...
}

// Interactions
func (e *Element) Click(ctx context.Context, opts *ActionOptions) error
func (e *Element) DblClick(ctx context.Context, opts *ActionOptions) error
func (e *Element) Type(ctx context.Context, text string, opts *ActionOptions) error
func (e *Element) Fill(ctx context.Context, value string, opts *ActionOptions) error
func (e *Element) Clear(ctx context.Context, opts *ActionOptions) error
func (e *Element) Press(ctx context.Context, key string, opts *ActionOptions) error
func (e *Element) Check(ctx context.Context, opts *ActionOptions) error
func (e *Element) Uncheck(ctx context.Context, opts *ActionOptions) error
func (e *Element) SelectOption(ctx context.Context, values SelectOptionValues, opts *ActionOptions) error
func (e *Element) Hover(ctx context.Context, opts *ActionOptions) error
func (e *Element) Focus(ctx context.Context, opts *ActionOptions) error

// State
func (e *Element) Text(ctx context.Context) (string, error)
func (e *Element) Value(ctx context.Context) (string, error)
func (e *Element) InnerHTML(ctx context.Context) (string, error)
func (e *Element) GetAttribute(ctx context.Context, name string) (string, error)
func (e *Element) BoundingBox(ctx context.Context) (*BoundingBox, error)
func (e *Element) IsVisible(ctx context.Context) (bool, error)
func (e *Element) IsEnabled(ctx context.Context) (bool, error)
func (e *Element) IsChecked(ctx context.Context) (bool, error)
```

### Options

```go
type LaunchOptions struct {
    Headless       bool
    Port           int
    ExecutablePath string
}

type FindOptions struct {
    Timeout     time.Duration
    Role        string
    Text        string
    Label       string
    Placeholder string
    TestID      string
}

type ActionOptions struct {
    Timeout time.Duration
}
```

## MCP Server

```go
import "github.com/grokify/vibium-go/mcp"

type Config struct {
    Headless       bool
    DefaultTimeout time.Duration
    Project        string
}

func NewServer(config Config) *Server
func (s *Server) Run(ctx context.Context) error
func (s *Server) Close(ctx context.Context) error
```

## Script Types

```go
import "github.com/grokify/vibium-go/script"

type Script struct {
    Name        string
    Description string
    Version     int
    Headless    bool
    BaseURL     string
    Timeout     string
    Variables   map[string]string
    Steps       []Step
}

type Step struct {
    Action   Action
    Selector string
    URL      string
    Value    string
    // ... see script/types.go
}

// Get JSON Schema
func Schema() []byte
```
