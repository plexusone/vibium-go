package vpat

import (
	"fmt"
	"strings"
	"time"

	"github.com/agentplexus/vibium-go/a11y"
	"github.com/agentplexus/vibium-go/vpat/criteria"
)

// Generator creates VPAT reports from accessibility test results.
type Generator struct {
	product    ProductInfo
	evaluation EvaluationInfo
	criteria   []criteria.Criterion
}

// NewGenerator creates a new VPAT generator for WCAG 2.2 AA.
func NewGenerator(product ProductInfo) *Generator {
	return &Generator{
		product:  product,
		criteria: criteria.WCAG22AA(),
		evaluation: EvaluationInfo{
			Date:    time.Now(),
			Methods: []string{"Automated testing"},
			Tools: []ToolInfo{
				{Name: "axe-core", Version: "4.8.4"},
				{Name: "Vibium", Version: "0.2.0"},
			},
		},
	}
}

// SetEvaluator sets the evaluator information.
func (g *Generator) SetEvaluator(evaluator string) {
	g.evaluation.Evaluator = evaluator
}

// SetScope sets the evaluation scope.
func (g *Generator) SetScope(scope string) {
	g.evaluation.Scope = scope
}

// AddURL adds a URL that was evaluated.
func (g *Generator) AddURL(url string) {
	g.evaluation.URLs = append(g.evaluation.URLs, url)
}

// Generate creates a VPAT report from axe-core results.
func (g *Generator) Generate(results []*a11y.Result) *Report {
	report := &Report{
		Product:     g.product,
		Evaluation:  g.evaluation,
		Standard:    "WCAG 2.2 Level AA",
		GeneratedAt: time.Now(),
		Notes: "This report was generated using automated accessibility testing tools. " +
			"Automated testing can only detect approximately 30-40% of accessibility issues. " +
			"A complete accessibility evaluation requires manual testing by accessibility experts.",
	}

	// Collect URLs from results
	for _, r := range results {
		if r.URL != "" && !containsString(report.Evaluation.URLs, r.URL) {
			report.Evaluation.URLs = append(report.Evaluation.URLs, r.URL)
		}
	}

	// Build a map of rule ID -> violations
	ruleViolations := make(map[string][]a11y.Violation)
	for _, r := range results {
		for _, v := range r.Violations {
			ruleViolations[v.ID] = append(ruleViolations[v.ID], v)
		}
	}

	// Evaluate each criterion
	for _, crit := range g.criteria {
		result := g.evaluateCriterion(crit, ruleViolations)
		report.Criteria = append(report.Criteria, result)
	}

	report.CalculateSummary()
	return report
}

// evaluateCriterion evaluates a single criterion based on axe-core results.
func (g *Generator) evaluateCriterion(crit criteria.Criterion, ruleViolations map[string][]a11y.Violation) CriterionResult {
	result := CriterionResult{
		ID:       crit.ID,
		Name:     crit.Name,
		Level:    crit.Level,
		AxeRules: crit.AxeRules,
	}

	// If no axe rules map to this criterion, mark as not evaluated
	if len(crit.AxeRules) == 0 {
		result.Conformance = ConformanceNotEvaluated
		result.EvaluationMethod = MethodNotTested
		result.Remarks = "Requires manual testing"
		return result
	}

	// Determine evaluation method
	if crit.CanAutomate {
		result.EvaluationMethod = MethodAutomated
	} else {
		result.EvaluationMethod = MethodHybrid
		result.Remarks = "Partial automation; manual testing recommended"
	}

	// Check for violations in mapped rules
	var allViolations []Violation
	for _, ruleID := range crit.AxeRules {
		if violations, ok := ruleViolations[ruleID]; ok {
			for _, v := range violations {
				viol := Violation{
					RuleID:      v.ID,
					Description: v.Help,
					Impact:      string(v.Impact),
					Count:       len(v.Nodes),
					HelpURL:     v.HelpURL,
				}
				// Add example elements
				for i, n := range v.Nodes {
					if i >= 3 {
						break
					}
					viol.Elements = append(viol.Elements, n.HTML)
				}
				allViolations = append(allViolations, viol)
			}
		}
	}

	// Determine conformance based on violations
	if len(allViolations) == 0 {
		result.Conformance = ConformanceSupports
		if !crit.CanAutomate {
			result.Remarks = "No automated violations detected; manual testing required for full evaluation"
		}
	} else {
		result.Violations = allViolations

		// Count total issues and determine severity
		totalIssues := 0
		hasCritical := false
		hasSerious := false
		for _, v := range allViolations {
			totalIssues += v.Count
			if v.Impact == string(a11y.ImpactCritical) {
				hasCritical = true
			}
			if v.Impact == string(a11y.ImpactSerious) {
				hasSerious = true
			}
		}

		if hasCritical || (hasSerious && totalIssues > 5) {
			result.Conformance = ConformanceDoesNotSupport
		} else {
			result.Conformance = ConformancePartiallySupports
		}

		// Build remarks
		var remarks []string
		remarks = append(remarks, fmt.Sprintf("%d instance(s) detected", totalIssues))
		if !crit.CanAutomate {
			remarks = append(remarks, "manual testing also required")
		}
		result.Remarks = strings.Join(remarks, "; ")
	}

	return result
}

// GenerateFromSingleResult creates a report from a single axe-core result.
func (g *Generator) GenerateFromSingleResult(result *a11y.Result) *Report {
	return g.Generate([]*a11y.Result{result})
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
