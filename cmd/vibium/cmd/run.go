package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	vibium "github.com/plexusone/vibium-go"
	"github.com/plexusone/vibium-go/script"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	runHeadless bool
	runTimeout  time.Duration
)

var runCmd = &cobra.Command{
	Use:   "run <script.yaml|script.json>",
	Short: "Run an automation script",
	Long: `Execute a series of browser automation commands from a YAML or JSON file.

Script format:
  name: My Test Script
  headless: true
  variables:
    baseUrl: https://example.com
  steps:
    - action: navigate
      url: ${baseUrl}
    - action: fill
      selector: "#email"
      value: "test@example.com"
    - action: click
      selector: "#submit"
    - action: assertAccessibility
      a11y:
        standard: wcag22aa
        failOn: serious
    - action: screenshot
      file: result.png

Available actions:
  Navigation: navigate, go, back, forward, reload
  Form: fill, type, clear, press, check, uncheck, select
  Mouse: click, dblclick, hover, focus, tap, dragTo
  Capture: screenshot, pdf
  Wait: wait, waitForSelector, waitForUrl, waitForLoad
  Assert: assertText, assertElement, assertVisible, assertHidden,
          assertUrl, assertTitle, assertAttribute, assertAccessibility
  Other: eval, setViewport, keyboardPress, keyboardType

Examples:
  vibium run test.yaml
  vibium run login.json --headless
  vibium run a11y-check.yaml --headless`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptFile := args[0]

		data, err := os.ReadFile(scriptFile)
		if err != nil {
			return fmt.Errorf("failed to read script: %w", err)
		}

		var scr script.Script
		if strings.HasSuffix(scriptFile, ".json") {
			if err := json.Unmarshal(data, &scr); err != nil {
				return fmt.Errorf("failed to parse JSON script: %w", err)
			}
		} else {
			if err := yaml.Unmarshal(data, &scr); err != nil {
				return fmt.Errorf("failed to parse YAML script: %w", err)
			}
		}

		// Override headless from CLI flag
		if cmd.Flags().Changed("headless") {
			scr.Headless = runHeadless
		}

		ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
		defer cancel()

		// Launch browser
		vibe, err := launchBrowser(ctx, scr.Headless)
		if err != nil {
			return err
		}
		defer func() {
			_ = vibe.Quit(context.Background())
			_ = clearSession()
		}()

		if scr.Name != "" {
			fmt.Printf("Running: %s\n", scr.Name)
		}

		// Execute steps
		for i, step := range scr.Steps {
			stepNum := i + 1
			stepName := step.Name
			if stepName == "" {
				stepName = describeStep(step)
			}
			if verbose {
				fmt.Printf("[%d] %s\n", stepNum, stepName)
			}

			// Substitute variables
			step = substituteVariables(step, scr.Variables)

			if err := executeStep(ctx, vibe, step); err != nil {
				if step.ContinueOnError {
					fmt.Printf("[%d] Warning: %v (continuing)\n", stepNum, err)
					continue
				}
				return fmt.Errorf("step %d (%s) failed: %w", stepNum, stepName, err)
			}
		}

		fmt.Printf("Completed %d steps\n", len(scr.Steps))
		return nil
	},
}

func substituteVariables(step script.Step, vars map[string]string) script.Step {
	if vars == nil {
		return step
	}

	subst := func(s string) string {
		for k, v := range vars {
			s = strings.ReplaceAll(s, "${"+k+"}", v)
		}
		return s
	}

	step.URL = subst(step.URL)
	step.Selector = subst(step.Selector)
	step.Value = subst(step.Value)
	step.Text = subst(step.Text)
	step.Expected = subst(step.Expected)
	step.Pattern = subst(step.Pattern)
	step.File = subst(step.File)
	step.Script = subst(step.Script)
	step.Target = subst(step.Target)

	return step
}

func describeStep(step script.Step) string {
	switch step.Action {
	case script.ActionNavigate, script.ActionGo:
		return fmt.Sprintf("navigate %s", step.URL)
	case script.ActionClick:
		return fmt.Sprintf("click %s", step.Selector)
	case script.ActionDblClick:
		return fmt.Sprintf("dblclick %s", step.Selector)
	case script.ActionType:
		return fmt.Sprintf("type %s", step.Selector)
	case script.ActionFill:
		return fmt.Sprintf("fill %s", step.Selector)
	case script.ActionClear:
		return fmt.Sprintf("clear %s", step.Selector)
	case script.ActionPress:
		return fmt.Sprintf("press %s on %s", step.Key, step.Selector)
	case script.ActionCheck:
		return fmt.Sprintf("check %s", step.Selector)
	case script.ActionUncheck:
		return fmt.Sprintf("uncheck %s", step.Selector)
	case script.ActionSelect:
		return fmt.Sprintf("select %s", step.Selector)
	case script.ActionHover:
		return fmt.Sprintf("hover %s", step.Selector)
	case script.ActionFocus:
		return fmt.Sprintf("focus %s", step.Selector)
	case script.ActionScreenshot:
		return fmt.Sprintf("screenshot %s", step.File)
	case script.ActionPDF:
		return fmt.Sprintf("pdf %s", step.File)
	case script.ActionEval:
		return "eval javascript"
	case script.ActionWait:
		return fmt.Sprintf("wait %s", step.Duration)
	case script.ActionWaitForSelector:
		return fmt.Sprintf("waitForSelector %s", step.Selector)
	case script.ActionWaitForURL:
		return fmt.Sprintf("waitForUrl %s", step.Pattern)
	case script.ActionWaitForLoad:
		return fmt.Sprintf("waitForLoad %s", step.LoadState)
	case script.ActionAssertText:
		return fmt.Sprintf("assertText %s", step.Selector)
	case script.ActionAssertElement:
		return fmt.Sprintf("assertElement %s", step.Selector)
	case script.ActionAssertVisible:
		return fmt.Sprintf("assertVisible %s", step.Selector)
	case script.ActionAssertHidden:
		return fmt.Sprintf("assertHidden %s", step.Selector)
	case script.ActionAssertURL:
		return fmt.Sprintf("assertUrl %s", step.Expected)
	case script.ActionAssertTitle:
		return fmt.Sprintf("assertTitle %s", step.Expected)
	case script.ActionAssertAccessibility:
		standard := "wcag22aa"
		if step.A11y != nil && step.A11y.Standard != "" {
			standard = step.A11y.Standard
		}
		return fmt.Sprintf("assertAccessibility (%s)", standard)
	default:
		return string(step.Action)
	}
}

func executeStep(ctx context.Context, vibe *vibium.Vibe, step script.Step) error {
	switch step.Action {
	case script.ActionNavigate, script.ActionGo:
		return vibe.Go(ctx, step.URL)

	case script.ActionBack:
		return vibe.Back(ctx)

	case script.ActionForward:
		return vibe.Forward(ctx)

	case script.ActionReload:
		return vibe.Reload(ctx)

	case script.ActionClick:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		return el.Click(ctx, nil)

	case script.ActionDblClick:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		return el.DblClick(ctx, nil)

	case script.ActionType:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		text := step.Text
		if text == "" {
			text = step.Value
		}
		return el.Type(ctx, text, nil)

	case script.ActionFill:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		value := step.Value
		if value == "" {
			value = step.Text
		}
		return el.Fill(ctx, value, nil)

	case script.ActionClear:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		return el.Clear(ctx, nil)

	case script.ActionPress:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		return el.Press(ctx, step.Key, nil)

	case script.ActionCheck:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		return el.Check(ctx, nil)

	case script.ActionUncheck:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		return el.Uncheck(ctx, nil)

	case script.ActionSelect:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		selectValues := vibium.SelectOptionValues{Values: []string{step.Value}}
		return el.SelectOption(ctx, selectValues, nil)

	case script.ActionHover:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		return el.Hover(ctx, nil)

	case script.ActionFocus:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		return el.Focus(ctx, nil)

	case script.ActionScrollIntoView:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		return el.ScrollIntoView(ctx, nil)

	case script.ActionTap:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		return el.Tap(ctx, nil)

	case script.ActionScreenshot:
		data, err := vibe.Screenshot(ctx)
		if err != nil {
			return err
		}
		return os.WriteFile(step.File, data, 0600)

	case script.ActionPDF:
		data, err := vibe.PDF(ctx, nil)
		if err != nil {
			return err
		}
		return os.WriteFile(step.File, data, 0600)

	case script.ActionEval:
		_, err := vibe.Evaluate(ctx, step.Script)
		return err

	case script.ActionWait:
		duration := step.Duration
		if duration == "" {
			duration = step.Timeout
		}
		d, err := time.ParseDuration(duration)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
		time.Sleep(d)
		return nil

	case script.ActionWaitForSelector:
		_, err := vibe.Find(ctx, step.Selector, nil)
		return err

	case script.ActionWaitForURL:
		timeout := 30 * time.Second
		if step.Timeout != "" {
			if d, err := time.ParseDuration(step.Timeout); err == nil {
				timeout = d
			}
		}
		return vibe.WaitForURL(ctx, step.Pattern, timeout)

	case script.ActionWaitForLoad:
		state := step.LoadState
		if state == "" {
			state = "load"
		}
		timeout := 30 * time.Second
		if step.Timeout != "" {
			if d, err := time.ParseDuration(step.Timeout); err == nil {
				timeout = d
			}
		}
		return vibe.WaitForLoad(ctx, state, timeout)

	case script.ActionSetViewport:
		viewport := vibium.Viewport{Width: step.Width, Height: step.Height}
		return vibe.SetViewport(ctx, viewport)

	case script.ActionKeyboardPress:
		kb, err := vibe.Keyboard(ctx)
		if err != nil {
			return err
		}
		return kb.Press(ctx, step.Key)

	case script.ActionKeyboardType:
		kb, err := vibe.Keyboard(ctx)
		if err != nil {
			return err
		}
		text := step.Text
		if text == "" {
			text = step.Value
		}
		return kb.Type(ctx, text)

	case script.ActionMouseClick:
		mouse, err := vibe.Mouse(ctx)
		if err != nil {
			return err
		}
		return mouse.Click(ctx, step.X, step.Y, nil)

	case script.ActionMouseMove:
		mouse, err := vibe.Mouse(ctx)
		if err != nil {
			return err
		}
		return mouse.Move(ctx, step.X, step.Y)

	// Assertions
	case script.ActionAssertText:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		text, err := el.Text(ctx)
		if err != nil {
			return err
		}
		if !strings.Contains(text, step.Expected) {
			return fmt.Errorf("text assertion failed: expected %q, got %q", step.Expected, text)
		}
		return nil

	case script.ActionAssertElement:
		_, err := vibe.Find(ctx, step.Selector, nil)
		return err

	case script.ActionAssertValue:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		value, err := el.Value(ctx)
		if err != nil {
			return err
		}
		if value != step.Expected {
			return fmt.Errorf("value assertion failed: expected %q, got %q", step.Expected, value)
		}
		return nil

	case script.ActionAssertVisible:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		visible, err := el.IsVisible(ctx)
		if err != nil {
			return err
		}
		if !visible {
			return fmt.Errorf("visibility assertion failed: element %s is not visible", step.Selector)
		}
		return nil

	case script.ActionAssertHidden:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			// Element not found is acceptable for assertHidden
			return nil
		}
		hidden, err := el.IsHidden(ctx)
		if err != nil {
			return err
		}
		if !hidden {
			return fmt.Errorf("hidden assertion failed: element %s is visible", step.Selector)
		}
		return nil

	case script.ActionAssertURL:
		url, err := vibe.URL(ctx)
		if err != nil {
			return err
		}
		if step.Pattern != "" {
			matched, err := regexp.MatchString(step.Pattern, url)
			if err != nil {
				return fmt.Errorf("invalid URL pattern: %w", err)
			}
			if !matched {
				return fmt.Errorf("URL assertion failed: %q does not match pattern %q", url, step.Pattern)
			}
		} else if !strings.Contains(url, step.Expected) {
			return fmt.Errorf("URL assertion failed: expected %q in %q", step.Expected, url)
		}
		return nil

	case script.ActionAssertTitle:
		title, err := vibe.Title(ctx)
		if err != nil {
			return err
		}
		if !strings.Contains(title, step.Expected) {
			return fmt.Errorf("title assertion failed: expected %q, got %q", step.Expected, title)
		}
		return nil

	case script.ActionAssertAttribute:
		el, err := vibe.Find(ctx, step.Selector, nil)
		if err != nil {
			return err
		}
		value, err := el.GetAttribute(ctx, step.Attribute)
		if err != nil {
			return err
		}
		if value != step.Expected {
			return fmt.Errorf("attribute assertion failed: expected %s=%q, got %q", step.Attribute, step.Expected, value)
		}
		return nil

	case script.ActionAssertAccessibility:
		return fmt.Errorf("assertAccessibility has moved to agent-a11y; use github.com/agentplexus/agent-a11y for accessibility testing")

	default:
		return fmt.Errorf("unknown action: %s", step.Action)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVar(&runHeadless, "headless", false, "Run browser in headless mode")
	runCmd.Flags().DurationVar(&runTimeout, "timeout", 5*time.Minute, "Total script timeout")
}
