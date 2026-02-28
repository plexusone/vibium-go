// Package render provides VPAT report renderers.
package render

import (
	"fmt"
	"strings"

	"github.com/agentplexus/vibium-go/vpat"
)

// Markdown renders a VPAT report as Markdown.
func Markdown(report *vpat.Report) string {
	var sb strings.Builder

	// Header
	sb.WriteString("# Voluntary Product Accessibility Template (VPAT)\n\n")
	sb.WriteString(fmt.Sprintf("**Standard:** %s\n\n", report.Standard))

	// Product Information
	sb.WriteString("## Product Information\n\n")
	sb.WriteString(fmt.Sprintf("| Field | Value |\n"))
	sb.WriteString(fmt.Sprintf("|-------|-------|\n"))
	sb.WriteString(fmt.Sprintf("| **Product Name** | %s |\n", report.Product.Name))
	if report.Product.Version != "" {
		sb.WriteString(fmt.Sprintf("| **Version** | %s |\n", report.Product.Version))
	}
	if report.Product.Vendor != "" {
		sb.WriteString(fmt.Sprintf("| **Vendor** | %s |\n", report.Product.Vendor))
	}
	if report.Product.URL != "" {
		sb.WriteString(fmt.Sprintf("| **Product URL** | %s |\n", report.Product.URL))
	}
	if report.Product.Description != "" {
		sb.WriteString(fmt.Sprintf("| **Description** | %s |\n", report.Product.Description))
	}
	sb.WriteString("\n")

	// Evaluation Information
	sb.WriteString("## Evaluation Information\n\n")
	sb.WriteString(fmt.Sprintf("| Field | Value |\n"))
	sb.WriteString(fmt.Sprintf("|-------|-------|\n"))
	sb.WriteString(fmt.Sprintf("| **Evaluation Date** | %s |\n", report.Evaluation.Date.Format("2006-01-02")))
	if report.Evaluation.Evaluator != "" {
		sb.WriteString(fmt.Sprintf("| **Evaluator** | %s |\n", report.Evaluation.Evaluator))
	}
	sb.WriteString(fmt.Sprintf("| **Methods** | %s |\n", strings.Join(report.Evaluation.Methods, ", ")))
	if len(report.Evaluation.Tools) > 0 {
		var tools []string
		for _, t := range report.Evaluation.Tools {
			if t.Version != "" {
				tools = append(tools, fmt.Sprintf("%s %s", t.Name, t.Version))
			} else {
				tools = append(tools, t.Name)
			}
		}
		sb.WriteString(fmt.Sprintf("| **Tools** | %s |\n", strings.Join(tools, ", ")))
	}
	if report.Evaluation.Scope != "" {
		sb.WriteString(fmt.Sprintf("| **Scope** | %s |\n", report.Evaluation.Scope))
	}
	sb.WriteString("\n")

	if len(report.Evaluation.URLs) > 0 {
		sb.WriteString("### URLs Evaluated\n\n")
		for _, url := range report.Evaluation.URLs {
			sb.WriteString(fmt.Sprintf("- %s\n", url))
		}
		sb.WriteString("\n")
	}

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("| Conformance Level | Count |\n"))
	sb.WriteString(fmt.Sprintf("|-------------------|-------|\n"))
	sb.WriteString(fmt.Sprintf("| Supports | %d |\n", report.Summary.Supports))
	sb.WriteString(fmt.Sprintf("| Partially Supports | %d |\n", report.Summary.PartiallySupports))
	sb.WriteString(fmt.Sprintf("| Does Not Support | %d |\n", report.Summary.DoesNotSupport))
	sb.WriteString(fmt.Sprintf("| Not Applicable | %d |\n", report.Summary.NotApplicable))
	sb.WriteString(fmt.Sprintf("| Not Evaluated | %d |\n", report.Summary.NotEvaluated))
	sb.WriteString(fmt.Sprintf("| **Total** | **%d** |\n", report.Summary.TotalCriteria))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("- **Automated Coverage:** %.1f%%\n", report.Summary.AutomatedCoverage))
	sb.WriteString(fmt.Sprintf("- **Total Violations Found:** %d\n\n", report.Summary.TotalViolations))

	// Detailed Results
	sb.WriteString("## Detailed Results\n\n")

	// Group by principle
	principles := []struct {
		name   string
		prefix string
	}{
		{"1. Perceivable", "1."},
		{"2. Operable", "2."},
		{"3. Understandable", "3."},
		{"4. Robust", "4."},
	}

	for _, p := range principles {
		sb.WriteString(fmt.Sprintf("### Principle %s\n\n", p.name))
		sb.WriteString("| Criteria | Conformance Level | Remarks |\n")
		sb.WriteString("|----------|-------------------|----------|\n")

		for _, c := range report.Criteria {
			if strings.HasPrefix(c.ID, p.prefix) {
				conformanceIcon := conformanceIcon(c.Conformance)
				remarks := c.Remarks
				if len(c.Violations) > 0 {
					var issues []string
					for _, v := range c.Violations {
						issues = append(issues, fmt.Sprintf("%s (%d)", v.RuleID, v.Count))
					}
					if remarks != "" {
						remarks += "; "
					}
					remarks += "Issues: " + strings.Join(issues, ", ")
				}
				sb.WriteString(fmt.Sprintf("| %s %s | %s %s | %s |\n",
					c.ID, c.Name, conformanceIcon, c.Conformance, remarks))
			}
		}
		sb.WriteString("\n")
	}

	// Violations Detail
	if report.Summary.TotalViolations > 0 {
		sb.WriteString("## Violations Detail\n\n")

		for _, c := range report.Criteria {
			if len(c.Violations) == 0 {
				continue
			}

			sb.WriteString(fmt.Sprintf("### %s %s\n\n", c.ID, c.Name))

			for _, v := range c.Violations {
				sb.WriteString(fmt.Sprintf("**%s** - %s\n\n", v.RuleID, v.Description))
				sb.WriteString(fmt.Sprintf("- Impact: %s\n", v.Impact))
				sb.WriteString(fmt.Sprintf("- Instances: %d\n", v.Count))
				if v.HelpURL != "" {
					sb.WriteString(fmt.Sprintf("- More info: %s\n", v.HelpURL))
				}
				if len(v.Elements) > 0 {
					sb.WriteString("- Example elements:\n")
					for _, el := range v.Elements {
						sb.WriteString(fmt.Sprintf("  ```html\n  %s\n  ```\n", truncate(el, 200)))
					}
				}
				sb.WriteString("\n")
			}
		}
	}

	// Notes
	if report.Notes != "" {
		sb.WriteString("## Notes\n\n")
		sb.WriteString(report.Notes)
		sb.WriteString("\n\n")
	}

	// Footer
	sb.WriteString("---\n\n")
	sb.WriteString(fmt.Sprintf("*Generated: %s*\n", report.GeneratedAt.Format("2006-01-02 15:04:05 MST")))

	return sb.String()
}

func conformanceIcon(c vpat.Conformance) string {
	switch c {
	case vpat.ConformanceSupports:
		return "[PASS]"
	case vpat.ConformancePartiallySupports:
		return "[PARTIAL]"
	case vpat.ConformanceDoesNotSupport:
		return "[FAIL]"
	case vpat.ConformanceNotApplicable:
		return "[N/A]"
	case vpat.ConformanceNotEvaluated:
		return "[?]"
	default:
		return ""
	}
}

func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.Join(strings.Fields(s), " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
