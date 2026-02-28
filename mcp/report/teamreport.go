package report

import (
	"io"
	"strings"
	"time"

	multiagentspec "github.com/plexusone/multi-agent-spec/sdk/go"
)

// ToTeamReport converts a TestResult to a multi-agent-spec TeamReport.
func ToTeamReport(tr *TestResult) *multiagentspec.TeamReport {
	teams := groupStepsIntoTeams(tr.Steps)

	report := &multiagentspec.TeamReport{
		Title:       "BROWSER TEST REPORT",
		Project:     tr.Project,
		Target:      tr.Target,
		Phase:       "TEST EXECUTION",
		Status:      convertStatus(tr.Status),
		GeneratedAt: tr.GeneratedAt,
		GeneratedBy: "vibium-mcp",
		Teams:       teams,
	}

	// Add browser info as summary block
	report.SummaryBlocks = []multiagentspec.ContentBlock{
		{
			Type:  multiagentspec.ContentBlockKVPairs,
			Title: "TEST INFO",
			Pairs: []multiagentspec.KVPair{
				{Key: "Project", Value: tr.Project},
				{Key: "Target", Value: tr.Target},
				{Key: "Browser", Value: formatBrowserInfo(tr.Browser)},
				{Key: "Duration", Value: formatDuration(tr.DurationMS)},
			},
		},
	}

	// Add recommendations as footer block if present
	if len(tr.Recommendations) > 0 {
		items := make([]multiagentspec.ListItem, len(tr.Recommendations))
		for i, rec := range tr.Recommendations {
			items[i] = multiagentspec.ListItem{Text: rec}
		}
		report.FooterBlocks = []multiagentspec.ContentBlock{
			{
				Type:  multiagentspec.ContentBlockList,
				Title: "RECOMMENDATIONS",
				Items: items,
			},
		}
	}

	return report
}

// groupStepsIntoTeams groups steps into logical team sections.
func groupStepsIntoTeams(steps []StepResult) []multiagentspec.TeamSection {
	// Group steps by action category
	categories := map[string][]StepResult{
		"navigation":  {},
		"interaction": {},
		"extraction":  {},
		"assertion":   {},
		"browser":     {},
		"other":       {},
	}

	categoryOrder := []string{"browser", "navigation", "interaction", "extraction", "assertion", "other"}

	for _, step := range steps {
		cat := categorizeAction(step.Action)
		categories[cat] = append(categories[cat], step)
	}

	// Build team sections
	var teams []multiagentspec.TeamSection
	var prevID string

	for _, cat := range categoryOrder {
		catSteps := categories[cat]
		if len(catSteps) == 0 {
			continue
		}

		tasks := make([]multiagentspec.TaskResult, len(catSteps))
		for i, step := range catSteps {
			tasks[i] = convertStepToTask(step)
		}

		section := multiagentspec.TeamSection{
			ID:     cat,
			Name:   cat,
			Status: convertStatus(ComputeOverallStatus(catSteps)),
			Tasks:  tasks,
		}

		// Add dependency on previous section
		if prevID != "" {
			section.DependsOn = []string{prevID}
		}
		prevID = cat

		teams = append(teams, section)
	}

	return teams
}

// categorizeAction maps action names to categories.
func categorizeAction(action string) string {
	switch action {
	case "browser_launch", "browser_quit":
		return "browser"
	case "navigate", "back", "forward", "reload":
		return "navigation"
	case "click", "type":
		return "interaction"
	case "get_text", "get_attribute", "screenshot", "evaluate", "find", "find_all":
		return "extraction"
	case "assert_text", "assert_element", "wait_for":
		return "assertion"
	default:
		return "other"
	}
}

// convertStepToTask converts a StepResult to a multi-agent-spec TaskResult.
func convertStepToTask(step StepResult) multiagentspec.TaskResult {
	detail := formatStepDetail(step)
	if len(detail) > 33 {
		detail = detail[:30] + "..."
	}

	return multiagentspec.TaskResult{
		ID:         step.ID,
		Status:     convertStatus(step.Status),
		Severity:   string(step.Severity),
		Detail:     detail,
		DurationMs: step.DurationMS,
	}
}

// formatStepDetail creates a brief detail string for a step.
func formatStepDetail(step StepResult) string {
	if step.Error != nil {
		return step.Error.Message
	}

	// Format based on action type
	switch step.Action {
	case "navigate":
		if url, ok := step.Args["url"].(string); ok {
			return url
		}
	case "click", "type", "get_text":
		if sel, ok := step.Args["selector"].(string); ok {
			return sel
		}
	case "screenshot":
		return "captured"
	}

	return ""
}

// convertStatus converts our Status to multi-agent-spec Status.
func convertStatus(s Status) multiagentspec.Status {
	switch s {
	case StatusGo:
		return multiagentspec.StatusGo
	case StatusWarn:
		return multiagentspec.StatusWarn
	case StatusNoGo:
		return multiagentspec.StatusNoGo
	case StatusSkip:
		return multiagentspec.StatusSkip
	default:
		return multiagentspec.StatusSkip
	}
}

// formatBrowserInfo formats browser info for display.
func formatBrowserInfo(bi BrowserInfo) string {
	mode := "headed"
	if bi.Headless {
		mode = "headless"
	}
	return bi.Name + " (" + mode + ")"
}

// formatDuration formats milliseconds as a human-readable duration.
func formatDuration(ms int64) string {
	d := time.Duration(ms) * time.Millisecond
	if d < time.Second {
		return d.String()
	}
	return strings.TrimSuffix(d.Round(time.Millisecond*100).String(), "0ms") + "s"
}

// RenderBox renders a TestResult as a box-format report.
func RenderBox(tr *TestResult, w io.Writer) error {
	teamReport := ToTeamReport(tr)
	renderer := multiagentspec.NewRenderer(w)
	return renderer.Render(teamReport)
}

// RenderBoxString renders a TestResult as a box-format string.
func RenderBoxString(tr *TestResult) (string, error) {
	var sb strings.Builder
	if err := RenderBox(tr, &sb); err != nil {
		return "", err
	}
	return sb.String(), nil
}
