// Package script defines the test script format for Vibium automation.
package script

// Script represents a Vibium automation test script.
// Scripts can be written in YAML or JSON format.
type Script struct {
	// Name is the human-readable name of the test script.
	Name string `json:"name" yaml:"name" jsonschema:"description=Human-readable name of the test script"`

	// Description provides additional context about what the script tests.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Additional context about what the script tests"`

	// Version is the schema version (currently 1).
	Version int `json:"version,omitempty" yaml:"version,omitempty" jsonschema:"description=Schema version (default: 1),default=1"`

	// Headless controls whether the browser runs in headless mode.
	Headless bool `json:"headless,omitempty" yaml:"headless,omitempty" jsonschema:"description=Run browser in headless mode"`

	// BaseURL is prepended to relative URLs in navigate actions.
	BaseURL string `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty" jsonschema:"description=Base URL prepended to relative URLs"`

	// Timeout is the default timeout for all steps (e.g., '30s', '1m').
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty" jsonschema:"description=Default timeout for all steps (e.g. 30s or 1m)"`

	// Variables defines reusable values that can be referenced in steps.
	Variables map[string]string `json:"variables,omitempty" yaml:"variables,omitempty" jsonschema:"description=Reusable values referenced in steps as ${varName}"`

	// Steps is the ordered list of automation steps to execute.
	Steps []Step `json:"steps" yaml:"steps" jsonschema:"description=Ordered list of automation steps,required"`
}

// Step represents a single automation action in a script.
type Step struct {
	// ID is an optional unique identifier for this step.
	ID string `json:"id,omitempty" yaml:"id,omitempty" jsonschema:"description=Optional unique identifier for this step"`

	// Name is an optional human-readable description of the step.
	Name string `json:"name,omitempty" yaml:"name,omitempty" jsonschema:"description=Human-readable description of the step"`

	// Action is the type of action to perform.
	Action Action `json:"action" yaml:"action" jsonschema:"description=Type of action to perform,required,enum=navigate,enum=go,enum=back,enum=forward,enum=reload,enum=click,enum=dblclick,enum=type,enum=fill,enum=clear,enum=press,enum=check,enum=uncheck,enum=select,enum=setFiles,enum=hover,enum=focus,enum=scrollIntoView,enum=dragTo,enum=tap,enum=screenshot,enum=pdf,enum=eval,enum=wait,enum=waitForSelector,enum=waitForUrl,enum=waitForLoad,enum=setViewport,enum=newPage,enum=closePage,enum=keyboardPress,enum=keyboardType,enum=mouseClick,enum=mouseMove,enum=assertText,enum=assertElement,enum=assertValue,enum=assertVisible,enum=assertHidden,enum=assertUrl,enum=assertTitle,enum=assertAttribute,enum=assertAccessibility,enum=getText,enum=getValue,enum=getAttribute,enum=getUrl,enum=getTitle"`

	// Selector is the CSS selector for element actions.
	Selector string `json:"selector,omitempty" yaml:"selector,omitempty" jsonschema:"description=CSS selector for element actions"`

	// URL is the target URL for navigation actions.
	URL string `json:"url,omitempty" yaml:"url,omitempty" jsonschema:"description=Target URL for navigation actions"`

	// Value is the input value for fill, type, select actions.
	Value string `json:"value,omitempty" yaml:"value,omitempty" jsonschema:"description=Input value for fill/type/select actions"`

	// Text is an alias for Value (for readability in type actions).
	Text string `json:"text,omitempty" yaml:"text,omitempty" jsonschema:"description=Alias for value (used in type actions)"`

	// Key is the key to press for press/keyboard actions.
	Key string `json:"key,omitempty" yaml:"key,omitempty" jsonschema:"description=Key to press (e.g. Enter or Tab or ArrowDown)"`

	// Script is the JavaScript code for eval actions.
	Script string `json:"script,omitempty" yaml:"script,omitempty" jsonschema:"description=JavaScript code for eval actions"`

	// File is the output file path for screenshot/pdf actions.
	File string `json:"file,omitempty" yaml:"file,omitempty" jsonschema:"description=Output file path for screenshot/pdf actions"`

	// Files is a list of file paths for file input actions.
	Files []string `json:"files,omitempty" yaml:"files,omitempty" jsonschema:"description=File paths for file input actions"`

	// Timeout overrides the default timeout for this step.
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty" jsonschema:"description=Timeout override for this step"`

	// Duration is the wait duration for wait actions.
	Duration string `json:"duration,omitempty" yaml:"duration,omitempty" jsonschema:"description=Duration for wait actions (e.g. 1s or 500ms)"`

	// FullPage captures the full page for screenshot actions.
	FullPage bool `json:"fullPage,omitempty" yaml:"fullPage,omitempty" jsonschema:"description=Capture full page for screenshots"`

	// Target is the destination element for drag actions.
	Target string `json:"target,omitempty" yaml:"target,omitempty" jsonschema:"description=Destination selector for drag actions"`

	// X is the X coordinate for mouse actions.
	X float64 `json:"x,omitempty" yaml:"x,omitempty" jsonschema:"description=X coordinate for mouse actions"`

	// Y is the Y coordinate for mouse actions.
	Y float64 `json:"y,omitempty" yaml:"y,omitempty" jsonschema:"description=Y coordinate for mouse actions"`

	// Width is the viewport width for setViewport actions.
	Width int `json:"width,omitempty" yaml:"width,omitempty" jsonschema:"description=Viewport width for setViewport actions"`

	// Height is the viewport height for setViewport actions.
	Height int `json:"height,omitempty" yaml:"height,omitempty" jsonschema:"description=Viewport height for setViewport actions"`

	// State is the expected state for wait actions (visible, hidden, attached, detached).
	State string `json:"state,omitempty" yaml:"state,omitempty" jsonschema:"description=Expected state for wait actions,enum=visible,enum=hidden,enum=attached,enum=detached"`

	// Pattern is the URL pattern for waitForUrl actions.
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty" jsonschema:"description=URL pattern for waitForUrl actions"`

	// LoadState is the load state for waitForLoad actions.
	LoadState string `json:"loadState,omitempty" yaml:"loadState,omitempty" jsonschema:"description=Load state for waitForLoad actions,enum=load,enum=domcontentloaded,enum=networkidle"`

	// Expected is the expected value for assertion actions.
	Expected string `json:"expected,omitempty" yaml:"expected,omitempty" jsonschema:"description=Expected value for assertion actions"`

	// Attribute is the attribute name for getAttribute actions.
	Attribute string `json:"attribute,omitempty" yaml:"attribute,omitempty" jsonschema:"description=Attribute name for getAttribute actions"`

	// Store saves the result to a variable for later use.
	Store string `json:"store,omitempty" yaml:"store,omitempty" jsonschema:"description=Variable name to store the result"`

	// ContinueOnError allows the script to continue if this step fails.
	ContinueOnError bool `json:"continueOnError,omitempty" yaml:"continueOnError,omitempty" jsonschema:"description=Continue script execution if this step fails"`

	// A11y specifies accessibility check options for assertAccessibility action.
	A11y *A11yOptions `json:"a11y,omitempty" yaml:"a11y,omitempty" jsonschema:"description=Accessibility check options for assertAccessibility action"`
}

// A11yOptions configures accessibility checking behavior.
type A11yOptions struct {
	// Standard is the WCAG standard to check against.
	// Supported values: "wcag2a", "wcag2aa", "wcag2aaa", "wcag21a", "wcag21aa", "wcag21aaa", "wcag22aa"
	// Default is "wcag22aa" (WCAG 2.2 Level AA).
	Standard string `json:"standard,omitempty" yaml:"standard,omitempty" jsonschema:"description=WCAG standard to check against,enum=wcag2a,enum=wcag2aa,enum=wcag2aaa,enum=wcag21a,enum=wcag21aa,enum=wcag21aaa,enum=wcag22aa,default=wcag22aa"`

	// IncludeSelector limits checking to elements matching this selector.
	IncludeSelector string `json:"include,omitempty" yaml:"include,omitempty" jsonschema:"description=CSS selector to limit checking scope"`

	// ExcludeSelector excludes elements matching this selector from checking.
	ExcludeSelector string `json:"exclude,omitempty" yaml:"exclude,omitempty" jsonschema:"description=CSS selector to exclude from checking"`

	// Rules specifies which axe-core rules to run. Empty means all applicable rules.
	Rules []string `json:"rules,omitempty" yaml:"rules,omitempty" jsonschema:"description=Specific axe-core rule IDs to run"`

	// DisabledRules specifies rules to skip.
	DisabledRules []string `json:"disabledRules,omitempty" yaml:"disabledRules,omitempty" jsonschema:"description=Rule IDs to skip"`

	// ReportFile saves the full accessibility report to this file.
	ReportFile string `json:"reportFile,omitempty" yaml:"reportFile,omitempty" jsonschema:"description=File path to save full accessibility report (JSON)"`

	// FailOn specifies which violation impacts cause failure.
	// Values: "any", "critical", "serious", "moderate", "minor"
	// Default is "serious" (fails on critical and serious violations).
	FailOn string `json:"failOn,omitempty" yaml:"failOn,omitempty" jsonschema:"description=Minimum impact level to fail on,enum=any,enum=critical,enum=serious,enum=moderate,enum=minor,default=serious"`
}

// Action represents the type of automation action.
type Action string

const (
	// Navigation actions
	ActionNavigate Action = "navigate"
	ActionGo       Action = "go" // Alias for navigate
	ActionBack     Action = "back"
	ActionForward  Action = "forward"
	ActionReload   Action = "reload"

	// Basic interactions
	ActionClick    Action = "click"
	ActionDblClick Action = "dblclick"
	ActionType     Action = "type"
	ActionFill     Action = "fill"
	ActionClear    Action = "clear"
	ActionPress    Action = "press"

	// Form controls
	ActionCheck    Action = "check"
	ActionUncheck  Action = "uncheck"
	ActionSelect   Action = "select"
	ActionSetFiles Action = "setFiles"

	// Element interactions
	ActionHover          Action = "hover"
	ActionFocus          Action = "focus"
	ActionScrollIntoView Action = "scrollIntoView"
	ActionDragTo         Action = "dragTo"
	ActionTap            Action = "tap"

	// Capture actions
	ActionScreenshot Action = "screenshot"
	ActionPDF        Action = "pdf"

	// JavaScript
	ActionEval Action = "eval"

	// Waiting
	ActionWait            Action = "wait"
	ActionWaitForSelector Action = "waitForSelector"
	ActionWaitForURL      Action = "waitForUrl"
	ActionWaitForLoad     Action = "waitForLoad"

	// Page actions
	ActionSetViewport Action = "setViewport"
	ActionNewPage     Action = "newPage"
	ActionClosePage   Action = "closePage"

	// Keyboard actions
	ActionKeyboardPress Action = "keyboardPress"
	ActionKeyboardType  Action = "keyboardType"

	// Mouse actions
	ActionMouseClick Action = "mouseClick"
	ActionMouseMove  Action = "mouseMove"

	// Assertions
	ActionAssertText          Action = "assertText"
	ActionAssertElement       Action = "assertElement"
	ActionAssertValue         Action = "assertValue"
	ActionAssertVisible       Action = "assertVisible"
	ActionAssertHidden        Action = "assertHidden"
	ActionAssertURL           Action = "assertUrl"
	ActionAssertTitle         Action = "assertTitle"
	ActionAssertAttribute     Action = "assertAttribute"
	ActionAssertAccessibility Action = "assertAccessibility"

	// Data extraction
	ActionGetText      Action = "getText"
	ActionGetValue     Action = "getValue"
	ActionGetAttribute Action = "getAttribute"
	ActionGetURL       Action = "getUrl"
	ActionGetTitle     Action = "getTitle"
)

// AllActions returns all valid action types.
func AllActions() []Action {
	return []Action{
		ActionNavigate, ActionGo, ActionBack, ActionForward, ActionReload,
		ActionClick, ActionDblClick, ActionType, ActionFill, ActionClear, ActionPress,
		ActionCheck, ActionUncheck, ActionSelect, ActionSetFiles,
		ActionHover, ActionFocus, ActionScrollIntoView, ActionDragTo, ActionTap,
		ActionScreenshot, ActionPDF,
		ActionEval,
		ActionWait, ActionWaitForSelector, ActionWaitForURL, ActionWaitForLoad,
		ActionSetViewport, ActionNewPage, ActionClosePage,
		ActionKeyboardPress, ActionKeyboardType,
		ActionMouseClick, ActionMouseMove,
		ActionAssertText, ActionAssertElement, ActionAssertValue, ActionAssertVisible,
		ActionAssertHidden, ActionAssertURL, ActionAssertTitle, ActionAssertAttribute,
		ActionAssertAccessibility,
		ActionGetText, ActionGetValue, ActionGetAttribute, ActionGetURL, ActionGetTitle,
	}
}
