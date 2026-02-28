// Package a11y provides accessibility testing using axe-core.
package a11y

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Standard represents a WCAG accessibility standard.
type Standard string

const (
	WCAG2A   Standard = "wcag2a"
	WCAG2AA  Standard = "wcag2aa"
	WCAG2AAA Standard = "wcag2aaa"
	WCAG21A  Standard = "wcag21a"
	WCAG21AA Standard = "wcag21aa"
	WCAG21AAA Standard = "wcag21aaa"
	WCAG22AA Standard = "wcag22aa"
)

// Impact represents the severity of an accessibility violation.
type Impact string

const (
	ImpactCritical Impact = "critical"
	ImpactSerious  Impact = "serious"
	ImpactModerate Impact = "moderate"
	ImpactMinor    Impact = "minor"
)

// Options configures accessibility checking.
type Options struct {
	// Standard is the WCAG standard to check against.
	Standard Standard

	// IncludeSelector limits checking to elements matching this selector.
	IncludeSelector string

	// ExcludeSelector excludes elements matching this selector.
	ExcludeSelector string

	// Rules specifies which rules to run.
	Rules []string

	// DisabledRules specifies rules to skip.
	DisabledRules []string

	// FailOn specifies which impact levels cause failure.
	// Default is "serious" (fails on critical and serious).
	FailOn Impact
}

// DefaultOptions returns sensible defaults for WCAG 2.2 AA.
func DefaultOptions() *Options {
	return &Options{
		Standard: WCAG22AA,
		FailOn:   ImpactSerious,
	}
}

// Result contains the accessibility check results.
type Result struct {
	Violations     []Violation `json:"violations"`
	Passes         []Rule      `json:"passes"`
	Incomplete     []Rule      `json:"incomplete"`
	Inapplicable   []Rule      `json:"inapplicable"`
	URL            string      `json:"url"`
	Timestamp      string      `json:"timestamp"`
	TestEngine     TestEngine  `json:"testEngine"`
	TestRunner     TestRunner  `json:"testRunner"`
	TestEnvironment TestEnvironment `json:"testEnvironment"`
}

// TestEngine contains axe-core version info.
type TestEngine struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// TestRunner contains runner info.
type TestRunner struct {
	Name string `json:"name"`
}

// TestEnvironment contains browser info.
type TestEnvironment struct {
	UserAgent    string `json:"userAgent"`
	WindowWidth  int    `json:"windowWidth"`
	WindowHeight int    `json:"windowHeight"`
}

// Violation represents an accessibility violation.
type Violation struct {
	ID          string   `json:"id"`
	Impact      Impact   `json:"impact"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Help        string   `json:"help"`
	HelpURL     string   `json:"helpUrl"`
	Nodes       []Node   `json:"nodes"`
}

// Rule represents an axe-core rule result.
type Rule struct {
	ID          string   `json:"id"`
	Impact      Impact   `json:"impact,omitempty"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Help        string   `json:"help"`
	HelpURL     string   `json:"helpUrl"`
	Nodes       []Node   `json:"nodes"`
}

// Node represents a DOM element that violated or passed a rule.
type Node struct {
	HTML          string   `json:"html"`
	Target        []string `json:"target"`
	FailureSummary string  `json:"failureSummary,omitempty"`
	Impact        Impact   `json:"impact,omitempty"`
}

// Evaluator is the interface for JavaScript evaluation.
type Evaluator interface {
	Evaluate(ctx context.Context, script string) (interface{}, error)
}

// Check runs accessibility checks against the current page.
func Check(ctx context.Context, evaluator Evaluator, opts *Options) (*Result, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	// Build axe-core run options
	axeOpts := buildAxeOptions(opts)

	// Inject axe-core and run analysis
	script := fmt.Sprintf(`
(async function() {
	// Inject axe-core from CDN if not already present
	if (typeof axe === 'undefined') {
		await new Promise((resolve, reject) => {
			const script = document.createElement('script');
			script.src = 'https://cdnjs.cloudflare.com/ajax/libs/axe-core/4.8.4/axe.min.js';
			script.onload = resolve;
			script.onerror = reject;
			document.head.appendChild(script);
		});
	}

	// Run axe-core analysis
	const results = await axe.run(%s);
	return JSON.stringify(results);
})()
`, axeOpts)

	resultRaw, err := evaluator.Evaluate(ctx, script)
	if err != nil {
		return nil, fmt.Errorf("axe-core execution failed: %w", err)
	}

	// Parse the JSON result
	var resultStr string
	switch v := resultRaw.(type) {
	case string:
		resultStr = v
	case map[string]interface{}:
		// If the result is already parsed, re-encode it
		data, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}
		resultStr = string(data)
	default:
		return nil, fmt.Errorf("unexpected result type: %T", resultRaw)
	}

	var result Result
	if err := json.Unmarshal([]byte(resultStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse axe-core results: %w", err)
	}

	return &result, nil
}

// buildAxeOptions creates the axe-core options object.
func buildAxeOptions(opts *Options) string {
	axeOpts := make(map[string]interface{})

	// Set runOnly based on standard
	tags := standardToTags(opts.Standard)
	if len(opts.Rules) > 0 {
		// Specific rules take precedence
		axeOpts["runOnly"] = map[string]interface{}{
			"type":   "rule",
			"values": opts.Rules,
		}
	} else if len(tags) > 0 {
		axeOpts["runOnly"] = map[string]interface{}{
			"type":   "tag",
			"values": tags,
		}
	}

	// Set include/exclude context
	if opts.IncludeSelector != "" || opts.ExcludeSelector != "" {
		context := make(map[string]interface{})
		if opts.IncludeSelector != "" {
			context["include"] = []string{opts.IncludeSelector}
		}
		if opts.ExcludeSelector != "" {
			context["exclude"] = []string{opts.ExcludeSelector}
		}
		axeOpts["context"] = context
	}

	// Disable specific rules
	if len(opts.DisabledRules) > 0 {
		rules := make(map[string]interface{})
		for _, rule := range opts.DisabledRules {
			rules[rule] = map[string]bool{"enabled": false}
		}
		axeOpts["rules"] = rules
	}

	data, _ := json.Marshal(axeOpts)
	return string(data)
}

// standardToTags converts a WCAG standard to axe-core tags.
func standardToTags(standard Standard) []string {
	switch standard {
	case WCAG2A:
		return []string{"wcag2a"}
	case WCAG2AA:
		return []string{"wcag2a", "wcag2aa"}
	case WCAG2AAA:
		return []string{"wcag2a", "wcag2aa", "wcag2aaa"}
	case WCAG21A:
		return []string{"wcag2a", "wcag21a"}
	case WCAG21AA:
		return []string{"wcag2a", "wcag2aa", "wcag21a", "wcag21aa"}
	case WCAG21AAA:
		return []string{"wcag2a", "wcag2aa", "wcag2aaa", "wcag21a", "wcag21aa", "wcag21aaa"}
	case WCAG22AA:
		return []string{"wcag2a", "wcag2aa", "wcag21a", "wcag21aa", "wcag22aa"}
	default:
		return []string{"wcag2a", "wcag2aa", "wcag21a", "wcag21aa", "wcag22aa"}
	}
}

// HasFailures checks if the result has violations at or above the specified impact level.
func (r *Result) HasFailures(failOn Impact) bool {
	for _, v := range r.Violations {
		if impactSeverity(v.Impact) >= impactSeverity(failOn) {
			return true
		}
	}
	return false
}

// FilterViolations returns violations at or above the specified impact level.
func (r *Result) FilterViolations(failOn Impact) []Violation {
	var filtered []Violation
	for _, v := range r.Violations {
		if impactSeverity(v.Impact) >= impactSeverity(failOn) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// impactSeverity returns a numeric severity for comparison.
func impactSeverity(impact Impact) int {
	switch impact {
	case ImpactCritical:
		return 4
	case ImpactSerious:
		return 3
	case ImpactModerate:
		return 2
	case ImpactMinor:
		return 1
	default:
		return 0
	}
}

// Summary returns a human-readable summary of the results.
func (r *Result) Summary() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Accessibility Results for %s\n", r.URL))
	sb.WriteString(fmt.Sprintf("Engine: %s %s\n\n", r.TestEngine.Name, r.TestEngine.Version))

	if len(r.Violations) == 0 {
		sb.WriteString("No accessibility violations found.\n")
	} else {
		sb.WriteString(fmt.Sprintf("Found %d violation(s):\n\n", len(r.Violations)))

		for i, v := range r.Violations {
			sb.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, strings.ToUpper(string(v.Impact)), v.Help))
			sb.WriteString(fmt.Sprintf("   Rule: %s\n", v.ID))
			sb.WriteString(fmt.Sprintf("   More info: %s\n", v.HelpURL))
			sb.WriteString(fmt.Sprintf("   Affected elements: %d\n", len(v.Nodes)))

			for j, node := range v.Nodes {
				if j >= 3 {
					sb.WriteString(fmt.Sprintf("   ... and %d more\n", len(v.Nodes)-3))
					break
				}
				sb.WriteString(fmt.Sprintf("   - %s\n", truncate(node.HTML, 80)))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString(fmt.Sprintf("Summary: %d violations, %d passes, %d incomplete, %d inapplicable\n",
		len(r.Violations), len(r.Passes), len(r.Incomplete), len(r.Inapplicable)))

	return sb.String()
}

// SaveReport saves the full results to a JSON file.
func (r *Result) SaveReport(filename string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}

func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.Join(strings.Fields(s), " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
