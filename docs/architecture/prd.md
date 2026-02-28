# Product Requirements Document: Vibium MCP Server

## Overview

Vibium MCP Server enables AI assistants (Claude Code, ChatGPT, etc.) to perform browser automation and testing through the Model Context Protocol (MCP). It exposes browser control primitives as MCP tools, allowing agents to navigate websites, interact with elements, and execute test plans.

## Vision

Provide a seamless bridge between AI assistants and browser automation, enabling:

- Natural language-driven browser testing
- Automated UI validation
- Interactive debugging of web applications
- Deterministic test reporting

## Target Users

1. **AI Assistants** - Claude Code, ChatGPT Desktop, custom agents
2. **Developers** - Using AI assistants for browser testing
3. **QA Engineers** - Automating test execution via AI
4. **DevOps** - CI/CD integration for automated browser tests

## Use Cases

### UC-1: Execute Browser Test Plan

**Actor**: AI Assistant (Claude Code)

**Flow**:

1. User provides test plan (Markdown or verbal description)
2. Assistant reads and interprets test steps
3. Assistant executes steps via MCP browser tools
4. Assistant collects results and generates report
5. User sees deterministic box report in terminal

### UC-2: Debug Test Failure

**Actor**: AI Assistant

**Flow**:

1. Test step fails (element not found, assertion failed)
2. Assistant requests diagnostic report
3. Assistant analyzes console logs, network errors, DOM state
4. Assistant suggests fixes or alternative selectors
5. Assistant retries or reports findings to user

### UC-3: Interactive Browser Exploration

**Actor**: Developer via AI Assistant

**Flow**:

1. User asks assistant to "check the login page"
2. Assistant navigates to URL
3. Assistant extracts page information (title, elements, text)
4. Assistant reports findings conversationally
5. User asks follow-up questions

### UC-4: CI/CD Test Execution

**Actor**: CI System

**Flow**:

1. CI triggers test run with headless browser
2. MCP server executes test plan
3. Results output as JSON for CI parsing
4. Box report saved as artifact
5. CI fails/passes based on status

## Features

### F-1: Browser Control Tools

Core MCP tools for browser automation:

| Tool | Description | Priority |
|------|-------------|----------|
| `browser_launch` | Launch browser (headless/headed) | P0 |
| `browser_quit` | Close browser and cleanup | P0 |
| `navigate` | Go to URL | P0 |
| `click` | Click element by selector | P0 |
| `type` | Type text into input element | P0 |
| `get_text` | Get text content of element | P0 |
| `screenshot` | Capture page screenshot | P0 |
| `find` | Find element, return info | P1 |
| `find_all` | Find all matching elements | P1 |
| `get_attribute` | Get element attribute value | P1 |
| `evaluate` | Execute JavaScript | P1 |
| `get_title` | Get page title | P1 |
| `get_url` | Get current URL | P1 |
| `wait_for` | Wait for element to appear | P1 |
| `assert_text` | Assert text exists on page | P1 |
| `assert_element` | Assert element exists | P1 |
| `back` | Navigate back | P2 |
| `forward` | Navigate forward | P2 |
| `reload` | Reload page | P2 |

### F-2: Report Generation

Two report formats for different consumers:

#### Box Report (Deterministic)

- Fixed 78-character width
- Unicode box drawing characters
- Status icons (🟢 🟡 🔴 ⚪)
- DAG-ordered team/step rendering
- Truncated details for readability
- Compatible with multi-agent-spec

**Consumer**: Humans (terminal), CI logs

#### Diagnostic Report (Rich JSON)

- Full error messages (not truncated)
- Browser console logs
- Network request/response errors
- DOM state at failure point
- Screenshot references
- Suggested fixes (similar selectors)
- Timing information

**Consumer**: AI agents for failure analysis

### F-3: Session Management

- Auto-launch browser on first tool call
- Reuse browser session across tool calls
- Configurable headless/headed mode
- Graceful cleanup on MCP disconnect
- Timeout protection for all operations

### F-4: Test Plan Support

- Read Markdown test plans from files
- Step-by-step execution
- Pass/fail tracking per step
- Aggregate status computation
- Skip downstream steps on failure

## Report Format Specifications

### Box Report Example (Success)

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                          BROWSER TEST REPORT                                 ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  Project: vibium-tests                                                       ║
║  Target:  Login Flow                                                         ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  🟢 navigation — GO                                                          ║
║     └── 🟢 GO navigate-to-app [info]                                         ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  🟢 authentication — GO                                                      ║
║     └── 🟢 GO type-email [info]                                              ║
║     └── 🟢 GO click-submit [info]                                            ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  🟢 assertions — GO                                                          ║
║     └── 🟢 GO assert-welcome [info]                                          ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                             🟢 TEST: GO 🟢                                   ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Box Report Example (Failure)

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                          BROWSER TEST REPORT                                 ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  Project: vibium-tests                                                       ║
║  Target:  Login Flow                                                         ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  🟢 navigation — GO                                                          ║
║     └── 🟢 GO navigate-to-app [info]                                         ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  🔴 authentication — NO-GO                                                   ║
║     └── 🟢 GO type-email [info]                                              ║
║     └── 🔴 NO-GO click-submit [critical] Element not found: #submit          ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  ⚪ assertions — SKIP                                                        ║
║     └── ⚪ SKIP assert-welcome [info] Skipped: upstream failed               ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                           🔴 TEST: NO-GO 🔴                                  ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Diagnostic Report Fields

| Field | Description |
|-------|-------------|
| `test_plan` | Source test plan file |
| `status` | Overall status (GO/WARN/NO-GO/SKIP) |
| `duration_ms` | Total execution time |
| `browser` | Browser info (name, headless, viewport) |
| `steps[]` | Array of step results |
| `steps[].id` | Step identifier |
| `steps[].action` | Tool/action name |
| `steps[].args` | Arguments passed |
| `steps[].status` | Step status |
| `steps[].duration_ms` | Step execution time |
| `steps[].result` | Success result data |
| `steps[].error` | Error details (on failure) |
| `steps[].error.type` | Error type name |
| `steps[].error.message` | Full error message |
| `steps[].error.suggestions` | Alternative selectors/fixes |
| `steps[].context` | Page state at execution |
| `steps[].context.page_url` | Current URL |
| `steps[].context.page_title` | Current title |
| `steps[].context.visible_buttons` | Visible interactive elements |
| `steps[].context.dom_snippet` | Relevant DOM fragment |
| `steps[].console_logs[]` | Browser console entries |
| `steps[].network_errors[]` | Failed network requests |
| `steps[].screenshot` | Screenshot path/base64 |
| `recommendations[]` | AI-friendly fix suggestions |

## Success Criteria

### MVP (v0.2.0)

- [ ] MCP server starts and accepts connections
- [ ] Core tools work: launch, quit, navigate, click, type, get_text, screenshot
- [ ] Box report renders correctly
- [ ] Diagnostic report contains error details
- [ ] Claude Code can execute simple test plan

### v0.3.0

- [ ] All P0 and P1 tools implemented
- [ ] Console log collection
- [ ] Network error collection
- [ ] Screenshot on failure
- [ ] Selector suggestions on ElementNotFound

### v1.0.0

- [ ] Full tool suite (P0, P1, P2)
- [ ] CI/CD integration documented
- [ ] Performance benchmarks
- [ ] Cross-platform testing (macOS, Linux, Windows)

## Non-Goals

- **Not a test framework** - No built-in test assertions, loops, or conditionals
- **Not a recorder** - No browser action recording
- **Not a proxy** - No network interception/modification
- **Not multi-browser** - Chromium only (via Clicker)

## Dependencies

- Vibium-Go client library (this repo)
- Clicker binary (external)
- multi-agent-spec SDK (for report rendering)
- MCP Go SDK (for protocol handling)

## References

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [multi-agent-spec](https://github.com/agentplexus/multi-agent-spec)
- [Vibium Clicker](https://github.com/aspect-build/aspect-cli) (example)
