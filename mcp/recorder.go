package mcp

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/plexusone/vibium-go/script"
)

// Recorder captures MCP tool calls and converts them to a script.
type Recorder struct {
	mu        sync.Mutex
	recording bool
	steps     []script.Step
	startTime time.Time
	metadata  RecorderMetadata
}

// RecorderMetadata contains metadata about the recording session.
type RecorderMetadata struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	BaseURL     string `json:"baseUrl,omitempty"`
}

// NewRecorder creates a new Recorder.
func NewRecorder() *Recorder {
	return &Recorder{
		steps: make([]script.Step, 0),
	}
}

// Start begins recording actions.
func (r *Recorder) Start(metadata RecorderMetadata) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recording = true
	r.steps = make([]script.Step, 0)
	r.startTime = time.Now()
	r.metadata = metadata
}

// Stop ends recording.
func (r *Recorder) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recording = false
}

// IsRecording returns whether recording is active.
func (r *Recorder) IsRecording() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.recording
}

// Clear removes all recorded steps.
func (r *Recorder) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.steps = make([]script.Step, 0)
}

// AddStep records a step if recording is active.
func (r *Recorder) AddStep(step script.Step) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.recording {
		r.steps = append(r.steps, step)
	}
}

// Steps returns a copy of the recorded steps.
func (r *Recorder) Steps() []script.Step {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]script.Step, len(r.steps))
	copy(result, r.steps)
	return result
}

// StepCount returns the number of recorded steps.
func (r *Recorder) StepCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.steps)
}

// Export returns the recorded session as a Script.
func (r *Recorder) Export() *script.Script {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := r.metadata.Name
	if name == "" {
		name = "Recorded Test"
	}

	return &script.Script{
		Name:        name,
		Description: r.metadata.Description,
		Version:     1,
		BaseURL:     r.metadata.BaseURL,
		Steps:       r.steps,
	}
}

// ExportJSON returns the recorded session as JSON.
func (r *Recorder) ExportJSON() ([]byte, error) {
	s := r.Export()
	return json.MarshalIndent(s, "", "  ")
}

// RecordNavigate records a navigation action.
func (r *Recorder) RecordNavigate(url string) {
	r.AddStep(script.Step{
		Action: script.ActionNavigate,
		URL:    url,
	})
}

// RecordClick records a click action.
func (r *Recorder) RecordClick(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionClick,
		Selector: selector,
	})
}

// RecordDblClick records a double-click action.
func (r *Recorder) RecordDblClick(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionDblClick,
		Selector: selector,
	})
}

// RecordType records a type action.
func (r *Recorder) RecordType(selector, text string) {
	r.AddStep(script.Step{
		Action:   script.ActionType,
		Selector: selector,
		Text:     text,
	})
}

// RecordFill records a fill action.
func (r *Recorder) RecordFill(selector, value string) {
	r.AddStep(script.Step{
		Action:   script.ActionFill,
		Selector: selector,
		Value:    value,
	})
}

// RecordClear records a clear action.
func (r *Recorder) RecordClear(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionClear,
		Selector: selector,
	})
}

// RecordPress records a press action.
func (r *Recorder) RecordPress(selector, key string) {
	r.AddStep(script.Step{
		Action:   script.ActionPress,
		Selector: selector,
		Key:      key,
	})
}

// RecordCheck records a check action.
func (r *Recorder) RecordCheck(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionCheck,
		Selector: selector,
	})
}

// RecordUncheck records an uncheck action.
func (r *Recorder) RecordUncheck(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionUncheck,
		Selector: selector,
	})
}

// RecordSelect records a select action.
func (r *Recorder) RecordSelect(selector, value string) {
	r.AddStep(script.Step{
		Action:   script.ActionSelect,
		Selector: selector,
		Value:    value,
	})
}

// RecordHover records a hover action.
func (r *Recorder) RecordHover(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionHover,
		Selector: selector,
	})
}

// RecordFocus records a focus action.
func (r *Recorder) RecordFocus(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionFocus,
		Selector: selector,
	})
}

// RecordScrollIntoView records a scroll action.
func (r *Recorder) RecordScrollIntoView(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionScrollIntoView,
		Selector: selector,
	})
}

// RecordScreenshot records a screenshot action.
func (r *Recorder) RecordScreenshot(file string, fullPage bool) {
	r.AddStep(script.Step{
		Action:   script.ActionScreenshot,
		File:     file,
		FullPage: fullPage,
	})
}

// RecordEval records an eval action.
func (r *Recorder) RecordEval(js string) {
	r.AddStep(script.Step{
		Action: script.ActionEval,
		Script: js,
	})
}

// RecordWait records a wait action.
func (r *Recorder) RecordWait(duration string) {
	r.AddStep(script.Step{
		Action:   script.ActionWait,
		Duration: duration,
	})
}

// RecordWaitForSelector records a waitForSelector action.
func (r *Recorder) RecordWaitForSelector(selector, state string) {
	r.AddStep(script.Step{
		Action:   script.ActionWaitForSelector,
		Selector: selector,
		State:    state,
	})
}

// RecordWaitForURL records a waitForUrl action.
func (r *Recorder) RecordWaitForURL(pattern string) {
	r.AddStep(script.Step{
		Action:  script.ActionWaitForURL,
		Pattern: pattern,
	})
}

// RecordWaitForLoad records a waitForLoad action.
func (r *Recorder) RecordWaitForLoad(state string) {
	r.AddStep(script.Step{
		Action:    script.ActionWaitForLoad,
		LoadState: state,
	})
}

// RecordSetViewport records a setViewport action.
func (r *Recorder) RecordSetViewport(width, height int) {
	r.AddStep(script.Step{
		Action: script.ActionSetViewport,
		Width:  width,
		Height: height,
	})
}

// RecordBack records a back action.
func (r *Recorder) RecordBack() {
	r.AddStep(script.Step{
		Action: script.ActionBack,
	})
}

// RecordForward records a forward action.
func (r *Recorder) RecordForward() {
	r.AddStep(script.Step{
		Action: script.ActionForward,
	})
}

// RecordReload records a reload action.
func (r *Recorder) RecordReload() {
	r.AddStep(script.Step{
		Action: script.ActionReload,
	})
}

// RecordAssertText records an assertText action.
func (r *Recorder) RecordAssertText(selector, expected string) {
	r.AddStep(script.Step{
		Action:   script.ActionAssertText,
		Selector: selector,
		Expected: expected,
	})
}

// RecordAssertElement records an assertElement action.
func (r *Recorder) RecordAssertElement(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionAssertElement,
		Selector: selector,
	})
}

// RecordAssertVisible records an assertVisible action.
func (r *Recorder) RecordAssertVisible(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionAssertVisible,
		Selector: selector,
	})
}

// RecordAssertURL records an assertUrl action.
func (r *Recorder) RecordAssertURL(expected string) {
	r.AddStep(script.Step{
		Action:   script.ActionAssertURL,
		Expected: expected,
	})
}

// RecordAssertTitle records an assertTitle action.
func (r *Recorder) RecordAssertTitle(expected string) {
	r.AddStep(script.Step{
		Action:   script.ActionAssertTitle,
		Expected: expected,
	})
}

// RecordMouseClick records a mouseClick action.
func (r *Recorder) RecordMouseClick(x, y float64) {
	r.AddStep(script.Step{
		Action: script.ActionMouseClick,
		X:      x,
		Y:      y,
	})
}

// RecordMouseMove records a mouseMove action.
func (r *Recorder) RecordMouseMove(x, y float64) {
	r.AddStep(script.Step{
		Action: script.ActionMouseMove,
		X:      x,
		Y:      y,
	})
}

// RecordKeyboardPress records a keyboardPress action.
func (r *Recorder) RecordKeyboardPress(key string) {
	r.AddStep(script.Step{
		Action: script.ActionKeyboardPress,
		Key:    key,
	})
}

// RecordKeyboardType records a keyboardType action.
func (r *Recorder) RecordKeyboardType(text string) {
	r.AddStep(script.Step{
		Action: script.ActionKeyboardType,
		Text:   text,
	})
}

// RecordDragTo records a dragTo action.
func (r *Recorder) RecordDragTo(selector, target string) {
	r.AddStep(script.Step{
		Action:   script.ActionDragTo,
		Selector: selector,
		Target:   target,
	})
}

// RecordTap records a tap action.
func (r *Recorder) RecordTap(selector string) {
	r.AddStep(script.Step{
		Action:   script.ActionTap,
		Selector: selector,
	})
}

// RecordSetFiles records a setFiles action.
func (r *Recorder) RecordSetFiles(selector string, files []string) {
	r.AddStep(script.Step{
		Action:   script.ActionSetFiles,
		Selector: selector,
		Files:    files,
	})
}

// RecordAccessibilityCheck records an assertAccessibility action.
func (r *Recorder) RecordAccessibilityCheck(standard, failOn string) {
	r.AddStep(script.Step{
		Action: script.ActionAssertAccessibility,
		A11y: &script.A11yOptions{
			Standard: standard,
			FailOn:   failOn,
		},
	})
}
