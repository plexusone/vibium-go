package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/agentplexus/vibium-go/a11y"
	"github.com/agentplexus/vibium-go/mcp/report"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CheckAccessibilityInput defines input for check_accessibility tool.
type CheckAccessibilityInput struct {
	// Standard is the WCAG standard to check against.
	// Supported: wcag2a, wcag2aa, wcag2aaa, wcag21a, wcag21aa, wcag21aaa, wcag22aa
	Standard string `json:"standard,omitempty" jsonschema:"description=WCAG standard: wcag2a or wcag2aa or wcag2aaa or wcag21a or wcag21aa or wcag21aaa or wcag22aa (default: wcag22aa)"`

	// Include limits checking to elements matching this CSS selector.
	Include string `json:"include,omitempty" jsonschema:"description=CSS selector to limit checking scope"`

	// Exclude excludes elements matching this CSS selector.
	Exclude string `json:"exclude,omitempty" jsonschema:"description=CSS selector to exclude from checking"`

	// FailOn specifies minimum impact level: any, critical, serious, moderate, minor.
	FailOn string `json:"failOn,omitempty" jsonschema:"description=Minimum impact level to report: any or critical or serious or moderate or minor (default: serious)"`
}

// CheckAccessibilityOutput contains the accessibility check results.
type CheckAccessibilityOutput struct {
	URL             string          `json:"url"`
	ViolationCount  int             `json:"violationCount"`
	PassCount       int             `json:"passCount"`
	IncompleteCount int             `json:"incompleteCount"`
	Violations      []ViolationInfo `json:"violations,omitempty"`
	Summary         string          `json:"summary"`
}

// ViolationInfo summarizes an accessibility violation.
type ViolationInfo struct {
	ID          string   `json:"id"`
	Impact      string   `json:"impact"`
	Description string   `json:"description"`
	Help        string   `json:"help"`
	HelpURL     string   `json:"helpUrl"`
	NodeCount   int      `json:"nodeCount"`
	Nodes       []string `json:"nodes,omitempty"`
}

func (s *Server) handleCheckAccessibility(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input CheckAccessibilityInput,
) (*mcp.CallToolResult, CheckAccessibilityOutput, error) {
	start := time.Now()

	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, CheckAccessibilityOutput{}, fmt.Errorf("browser not launched: %w", err)
	}

	opts := a11y.DefaultOptions()
	if input.Standard != "" {
		opts.Standard = a11y.Standard(input.Standard)
	}
	if input.Include != "" {
		opts.IncludeSelector = input.Include
	}
	if input.Exclude != "" {
		opts.ExcludeSelector = input.Exclude
	}
	if input.FailOn != "" {
		opts.FailOn = a11y.Impact(input.FailOn)
	}

	result, err := a11y.Check(ctx, vibe, opts)
	duration := time.Since(start)

	stepResult := report.StepResult{
		ID:         s.session.NextStepID("check_accessibility"),
		Action:     "check_accessibility",
		Args:       map[string]any{"standard": string(opts.Standard), "failOn": string(opts.FailOn)},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		stepResult.Status = report.StatusNoGo
		stepResult.Severity = report.SeverityHigh
		stepResult.Error = &report.StepError{
			Type:    "AccessibilityCheckError",
			Message: err.Error(),
		}
		s.session.RecordStep(stepResult)
		return nil, CheckAccessibilityOutput{}, fmt.Errorf("accessibility check failed: %w", err)
	}

	// Convert to output format
	output := CheckAccessibilityOutput{
		URL:             result.URL,
		ViolationCount:  len(result.Violations),
		PassCount:       len(result.Passes),
		IncompleteCount: len(result.Incomplete),
		Summary:         result.Summary(),
	}

	// Include violation details (limit nodes to avoid huge output)
	for _, v := range result.Violations {
		vi := ViolationInfo{
			ID:          v.ID,
			Impact:      string(v.Impact),
			Description: v.Description,
			Help:        v.Help,
			HelpURL:     v.HelpURL,
			NodeCount:   len(v.Nodes),
		}
		// Include first few nodes
		for i, n := range v.Nodes {
			if i >= 3 {
				break
			}
			vi.Nodes = append(vi.Nodes, n.HTML)
		}
		output.Violations = append(output.Violations, vi)
	}

	// Determine severity based on violations
	if output.ViolationCount > 0 {
		stepResult.Status = report.StatusNoGo
		stepResult.Severity = report.SeverityMedium
		stepResult.Error = &report.StepError{
			Type:    "AccessibilityViolations",
			Message: fmt.Sprintf("%d accessibility violations found", output.ViolationCount),
		}
	} else {
		stepResult.Status = report.StatusGo
		stepResult.Severity = report.SeverityInfo
	}
	s.session.RecordStep(stepResult)

	// Record if recording
	if s.session.Recorder().IsRecording() {
		standard := string(opts.Standard)
		if standard == "" {
			standard = "wcag22aa"
		}
		s.session.Recorder().RecordAccessibilityCheck(standard, string(opts.FailOn))
	}

	return nil, output, nil
}

type GetA11yTreeInput struct{}

type GetA11yTreeOutput struct {
	Tree string `json:"tree"`
}

func (s *Server) handleGetA11yTree(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetA11yTreeInput,
) (*mcp.CallToolResult, GetA11yTreeOutput, error) {
	start := time.Now()

	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetA11yTreeOutput{}, fmt.Errorf("browser not launched: %w", err)
	}

	tree, err := vibe.A11yTree(ctx)
	duration := time.Since(start)

	stepResult := report.StepResult{
		ID:         s.session.NextStepID("get_a11y_tree"),
		Action:     "get_a11y_tree",
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		stepResult.Status = report.StatusNoGo
		stepResult.Severity = report.SeverityMedium
		stepResult.Error = &report.StepError{
			Type:    "A11yTreeError",
			Message: err.Error(),
		}
		s.session.RecordStep(stepResult)
		return nil, GetA11yTreeOutput{}, fmt.Errorf("failed to get a11y tree: %w", err)
	}

	// Convert interface{} to string
	var treeStr string
	switch v := tree.(type) {
	case string:
		treeStr = v
	default:
		// Try to JSON encode it
		data, err := json.MarshalIndent(tree, "", "  ")
		if err != nil {
			return nil, GetA11yTreeOutput{}, fmt.Errorf("failed to format a11y tree: %w", err)
		}
		treeStr = string(data)
	}

	stepResult.Status = report.StatusGo
	stepResult.Severity = report.SeverityInfo
	s.session.RecordStep(stepResult)

	return nil, GetA11yTreeOutput{Tree: treeStr}, nil
}
