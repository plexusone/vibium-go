package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	w3pilot "github.com/plexusone/w3pilot"
)

var (
	assertTimeout  time.Duration
	assertSelector string
)

// AssertResult represents the result of an assertion.
type AssertResult struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

// assertTextCmd asserts that text exists on the page
var assertTextCmd = &cobra.Command{
	Use:   "assert-text <text>",
	Short: "Assert text exists on page",
	Long: `Assert that the specified text exists on the page.

By default, searches the entire page. Use --selector to scope the search.

Examples:
  w3pilot test assert-text "Welcome"
  w3pilot test assert-text "Login" --selector "#header"
  w3pilot test assert-text "Success" --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), assertTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		opts := &w3pilot.AssertOptions{
			Timeout:  assertTimeout,
			Selector: assertSelector,
		}

		err := pilot.AssertText(ctx, text, opts)

		result := AssertResult{Passed: err == nil}
		if err != nil {
			var aErr *w3pilot.AssertionError
			if errors.As(err, &aErr) {
				result.Message = aErr.Message
			} else {
				result.Message = err.Error()
			}
		} else {
			result.Message = fmt.Sprintf("Text found: %q", text)
		}

		Output(result, func(data interface{}) string {
			r := data.(AssertResult)
			if r.Passed {
				return fmt.Sprintf("PASS: %s", r.Message)
			}
			return fmt.Sprintf("FAIL: %s", r.Message)
		})

		if !result.Passed {
			return fmt.Errorf("assertion failed")
		}
		return nil
	},
}

// assertElementCmd asserts that an element exists
var assertElementCmd = &cobra.Command{
	Use:   "assert-element <selector>",
	Short: "Assert element exists",
	Long: `Assert that an element matching the selector exists on the page.

Examples:
  w3pilot test assert-element "#login-button"
  w3pilot test assert-element ".submit-form"
  w3pilot test assert-element "[data-testid=submit]" --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), assertTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		opts := &w3pilot.AssertOptions{
			Timeout: assertTimeout,
		}

		err := pilot.AssertElement(ctx, selector, opts)

		result := AssertResult{Passed: err == nil}
		if err != nil {
			var aErr *w3pilot.AssertionError
			if errors.As(err, &aErr) {
				result.Message = aErr.Message
			} else {
				result.Message = err.Error()
			}
		} else {
			result.Message = fmt.Sprintf("Element found: %s", selector)
		}

		Output(result, func(data interface{}) string {
			r := data.(AssertResult)
			if r.Passed {
				return fmt.Sprintf("PASS: %s", r.Message)
			}
			return fmt.Sprintf("FAIL: %s", r.Message)
		})

		if !result.Passed {
			return fmt.Errorf("assertion failed")
		}
		return nil
	},
}

// assertURLCmd asserts that the URL matches a pattern
var assertURLCmd = &cobra.Command{
	Use:   "assert-url <pattern>",
	Short: "Assert URL matches pattern",
	Long: `Assert that the current URL matches the specified pattern.

The pattern can be:
  - Exact URL: "https://example.com/page"
  - Glob pattern: "https://example.com/users/*" or "**/dashboard"
  - Regex (wrapped in /): "/https://example\.com/users/\d+/"

Examples:
  w3pilot test assert-url "https://example.com/dashboard"
  w3pilot test assert-url "**/login"
  w3pilot test assert-url "/.*\/users\/\d+/"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), assertTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		err := pilot.AssertURL(ctx, pattern, nil)

		result := AssertResult{Passed: err == nil}
		if err != nil {
			var aErr *w3pilot.AssertionError
			if errors.As(err, &aErr) {
				result.Message = aErr.Message
			} else {
				result.Message = err.Error()
			}
		} else {
			currentURL, _ := pilot.URL(ctx)
			result.Message = fmt.Sprintf("URL matches: %s", currentURL)
		}

		Output(result, func(data interface{}) string {
			r := data.(AssertResult)
			if r.Passed {
				return fmt.Sprintf("PASS: %s", r.Message)
			}
			return fmt.Sprintf("FAIL: %s", r.Message)
		})

		if !result.Passed {
			return fmt.Errorf("assertion failed")
		}
		return nil
	},
}

func init() {
	testCmd.AddCommand(assertTextCmd)
	testCmd.AddCommand(assertElementCmd)
	testCmd.AddCommand(assertURLCmd)

	// Common flags
	assertTextCmd.Flags().DurationVar(&assertTimeout, "timeout", 5*time.Second, "Timeout")
	assertTextCmd.Flags().StringVar(&assertSelector, "selector", "", "Scope search to selector")

	assertElementCmd.Flags().DurationVar(&assertTimeout, "timeout", 5*time.Second, "Timeout")

	assertURLCmd.Flags().DurationVar(&assertTimeout, "timeout", 5*time.Second, "Timeout")
}
