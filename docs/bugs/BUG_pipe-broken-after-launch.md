# Bug: Broken Pipe After browser_launch in Pipe Mode

**ID**: BUG_pipe-broken-after-launch
**Severity**: P0 — blocks all MCP usage
**Status**: ✅ Fixed (commit 447bb95)
**Date**: 2026-03-29
**Version**: w3pilot-mcp v0.2.0 (binary built 2026-03-29 13:46)
**Clicker**: `clicker vdev` (built 2026-03-26)
**OS**: macOS (arm64)

## Fix

**Root cause**: `exec.CommandContext(ctx)` was used to spawn clicker, which kills the child process when the context is cancelled. In MCP, each tool call has its own request context that gets cancelled after the tool returns. This killed clicker after `browser_launch`, causing subsequent commands to fail.

**Solution** (transport_pipe.go):
- Use `exec.Command()` instead of `exec.CommandContext()` for clicker process
- Capture `browsingContext` from `contextCreated` event during startup
- Log clicker stderr when `W3PILOT_DEBUG=1` is set

## Summary

`browser_launch` succeeds but every subsequent command (`page_navigate`, `page_get_url`, etc.) fails with:

```
failed to get browsing context: failed to send command: write |1: broken pipe
```

This blocks all MCP server usage.

## Reproduction

### Via MCP (Kiro CLI)

```json
// Step 1: succeeds
{"method": "tools/call", "params": {"name": "browser_launch", "arguments": {"headless": true}}}
// Response: {"message": "Browser launched successfully"}

// Step 2: fails
{"method": "tools/call", "params": {"name": "page_navigate", "arguments": {"url": "https://example.com"}}}
// Response: "navigation failed: failed to get browsing context: failed to send command: write |1: broken pipe"
```

### Via FIFO (standalone, no Kiro)

```bash
mkfifo /tmp/w3pilot-fifo
/Users/johnwang/go/bin/w3pilot-mcp -headless=true -timeout=10s < /tmp/w3pilot-fifo > /tmp/out.log 2>/tmp/err.log &

(
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.1"}}}'
sleep 2
echo '{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}'
sleep 1
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"browser_launch","arguments":{"headless":true}}}'
sleep 3
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"page_navigate","arguments":{"url":"https://example.com"}}}'
sleep 5
) > /tmp/w3pilot-fifo

cat /tmp/out.log
```

Result: `browser_launch` returns success, `page_navigate` returns broken pipe.

### Clicker works standalone

```bash
echo "test" | clicker pipe 2>&1
```

Output confirms clicker launches fine:
```
[router] Launching browser for client 1...
[router] Browser launched for client 1 (BiDi session)
{"method":"browsingContext.contextCreated","params":{"context":"FD5508EE...","url":"about:blank",...}}
{"method":"vibium:lifecycle.ready","params":{"version":"dev"}}
[router] Closing browser session for client 1
```

## Root Cause Analysis

The issue appears to be in `transport_pipe.go` / `pilot.go` pipe mode initialization.

### Observation 1: Browsing context not captured in pipe mode

In `launchPipe()` (`pilot.go:69-91`), after `transport.Start()` returns:
- The `browsingContext.contextCreated` event has already been received (clicker sends it before `vibium:lifecycle.ready`)
- But `launchPipe` never registers a handler for it
- `pilot.browsingContext` remains empty

Compare with `transport_ws.go:62-63` which explicitly listens for `browsingContext.contextCreated`.

### Observation 2: getContext triggers the broken pipe

When `page_navigate` is called, `getContext()` (`pilot.go:210+`) finds `browsingContext` is empty and sends `browsingContext.getTree` over the pipe. By this point the pipe may already be in a bad state, causing the write to fail.

### Observation 3: readStderr discards clicker output

`readStderr()` (`transport_pipe.go:119-124`) silently discards all stderr. If clicker is logging errors about why it's closing, we'd never see them.

## Suggested Fix

In `launchPipe()`, register a handler for `browsingContext.contextCreated` before calling `transport.Start()`, similar to how `waitForReady` registers for `vibium:lifecycle.ready`:

```go
func (b *browserLauncher) launchPipe(ctx context.Context, opts *LaunchOptions) (*Pilot, error) {
    transport := newPipeTransport()

    // Capture browsing context from the contextCreated event
    contextCh := make(chan string, 1)
    transport.OnEvent("browsingContext.contextCreated", func(event *BiDiEvent) {
        var params struct {
            Context string `json:"context"`
        }
        if err := json.Unmarshal(event.Params, &params); err == nil {
            select {
            case contextCh <- params.Context:
            default:
            }
        }
    })

    if err := transport.Start(ctx, pipeOpts); err != nil {
        return nil, err
    }

    // Wait for context (should arrive before or with lifecycle.ready)
    select {
    case ctx := <-contextCh:
        pilot.browsingContext = ctx
    case <-time.After(5 * time.Second):
        // Fall back to getTree
    }

    // ... rest of initialization
}
```

Also consider logging stderr instead of discarding it, at least when debug mode is enabled.

## Environment

```
w3pilot-mcp: v0.2.0 (17,860,338 bytes, Mach-O arm64)
clicker:     vdev (15,272,146 bytes, ~/go/bin/clicker)
Go:          (check go version)
macOS:       arm64
MCP config:  ~/.kiro/settings/mcp.json
  args: ["-headless=false", "-timeout=60s", "-project=achmea-dast"]
```

## Impact

All w3pilot MCP tool calls fail after `browser_launch`. The MCP server is unusable for any browser automation task.

Previously working config with `vibium-mcp` binary (same architecture) at:
`~/go/src/gitlab.com/saviynt/product/appsec/analysis_CNP-CC-EIC-TRUNK/.kiro/settings/mcp.json`
