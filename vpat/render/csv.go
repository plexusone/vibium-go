package render

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/agentplexus/vibium-go/vpat"
)

// CSV renders a VPAT report as CSV for spreadsheet import.
func CSV(report *vpat.Report) (string, error) {
	var sb strings.Builder
	w := csv.NewWriter(&sb)

	// Header row
	header := []string{
		"Criteria ID",
		"Criteria Name",
		"Level",
		"Conformance",
		"Evaluation Method",
		"Remarks",
		"Violations Count",
		"Violation Rules",
	}
	if err := w.Write(header); err != nil {
		return "", err
	}

	// Data rows
	for _, c := range report.Criteria {
		violationCount := 0
		var violationRules []string
		for _, v := range c.Violations {
			violationCount += v.Count
			violationRules = append(violationRules, v.RuleID)
		}

		row := []string{
			c.ID,
			c.Name,
			c.Level,
			string(c.Conformance),
			string(c.EvaluationMethod),
			c.Remarks,
			fmt.Sprintf("%d", violationCount),
			strings.Join(violationRules, "; "),
		}
		if err := w.Write(row); err != nil {
			return "", err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

// CSVSummary renders just the summary as CSV.
func CSVSummary(report *vpat.Report) (string, error) {
	var sb strings.Builder
	w := csv.NewWriter(&sb)

	// Header
	if err := w.Write([]string{"Metric", "Value"}); err != nil {
		return "", err
	}

	// Data
	rows := [][]string{
		{"Product", report.Product.Name},
		{"Version", report.Product.Version},
		{"Standard", report.Standard},
		{"Evaluation Date", report.Evaluation.Date.Format("2006-01-02")},
		{"Total Criteria", fmt.Sprintf("%d", report.Summary.TotalCriteria)},
		{"Supports", fmt.Sprintf("%d", report.Summary.Supports)},
		{"Partially Supports", fmt.Sprintf("%d", report.Summary.PartiallySupports)},
		{"Does Not Support", fmt.Sprintf("%d", report.Summary.DoesNotSupport)},
		{"Not Applicable", fmt.Sprintf("%d", report.Summary.NotApplicable)},
		{"Not Evaluated", fmt.Sprintf("%d", report.Summary.NotEvaluated)},
		{"Automated Coverage (%)", fmt.Sprintf("%.1f", report.Summary.AutomatedCoverage)},
		{"Total Violations", fmt.Sprintf("%d", report.Summary.TotalViolations)},
	}

	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return "", err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

// CSVViolations renders just violations as CSV.
func CSVViolations(report *vpat.Report) (string, error) {
	var sb strings.Builder
	w := csv.NewWriter(&sb)

	// Header
	header := []string{
		"Criteria ID",
		"Criteria Name",
		"Rule ID",
		"Impact",
		"Count",
		"Description",
		"Help URL",
	}
	if err := w.Write(header); err != nil {
		return "", err
	}

	// Data
	for _, c := range report.Criteria {
		for _, v := range c.Violations {
			row := []string{
				c.ID,
				c.Name,
				v.RuleID,
				v.Impact,
				fmt.Sprintf("%d", v.Count),
				v.Description,
				v.HelpURL,
			}
			if err := w.Write(row); err != nil {
				return "", err
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}

	return sb.String(), nil
}
