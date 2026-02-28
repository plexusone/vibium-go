package mcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	vibium "github.com/agentplexus/vibium-go"
	"github.com/agentplexus/vibium-go/mcp/report"
)

// Session manages a browser session and collects test results.
type Session struct {
	mu       sync.Mutex
	vibe     *vibium.Vibe
	config   SessionConfig
	results  []report.StepResult
	stepNum  int
	recorder *Recorder
}

// SessionConfig holds session configuration.
type SessionConfig struct {
	Headless       bool
	DefaultTimeout time.Duration
	Project        string
	Target         string
}

// NewSession creates a new Session.
func NewSession(config SessionConfig) *Session {
	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = 30 * time.Second
	}
	if config.Project == "" {
		config.Project = "vibium-tests"
	}
	return &Session{
		config:   config,
		results:  make([]report.StepResult, 0),
		recorder: NewRecorder(),
	}
}

// Recorder returns the session's recorder.
func (s *Session) Recorder() *Recorder {
	return s.recorder
}

// LaunchIfNeeded launches the browser if not already running.
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

// Vibe returns the browser controller, launching if needed.
func (s *Session) Vibe(ctx context.Context) (*vibium.Vibe, error) {
	if err := s.LaunchIfNeeded(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.vibe, nil
}

// RecordStep records a step result.
func (s *Session) RecordStep(result report.StepResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results = append(s.results, result)
}

// NextStepID returns the next step ID.
func (s *Session) NextStepID(action string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stepNum++
	return fmt.Sprintf("%s-%d", action, s.stepNum)
}

// GetTestResult returns the current test result.
func (s *Session) GetTestResult() *report.TestResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	steps := make([]report.StepResult, len(s.results))
	copy(steps, s.results)

	tr := &report.TestResult{
		Project:     s.config.Project,
		Target:      s.config.Target,
		Status:      report.ComputeOverallStatus(steps),
		DurationMS:  report.ComputeTotalDuration(steps),
		Steps:       steps,
		GeneratedAt: time.Now().UTC(),
	}

	// Set browser info
	tr.Browser.Name = "chromium"
	tr.Browser.Headless = s.config.Headless
	tr.Browser.Viewport.Width = 1280
	tr.Browser.Viewport.Height = 720

	return tr
}

// Reset clears the session results.
func (s *Session) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results = make([]report.StepResult, 0)
	s.stepNum = 0
}

// Close closes the browser session.
func (s *Session) Close(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.vibe != nil {
		err := s.vibe.Quit(ctx)
		s.vibe = nil
		return err
	}
	return nil
}

// SetTarget sets the test target description.
func (s *Session) SetTarget(target string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.Target = target
}

// CaptureScreenshot captures a screenshot and returns a ScreenshotRef.
func (s *Session) CaptureScreenshot(ctx context.Context) *report.ScreenshotRef {
	s.mu.Lock()
	vibe := s.vibe
	s.mu.Unlock()

	if vibe == nil {
		return nil
	}

	data, err := vibe.Screenshot(ctx)
	if err != nil {
		return nil
	}

	return &report.ScreenshotRef{
		Base64: base64.StdEncoding.EncodeToString(data),
	}
}

// CaptureContext captures the current page context.
func (s *Session) CaptureContext(ctx context.Context) *report.StepContext {
	s.mu.Lock()
	vibe := s.vibe
	s.mu.Unlock()

	if vibe == nil {
		return nil
	}

	pageURL, _ := vibe.URL(ctx)
	pageTitle, _ := vibe.Title(ctx)

	stepContext := &report.StepContext{
		PageURL:   pageURL,
		PageTitle: pageTitle,
	}

	// Get visible interactive elements
	script := `return Array.from(document.querySelectorAll('button, input[type="submit"], a[href]'))
		.filter(el => el.offsetParent !== null)
		.map(el => el.id ? '#' + el.id : (el.className ? '.' + el.className.split(' ')[0] : el.tagName))
		.slice(0, 10)`
	if result, err := vibe.Evaluate(ctx, script); err == nil {
		if elems, ok := result.([]any); ok {
			for _, e := range elems {
				if str, ok := e.(string); ok {
					stepContext.VisibleButtons = append(stepContext.VisibleButtons, str)
				}
			}
		}
	}

	return stepContext
}

// FindSimilarSelectors attempts to find similar selectors to the given one.
func (s *Session) FindSimilarSelectors(ctx context.Context, selector string) []string {
	s.mu.Lock()
	vibe := s.vibe
	s.mu.Unlock()

	if vibe == nil {
		return nil
	}

	// Extract the base selector name for variations
	baseName := selector
	if len(baseName) > 0 && (baseName[0] == '#' || baseName[0] == '.') {
		baseName = baseName[1:]
	}

	script := fmt.Sprintf(`
		(function() {
			const suggestions = [];
			const base = %q;

			// Try ID variations
			['#' + base, '#' + base + '-btn', '#' + base + '-button', '#' + base + 'Btn'].forEach(sel => {
				try { if (document.querySelector(sel)) suggestions.push(sel); } catch {}
			});

			// Try class variations
			['.' + base, '.' + base + '-btn', '.' + base + '-button'].forEach(sel => {
				try { if (document.querySelector(sel)) suggestions.push(sel); } catch {}
			});

			// Try data-testid
			try {
				const testId = document.querySelector('[data-testid="' + base + '"]');
				if (testId) suggestions.push('[data-testid="' + base + '"]');
			} catch {}

			// Find buttons/inputs with similar text
			document.querySelectorAll('button, input[type="submit"], a').forEach(el => {
				const text = (el.textContent || el.value || '').toLowerCase();
				if (text.includes(base.toLowerCase())) {
					const id = el.id ? '#' + el.id : '';
					const cls = el.className ? '.' + el.className.split(' ')[0] : '';
					if (id) suggestions.push(id);
					else if (cls) suggestions.push(cls);
				}
			});

			return [...new Set(suggestions)].slice(0, 5);
		})()
	`, baseName)

	result, err := vibe.Evaluate(ctx, script)
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
