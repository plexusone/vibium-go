// Package vpat provides VPAT (Voluntary Product Accessibility Template) report generation.
package vpat

import (
	"time"
)

// Report represents a complete VPAT accessibility conformance report.
type Report struct {
	// Product contains information about the product being evaluated.
	Product ProductInfo `json:"product" jsonschema:"description=Information about the product being evaluated"`

	// Evaluation contains information about the evaluation process.
	Evaluation EvaluationInfo `json:"evaluation" jsonschema:"description=Information about the evaluation process"`

	// Standard is the accessibility standard used (e.g., WCAG 2.2 AA).
	Standard string `json:"standard" jsonschema:"description=Accessibility standard (e.g. WCAG 2.2 Level AA)"`

	// Criteria contains results for each success criterion.
	Criteria []CriterionResult `json:"criteria" jsonschema:"description=Results for each success criterion"`

	// Summary provides aggregate statistics.
	Summary Summary `json:"summary" jsonschema:"description=Aggregate conformance statistics"`

	// Notes contains additional notes or disclaimers.
	Notes string `json:"notes,omitempty" jsonschema:"description=Additional notes or disclaimers"`

	// GeneratedAt is when the report was generated.
	GeneratedAt time.Time `json:"generatedAt" jsonschema:"description=Report generation timestamp"`
}

// ProductInfo describes the product being evaluated.
type ProductInfo struct {
	// Name is the product name.
	Name string `json:"name" jsonschema:"description=Product name"`

	// Version is the product version.
	Version string `json:"version,omitempty" jsonschema:"description=Product version"`

	// Description provides a brief description of the product.
	Description string `json:"description,omitempty" jsonschema:"description=Brief product description"`

	// Vendor is the company or organization that created the product.
	Vendor string `json:"vendor,omitempty" jsonschema:"description=Vendor or organization name"`

	// URL is the product's website or documentation URL.
	URL string `json:"url,omitempty" jsonschema:"description=Product URL"`
}

// EvaluationInfo describes how the evaluation was conducted.
type EvaluationInfo struct {
	// Date is when the evaluation was performed.
	Date time.Time `json:"date" jsonschema:"description=Evaluation date"`

	// Evaluator is who performed the evaluation.
	Evaluator string `json:"evaluator,omitempty" jsonschema:"description=Person or organization performing evaluation"`

	// Methods describes the evaluation methods used.
	Methods []string `json:"methods" jsonschema:"description=Evaluation methods used (e.g. Automated testing or Manual testing)"`

	// Tools lists the tools used for evaluation.
	Tools []ToolInfo `json:"tools,omitempty" jsonschema:"description=Tools used for evaluation"`

	// URLs lists the URLs that were evaluated.
	URLs []string `json:"urls,omitempty" jsonschema:"description=URLs that were evaluated"`

	// Scope describes what was included/excluded from evaluation.
	Scope string `json:"scope,omitempty" jsonschema:"description=Evaluation scope description"`
}

// ToolInfo describes an evaluation tool.
type ToolInfo struct {
	// Name is the tool name.
	Name string `json:"name" jsonschema:"description=Tool name"`

	// Version is the tool version.
	Version string `json:"version,omitempty" jsonschema:"description=Tool version"`
}

// CriterionResult contains the evaluation result for a single WCAG criterion.
type CriterionResult struct {
	// ID is the criterion identifier (e.g., "1.1.1").
	ID string `json:"id" jsonschema:"description=WCAG criterion ID (e.g. 1.1.1)"`

	// Name is the criterion name (e.g., "Non-text Content").
	Name string `json:"name" jsonschema:"description=Criterion name"`

	// Level is the WCAG level (A, AA, or AAA).
	Level string `json:"level" jsonschema:"description=WCAG level,enum=A,enum=AA,enum=AAA"`

	// Conformance is the conformance level achieved.
	Conformance Conformance `json:"conformance" jsonschema:"description=Conformance level achieved"`

	// EvaluationMethod indicates how this criterion was evaluated.
	EvaluationMethod EvaluationMethod `json:"evaluationMethod" jsonschema:"description=How this criterion was evaluated"`

	// Remarks provides additional explanation.
	Remarks string `json:"remarks,omitempty" jsonschema:"description=Additional explanation or context"`

	// Violations lists specific violations found.
	Violations []Violation `json:"violations,omitempty" jsonschema:"description=Specific violations found"`

	// AxeRules lists the axe-core rules that map to this criterion.
	AxeRules []string `json:"axeRules,omitempty" jsonschema:"description=Axe-core rules mapped to this criterion"`
}

// Conformance represents the conformance level for a criterion.
type Conformance string

const (
	// ConformanceSupports means the product fully supports this criterion.
	ConformanceSupports Conformance = "Supports"

	// ConformancePartiallySupports means the product partially supports this criterion.
	ConformancePartiallySupports Conformance = "Partially Supports"

	// ConformanceDoesNotSupport means the product does not support this criterion.
	ConformanceDoesNotSupport Conformance = "Does Not Support"

	// ConformanceNotApplicable means this criterion is not applicable to the product.
	ConformanceNotApplicable Conformance = "Not Applicable"

	// ConformanceNotEvaluated means this criterion was not evaluated.
	ConformanceNotEvaluated Conformance = "Not Evaluated"
)

// EvaluationMethod indicates how a criterion was evaluated.
type EvaluationMethod string

const (
	// MethodAutomated means the criterion was evaluated using automated testing.
	MethodAutomated EvaluationMethod = "Automated"

	// MethodManual means the criterion was evaluated using manual testing.
	MethodManual EvaluationMethod = "Manual"

	// MethodHybrid means the criterion was evaluated using both methods.
	MethodHybrid EvaluationMethod = "Hybrid"

	// MethodNotTested means the criterion was not tested.
	MethodNotTested EvaluationMethod = "Not Tested"
)

// Violation represents a specific accessibility violation.
type Violation struct {
	// RuleID is the axe-core rule identifier.
	RuleID string `json:"ruleId,omitempty" jsonschema:"description=Axe-core rule ID"`

	// Description describes the violation.
	Description string `json:"description" jsonschema:"description=Violation description"`

	// Impact is the severity of the violation.
	Impact string `json:"impact,omitempty" jsonschema:"description=Violation impact (critical or serious or moderate or minor)"`

	// Count is the number of instances found.
	Count int `json:"count" jsonschema:"description=Number of instances found"`

	// Elements lists example affected elements (HTML snippets).
	Elements []string `json:"elements,omitempty" jsonschema:"description=Example affected elements (HTML)"`

	// HelpURL is a link to more information about the violation.
	HelpURL string `json:"helpUrl,omitempty" jsonschema:"description=URL with more information"`
}

// Summary provides aggregate conformance statistics.
type Summary struct {
	// TotalCriteria is the total number of criteria evaluated.
	TotalCriteria int `json:"totalCriteria" jsonschema:"description=Total criteria in scope"`

	// Supports is the count of criteria that fully support accessibility.
	Supports int `json:"supports" jsonschema:"description=Criteria fully supported"`

	// PartiallySupports is the count of criteria with partial support.
	PartiallySupports int `json:"partiallySupports" jsonschema:"description=Criteria partially supported"`

	// DoesNotSupport is the count of criteria that are not supported.
	DoesNotSupport int `json:"doesNotSupport" jsonschema:"description=Criteria not supported"`

	// NotApplicable is the count of criteria that are not applicable.
	NotApplicable int `json:"notApplicable" jsonschema:"description=Criteria not applicable"`

	// NotEvaluated is the count of criteria that were not evaluated.
	NotEvaluated int `json:"notEvaluated" jsonschema:"description=Criteria not evaluated"`

	// AutomatedCoverage is the percentage of criteria covered by automated testing.
	AutomatedCoverage float64 `json:"automatedCoverage" jsonschema:"description=Percentage of criteria covered by automated testing"`

	// TotalViolations is the total number of violations found.
	TotalViolations int `json:"totalViolations" jsonschema:"description=Total violations found"`
}

// CalculateSummary computes the summary from criteria results.
func (r *Report) CalculateSummary() {
	r.Summary = Summary{
		TotalCriteria: len(r.Criteria),
	}

	automatedCount := 0
	for _, c := range r.Criteria {
		switch c.Conformance {
		case ConformanceSupports:
			r.Summary.Supports++
		case ConformancePartiallySupports:
			r.Summary.PartiallySupports++
		case ConformanceDoesNotSupport:
			r.Summary.DoesNotSupport++
		case ConformanceNotApplicable:
			r.Summary.NotApplicable++
		case ConformanceNotEvaluated:
			r.Summary.NotEvaluated++
		}

		if c.EvaluationMethod == MethodAutomated || c.EvaluationMethod == MethodHybrid {
			automatedCount++
		}

		for _, v := range c.Violations {
			r.Summary.TotalViolations += v.Count
		}
	}

	if r.Summary.TotalCriteria > 0 {
		r.Summary.AutomatedCoverage = float64(automatedCount) / float64(r.Summary.TotalCriteria) * 100
	}
}
