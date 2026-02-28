package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/agentplexus/vibium-go/a11y"
	"github.com/agentplexus/vibium-go/vpat"
	"github.com/agentplexus/vibium-go/vpat/render"
	"github.com/spf13/cobra"
)

var (
	vpatFormat    string
	vpatOutput    string
	vpatProduct   string
	vpatVersion   string
	vpatVendor    string
	vpatEvaluator string
	vpatScope     string
	vpatStandard  string
	vpatTimeout   time.Duration
)

var vpatCmd = &cobra.Command{
	Use:   "vpat <url>...",
	Short: "Generate a VPAT accessibility conformance report",
	Long: `Generate a VPAT (Voluntary Product Accessibility Template) report by
running automated accessibility checks against one or more URLs.

The report maps axe-core findings to WCAG 2.2 AA success criteria and
outputs in the specified format.

Output formats:
  json     - JSON intermediate representation
  markdown - Markdown (ITI VPAT format)
  html     - HTML (ITI VPAT format)
  csv      - CSV for spreadsheet import

Examples:
  # Generate HTML VPAT for a website
  vibium vpat https://example.com --product "Example Site" --format html -o vpat.html

  # Check multiple pages and generate Markdown report
  vibium vpat https://example.com https://example.com/about --format markdown

  # Generate JSON for further processing
  vibium vpat https://example.com --format json -o vpat.json`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		urls := args

		ctx, cancel := context.WithTimeout(context.Background(), vpatTimeout)
		defer cancel()

		// Launch browser
		vibe, err := launchBrowser(ctx, true) // Always headless for VPAT
		if err != nil {
			return err
		}
		defer func() {
			_ = vibe.Quit(context.Background())
			_ = clearSession()
		}()

		// Collect accessibility results from all URLs
		var results []*a11y.Result
		for _, url := range urls {
			fmt.Printf("Checking: %s\n", url)

			if err := vibe.Go(ctx, url); err != nil {
				return fmt.Errorf("failed to navigate to %s: %w", url, err)
			}

			// Wait for page to load
			if err := vibe.WaitForLoad(ctx, "networkidle", 30*time.Second); err != nil {
				// Non-fatal - continue with check
				fmt.Printf("  Warning: page may not be fully loaded\n")
			}

			opts := &a11y.Options{
				Standard: a11y.Standard(vpatStandard),
			}

			result, err := a11y.Check(ctx, vibe, opts)
			if err != nil {
				return fmt.Errorf("accessibility check failed for %s: %w", url, err)
			}

			results = append(results, result)
			fmt.Printf("  Found %d violations, %d passes\n", len(result.Violations), len(result.Passes))
		}

		// Generate VPAT report
		productInfo := vpat.ProductInfo{
			Name:    vpatProduct,
			Version: vpatVersion,
			Vendor:  vpatVendor,
		}
		if productInfo.Name == "" {
			// Use first URL as product name
			productInfo.Name = urls[0]
		}

		generator := vpat.NewGenerator(productInfo)
		if vpatEvaluator != "" {
			generator.SetEvaluator(vpatEvaluator)
		}
		if vpatScope != "" {
			generator.SetScope(vpatScope)
		}

		report := generator.Generate(results)

		// Render output
		var output string
		switch strings.ToLower(vpatFormat) {
		case "json":
			data, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			output = string(data)
		case "markdown", "md":
			output = render.Markdown(report)
		case "html":
			output = render.HTML(report)
		case "csv":
			var err error
			output, err = render.CSV(report)
			if err != nil {
				return fmt.Errorf("failed to render CSV: %w", err)
			}
		default:
			return fmt.Errorf("unknown format: %s (use json, markdown, html, or csv)", vpatFormat)
		}

		// Write output
		if vpatOutput != "" {
			if err := os.WriteFile(vpatOutput, []byte(output), 0600); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			fmt.Printf("VPAT report written to: %s\n", vpatOutput)
		} else {
			fmt.Println(output)
		}

		// Print summary
		fmt.Printf("\nSummary:\n")
		fmt.Printf("  Supports:           %d\n", report.Summary.Supports)
		fmt.Printf("  Partially Supports: %d\n", report.Summary.PartiallySupports)
		fmt.Printf("  Does Not Support:   %d\n", report.Summary.DoesNotSupport)
		fmt.Printf("  Not Evaluated:      %d\n", report.Summary.NotEvaluated)
		fmt.Printf("  Automated Coverage: %.1f%%\n", report.Summary.AutomatedCoverage)
		fmt.Printf("  Total Violations:   %d\n", report.Summary.TotalViolations)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(vpatCmd)

	vpatCmd.Flags().StringVarP(&vpatFormat, "format", "f", "markdown", "Output format: json, markdown, html, csv")
	vpatCmd.Flags().StringVarP(&vpatOutput, "output", "o", "", "Output file (default: stdout)")
	vpatCmd.Flags().StringVar(&vpatProduct, "product", "", "Product name for the report")
	vpatCmd.Flags().StringVar(&vpatVersion, "version", "", "Product version")
	vpatCmd.Flags().StringVar(&vpatVendor, "vendor", "", "Vendor/organization name")
	vpatCmd.Flags().StringVar(&vpatEvaluator, "evaluator", "", "Evaluator name")
	vpatCmd.Flags().StringVar(&vpatScope, "scope", "", "Evaluation scope description")
	vpatCmd.Flags().StringVar(&vpatStandard, "standard", "wcag22aa", "WCAG standard: wcag2a, wcag2aa, wcag21aa, wcag22aa")
	vpatCmd.Flags().DurationVar(&vpatTimeout, "timeout", 10*time.Minute, "Total timeout for all checks")
}
