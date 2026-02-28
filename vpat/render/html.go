package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/agentplexus/vibium-go/vpat"
)

// HTML renders a VPAT report as HTML following the ITI VPAT 2.5 format.
func HTML(report *vpat.Report) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>VPAT - `)
	sb.WriteString(html.EscapeString(report.Product.Name))
	sb.WriteString(`</title>
<style>
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  line-height: 1.6;
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  color: #333;
}
h1, h2, h3 { color: #1a1a2e; }
table {
  width: 100%;
  border-collapse: collapse;
  margin: 1em 0;
}
th, td {
  border: 1px solid #ddd;
  padding: 12px;
  text-align: left;
}
th {
  background-color: #f5f5f5;
  font-weight: 600;
}
tr:nth-child(even) { background-color: #fafafa; }
.supports { color: #22863a; font-weight: bold; }
.partially-supports { color: #b08800; font-weight: bold; }
.does-not-support { color: #cb2431; font-weight: bold; }
.not-applicable { color: #6a737d; }
.not-evaluated { color: #6a737d; font-style: italic; }
.summary-box {
  background: #f6f8fa;
  border: 1px solid #e1e4e8;
  border-radius: 6px;
  padding: 16px;
  margin: 1em 0;
}
.violation {
  background: #ffeef0;
  border-left: 4px solid #cb2431;
  padding: 12px;
  margin: 1em 0;
}
.violation code {
  background: #f6f8fa;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 0.9em;
}
.note {
  background: #fff8c5;
  border: 1px solid #f9c513;
  border-radius: 6px;
  padding: 16px;
  margin: 1em 0;
}
footer {
  margin-top: 2em;
  padding-top: 1em;
  border-top: 1px solid #e1e4e8;
  color: #6a737d;
  font-size: 0.9em;
}
</style>
</head>
<body>
`)

	// Header
	sb.WriteString("<h1>Voluntary Product Accessibility Template (VPAT)</h1>\n")
	sb.WriteString(fmt.Sprintf("<p><strong>Standard:</strong> %s</p>\n", html.EscapeString(report.Standard)))

	// Product Information
	sb.WriteString("<h2>Product Information</h2>\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>Field</th><th>Value</th></tr>\n")
	sb.WriteString(fmt.Sprintf("<tr><td>Product Name</td><td>%s</td></tr>\n", html.EscapeString(report.Product.Name)))
	if report.Product.Version != "" {
		sb.WriteString(fmt.Sprintf("<tr><td>Version</td><td>%s</td></tr>\n", html.EscapeString(report.Product.Version)))
	}
	if report.Product.Vendor != "" {
		sb.WriteString(fmt.Sprintf("<tr><td>Vendor</td><td>%s</td></tr>\n", html.EscapeString(report.Product.Vendor)))
	}
	if report.Product.URL != "" {
		sb.WriteString(fmt.Sprintf("<tr><td>Product URL</td><td><a href=\"%s\">%s</a></td></tr>\n",
			html.EscapeString(report.Product.URL), html.EscapeString(report.Product.URL)))
	}
	if report.Product.Description != "" {
		sb.WriteString(fmt.Sprintf("<tr><td>Description</td><td>%s</td></tr>\n", html.EscapeString(report.Product.Description)))
	}
	sb.WriteString("</table>\n")

	// Evaluation Information
	sb.WriteString("<h2>Evaluation Information</h2>\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>Field</th><th>Value</th></tr>\n")
	sb.WriteString(fmt.Sprintf("<tr><td>Evaluation Date</td><td>%s</td></tr>\n", report.Evaluation.Date.Format("2006-01-02")))
	if report.Evaluation.Evaluator != "" {
		sb.WriteString(fmt.Sprintf("<tr><td>Evaluator</td><td>%s</td></tr>\n", html.EscapeString(report.Evaluation.Evaluator)))
	}
	sb.WriteString(fmt.Sprintf("<tr><td>Methods</td><td>%s</td></tr>\n", html.EscapeString(strings.Join(report.Evaluation.Methods, ", "))))
	if len(report.Evaluation.Tools) > 0 {
		var tools []string
		for _, t := range report.Evaluation.Tools {
			if t.Version != "" {
				tools = append(tools, fmt.Sprintf("%s %s", t.Name, t.Version))
			} else {
				tools = append(tools, t.Name)
			}
		}
		sb.WriteString(fmt.Sprintf("<tr><td>Tools</td><td>%s</td></tr>\n", html.EscapeString(strings.Join(tools, ", "))))
	}
	if report.Evaluation.Scope != "" {
		sb.WriteString(fmt.Sprintf("<tr><td>Scope</td><td>%s</td></tr>\n", html.EscapeString(report.Evaluation.Scope)))
	}
	sb.WriteString("</table>\n")

	if len(report.Evaluation.URLs) > 0 {
		sb.WriteString("<h3>URLs Evaluated</h3>\n<ul>\n")
		for _, url := range report.Evaluation.URLs {
			sb.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a></li>\n",
				html.EscapeString(url), html.EscapeString(url)))
		}
		sb.WriteString("</ul>\n")
	}

	// Summary
	sb.WriteString("<h2>Summary</h2>\n")
	sb.WriteString("<div class=\"summary-box\">\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>Conformance Level</th><th>Count</th></tr>\n")
	sb.WriteString(fmt.Sprintf("<tr><td class=\"supports\">Supports</td><td>%d</td></tr>\n", report.Summary.Supports))
	sb.WriteString(fmt.Sprintf("<tr><td class=\"partially-supports\">Partially Supports</td><td>%d</td></tr>\n", report.Summary.PartiallySupports))
	sb.WriteString(fmt.Sprintf("<tr><td class=\"does-not-support\">Does Not Support</td><td>%d</td></tr>\n", report.Summary.DoesNotSupport))
	sb.WriteString(fmt.Sprintf("<tr><td class=\"not-applicable\">Not Applicable</td><td>%d</td></tr>\n", report.Summary.NotApplicable))
	sb.WriteString(fmt.Sprintf("<tr><td class=\"not-evaluated\">Not Evaluated</td><td>%d</td></tr>\n", report.Summary.NotEvaluated))
	sb.WriteString(fmt.Sprintf("<tr><th>Total</th><th>%d</th></tr>\n", report.Summary.TotalCriteria))
	sb.WriteString("</table>\n")
	sb.WriteString(fmt.Sprintf("<p><strong>Automated Coverage:</strong> %.1f%%</p>\n", report.Summary.AutomatedCoverage))
	sb.WriteString(fmt.Sprintf("<p><strong>Total Violations Found:</strong> %d</p>\n", report.Summary.TotalViolations))
	sb.WriteString("</div>\n")

	// Detailed Results
	sb.WriteString("<h2>Detailed Results</h2>\n")

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
		sb.WriteString(fmt.Sprintf("<h3>Principle %s</h3>\n", html.EscapeString(p.name)))
		sb.WriteString("<table>\n")
		sb.WriteString("<tr><th>Criteria</th><th>Conformance Level</th><th>Remarks</th></tr>\n")

		for _, c := range report.Criteria {
			if strings.HasPrefix(c.ID, p.prefix) {
				conformanceClass := conformanceClass(c.Conformance)
				remarks := html.EscapeString(c.Remarks)
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
				sb.WriteString(fmt.Sprintf("<tr><td>%s %s</td><td class=\"%s\">%s</td><td>%s</td></tr>\n",
					html.EscapeString(c.ID), html.EscapeString(c.Name), conformanceClass, c.Conformance, remarks))
			}
		}
		sb.WriteString("</table>\n")
	}

	// Violations Detail
	if report.Summary.TotalViolations > 0 {
		sb.WriteString("<h2>Violations Detail</h2>\n")

		for _, c := range report.Criteria {
			if len(c.Violations) == 0 {
				continue
			}

			sb.WriteString(fmt.Sprintf("<h3>%s %s</h3>\n", html.EscapeString(c.ID), html.EscapeString(c.Name)))

			for _, v := range c.Violations {
				sb.WriteString("<div class=\"violation\">\n")
				sb.WriteString(fmt.Sprintf("<p><strong>%s</strong> - %s</p>\n",
					html.EscapeString(v.RuleID), html.EscapeString(v.Description)))
				sb.WriteString("<ul>\n")
				sb.WriteString(fmt.Sprintf("<li>Impact: %s</li>\n", html.EscapeString(v.Impact)))
				sb.WriteString(fmt.Sprintf("<li>Instances: %d</li>\n", v.Count))
				if v.HelpURL != "" {
					sb.WriteString(fmt.Sprintf("<li>More info: <a href=\"%s\">%s</a></li>\n",
						html.EscapeString(v.HelpURL), html.EscapeString(v.HelpURL)))
				}
				sb.WriteString("</ul>\n")
				if len(v.Elements) > 0 {
					sb.WriteString("<p>Example elements:</p>\n")
					for _, el := range v.Elements {
						sb.WriteString(fmt.Sprintf("<pre><code>%s</code></pre>\n", html.EscapeString(truncate(el, 200))))
					}
				}
				sb.WriteString("</div>\n")
			}
		}
	}

	// Notes
	if report.Notes != "" {
		sb.WriteString("<div class=\"note\">\n")
		sb.WriteString("<h2>Notes</h2>\n")
		sb.WriteString(fmt.Sprintf("<p>%s</p>\n", html.EscapeString(report.Notes)))
		sb.WriteString("</div>\n")
	}

	// Footer
	sb.WriteString("<footer>\n")
	sb.WriteString(fmt.Sprintf("<p>Generated: %s</p>\n", report.GeneratedAt.Format("2006-01-02 15:04:05 MST")))
	sb.WriteString("</footer>\n")

	sb.WriteString("</body>\n</html>\n")

	return sb.String()
}

func conformanceClass(c vpat.Conformance) string {
	switch c {
	case vpat.ConformanceSupports:
		return "supports"
	case vpat.ConformancePartiallySupports:
		return "partially-supports"
	case vpat.ConformanceDoesNotSupport:
		return "does-not-support"
	case vpat.ConformanceNotApplicable:
		return "not-applicable"
	case vpat.ConformanceNotEvaluated:
		return "not-evaluated"
	default:
		return ""
	}
}
