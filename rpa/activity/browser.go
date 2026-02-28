package activity

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	vibium "github.com/plexusone/vibium-go"
)

// NavigateActivity navigates to a URL.
type NavigateActivity struct{}

func (a *NavigateActivity) Name() string { return "browser.navigate" }

func (a *NavigateActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	url := GetString(params, "url")
	if url == "" {
		return nil, fmt.Errorf("url parameter is required")
	}

	if err := env.Vibe.Go(ctx, url); err != nil {
		return nil, fmt.Errorf("navigation failed: %w", err)
	}

	// Wait for load if specified
	if wait := GetString(params, "wait"); wait != "" {
		timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
		if err := env.Vibe.WaitForLoad(ctx, wait, timeout); err != nil {
			return nil, fmt.Errorf("wait for load failed: %w", err)
		}
	}

	return nil, nil
}

// ClickActivity clicks on an element.
type ClickActivity struct{}

func (a *ClickActivity) Name() string { return "browser.click" }

func (a *ClickActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	actionOpts := &vibium.ActionOptions{Timeout: timeout}
	if err := el.Click(ctx, actionOpts); err != nil {
		return nil, fmt.Errorf("click failed: %w", err)
	}

	return nil, nil
}

// FillActivity fills an input with a value (clears first).
type FillActivity struct{}

func (a *FillActivity) Name() string { return "browser.fill" }

func (a *FillActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	value := GetString(params, "value")

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	actionOpts := &vibium.ActionOptions{Timeout: timeout}
	if err := el.Fill(ctx, value, actionOpts); err != nil {
		return nil, fmt.Errorf("fill failed: %w", err)
	}

	return nil, nil
}

// TypeActivity types text into an element (without clearing).
type TypeActivity struct{}

func (a *TypeActivity) Name() string { return "browser.type" }

func (a *TypeActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	text := GetString(params, "text")
	if text == "" {
		return nil, fmt.Errorf("text parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	actionOpts := &vibium.ActionOptions{Timeout: timeout}
	if err := el.Type(ctx, text, actionOpts); err != nil {
		return nil, fmt.Errorf("type failed: %w", err)
	}

	return nil, nil
}

// SelectOptionActivity selects an option in a dropdown.
type SelectOptionActivity struct{}

func (a *SelectOptionActivity) Name() string { return "browser.select" }

func (a *SelectOptionActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	selectOpts := vibium.SelectOptionValues{}
	if value := GetString(params, "value"); value != "" {
		selectOpts.Values = []string{value}
	}
	if label := GetString(params, "label"); label != "" {
		selectOpts.Labels = []string{label}
	}
	if index := GetInt(params, "index"); index > 0 {
		selectOpts.Indexes = []int{index}
	}

	actionOpts := &vibium.ActionOptions{Timeout: timeout}
	if err := el.SelectOption(ctx, selectOpts, actionOpts); err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	return nil, nil
}

// CheckActivity checks a checkbox.
type CheckActivity struct{}

func (a *CheckActivity) Name() string { return "browser.check" }

func (a *CheckActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	actionOpts := &vibium.ActionOptions{Timeout: timeout}
	if err := el.Check(ctx, actionOpts); err != nil {
		return nil, fmt.Errorf("check failed: %w", err)
	}

	return nil, nil
}

// UncheckActivity unchecks a checkbox.
type UncheckActivity struct{}

func (a *UncheckActivity) Name() string { return "browser.uncheck" }

func (a *UncheckActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	actionOpts := &vibium.ActionOptions{Timeout: timeout}
	if err := el.Uncheck(ctx, actionOpts); err != nil {
		return nil, fmt.Errorf("uncheck failed: %w", err)
	}

	return nil, nil
}

// ScrollActivity scrolls the page or an element.
type ScrollActivity struct{}

func (a *ScrollActivity) Name() string { return "browser.scroll" }

func (a *ScrollActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond

	if selector != "" {
		opts := &vibium.FindOptions{Timeout: timeout}
		el, err := env.Vibe.Find(ctx, selector, opts)
		if err != nil {
			return nil, fmt.Errorf("element not found: %w", err)
		}

		actionOpts := &vibium.ActionOptions{Timeout: timeout}
		if err := el.ScrollIntoView(ctx, actionOpts); err != nil {
			return nil, fmt.Errorf("scroll failed: %w", err)
		}
	} else {
		// Scroll by delta
		deltaX := GetInt(params, "deltaX")
		deltaY := GetInt(params, "deltaY")
		script := fmt.Sprintf("window.scrollBy(%d, %d)", deltaX, deltaY)
		if _, err := env.Vibe.Evaluate(ctx, script); err != nil {
			return nil, fmt.Errorf("scroll failed: %w", err)
		}
	}

	return nil, nil
}

// ScreenshotActivity captures a screenshot.
type ScreenshotActivity struct{}

func (a *ScreenshotActivity) Name() string { return "browser.screenshot" }

func (a *ScreenshotActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")

	var data []byte
	var err error

	if selector != "" {
		timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
		opts := &vibium.FindOptions{Timeout: timeout}
		el, findErr := env.Vibe.Find(ctx, selector, opts)
		if findErr != nil {
			return nil, fmt.Errorf("element not found: %w", findErr)
		}
		data, err = el.Screenshot(ctx)
	} else {
		data, err = env.Vibe.Screenshot(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("screenshot failed: %w", err)
	}

	// Return base64 encoded data
	return base64.StdEncoding.EncodeToString(data), nil
}

// PDFActivity generates a PDF of the page.
type PDFActivity struct{}

func (a *PDFActivity) Name() string { return "browser.pdf" }

func (a *PDFActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	opts := &vibium.PDFOptions{
		PrintBackground: GetBool(params, "printBackground"),
		Landscape:       GetBool(params, "landscape"),
		DisplayHeader:   GetBool(params, "displayHeader"),
		DisplayFooter:   GetBool(params, "displayFooter"),
	}

	if scale := GetFloat(params, "scale"); scale > 0 {
		opts.Scale = scale
	}
	if format := GetString(params, "format"); format != "" {
		opts.Format = format
	}

	data, err := env.Vibe.PDF(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("PDF generation failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
