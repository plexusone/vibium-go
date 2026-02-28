# Technical Requirements Document: Vibium MCP Server

## Architecture Overview

```
┌────────────────────────────────────────────────────────────────────┐
│                      AI ASSISTANT                                  │
│                 (Claude Code, ChatGPT, etc.)                       │
└────────────────────────────────────────────────────────────────────┘
                               │
                               │ MCP Protocol (stdio)
                               │ JSON-RPC 2.0
                               ▼
┌────────────────────────────────────────────────────────────────────┐
│  cmd/vibium-mcp/main.go                                           │
│  └─► mcp.NewServer(config).Run()                                  │
└────────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌────────────────────────────────────────────────────────────────────┐
│  mcp/  (github.com/grokify/vibium-go/mcp)                         │
│                                                                    │
│  ├── server.go      MCP protocol handling, tool dispatch          │
│  ├── tools.go       Tool definitions and handlers                 │
│  ├── session.go     Browser session lifecycle                     │
│  └── report/                                                       │
│      ├── result.go      TestResult internal struct                │
│      ├── diagnostic.go  DiagnosticReport for agents               │
│      ├── teamreport.go  Convert to multi-agent-spec format        │
│      └── collector.go   Console/network log collection            │
└────────────────────────────────────────────────────────────────────┘
                               │
                               │ imports
                               ▼
┌────────────────────────────────────────────────────────────────────┐
│  vibium (github.com/grokify/vibium-go)  [PUBLIC API]              │
│                                                                    │
│  vibium.go     Launch(ctx) / LaunchHeadless(ctx)                  │
│  vibe.go       Vibe.Go() / Find() / Screenshot() / Evaluate()     │
│  element.go    Element.Click() / Type() / Text()                  │
│  options.go    LaunchOptions / FindOptions / ActionOptions        │
│  errors.go     ErrElementNotFound / ErrTimeout / ...              │
└────────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌────────────────────────────────────────────────────────────────────┐
│  internal/bidi/     BiDi WebSocket client (private)               │
│  internal/clicker/  Clicker process management (private)          │
└────────────────────────────────────────────────────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │   Clicker Binary    │
                    │   (WebDriver BiDi)  │
                    └─────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │   Chromium Browser  │
                    └─────────────────────┘
```

## Directory Structure

```
vibium-go/
├── go.mod
├── go.sum
├── doc.go
│
├── vibium.go                   # Public entry: Launch, LaunchHeadless
├── vibe.go                     # Vibe browser controller
├── element.go                  # Element interaction
├── options.go                  # LaunchOptions, FindOptions, ActionOptions
├── errors.go                   # Error types
├── types.go                    # BoundingBox, ElementInfo
│
├── internal/
│   ├── bidi/                   # BiDi protocol (private)
│   │   ├── client.go           # BiDiClient WebSocket handling
│   │   └── messages.go         # BiDi message types
│   └── clicker/                # Clicker process (private)
│       ├── process.go          # ClickerProcess lifecycle
│       └── locate.go           # Binary discovery
│
├── mcp/                        # MCP server package
│   ├── server.go               # MCP protocol, tool registration
│   ├── tools.go                # Tool definitions & handlers
│   ├── session.go              # Browser session management
│   ├── config.go               # Server configuration
│   │
│   └── report/                 # Report generation
│       ├── result.go           # TestResult, StepResult structs
│       ├── diagnostic.go       # DiagnosticReport struct & generation
│       ├── teamreport.go       # multi-agent-spec TeamReport conversion
│       ├── collector.go        # Console/network log collector
│       └── render.go           # Format selection & rendering
│
├── cmd/
│   └── vibium-mcp/             # MCP server executable
│       └── main.go
│
├── testplans/                  # Example test plans
│   ├── example_login.md
│   └── example_navigation.md
│
├── docs/
│   ├── PRD.md                  # Product requirements
│   └── TRD.md                  # Technical requirements (this file)
│
└── integration/                # Integration tests
    └── ...
```

## Module Dependencies

### go.mod

```go
module github.com/grokify/vibium-go

go 1.22

require (
    github.com/gorilla/websocket v1.5.3
    github.com/agentplexus/multi-agent-spec v0.x.x
    github.com/modelcontextprotocol/go-sdk v1.3.1
)
```

### Dependency Purposes

| Module | Purpose |
|--------|---------|
| `gorilla/websocket` | BiDi WebSocket communication |
| `multi-agent-spec` | TeamReport rendering (box format) |
| `modelcontextprotocol/go-sdk` | Official MCP protocol SDK |

## MCP Tool Specifications

### Tool: browser_launch

```json
{
  "name": "browser_launch",
  "description": "Launch a browser instance. Call this before any other browser operations.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "headless": {
        "type": "boolean",
        "description": "Run browser without GUI (default: true)",
        "default": true
      }
    }
  }
}
```

### Tool: browser_quit

```json
{
  "name": "browser_quit",
  "description": "Close the browser and cleanup resources."
}
```

### Tool: navigate

```json
{
  "name": "navigate",
  "description": "Navigate to a URL.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "url": {
        "type": "string",
        "description": "The URL to navigate to"
      }
    },
    "required": ["url"]
  }
}
```

### Tool: click

```json
{
  "name": "click",
  "description": "Click an element by CSS selector.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "selector": {
        "type": "string",
        "description": "CSS selector for the element to click"
      },
      "timeout_ms": {
        "type": "integer",
        "description": "Timeout in milliseconds (default: 5000)",
        "default": 5000
      }
    },
    "required": ["selector"]
  }
}
```

### Tool: type

```json
{
  "name": "type",
  "description": "Type text into an input element.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "selector": {
        "type": "string",
        "description": "CSS selector for the input element"
      },
      "text": {
        "type": "string",
        "description": "Text to type"
      },
      "timeout_ms": {
        "type": "integer",
        "description": "Timeout in milliseconds (default: 5000)",
        "default": 5000
      }
    },
    "required": ["selector", "text"]
  }
}
```

### Tool: get_text

```json
{
  "name": "get_text",
  "description": "Get the text content of an element.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "selector": {
        "type": "string",
        "description": "CSS selector for the element"
      },
      "timeout_ms": {
        "type": "integer",
        "description": "Timeout in milliseconds (default: 5000)",
        "default": 5000
      }
    },
    "required": ["selector"]
  }
}
```

### Tool: screenshot

```json
{
  "name": "screenshot",
  "description": "Capture a screenshot of the current page.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "format": {
        "type": "string",
        "enum": ["base64", "file"],
        "description": "Output format (default: base64)",
        "default": "base64"
      },
      "path": {
        "type": "string",
        "description": "File path (required if format is 'file')"
      }
    }
  }
}
```

### Tool: get_test_report

```json
{
  "name": "get_test_report",
  "description": "Get the test execution report in specified format.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "format": {
        "type": "string",
        "enum": ["box", "diagnostic", "json"],
        "description": "Report format: box (terminal), diagnostic (full JSON for agents), json (multi-agent-spec)",
        "default": "box"
      }
    }
  }
}
```

### Tool: assert_text

```json
{
  "name": "assert_text",
  "description": "Assert that text exists on the page.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "text": {
        "type": "string",
        "description": "Text to search for"
      },
      "selector": {
        "type": "string",
        "description": "Optional: limit search to element matching selector"
      }
    },
    "required": ["text"]
  }
}
```

### Tool: assert_element

```json
{
  "name": "assert_element",
  "description": "Assert that an element exists on the page.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "selector": {
        "type": "string",
        "description": "CSS selector for the element"
      },
      "timeout_ms": {
        "type": "integer",
        "description": "Timeout in milliseconds (default: 5000)",
        "default": 5000
      }
    },
    "required": ["selector"]
  }
}
```

## Data Structures

### StepResult

```go
// mcp/report/result.go

type StepResult struct {
    ID         string            `json:"id"`
    Action     string            `json:"action"`
    Args       map[string]any    `json:"args"`
    Status     Status            `json:"status"`     // GO, WARN, NO-GO, SKIP
    Severity   Severity          `json:"severity"`   // critical, high, medium, low, info
    DurationMS int64             `json:"duration_ms"`
    Result     any               `json:"result,omitempty"`
    Error      *StepError        `json:"error,omitempty"`
    Context    *StepContext      `json:"context,omitempty"`
    Console    []ConsoleEntry    `json:"console_logs,omitempty"`
    Network    []NetworkError    `json:"network_errors,omitempty"`
    Screenshot *ScreenshotRef    `json:"screenshot,omitempty"`
}

type StepError struct {
    Type        string   `json:"type"`
    Message     string   `json:"message"`
    Selector    string   `json:"selector,omitempty"`
    TimeoutMS   int64    `json:"timeout_ms,omitempty"`
    Suggestions []string `json:"suggestions,omitempty"`
}

type StepContext struct {
    PageURL        string   `json:"page_url"`
    PageTitle      string   `json:"page_title"`
    VisibleButtons []string `json:"visible_buttons,omitempty"`
    DOMSnippet     string   `json:"dom_snippet,omitempty"`
}

type ConsoleEntry struct {
    Level   string `json:"level"`   // error, warn, info, log
    Message string `json:"message"`
    Source  string `json:"source"`  // javascript, network
    URL     string `json:"url,omitempty"`
}

type NetworkError struct {
    URL        string `json:"url"`
    Method     string `json:"method"`
    StatusCode int    `json:"status"`
}

type ScreenshotRef struct {
    Path   string `json:"path,omitempty"`
    Base64 string `json:"base64,omitempty"`
}
```

### TestResult

```go
// mcp/report/result.go

type TestResult struct {
    TestPlan   string        `json:"test_plan,omitempty"`
    Project    string        `json:"project"`
    Target     string        `json:"target"`
    Status     Status        `json:"status"`
    DurationMS int64         `json:"duration_ms"`
    Browser    BrowserInfo   `json:"browser"`
    Steps      []StepResult  `json:"steps"`
    Recommendations []string `json:"recommendations,omitempty"`
    GeneratedAt time.Time    `json:"generated_at"`
}

type BrowserInfo struct {
    Name     string `json:"name"`
    Headless bool   `json:"headless"`
    Viewport struct {
        Width  int `json:"width"`
        Height int `json:"height"`
    } `json:"viewport"`
}
```

### DiagnosticReport

```go
// mcp/report/diagnostic.go

type DiagnosticReport struct {
    TestResult
    // DiagnosticReport embeds TestResult and adds no extra fields
    // The full StepError, StepContext, Console, Network data
    // is already in TestResult.Steps
}

func (r *DiagnosticReport) JSON() ([]byte, error) {
    return json.MarshalIndent(r, "", "  ")
}
```

### TeamReport Conversion

```go
// mcp/report/teamreport.go

import masreport "github.com/agentplexus/multi-agent-spec/sdk/go"

func ToTeamReport(tr *TestResult) *masreport.TeamReport {
    report := &masreport.TeamReport{
        Title:       "BROWSER TEST REPORT",
        Project:     tr.Project,
        Target:      tr.Target,
        Status:      convertStatus(tr.Status),
        GeneratedAt: tr.GeneratedAt,
        Teams:       groupStepsIntoTeams(tr.Steps),
    }
    return report
}

func groupStepsIntoTeams(steps []StepResult) []masreport.TeamSection {
    // Group steps by category (navigation, interaction, assertion)
    // Return as TeamSections with proper dependencies
}
```

## Session Management

```go
// mcp/session.go

type Session struct {
    mu       sync.Mutex
    vibe     *vibium.Vibe
    config   SessionConfig
    results  []StepResult
    console  []ConsoleEntry
    network  []NetworkError
}

type SessionConfig struct {
    Headless       bool
    DefaultTimeout time.Duration
    Project        string
}

func (s *Session) LaunchIfNeeded(ctx context.Context) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if s.vibe != nil && !s.vibe.IsClosed() {
        return nil
    }

    var err error
    if s.config.Headless {
        s.vibe, err = vibium.LaunchHeadless(ctx)
    } else {
        s.vibe, err = vibium.Launch(ctx)
    }
    return err
}

func (s *Session) RecordStep(result StepResult) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.results = append(s.results, result)
}

func (s *Session) GetTestResult() *TestResult {
    s.mu.Lock()
    defer s.mu.Unlock()

    return &TestResult{
        Project:     s.config.Project,
        Status:      computeOverallStatus(s.results),
        DurationMS:  computeTotalDuration(s.results),
        Steps:       s.results,
        GeneratedAt: time.Now(),
    }
}

func (s *Session) Close(ctx context.Context) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if s.vibe != nil {
        return s.vibe.Quit(ctx)
    }
    return nil
}
```

## Server Implementation

```go
// mcp/server.go

import (
    "context"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
    session   *Session
    mcpServer *mcp.Server
    config    Config
}

func NewServer(config Config) *Server {
    s := &Server{
        config: config,
        session: NewSession(SessionConfig{
            Headless:       config.Headless,
            DefaultTimeout: config.DefaultTimeout,
            Project:        config.Project,
        }),
    }

    s.mcpServer = mcp.NewServer(
        &mcp.Implementation{
            Name:    "vibium-mcp",
            Version: "0.2.0",
        },
        nil,
    )

    s.registerTools()
    return s
}

func (s *Server) registerTools() {
    // browser_launch
    mcp.AddTool(s.mcpServer, &mcp.Tool{
        Name:        "browser_launch",
        Description: "Launch a browser instance. Call this before any other browser operations.",
    }, s.handleBrowserLaunch)

    // navigate
    mcp.AddTool(s.mcpServer, &mcp.Tool{
        Name:        "navigate",
        Description: "Navigate to a URL.",
    }, s.handleNavigate)

    // click
    mcp.AddTool(s.mcpServer, &mcp.Tool{
        Name:        "click",
        Description: "Click an element by CSS selector.",
    }, s.handleClick)

    // type
    mcp.AddTool(s.mcpServer, &mcp.Tool{
        Name:        "type",
        Description: "Type text into an input element.",
    }, s.handleType)

    // get_text
    mcp.AddTool(s.mcpServer, &mcp.Tool{
        Name:        "get_text",
        Description: "Get the text content of an element.",
    }, s.handleGetText)

    // screenshot
    mcp.AddTool(s.mcpServer, &mcp.Tool{
        Name:        "screenshot",
        Description: "Capture a screenshot of the current page.",
    }, s.handleScreenshot)

    // get_test_report
    mcp.AddTool(s.mcpServer, &mcp.Tool{
        Name:        "get_test_report",
        Description: "Get the test execution report in specified format.",
    }, s.handleGetTestReport)
}

func (s *Server) Run(ctx context.Context) error {
    return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}
```

## Tool Handlers

The official MCP SDK uses typed input/output structs with JSON schema tags:

```go
// mcp/tools.go

// Input/Output types with jsonschema tags for automatic schema generation

type BrowserLaunchInput struct {
    Headless bool `json:"headless" jsonschema:"description=Run browser without GUI,default=true"`
}

type BrowserLaunchOutput struct {
    Message string `json:"message"`
}

func (s *Server) handleBrowserLaunch(
    ctx context.Context,
    req *mcp.CallToolRequest,
    input BrowserLaunchInput,
) (*mcp.CallToolResult, BrowserLaunchOutput, error) {
    s.session.config.Headless = input.Headless

    start := time.Now()
    err := s.session.LaunchIfNeeded(ctx)
    duration := time.Since(start)

    result := StepResult{
        ID:         "browser-launch",
        Action:     "browser_launch",
        Args:       map[string]any{"headless": input.Headless},
        DurationMS: duration.Milliseconds(),
    }

    if err != nil {
        result.Status = StatusNoGo
        result.Severity = SeverityCritical
        result.Error = &StepError{
            Type:    "LaunchError",
            Message: err.Error(),
        }
        s.session.RecordStep(result)
        return nil, BrowserLaunchOutput{}, err
    }

    result.Status = StatusGo
    result.Severity = SeverityInfo
    s.session.RecordStep(result)

    return nil, BrowserLaunchOutput{Message: "Browser launched successfully"}, nil
}

type NavigateInput struct {
    URL string `json:"url" jsonschema:"description=The URL to navigate to,required"`
}

type NavigateOutput struct {
    URL   string `json:"url"`
    Title string `json:"title"`
}

func (s *Server) handleNavigate(
    ctx context.Context,
    req *mcp.CallToolRequest,
    input NavigateInput,
) (*mcp.CallToolResult, NavigateOutput, error) {
    if err := s.session.LaunchIfNeeded(ctx); err != nil {
        return nil, NavigateOutput{}, fmt.Errorf("browser not launched: %w", err)
    }

    start := time.Now()
    err := s.session.vibe.Go(ctx, input.URL)
    duration := time.Since(start)

    result := StepResult{
        ID:         fmt.Sprintf("navigate-%d", len(s.session.results)),
        Action:     "navigate",
        Args:       map[string]any{"url": input.URL},
        DurationMS: duration.Milliseconds(),
    }

    if err != nil {
        result.Status = StatusNoGo
        result.Severity = SeverityCritical
        result.Error = &StepError{
            Type:    "NavigationError",
            Message: err.Error(),
        }
        result.Screenshot = s.captureScreenshot(ctx)
        s.session.RecordStep(result)
        return nil, NavigateOutput{}, err
    }

    result.Status = StatusGo
    result.Severity = SeverityInfo
    result.Result = map[string]any{
        "url":   s.session.vibe.URL(),
        "title": s.session.vibe.Title(),
    }
    s.session.RecordStep(result)

    return nil, NavigateOutput{
        URL:   s.session.vibe.URL(),
        Title: s.session.vibe.Title(),
    }, nil
}

type ClickInput struct {
    Selector  string `json:"selector" jsonschema:"description=CSS selector for the element to click,required"`
    TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds,default=5000"`
}

type ClickOutput struct {
    Message string `json:"message"`
}

func (s *Server) handleClick(
    ctx context.Context,
    req *mcp.CallToolRequest,
    input ClickInput,
) (*mcp.CallToolResult, ClickOutput, error) {
    if input.TimeoutMS == 0 {
        input.TimeoutMS = 5000
    }
    timeout := time.Duration(input.TimeoutMS) * time.Millisecond

    start := time.Now()
    elem, err := s.session.vibe.Find(ctx, input.Selector, vibium.FindOptions{Timeout: timeout})

    result := StepResult{
        ID:     fmt.Sprintf("click-%d", len(s.session.results)),
        Action: "click",
        Args:   map[string]any{"selector": input.Selector},
    }

    if err != nil {
        result.DurationMS = time.Since(start).Milliseconds()
        result.Status = StatusNoGo
        result.Severity = SeverityCritical
        result.Error = &StepError{
            Type:      "ElementNotFoundError",
            Message:   err.Error(),
            Selector:  input.Selector,
            TimeoutMS: int64(input.TimeoutMS),
        }
        result.Error.Suggestions = s.findSimilarSelectors(ctx, input.Selector)
        result.Context = s.captureContext(ctx)
        result.Screenshot = s.captureScreenshot(ctx)
        s.session.RecordStep(result)
        return nil, ClickOutput{}, err
    }

    err = elem.Click(ctx, vibium.ActionOptions{Timeout: timeout})
    result.DurationMS = time.Since(start).Milliseconds()

    if err != nil {
        result.Status = StatusNoGo
        result.Severity = SeverityCritical
        result.Error = &StepError{
            Type:     "ClickError",
            Message:  err.Error(),
            Selector: input.Selector,
        }
        result.Screenshot = s.captureScreenshot(ctx)
        s.session.RecordStep(result)
        return nil, ClickOutput{}, err
    }

    result.Status = StatusGo
    result.Severity = SeverityInfo
    s.session.RecordStep(result)

    return nil, ClickOutput{Message: fmt.Sprintf("Clicked %s", input.Selector)}, nil
}

type GetTestReportInput struct {
    Format string `json:"format" jsonschema:"description=Report format: box (terminal) or diagnostic (full JSON for agents) or json (multi-agent-spec),enum=box,enum=diagnostic,enum=json,default=box"`
}

type GetTestReportOutput struct {
    Report string `json:"report"`
}

func (s *Server) handleGetTestReport(
    ctx context.Context,
    req *mcp.CallToolRequest,
    input GetTestReportInput,
) (*mcp.CallToolResult, GetTestReportOutput, error) {
    if input.Format == "" {
        input.Format = "box"
    }

    testResult := s.session.GetTestResult()

    switch input.Format {
    case "box":
        teamReport := report.ToTeamReport(testResult)
        rendered := masreport.Render(teamReport)
        return nil, GetTestReportOutput{Report: rendered}, nil

    case "diagnostic":
        diag := &report.DiagnosticReport{TestResult: *testResult}
        jsonBytes, err := diag.JSON()
        if err != nil {
            return nil, GetTestReportOutput{}, err
        }
        return nil, GetTestReportOutput{Report: string(jsonBytes)}, nil

    case "json":
        teamReport := report.ToTeamReport(testResult)
        jsonBytes, err := json.MarshalIndent(teamReport, "", "  ")
        if err != nil {
            return nil, GetTestReportOutput{}, err
        }
        return nil, GetTestReportOutput{Report: string(jsonBytes)}, nil

    default:
        return nil, GetTestReportOutput{}, fmt.Errorf("unknown format: %s", input.Format)
    }
}
```

## Helper Functions

```go
// mcp/tools.go

func (s *Server) findSimilarSelectors(ctx context.Context, selector string) []string {
    // Execute JS to find similar elements
    script := fmt.Sprintf(`
        (function() {
            const suggestions = [];
            // Try common variations
            const variations = [
                '%s-btn', '%s-button', 'button%s',
                '#%s', '.%s', '[data-testid="%s"]'
            ];
            for (const sel of variations) {
                try {
                    if (document.querySelector(sel)) {
                        suggestions.push(sel);
                    }
                } catch {}
            }
            // Find buttons/inputs with similar text
            document.querySelectorAll('button, input[type="submit"], a').forEach(el => {
                const text = el.textContent || el.value || '';
                if (text.toLowerCase().includes('%s'.toLowerCase())) {
                    const id = el.id ? '#' + el.id : '';
                    const cls = el.className ? '.' + el.className.split(' ')[0] : '';
                    suggestions.push(id || cls || el.tagName.toLowerCase());
                }
            });
            return [...new Set(suggestions)].slice(0, 5);
        })()
    `, selector, selector, selector, selector, selector, selector, selector)

    result, err := s.session.vibe.Evaluate(ctx, script)
    if err != nil {
        return nil
    }

    if suggestions, ok := result.([]any); ok {
        var strs []string
        for _, s := range suggestions {
            if str, ok := s.(string); ok {
                strs = append(strs, str)
            }
        }
        return strs
    }
    return nil
}

func (s *Server) captureContext(ctx context.Context) *StepContext {
    context := &StepContext{
        PageURL:   s.session.vibe.URL(),
        PageTitle: s.session.vibe.Title(),
    }

    // Get visible interactive elements
    script := `
        Array.from(document.querySelectorAll('button, input[type="submit"], a[href]'))
            .filter(el => el.offsetParent !== null)
            .map(el => el.id ? '#' + el.id : (el.className ? '.' + el.className.split(' ')[0] : el.tagName))
            .slice(0, 10)
    `
    if result, err := s.session.vibe.Evaluate(ctx, script); err == nil {
        if elems, ok := result.([]any); ok {
            for _, e := range elems {
                if str, ok := e.(string); ok {
                    context.VisibleButtons = append(context.VisibleButtons, str)
                }
            }
        }
    }

    return context
}

func (s *Server) captureScreenshot(ctx context.Context) *ScreenshotRef {
    data, err := s.session.vibe.Screenshot()
    if err != nil {
        return nil
    }
    return &ScreenshotRef{
        Base64: base64.StdEncoding.EncodeToString(data),
    }
}
```

## CLI Entry Point

```go
// cmd/vibium-mcp/main.go

package main

import (
    "flag"
    "log"
    "os"

    "github.com/grokify/vibium-go/mcp"
)

func main() {
    headless := flag.Bool("headless", true, "Run browser in headless mode")
    project := flag.String("project", "vibium-tests", "Project name for reports")
    flag.Parse()

    config := mcp.Config{
        Headless: *headless,
        Project:  *project,
    }

    server := mcp.NewServer(config)

    if err := server.Run(); err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
}
```

## Claude Code Configuration

### ~/.claude/claude_desktop_config.json

```json
{
  "mcpServers": {
    "vibium": {
      "command": "vibium-mcp",
      "args": ["--headless", "--project", "my-app"]
    }
  }
}
```

### Alternative: Local development

```json
{
  "mcpServers": {
    "vibium": {
      "command": "go",
      "args": ["run", "./cmd/vibium-mcp", "--headless"],
      "cwd": "/path/to/vibium-go"
    }
  }
}
```

## Implementation Phases

### Phase 1: Core Infrastructure

1. Create `internal/bidi/` and `internal/clicker/` from existing code
2. Create `mcp/` package structure
3. Implement `mcp/server.go` with MCP protocol handling
4. Implement `mcp/session.go` for browser lifecycle
5. Create `cmd/vibium-mcp/main.go` entry point

### Phase 2: Core Tools

1. Implement `browser_launch` / `browser_quit`
2. Implement `navigate`
3. Implement `click`
4. Implement `type`
5. Implement `get_text`
6. Implement `screenshot`

### Phase 3: Report System

1. Create `mcp/report/result.go` structs
2. Create `mcp/report/diagnostic.go`
3. Create `mcp/report/teamreport.go` (multi-agent-spec integration)
4. Implement `get_test_report` tool

### Phase 4: Advanced Tools

1. Implement `find` / `find_all`
2. Implement `assert_text` / `assert_element`
3. Implement `evaluate`
4. Implement `get_attribute`
5. Implement `wait_for`

### Phase 5: Diagnostics

1. Implement console log collection
2. Implement network error collection
3. Implement selector suggestions
4. Implement context capture

## Testing Strategy

### Unit Tests

- `mcp/report/*_test.go` - Report generation
- `mcp/session_test.go` - Session lifecycle (mock Vibe)
- `mcp/tools_test.go` - Tool handlers (mock Session)

### Integration Tests

- `mcp/integration_test.go` - Full MCP protocol test
- Test against real browser with example.com

### Manual Testing

```bash
# Build and run
go build -o vibium-mcp ./cmd/vibium-mcp
./vibium-mcp

# Test with MCP inspector
npx @anthropic-ai/mcp-inspector vibium-mcp
```

## Error Handling

All tool handlers follow this pattern:

1. Validate inputs
2. Execute operation with timeout
3. On success: record GO step, return success message
4. On failure: record NO-GO step with full error context, return error message

Errors are **never** propagated as MCP protocol errors. They are returned as tool results with error messages, allowing the agent to handle them gracefully.

## Security Considerations

- No arbitrary code execution (evaluate is limited to read-only inspection)
- No file system access except screenshot output
- No network interception
- Timeout protection on all operations
- Session isolation (one browser per MCP connection)
