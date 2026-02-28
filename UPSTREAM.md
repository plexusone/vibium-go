# Upstream Reference

This file tracks the upstream VibiumDev/vibium repository version that the Go client is aligned with.

## Current Alignment

| Field | Value |
|-------|-------|
| Repository | github.com/VibiumDev/vibium |
| Commit | `256d874e5582c4534d920a7afeae54f0e1df9a8f` |
| Date | 2026-02-24 |
| Message | Merge pull request #69 from evanniedojadlo/fix/actionability-reason-field-dropped |

## Clients Reviewed

- `clients/javascript` - v26.2.0 (package.json)
- `clients/python` - v26.2.0 (pyproject.toml)

## Feature Parity Checklist

### Phase 1: Essential Form Interactions

- [ ] `Element.Fill()` - Clear and fill form field
- [ ] `Element.Press()` - Press specific key
- [ ] `Element.Clear()` - Clear text input
- [ ] `Element.Check()` / `Uncheck()` - Checkbox interaction
- [ ] `Element.SelectOption()` - Dropdown selection
- [ ] `Element.Focus()` - Focus element
- [ ] `Element.Hover()` - Mouse hover
- [ ] `Element.ScrollIntoView()` - Scroll into viewport
- [ ] `Element.Value()` - Get input value
- [ ] `Element.InnerHTML()` - Get innerHTML
- [ ] `Element.IsVisible()` / `IsHidden()` - Visibility checks
- [ ] `Element.IsEnabled()` / `IsEditable()` - State checks
- [ ] `Element.IsChecked()` - Checkbox state
- [ ] `Element.WaitUntil(state)` - Wait for state

### Phase 2: Semantic Selectors

- [ ] Find by `role`
- [ ] Find by `text`
- [ ] Find by `label`
- [ ] Find by `placeholder`
- [ ] Find by `testid`
- [ ] Find by `xpath`
- [ ] Find by `alt`
- [ ] Find by `title`
- [ ] Find by `near`

### Phase 3: Advanced Element Actions

- [ ] `Element.DblClick()` - Double click
- [ ] `Element.DragTo()` - Drag and drop
- [ ] `Element.Tap()` - Touch tap
- [ ] `Element.DispatchEvent()` - Dispatch DOM event
- [ ] `Element.SetFiles()` - Set file input
- [ ] `Element.Screenshot()` - Element screenshot

### Phase 4: Page Enhancements

- [ ] `Vibe.Content()` - Get page HTML
- [ ] `Vibe.SetContent()` - Set page HTML
- [ ] `Vibe.Viewport()` / `SetViewport()` - Viewport control
- [ ] `Vibe.NewPage()` - Create new tab
- [ ] `Vibe.Close()` - Close page (vs Quit which closes browser)
- [ ] `Vibe.BringToFront()` - Activate page
- [ ] `Vibe.Frames()` / `Frame()` - Frame handling
- [ ] `Vibe.PDF()` - Print to PDF

### Phase 5: Input Controllers

- [ ] `Keyboard` - press/down/up/type
- [ ] `Mouse` - click/move/down/up/wheel
- [ ] `Touch` - tap

### Phase 6: Network & Events

- [ ] `Route()` / `Unroute()` - Network interception
- [ ] `OnRequest()` / `OnResponse()` - Network events
- [ ] `OnConsole()` / `ConsoleMessages()` - Console collection
- [ ] `OnDialog()` - Dialog handling
- [ ] `OnDownload()` - Download handling
- [ ] `SetHeaders()` - Set extra HTTP headers

### Phase 7: Browser Context

- [ ] `NewContext()` - Create isolated context
- [ ] `Context.Cookies()` / `SetCookies()` / `ClearCookies()`
- [ ] `Context.StorageState()`
- [ ] `Context.AddInitScript()`

### Phase 8: Advanced Features

- [ ] `Clock` control - Timer manipulation for testing
- [ ] `Tracing` - Browser trace recording
- [ ] `A11yTree()` - Accessibility tree
- [ ] `EmulateMedia()` - Media emulation
- [ ] `SetGeolocation()` - Geolocation override
- [ ] `AddScript()` / `AddStyle()` - Inject scripts/styles
- [ ] `Expose()` - Expose function to page

## How to Update

When syncing with upstream:

1. Get latest commit: `cd /path/to/VibiumDev/vibium && git pull && git rev-parse HEAD`
2. Compare JS client: `diff` or review `clients/javascript/src/`
3. Compare Python client: review `clients/python/vibium/`
4. Update this file with new commit hash
5. Implement missing features
6. Update checklist above
